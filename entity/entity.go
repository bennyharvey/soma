package entity

import (
	"time"

	"git.tattelecom.ru/dimuls/face"
)

type FaceDetection = face.Detection

type DetectedFace struct {
	StreamID         string
	Photo            []byte
	DetectConfidence float64
	DetectTime       time.Time
	FrameTime        time.Time
}

type RecognizedFace struct {
	DetectedFace
	Descriptor    FaceDescriptor
	RecognizeTime time.Time
}

type PassageType string

const (
	Z5R   PassageType = "z5r"
	Sigur PassageType = "sigur"
	Dummy PassageType = "dummy"
	Beward PassageType = "beward"
)

type Direction string

const (
	In  Direction = "in"
	Out Direction = "out"
)

var EqualFacesMaxDistance float32 = 0.5

type Role string

const (
	Admin    = "admin"
	Security = "security"
)

type User struct {
	Login        string `json:"login" db:"login"`
	Password     string `json:"password" db:"-"`
	PasswordHash []byte `json:"-" db:"password_hash"`
	Role         Role   `json:"role" db:"role"`
}

type Person struct {
	ID       int64  `json:"id" db:"id"`
	Name     string `json:"name" db:"name"`
	Position string `json:"position" db:"position"`
	Unit     string `json:"unit" db:"unit"`
}

type PersonFace struct {
	ID         int64          `json:"id" db:"id"`
	PersonID   int64          `json:"person_id" db:"person_id"`
	Descriptor FaceDescriptor `json:"descriptor" db:"descriptor"`
	PhotoID    string         `json:"photo_id" db:"photo_id"`
}
