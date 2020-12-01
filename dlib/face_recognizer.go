package dlib

import (
	"git.tattelecom.ru/dimuls/face"
	"github.com/bennyharvey/soma/entity"
	"gocv.io/x/gocv"
)

type FaceRecognizer struct {
	recognizer *face.Recognizer
}

func NewFaceRecognizer(shaperModelPath, recognizerModelPath string, jittering int) (*FaceRecognizer, error) {
	r, err := face.NewRecognizer(shaperModelPath, recognizerModelPath, jittering)
	if err != nil {
		return nil, err
	}

	return &FaceRecognizer{recognizer: r}, nil
}

func (fr *FaceRecognizer) Close() {
	fr.recognizer.Close()
}

func (fr *FaceRecognizer) RecognizeFace(face gocv.Mat) (entity.FaceDescriptor, error) {
	d, err := fr.recognizer.Recognize(face)
	return entity.FaceDescriptor(d), err
}
