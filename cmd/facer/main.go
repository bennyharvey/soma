package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"../../facer"
	"../../rmq"

	"../../dlib"

	"github.com/sirupsen/logrus"
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
	
	// logrus.Info(configPath)
	
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

	fr, err := dlib.NewFaceRecognizer(c.ShaperModelPath, c.RecognizerModelPath, c.Jittering)
	if err != nil {
		logrus.WithError(err).Fatal("failed to create dlib_face_recognizer")
	}
	defer func() {
		fr.Close()
		logrus.Info("dlib_face_recognizer closed")
	}()

	logrus.Info("dlib_face_recognizer created")

	rfp := rmq.NewRecognizedFacePublisher(c.RabbitMQURI, c.RabbitMQExchange, c.PassageID)
	defer func() {
		rfp.Stop()
		logrus.Info("rmq_recognized_face_publisher stopped")
	}()

	logrus.Info("rmq_recognized_face_publisher created and started")

	dfh := facer.NewDetectedFaceHandler(c.ConfidenceLimit, fr, rfp)

	logrus.Info("facer_detected_face_handler created")

	dfc := rmq.NewDetectedFaceConsumer(c.RabbitMQURI, c.RabbitMQExchange, c.PassageID, dfh)
	defer func() {
		dfc.Stop()
		logrus.Info("rmq_recognized_face_publisher stopped")
	}()

	logrus.Info("rmq_recognized_face_publisher created and started")

	logrus.Info("started")

	signals := make(chan os.Signal)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	logrus.Infof("captured %v signal, stopping", <-signals)

	st = time.Now()
}
