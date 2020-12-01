//go:generate rice embed-go
package web

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"gocv.io/x/gocv"

	rice "github.com/GeertJohan/go.rice"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/sirupsen/logrus"

	"github.com/bennyharvey/soma/entity"
)

type DBStorage interface {
	User(login string) (entity.User, error)
	AddUser(entity.User) error
	SetUser(entity.User) error
	RemoveUser(login string) error
	Users() ([]entity.User, error)

	Person(personID int64) (entity.Person, error)
	AddPerson(entity.Person) (entity.Person, error)
	SetPerson(entity.Person) error
	RemovePerson(personID int64) error
	Persons() ([]entity.Person, error)

	PersonFaces(personID int64) ([]entity.PersonFace, error)
	AddPersonFace(entity.PersonFace) (entity.PersonFace, error)
	RemovePersonFace(personFaceID int64) error

	Events(...entity.EventsFilter) ([]entity.Event, error)
}

type PhotoStorage interface {
	AddPhoto(photoID string, photo []byte) error
	PhotoPath(photoID string) string
}

type FaceDetector interface {
	DetectFaces(photo gocv.Mat) ([]entity.FaceDetection, error)
}

type FaceRecognizer interface {
	RecognizeFace(face gocv.Mat) (entity.FaceDescriptor, error)
}

type ServerConfig struct {
	BindAddr       string            `yaml:"bind_addr"`
	JWTSigningKey  string            `yaml:"jwt_signing_key"`
	TLSCrtFilePath string            `yaml:"tls_crt_file_path"`
	TLSKeyFilePath string            `yaml:"tls_key_file_path"`
	Debug          bool              `yaml:"debug"`
	PassageNames   map[string]string `yaml:"passage_names"`
}

type Server struct {
	jwtSigningKey         string
	detectConfidenseLimit float64
	passageNames          map[string]string

	dbStorage      DBStorage
	photoStorage   PhotoStorage
	faceDetector   FaceDetector
	faceRecognizer FaceRecognizer

	echo *echo.Echo

	log  *logrus.Entry
	stop chan struct{}
	wg   sync.WaitGroup
}

func NewServer(bindAddr, jwtSigningKey, tlsCrtFilePath, tlsKeyFilePath string, detectConfidenseLimit float64,
	passageNames map[string]string, debug bool,
	dbs DBStorage, ps PhotoStorage, fd FaceDetector, fr FaceRecognizer) (*Server, error) {

	s := &Server{
		jwtSigningKey:         jwtSigningKey,
		detectConfidenseLimit: detectConfidenseLimit,
		passageNames:          passageNames,
		dbStorage:             dbs,
		photoStorage:          ps,
		faceDetector:          fd,
		faceRecognizer:        fr,
		log:                   logrus.WithField("subsystem", "web_server"),
	}

	s.stop = make(chan struct{})

	uiBox, err := rice.FindBox("ui/build")
	if err != nil {
		return nil, fmt.Errorf("rice.FindBox: %w", err)
	}

	indexHTML, err := uiBox.String("index.html")
	if err != nil {
		return nil, fmt.Errorf("uiBox.String index.html: %w", err)
	}

	uiHandler := echo.WrapHandler(http.FileServer(uiBox.HTTPBox()))

	e := echo.New()

	e.HideBanner = true
	e.HidePort = true
	e.HTTPErrorHandler = s.httpErrorHandler
	e.Debug = debug

	e.Use(middleware.Recover())
	e.Use(logrusLogger)

	corsCfg := middleware.DefaultCORSConfig
	corsCfg.AllowCredentials = true

	e.Use(middleware.CORSWithConfig(corsCfg))

	e.POST("/api/login", s.postLogin)

	aa := e.Group("/api", middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:  []byte(jwtSigningKey),
		Claims:      &jwtClaims{},
		TokenLookup: "cookie:auth",
		ErrorHandler: func(err error) error {
			return &echo.HTTPError{
				Code:    http.StatusUnauthorized,
				Message: err.Error(),
			}
		},
	}), func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tokenI := c.Get("user")
			if tokenI == nil {
				return echo.NewHTTPError(http.StatusUnauthorized)
			}

			token, ok := tokenI.(*jwt.Token)
			if !ok {
				return fmt.Errorf("expected %T but got %T", &jwt.Token{}, tokenI)
			}

			claims, ok := token.Claims.(*jwtClaims)
			if !ok {
				return fmt.Errorf("expected %T but got %T", &jwtClaims{}, token.Claims)
			}

			user, err := s.dbStorage.User(claims.Login)
			if err != nil {
				if err == entity.ErrUserNotFound {
					return echo.NewHTTPError(http.StatusUnauthorized)
				}
				return fmt.Errorf("dbStorage.User: %w", err)
			}

			c.Set("user", user)

			return next(c)
		}
	})

	adminOnly := withRoles(entity.Admin)
	adminWithSecurity := withRoles(entity.Security, entity.Admin)

	aa.GET("/users", s.getAPIUsers, adminOnly)
	aa.GET("/users/:user_login", s.getAPIUser, adminOnly)
	aa.POST("/users", s.postAPIUsers, adminOnly)
	aa.PUT("/users", s.putAPIUsers, adminOnly)
	aa.DELETE("/users/:user_login", s.deleteAPIUser, adminOnly)

	aa.GET("/photos/:photo_id", s.getAPIPhoto, adminWithSecurity)

	aa.GET("/persons", s.getAPIPersons, adminWithSecurity)
	aa.POST("/persons", s.postAPIPersons, adminWithSecurity)
	aa.PUT("/persons", s.putAPIPersons, adminWithSecurity)
	aa.GET("/persons/:person_id", s.getAPIPerson, adminWithSecurity)
	aa.DELETE("/persons/:person_id", s.deleteAPIPerson, adminWithSecurity)
	aa.GET("/persons/:person_id/faces", s.getAPIPersonFaces, adminWithSecurity)
	aa.POST("/persons/:person_id/faces", s.postAPIPersonFaces, adminWithSecurity)

	aa.DELETE("/person_faces/:person_face_id", s.deleteAPIPersonFace, adminWithSecurity)

	aa.GET("/passage_names", s.getAPIPassageNames, adminWithSecurity)

	aa.GET("/events", s.getAPIEvents, adminWithSecurity)

	e.GET("/", func(c echo.Context) error {
		return c.HTML(http.StatusOK, indexHTML)
	})

	e.GET("/*", func(c echo.Context) error {
		_, err := uiBox.Bytes(c.Request().URL.Path)
		if err == nil {
			return uiHandler(c)
		}
		return c.HTML(http.StatusOK, indexHTML)
	})

	s.echo = e

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		for {
			select {
			case <-s.stop:
				return
			default:
			}

			err := s.echo.StartTLS(bindAddr, tlsCrtFilePath, tlsKeyFilePath)
			if err != nil {
				if err == http.ErrServerClosed {
					return
				}
				s.log.WithError(err).Error("failed to start")
				time.Sleep(3 * time.Second)
			}
		}
	}()

	return s, nil
}

