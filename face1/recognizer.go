package face1

// #cgo pkg-config: dlib-1 opencv4
// #cgo CXXFLAGS: -std=c++1z -Wall -O3 -DNDEBUG -march=native
// #cgo LDFLAGS: -ldlib -lcblas -lblas -llapack -lopencv_core
// #include <stdlib.h>
// #include "recognizer.h"
import "C"

import (
	"errors"
	"unsafe"

	"gocv.io/x/gocv"
)

const DescriptorSize = 128

type Descriptor [DescriptorSize]float32

type Recognizer struct {
	recognizer unsafe.Pointer
}

func NewRecognizer(shaperModelPath, faceRecognizerModelPath string, jittering int) (*Recognizer, error) {
	cShaperModelPath := C.CString(shaperModelPath)
	defer C.free(unsafe.Pointer(cShaperModelPath))

	cRecognizerModelPath := C.CString(faceRecognizerModelPath)
	defer C.free(unsafe.Pointer(cRecognizerModelPath))

	result := C.recognizer_init(cShaperModelPath, cRecognizerModelPath, C.int(jittering))
	defer C.free(unsafe.Pointer(result))

	if result.err_str != nil {
		defer C.free(unsafe.Pointer(result.err_str))
		return nil, errors.New(C.GoString(result.err_str))
	}

	return &Recognizer{recognizer: result.recognizer}, nil
}

func (r *Recognizer) Close() {
	C.recognizer_free(r.recognizer)
	r.recognizer = nil
}

func (r *Recognizer) Recognize(img gocv.Mat) (Descriptor, error) {
	result := C.recognizer_recognize(r.recognizer, unsafe.Pointer(img.Ptr()))

	if result.err_str != nil {
		defer C.free(unsafe.Pointer(result.err_str))
		return Descriptor{}, errors.New(C.GoString(result.err_str))
	}

	defer C.free(unsafe.Pointer(result.descriptor))

	var descriptorArray *C.float = result.descriptor
	descriptorSlice := (*[1 << 28]float32)(unsafe.Pointer(descriptorArray))[:DescriptorSize:DescriptorSize]

	var fd Descriptor

	copy(fd[:], descriptorSlice)

	return fd, nil
}
