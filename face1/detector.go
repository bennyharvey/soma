package face1

// #cgo pkg-config: dlib-1 opencv4
// #cgo CXXFLAGS: -std=c++1z -Wall -O3 -DNDEBUG -march=native
// #cgo LDFLAGS: -ldlib -lcblas -lblas -llapack -lopencv_core
// #include <stdlib.h>
// #include "detector.h"
import "C"

import (
	"errors"
	"image"
	"unsafe"

	"github.com/sirupsen/logrus"

	"gocv.io/x/gocv"
)

type Detection struct {
	Rectangle  image.Rectangle
	Confidence float64
}

type Detector struct {
	faceDetector unsafe.Pointer
}

func NewDetector(modelPath string) (*Detector, error) {
	cModelPath := C.CString(modelPath)
	defer C.free(unsafe.Pointer(cModelPath))

	result := C.detector_init(cModelPath)
	defer C.free(unsafe.Pointer(result))

	if result.err_str != nil {
		defer C.free(unsafe.Pointer(result.err_str))
		return nil, errors.New(C.GoString(result.err_str))
	}

	return &Detector{faceDetector: result.detector}, nil
}

func (d *Detector) Close() {
	C.detector_free(d.faceDetector)
	d.faceDetector = nil
}

func (d *Detector) Detect(img gocv.Mat) ([]Detection, error) {
	result := C.detector_detect(d.faceDetector, unsafe.Pointer(img.Ptr()))

	logrus.Info("grettings from detector")

	if result.err_str != nil {
		defer C.free(unsafe.Pointer(result.err_str))
		return nil, errors.New(C.GoString(result.err_str))
	}

	defer C.free(unsafe.Pointer(result.detections))

	var detectionsArray *C.Detection = result.detections

	detections := (*[1 << 28]C.Detection)(unsafe.Pointer(detectionsArray))[:result.detections_count:result.detections_count]

	var ds []Detection

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

	return ds, nil
}
