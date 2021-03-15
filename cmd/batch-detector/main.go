package main

import (
	"fmt"
	"image"
	"image/color"
	"os"
	"os/signal"
	"soma/face"

	"strconv"
	"syscall"
	"time"
	"github.com/sirupsen/logrus"
	"gocv.io/x/gocv"

	//"git.tattelecom.ru/dimuls/face"

	zurabiyGocv "github.com/bennyharvey/soma/gocv"
)

type frameHandler struct {
	window   *gocv.Window
	detector *face1.BatchDetector
	log      *logrus.Entry
}

func newFrameHandler(streamID string, d *face1.BatchDetector) *frameHandler {
	return &frameHandler{
		window:   gocv.NewWindow(streamID),
		detector: d,
		log:      logrus.WithField("stream_id", streamID),
	}
}

func (fh *frameHandler) close() error {
	return fh.window.Close()
}

func (fh *frameHandler) HandleFrame(frame gocv.Mat) {
	detections, errs, err := fh.detector.Detect(frame)
	if err != nil {
		fh.log.WithError(err).Error("failed to prepare detect")
		return
	}

	var (
		ds   []face1.Detection
		more bool
	)

	select {
	case ds, more = <-detections:
		if !more {
			break
		}
	case err, more = <-errs:
		if !more {
			break
		}
		fh.log.WithError(err).Error("failed to detect")
	}

	for _, d := range ds {
		if d.Rectangle.Min.X < 0 {
			d.Rectangle.Min.X = 0
		}
		if d.Rectangle.Min.Y < 0 {
			d.Rectangle.Min.Y = 0
		}
		if d.Rectangle.Max.X > 1280 {
			d.Rectangle.Max.X = 1280
		}
		if d.Rectangle.Max.Y > 720 {
			d.Rectangle.Max.Y = 720
		}
		gocv.Rectangle(&frame, d.Rectangle, color.RGBA{0, 0, 255, 0}, 3)

		if d.Rectangle.Min.Y-10 < 0 {
			continue
		}
		gocv.PutText(&frame, fmt.Sprintf("%0.2g", d.Confidence), image.Point{
			X: d.Rectangle.Min.X,
			Y: d.Rectangle.Min.Y - 10,
		}, gocv.FontHersheyPlain, 2, color.RGBA{0, 0, 255, 0}, 2)
	}

	fh.window.IMShow(frame)
	fh.window.WaitKey(1)
}

func main() {
	detector, err := face1.NewBatchDetector(os.Args[1], image.Point{X: 640, Y: 480}, len(os.Args)-2, 100*time.Millisecond)
	if err != nil {
		logrus.WithError(err).Fatal("failed to create detector")
	}

	defer detector.Close()

	var fhs []*frameHandler
	var ss []*zurabiyGocv.Stream

	for i, _ := range os.Args[2:] {
		id := strconv.Itoa(i)
		fh := newFrameHandler(id, detector)
		fhs = append(fhs, fh)
		//ss = append(ss, zurabiyGocv.NewStream(id, uri, 10*time.Second, 25, fh))
	}

	signals := make(chan os.Signal)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	<-signals

	for _, s := range ss {
		s.Stop()
	}

	for _, fh := range fhs {
		err := fh.close()
		if err != nil {
			logrus.WithError(err).Error("failed to close frame handler")
		}
	}
}
