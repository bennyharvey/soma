package dlib

import (
	"image"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"gocv.io/x/gocv"

	"git.tattelecom.ru/dimuls/face"

	"github.com/bennyharvey/soma/entity"
)

type BatchFaceDetector struct {
	detector *face.BatchDetector
	log      *logrus.Entry
	wg       sync.WaitGroup
}

func NewBatchFaceDetector(modelPath string, imageSize image.Point, batchSize int, waitTime time.Duration) (*BatchFaceDetector, error) {
	d, err := face.NewBatchDetector(modelPath, imageSize, batchSize, waitTime)
	if err != nil {
		return nil, err
	}

	return &BatchFaceDetector{detector: d, log: logrus.WithField("subsystem", "dlib_batch_face_detector")}, nil
}

func (fd *BatchFaceDetector) Close() {
	fd.wg.Wait()
	fd.detector.Close()
}

func (fd *BatchFaceDetector) DetectFaces(img gocv.Mat) (chan []entity.FaceDetection, error) {
	detects, detectErr, err := fd.detector.Detect(img)
	if err != nil {
		return nil, err
	}

	fd.wg.Add(1)
	go func() {
		defer fd.wg.Done()
		err := <-detectErr
		if err != nil {
			fd.log.WithError(err).Error("failed to detect")
		}
	}()

	return detects, nil
}
