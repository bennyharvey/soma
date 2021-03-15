package dlib

import (
	"sync"

	"github.com/sirupsen/logrus"
	"gocv.io/x/gocv"

	"github.com/bennyharvey/soma/face"

	"github.com/bennyharvey/soma/entity"
)

type FaceDetector struct {
	detector *face.Detector
	log      *logrus.Entry
	wg       sync.WaitGroup
}

func NewFaceDetector(modelPath string) (*FaceDetector, error) {
	d, err := face.NewDetector(modelPath)
	if err != nil {
		return nil, err
	}

	return &FaceDetector{detector: d, log: logrus.WithField("subsystem", "dlib_face_detector")}, nil
}

func (fd *FaceDetector) Close() {
	fd.wg.Wait()
	fd.detector.Close()
}

func (fd *FaceDetector) DetectFaces(img gocv.Mat) ([]entity.FaceDetection, error) {
	return fd.detector.Detect(img)
}
