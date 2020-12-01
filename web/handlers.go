package web

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"gocv.io/x/gocv"

	"github.com/dgrijalva/jwt-go"
	"github.com/iancoleman/strcase"
	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"

	"github.com/bennyharvey/soma/entity"
)

func (s *Server) postLogin(c echo.Context) error {
	var params struct {
		Login    string
		Password string
	}

	err := c.Bind(&params)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("bind user: %w", err))
	}

	u, err := s.dbStorage.User(params.Login)
	if err != nil {
		if err == entity.ErrUserNotFound {
			return echo.NewHTTPError(http.StatusUnauthorized)
		}
		return fmt.Errorf("dbStorage.User: %w", err)
	}

	err = bcrypt.CompareHashAndPassword(u.PasswordHash, []byte(params.Password))
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	claims := &jwtClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(12 * time.Hour).Unix(),
		},
		Login: u.Login,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(s.jwtSigningKey))
	if err != nil {
		return fmt.Errorf("token.SignedString: %w", err)
	}

	cookie := new(http.Cookie)
	cookie.Name = "auth"
	cookie.Value = signedToken
	cookie.Expires = time.Now().Add(12 * time.Hour)

	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, echo.Map{
		"user":  u,
		"token": signedToken,
	})
}

func (s *Server) getAPIUsers(c echo.Context) error {
	us, err := s.dbStorage.Users()
	if err != nil {
		return fmt.Errorf("dbStorage.Users: %w", err)
	}
	if us == nil {
		us = []entity.User{}
	}
	return c.JSON(http.StatusOK, us)
}

func (s *Server) getAPIUser(c echo.Context) error {
	login := c.Param("user_login")

	u, err := s.dbStorage.User(login)
	if err != nil {
		if err == entity.ErrUserNotFound {
			return echo.NewHTTPError(http.StatusNotFound, err)
		}
		return fmt.Errorf("dbStorage.User: %w", err)
	}

	return c.JSON(http.StatusOK, u)
}

func (s *Server) postAPIUsers(c echo.Context) error {
	var u entity.User

	err := c.Bind(&u)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("bind user: %w", err))
	}

	u.Login = strings.TrimSpace(u.Login)
	if u.Login == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "empty login")
	}

	switch u.Role {
	case entity.Admin, entity.Security:
	default:
		return echo.NewHTTPError(http.StatusBadRequest, "invalid role")
	}

	if len(u.Password) < 16 {
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("password length is less than 16"))
	}

	u.PasswordHash, err = bcrypt.GenerateFromPassword(
		[]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("bcrypt.GenerateFromPassword: %w", err)
	}

	err = s.dbStorage.AddUser(u)
	if err != nil {
		return fmt.Errorf("dbStorage.AddUser: %w", err)
	}

	return c.NoContent(http.StatusOK)
}

func (s *Server) putAPIUsers(c echo.Context) error {
	var u entity.User

	err := c.Bind(&u)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("bind user: %w", err))
	}

	u.Login = strings.TrimSpace(u.Login)
	if u.Login == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "empty login")
	}

	switch u.Role {
	case entity.Admin, entity.Security:
	default:
		return echo.NewHTTPError(http.StatusBadRequest, "invalid role")
	}

	if len(u.Password) > 0 {
		if len(u.Password) < 16 {
			return echo.NewHTTPError(http.StatusBadRequest, errors.New("password length is less than 16"))
		}
		u.PasswordHash, err = bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("bcrypt.GenerateFromPassword: %w", err)
		}
	}

	err = s.dbStorage.SetUser(u)
	if err != nil {
		return fmt.Errorf("dbStorage.SetUser: %w", err)
	}

	return c.NoContent(http.StatusOK)
}

func (s *Server) deleteAPIUser(c echo.Context) error {
	login := c.Param("user_login")

	err := s.dbStorage.RemoveUser(login)
	if err != nil {
		return fmt.Errorf("dbStorage.RemoveUser: %w", err)
	}

	return c.NoContent(http.StatusOK)
}

