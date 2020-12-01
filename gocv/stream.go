package gocv

import (
	"runtime"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"gocv.io/x/gocv"
)

type FrameHandler interface {
	HandleFrame(frame gocv.Mat)
}

type Stream struct {
	log  *logrus.Entry
	stop chan struct{}
	wg   sync.WaitGroup
}

func NewStream(passageID, uri string, closedDuration time.Duration, frameRate int, fh FrameHandler) *Stream {
	s := &Stream{
		log: logrus.WithFields(logrus.Fields{
			"subsystem":  "gocv_stream",
			"passage_id": passageID,
		}),
		stop: make(chan struct{}),
	}

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()

		runtime.LockOSThread()
		defer runtime.UnlockOSThread()

		var (
			err    error
			stream *gocv.VideoCapture

			successTime = time.Now()
			readTries   int

			frame = gocv.NewMat()

			frameDuration = time.Second / time.Duration(frameRate)
			framesToSkip  int
		)

		defer func() {
			err := frame.Close()
			if err != nil {
				s.log.WithError(err).Error("failed to close gocv.Mat")
			}

			if stream != nil {
				err = stream.Close()
				if err != nil {
					s.log.WithError(err).Error("failed to close stream")
				}
			}
		}()

	loop:
		for {
			select {
			case <-s.stop:
				break loop
			default:
			}

			if stream == nil {
				stream, err = gocv.OpenVideoCapture(uri)
				if err != nil {
					s.log.WithError(err).Error("failed to open stream")
					time.Sleep(3 * time.Second)
					continue
				}
			}

			if readTries > 0 && time.Now().Sub(successTime) > closedDuration {

				s.log.Warning("failed to read too long, reconnecting")

				err = stream.Close()
				if err != nil {
					s.log.WithError(err).Error("failed to close stream")
				}

				stream = nil
				successTime = time.Now()
				readTries = 0
				framesToSkip = 0

				continue
			}

			stream.Grab(framesToSkip)

			if !stream.Read(&frame) {
				readTries++
				time.Sleep(time.Second)
				continue
			}

			handleStartTime := time.Now()

			fh.HandleFrame(frame)

			framesToSkip = int(time.Now().Sub(handleStartTime) / frameDuration)

			successTime, readTries = time.Now(), 0
		}
	}()

	return s
}

func (s *Stream) Stop() {
	close(s.stop)
	s.wg.Wait()
}
