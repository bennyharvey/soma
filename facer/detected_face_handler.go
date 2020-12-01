package facer

import (
	"time"

	"github.com/sirupsen/logrus"
	"gocv.io/x/gocv"

	"github.com/bennyharvey/soma/entity"
)

type FaceRecognizer interface {
	RecognizeFace(photo gocv.Mat) (entity.FaceDescriptor, error)
}

type RecognizedFacePublisher interface {
	PublishRecognizedFace(entity.RecognizedFace)
}

type DetectedFaceHandler struct {
	confidenceLimit         float64
	faceRecognizer          FaceRecognizer
	recognizedFacePublisher RecognizedFacePublisher
	log                     *logrus.Entry
}

func NewDetectedFaceHandler(confidenceLimit float64, fr FaceRecognizer, rfp RecognizedFacePublisher) *DetectedFaceHandler {
	return &DetectedFaceHandler{
		confidenceLimit:         confidenceLimit,
		faceRecognizer:          fr,
		recognizedFacePublisher: rfp,
		log:                     logrus.WithField("subsystem", "facer_detected_face_handler"),
	}
}

func (dfh *DetectedFaceHandler) HandleDetectedFace(df entity.DetectedFace) {
	if df.DetectConfidence < dfh.confidenceLimit {
		return
	}

	photo, err := gocv.IMDecode(df.Photo, gocv.IMReadUnchanged)
	if err != nil {
		dfh.log.WithError(err).Error("failed to decode photo")
		return
	}

	defer func() {
		err := photo.Close()
		if err != nil {
			dfh.log.WithError(err).Error("failed to close gocv.Mat")
		}
	}()

	descr, err := dfh.faceRecognizer.RecognizeFace(photo)
	if err != nil {
		dfh.log.WithError(err).Error("failed to recognize face")
		return
	}

	dfh.recognizedFacePublisher.PublishRecognizedFace(entity.RecognizedFace{
		DetectedFace:  df,
		Descriptor:    descr,
		RecognizeTime: time.Now(),
	})
}