func (s *Server) getAPIPhoto(c echo.Context) error {
	return c.File(s.photoStorage.PhotoPath(c.Param("photo_id")))
}

func (s *Server) getAPIPersons(c echo.Context) error {
	ps, err := s.dbStorage.Persons()
	if err != nil {
		return fmt.Errorf("dbStorage.Persons: %w", err)
	}
	if ps == nil {
		ps = []entity.Person{}
	}
	return c.JSON(http.StatusOK, ps)
}

func (s *Server) postAPIPersons(c echo.Context) error {
	var p entity.Person

	err := c.Bind(&p)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("bind person: %w", err))
	}

	p, err = s.dbStorage.AddPerson(p)
	if err != nil {
		return fmt.Errorf("dbStorage.AddPerson: %w", err)
	}

	return c.JSON(http.StatusOK, p)
}

func (s *Server) putAPIPersons(c echo.Context) error {
	var p entity.Person

	err := c.Bind(&p)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("bind person: %w", err))
	}

	err = s.dbStorage.SetPerson(p)
	if err != nil {
		return fmt.Errorf("dbStorage.SetPerson: %w", err)
	}

	return c.NoContent(http.StatusOK)
}

func (s *Server) getAPIPerson(c echo.Context) error {
	personID, err := strconv.ParseInt(c.Param("person_id"),
		10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "parse person_id: "+err.Error())
	}

	ps, err := s.dbStorage.Person(personID)
	if err != nil {
		return fmt.Errorf("dbStorage.Person: %w", err)
	}

	return c.JSON(http.StatusOK, ps)
}

func (s *Server) deleteAPIPerson(c echo.Context) error {
	personID, err := strconv.ParseInt(c.Param("person_id"),
		10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "parse person_id: "+err.Error())
	}

	err = s.dbStorage.RemovePerson(personID)
	if err != nil {
		return fmt.Errorf("dbStorage.RemovePerson: %w", err)
	}

	return c.NoContent(http.StatusOK)
}

func (s *Server) getAPIPersonFaces(c echo.Context) error {
	personID, err := strconv.ParseInt(c.Param("person_id"),
		10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "parse person_id: "+err.Error())
	}

	pfs, err := s.dbStorage.PersonFaces(personID)
	if err != nil {
		return fmt.Errorf("dbStorage.PersonFaces: %w", err)
	}

	if pfs == nil {
		pfs = []entity.PersonFace{}
	}

	return c.JSON(http.StatusOK, pfs)
}

func (s *Server) postAPIPersonFaces(c echo.Context) error {
	personID, err := strconv.ParseInt(c.Param("person_id"),
		10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "parse person_id: "+err.Error())
	}

	photoBytes, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return fmt.Errorf("read all request body: %w", err)
	}

	photo, err := gocv.IMDecode(photoBytes, gocv.IMReadUnchanged)
	if err != nil {
		return fmt.Errorf("decode photo: %w", err)
	}

	defer func() {
		err := photo.Close()
		if err != nil {
			s.log.WithError(err).Error("failed to close gocv.Mat")
		}
	}()

	faceDets, err := s.faceDetector.DetectFaces(photo)
	if err != nil {
		return fmt.Errorf("face detect: %w", err)
	}

	var faceDet entity.FaceDetection

	switch len(faceDets) {
	case 0:
		return echo.NewHTTPError(http.StatusBadRequest, "no face detected on photo")
	case 1:
		faceDet = faceDets[0]
	default:
		maxSize := 0
		for _, fd := range faceDets {
			sizeP := fd.Rectangle.Size()
			size := sizeP.X + sizeP.Y
			if size > maxSize {
				maxSize = size
				faceDet = fd
			}
		}
	}

	if faceDet.Confidence < s.detectConfidenseLimit {
		return echo.NewHTTPError(http.StatusBadRequest, "face detect confidence too low")
	}

	if faceDet.Rectangle.Min.X < 0 {
		faceDet.Rectangle.Min.X = 0
	}

	if faceDet.Rectangle.Min.Y < 0 {
		faceDet.Rectangle.Min.Y = 0
	}

	if faceDet.Rectangle.Max.X > photo.Cols() {
		faceDet.Rectangle.Max.X = photo.Cols()
	}

	if faceDet.Rectangle.Max.Y > photo.Rows() {
		faceDet.Rectangle.Max.Y = photo.Rows()
	}

	descriptor, err := s.faceRecognizer.RecognizeFace(photo.Region(faceDet.Rectangle))
	if err != nil {
		return fmt.Errorf("faceRecongizer.RecognizeFace: %w", err)
	}

	photoMD5Sum := md5.Sum(photoBytes)
	photoID := hex.EncodeToString(photoMD5Sum[:])

	err = s.photoStorage.AddPhoto(photoID, photoBytes)
	if err != nil {
		return fmt.Errorf(
			"photoStorage.AddPhoto: %w", err)
	}

	pf, err := s.dbStorage.AddPersonFace(entity.PersonFace{
		PersonID:   personID,
		Descriptor: descriptor,
		PhotoID:    photoID,
	})
	if err != nil {
		return fmt.Errorf("dbStorage.AddPersonFace: %w", err)
	}

	return c.JSON(http.StatusOK, pf)
}

