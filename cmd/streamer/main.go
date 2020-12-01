package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/bennyharvey/soma/dlib"
	"github.com/bennyharvey/soma/gocv"
	"github.com/bennyharvey/soma/rmq"
	"github.com/bennyharvey/soma/streamer"
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

	bfd, err := dlib.NewBatchFaceDetector(c.DetectorModelPath, c.StreamsResolution, len(c.Streams),
		c.FaceDetectorWaitTime)
	if err != nil {
		logrus.WithError(err).Fatal("failed to create dlib_batch_face_detector")
	}
	defer func() {
		bfd.Close()
		logrus.Info("dlib_batch_face_detector closed")
	}()

	logrus.Info("dlib_batch_face_detector created")

	for i, sc := range c.Streams {
		log := logrus.WithFields(logrus.Fields{
			"stream_index": i,
			"passage_id":   sc.PassageID,
		})

		dfp := rmq.NewDetectedFacePublisher(c.RabbitMQURI, c.RabbitMQExchange, sc.PassageID)
		defer func() {
			dfp.Stop()
			log.Info("rmq_detected_face_publisher stopped")
		}()

		log.Info("rmq_detected_faced_publisher created and started")

		dfh := streamer.NewFrameHandler(bfd, dfp)

		log.Info("streamer_frame_handler created")

		s := gocv.NewStream(sc.PassageID, sc.URI, sc.ClosedDuration, sc.FrameRate, dfh)
		defer func() {
			s.Stop()
			log.Info("gocv_stream stopped")
		}()

		log.Info("gocv_stream created and started")
	}

	logrus.Info("started")

	signals := make(chan os.Signal)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	logrus.Infof("captured %v signal, stopping", <-signals)

	st = time.Now()
}
