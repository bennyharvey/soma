package face1

// #cgo pkg-config: dlib-1 opencv4
// #cgo CXXFLAGS: -std=c++1z -Wall -O3 -DNDEBUG -march=native
// #cgo LDFLAGS: -ldlib -lcblas -lblas -llapack -lopencv_core
// #include <stdlib.h>
// #include "detector.h"
import "C"

import (
	"errors"
	"fmt"
	"image"
	"sync"
	"time"
	"unsafe"

	"gocv.io/x/gocv"
)

type BatchDetector struct {
	detector  unsafe.Pointer
	imageSize image.Point
	batchSize int
	waitTime  time.Duration

	images     []gocv.Mat
	detections []chan []Detection
	error      chan error
	hardDetect chan struct{}
	mx         sync.Mutex
}

func NewBatchDetector(modelPath string, imageSize image.Point, batchSize int, waitTime time.Duration) (*BatchDetector, error) {
	cModelPath := C.CString(modelPath)
	defer C.free(unsafe.Pointer(cModelPath))

	result := C.detector_init(cModelPath)
	defer C.free(unsafe.Pointer(result))

	if result.err_str != nil {
		defer C.free(unsafe.Pointer(result.err_str))
		return nil, errors.New(C.GoString(result.err_str))
	}

	return &BatchDetector{
		detector:  result.detector,
		imageSize: imageSize,
		batchSize: batchSize,
		waitTime:  waitTime,
	}, nil
}

func (bd *BatchDetector) Close() {
	C.detector_free(bd.detector)
	bd.detector = nil
}

func (bd *BatchDetector) Detect(img gocv.Mat) (chan []Detection, chan error, error) {
	bd.mx.Lock()
	defer bd.mx.Unlock()

	if img.Cols() != bd.imageSize.X && img.Rows() != bd.imageSize.Y {
		return nil, nil, fmt.Errorf("invalid image size %dx%d, expected %dx%d",
			img.Cols(), img.Rows(), bd.imageSize.X, bd.imageSize.Y)
	}

	ds := make(chan []Detection)

	bd.images = append(bd.images, img)
	bd.detections = append(bd.detections, ds)
	bd.error = make(chan error)

	switch len(bd.images) {
	case 1:
		bd.hardDetect = make(chan struct{})
		go bd.waitDetect()
	case bd.batchSize:
		close(bd.hardDetect)
	}

	return ds, bd.error, nil
}

func (bd *BatchDetector) waitDetect() {
	t := time.NewTimer(bd.waitTime)

	select {
	case <-bd.hardDetect:
		if !t.Stop() {
			<-t.C
		}
	case <-t.C:
		close(bd.hardDetect)
	}

	bd.hardDetect = nil

	bd.detect()
}

func (bd *BatchDetector) detect() {
	bd.mx.Lock()
	defer bd.mx.Unlock()

	defer func() {
		bd.images = nil

		for _, d := range bd.detections {
			close(d)
		}
		bd.detections = nil

		close(bd.error)
		bd.error = nil
	}()

	var images []unsafe.Pointer
	for _, i := range bd.images {
		images = append(images, unsafe.Pointer(i.Ptr()))
	}

	result := C.detector_batch_detect(bd.detector, unsafe.Pointer(&images[0]), C.int(len(images)))

	if result.err_str != nil {
		defer C.free(unsafe.Pointer(result.err_str))
		bd.error <- errors.New(C.GoString(result.err_str))
		return
	}

	defer C.free(unsafe.Pointer(result.detections))

	var batchDetectionsArray *C.BatchDetection = result.detections
	batchDetections := (*[1 << 28]C.BatchDetection)(unsafe.Pointer(batchDetectionsArray))[:len(images):len(images)]

	for i, batchDetection := range batchDetections {
		var ds []Detection

		detections := (*[1 << 28]C.Detection)(unsafe.Pointer(batchDetection.detections))[:batchDetection.detections_count:batchDetection.detections_count]

		for _, detection := range detections {
			ds = append(ds, Detection{
				Rectangle: image.Rectangle{
					Min: image.Point{
						X: int(detection.rectangle.left),
						Y: int(detection.rectangle.top),
					},
					Max: image.Point{
						X: int(detection.rectangle.right),
						Y: int(detection.rectangle.bottom),
					},
				},
				Confidence: float64(detection.confidence),
			})
		}

		bd.detections[i] <- ds

		C.free(unsafe.Pointer(batchDetection.detections))
	}
}
