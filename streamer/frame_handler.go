package streamer

import (
	"time"

	"github.com/sirupsen/logrus"
	"gocv.io/x/gocv"

	"github.com/bennyharvey/soma/entity"
)

type FaceDetector interface {
	DetectFaces(gocv.Mat) (chan []entity.FaceDetection, error)
}

type DetectedFacePublisher interface {
	PublishDetectedFace(df entity.DetectedFace)
}

type FrameHandler struct {
	faceDetector          FaceDetector
	detectedFacePublisher DetectedFacePublisher
	log                   *logrus.Entry
}

func NewFrameHandler(fd FaceDetector, dfp DetectedFacePublisher) *FrameHandler {
	return &FrameHandler{
		faceDetector:          fd,
		detectedFacePublisher: dfp,
		log:                   logrus.WithField("subsystem", "streamer_frame_handler"),
	}
}

func (fh *FrameHandler) HandleFrame(frame gocv.Mat) {
	frameTime := time.Now()

	ds, err := fh.faceDetector.DetectFaces(frame)
	if err != nil {
		fh.log.WithError(err).Error("failed to detect faces")
		return
	}

	detectTime := time.Now()

	for _, d := range <-ds {
		if d.Rectangle.Min.X < 0 || d.Rectangle.Min.Y < 0 ||
			d.Rectangle.Max.X > frame.Cols() || d.Rectangle.Max.Y > frame.Rows() {
			continue
		}

		photo, err := gocv.IMEncode(gocv.JPEGFileExt, frame.Region(d.Rectangle))
		if err != nil {
			fh.log.WithError(err).Error("failed to encode photo")
			continue
		}

		fh.detectedFacePublisher.PublishDetectedFace(entity.DetectedFace{
			Photo:            photo,
			DetectConfidence: d.Confidence,
			DetectTime:       detectTime,
			FrameTime:        frameTime,
		})
	}
}