func (s *Server) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	close(s.stop)

	err := s.echo.Shutdown(ctx)
	if err != nil {
		s.log.WithError(err).Error("failed to graceful shutdown")
	}

	s.wg.Wait()
}

func (s *Server) httpErrorHandler(err error, c echo.Context) {
	var (
		code = http.StatusInternalServerError
		msg  interface{}
	)

	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		msg = he.Message
	} else if s.echo.Debug {
		msg = err.Error()
	} else {
		msg = http.StatusText(code)
	}
	if _, ok := msg.(string); !ok {
		msg = fmt.Sprintf("%v", msg)
	}

	// Send response
	if !c.Response().Committed {
		if c.Request().Method == http.MethodHead { // Issue #608
			err = c.NoContent(code)
		} else {
			err = c.String(code, msg.(string))
		}
		if err != nil {
			s.log.WithError(err).Error("failed to error response")
		}
	}
}

func withRoles(roles ...string) func(echo.HandlerFunc) echo.HandlerFunc {

	rolesMap := map[string]struct{}{}
	for _, r := range roles {
		rolesMap[r] = struct{}{}
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userI := c.Get("user")
			if userI == nil {
				return fmt.Errorf("expected not nil user in context")
			}

			user, ok := userI.(entity.User)
			if !ok {
				return fmt.Errorf("expected %T but got %T", entity.User{}, userI)
			}

			if _, exists := rolesMap[string(user.Role)]; exists {
				return next(c)
			}

			return echo.NewHTTPError(http.StatusForbidden)
		}
	}
}

func logrusLogger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()

		err := next(c)

		stop := time.Now()

		if err != nil {
			c.Error(err)
		}

		req := c.Request()
		res := c.Response()

		p := req.URL.Path
		if p == "" {
			p = "/"
		}

		bytesIn := req.Header.Get(echo.HeaderContentLength)
		if bytesIn == "" {
			bytesIn = "0"
		}

		entry := logrus.WithFields(map[string]interface{}{
			"subsystem":    "web_server",
			"remote_ip":    c.RealIP(),
			"host":         req.Host,
			"query_params": c.QueryParams(),
			"uri":          req.RequestURI,
			"method":       req.Method,
			"path":         p,
			"referer":      req.Referer(),
			"user_agent":   req.UserAgent(),
			"status":       res.Status,
			"latency":      stop.Sub(start).String(),
			"bytes_in":     bytesIn,
			"bytes_out":    strconv.FormatInt(res.Size, 10),
		})

		const msg = "request handled"

		if res.Status >= 500 {
			if err != nil {
				entry = entry.WithError(err)
			}
			entry.Error(msg)
		} else if res.Status >= 400 {
			if err != nil {
				entry = entry.WithError(err)
			}
			entry.Warn(msg)
		} else {
			entry.Info(msg)
		}

		return nil
	}
}

type jwtClaims struct {
	jwt.StandardClaims
	Login string
}