func (s *Server) deleteAPIPersonFace(c echo.Context) error {
	personFaceID, err := strconv.ParseInt(c.Param("person_face_id"),
		10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "parse person_id: "+err.Error())
	}

	err = s.dbStorage.RemovePersonFace(personFaceID)
	if err != nil {
		return fmt.Errorf("dbStorage.RemovePersonFace: %w", err)
	}

	return c.NoContent(http.StatusOK)
}

func (s *Server) getAPIPassageNames(c echo.Context) error {
	return c.JSON(http.StatusOK, s.passageNames)
}

func (s *Server) getAPIEvents(c echo.Context) error {
	var eventsFilters []entity.EventsFilter

	for key, values := range c.QueryParams() {
		switch key {
		case "from":
			from, err := time.Parse(time.RFC3339, values[0])
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "invalid from")
			}
			eventsFilters = append(eventsFilters, entity.EventsFrom(from))
		case "to":
			to, err := time.Parse(time.RFC3339, values[0])
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "invalid to")
			}
			eventsFilters = append(eventsFilters, entity.EventsTo(to.Add(time.Second-1)))
		case "passage_id":
			eventsFilters = append(eventsFilters, entity.EventsPassageID(values[0]))
		case "person_id":
			personID, err := strconv.ParseInt(values[0], 10, 64)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "invalid person_id")
			}
			eventsFilters = append(eventsFilters, entity.EventsPersonID(personID))
		case "person_name":
			eventsFilters = append(eventsFilters, entity.EventsPersonName(values[0]))
		case "order_by":
			eventsFilters = append(eventsFilters, entity.EventsOrderBy(values[0]))
		case "order_direction":
			eventsFilters = append(eventsFilters, entity.EventsOrderDirection(values[0]))
		case "limit":
			limit, err := strconv.Atoi(values[0])
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "invalid limit")
			}
			eventsFilters = append(eventsFilters, entity.EventsLimit(limit))
		case "offset":
			offset, err := strconv.Atoi(values[0])
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "invalid offset")
			}
			eventsFilters = append(eventsFilters, entity.EventsOffset(offset))
		}
	}

	e, err := s.dbStorage.Events(eventsFilters...)
	if err != nil {
		switch tErr := err.(type) {
		case entity.InvalidParamErr:
			return echo.NewHTTPError(http.StatusBadRequest, "invalid "+strcase.ToSnake(tErr.Param))
		default:
			return fmt.Errorf("dbStorage.Event: %w", err)
		}
	}

	if e == nil {
		e = []entity.Event{}
	}

	return c.JSON(http.StatusOK, e)
}
