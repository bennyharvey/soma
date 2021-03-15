package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/bennyharvey/soma/dlib"
	"github.com/bennyharvey/soma/entity"
	"github.com/bennyharvey/soma/file"
	"../../pg"
	"github.com/bennyharvey/soma/rmq"
	"github.com/bennyharvey/soma/sigur"
	"github.com/bennyharvey/soma/skuder"
	// "github.com/bennyharvey/soma/web"
	"../../web"
	"soma/z5r"
	"soma/beward"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)

	logrus.Info("starting")

	var st time.Time
	defer func() {
		logrus.WithField("shutdown_time", time.Now().Sub(st)).Info("stopped")
	}()

	var configPath string

	flag.StringVar(&configPath, "c", "", "config file path")
	flag.Parse()

	c, err := loadConfig(configPath)
	if err != nil {
		logrus.WithError(err).Fatal("failed to load config")
	}

	err = c.Validate()
	if err != nil {
		logrus.WithError(err).Fatal("invalid config")
	}

	logrus.Info("config loaded")

	if c.CUDAVisibleDevices != "" {
		err = os.Setenv("CUDA_VISIBLE_DEVICES", c.CUDAVisibleDevices)
		if err != nil {
			logrus.WithError(err).Fatal("failed to set CUDA_VISIBLE_DEVICES environment variable")
		}
		logrus.Info("environment variables set")
	}

	pgStorage, err := pg.NewStorage(c.PostgresURI)
	if err != nil {
		logrus.WithError(err).Fatal("failed to create pg_storage")
	}
	defer func() {
		err := pgStorage.Close()
		if err != nil {
			logrus.WithError(err).Error("failed to close pg_storage")
		} else {
			logrus.Info("pg_storage closed")
		}
	}()

	logrus.Info("pg_storage created")

	// err = pgStorage.Migrate()
	// if err != nil {
	// 	logrus.WithError(err).Fatal("failed to migrate pg_storage")
	// }

	// logrus.Info("pg_storage migrated")

	err = pgStorage.LoadPersonFaces()
	if err != nil {
		logrus.WithError(err).Fatal("failed to load person faces")
	}

	logrus.Info("person faces loaded")

	photoStorage := file.NewPhotoStorage(c.PhotoStoragePath)

	logrus.Info("photo_storage created")

	for passageID, poc := range c.PassageOpeners {
		log := logrus.WithField("passage_id", passageID)

		var po skuder.PassageOpener

		switch poc.Type {
			case entity.Sigur:
				po = sigur.NewPassageOpener(poc.Address, poc.Direction)
			case entity.Z5R:
				po = z5r.NewPassageOpener(poc.Address, poc.Direction)
			case entity.Beward:
				po = beward.NewPassageOpener(poc.Address, poc.Direction)
			case entity.Dummy:
				po = newDummyPassageOpener()
		}

		rfh := skuder.NewRecognizedFaceHandler(passageID, poc.WaitAfterOpen, c.DescriptorsMatchDistance,
			c.DetectConfidenceLimit, pgStorage, photoStorage, po)

		log.Info("recognized_face_handler created")

		rfc := rmq.NewRecognizedFaceConsumer(c.RabbitMQURI, c.RabbitMQExchange, passageID, rfh)
		defer func() {
			rfc.Stop()
			logrus.Info("recognized_face_consumer stopped")
		}()

		log.Info("recognized_face_consumer created and started")
	}

	fd, err := dlib.NewFaceDetector(c.DetectorModelPath)
	if err != nil {
		logrus.WithError(err).Error("failed to create dlib_face_detector")
	}
	defer func() {
		fd.Close()
		logrus.Info("dlib_face_detector closed")
	}()

	logrus.Info("dlib_face_detector created")

	fr, err := dlib.NewFaceRecognizer(c.ShaperModelPath, c.RecognizerModelPath, c.Jittering)
	if err != nil {
		logrus.WithError(err).Error("failed to create dlib_face_recognizer")
	}
	defer func() {
		fd.Close()
		logrus.Info("dlib_face_recognizer closed")
	}()

	logrus.Info("dlib_face_recognizer created")

	ws, err := web.NewServer(c.WebServer.BindAddr, c.WebServer.JWTSigningKey, c.WebServer.TLSCrtFilePath,
		c.WebServer.TLSKeyFilePath, c.DetectConfidenceLimit, c.WebServer.PassageNames, c.WebServer.Debug,
		pgStorage, photoStorage, fd, fr)
	if err != nil {
		logrus.WithError(err).Fatal("failed to create web_server")
	}
	defer func() {
		ws.Stop()
		logrus.Info("web_server stopped")
	}()

	logrus.Info("web_server created and started")

	logrus.Info("started")

	signals := make(chan os.Signal)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	logrus.Infof("captured %v signal, stopping", <-signals)

	st = time.Now()
}

type dummyPassageOpener struct {
	lastOpenTime time.Time
}

func newDummyPassageOpener() *dummyPassageOpener {
	return &dummyPassageOpener{}
}

func (po *dummyPassageOpener) OpenPassage() error {
	po.lastOpenTime = time.Now()
	return nil
}

func (po *dummyPassageOpener) LastOpenTime() time.Time {
	return po.lastOpenTime
}
