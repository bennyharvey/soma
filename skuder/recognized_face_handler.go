package skuder

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/bennyharvey/soma/entity"
)

type DBStorage interface {
	Person(personID int64) (entity.Person, error)
	FindClosestPersonFace(entity.FaceDescriptor) (entity.PersonFace, float64, bool)
	AddEvent(entity.Event) error
}

type PassageOpener interface {
	OpenPassage() error
	LastOpenTime() time.Time
}

type PhotoStorage interface {
	AddPhoto(photoID string, photo []byte) error
}

type RecognizedFaceHandler struct {
	passageID             string
	waitAfterOpen         time.Duration
	matchDistance         float64
	detectConfidenceLimit float64
	dbStorage             DBStorage
	photoStorage          PhotoStorage
	passageOpener         PassageOpener
	log                   *logrus.Entry
}

func NewRecognizedFaceHandler(passageID string, waitAfterOpen time.Duration, matchDistance float64,
	detectConfidenceLimit float64, dbs DBStorage, ps PhotoStorage, po PassageOpener) *RecognizedFaceHandler {

	return &RecognizedFaceHandler{
		passageID:             passageID,
		waitAfterOpen:         waitAfterOpen,
		matchDistance:         matchDistance,
		detectConfidenceLimit: detectConfidenceLimit,
		dbStorage:             dbs,
		photoStorage:          ps,
		passageOpener:         po,
		log:                   logrus.WithField("subsystem", "facer_recognized_face_handler"),
	}
}

func (rfh *RecognizedFaceHandler) HandleRecognizedFace(rf entity.RecognizedFace) {
	photoMD5Sum := md5.Sum(rf.Photo)
	photoID := hex.EncodeToString(photoMD5Sum[:])

	log := logrus.WithFields(logrus.Fields{
		"photo_id":           photoID,
		"frame_time":         rf.FrameTime,
		"detect_time":        rf.DetectTime,
		"detect_duration":    rf.DetectTime.Sub(rf.FrameTime),
		"detect_confidence":  rf.DetectConfidence,
		"recognize_time":     rf.RecognizeTime,
		"recognize_duration": rf.RecognizeTime.Sub(rf.DetectTime),
	})

	err := rfh.photoStorage.AddPhoto(photoID, rf.Photo)
	if err != nil {
		log.WithError(err).Error("failed to add photo to photo storage")
	}

	if rf.DetectConfidence < rfh.detectConfidenceLimit {
		log.Info("too low confidence, skipping")
		return
	}

	pf, distance, found := rfh.dbStorage.FindClosestPersonFace(rf.Descriptor)

	if found && distance < rfh.matchDistance {
		p, err := rfh.dbStorage.Person(pf.PersonID)
		if err != nil {
			if err == entity.ErrPersonNotFound {
				rfh.addFaceRecognizeEvent(log, rf, photoID)
			} else {
				log.WithError(err).Error("failed to get person from DB storage")
			}
			return
		}

		rfh.addPersonRecognizeEvent(log, rf, pf, p, photoID, distance)
		rfh.openPassage(log, pf, p)
	} else {
		rfh.addFaceRecognizeEvent(log, rf, photoID)
	}
}

func (rfh *RecognizedFaceHandler) addFaceRecognizeEvent(log *logrus.Entry, rf entity.RecognizedFace, photoID string) {
	data, err := json.Marshal(entity.FaceRecognizedData{
		PhotoID:          photoID,
		FaceDescriptor:   rf.Descriptor,
		DetectConfidence: rf.DetectConfidence,
	})
	if err != nil {
		log.WithError(err).Error("failed to JSON marshal face recognized data")
		return
	}

	err = rfh.dbStorage.AddEvent(entity.Event{
		Time: time.Now(),
		Type: entity.FaceRecognize,
		Data: data,
	})
	if err != nil {
		log.WithError(err).Error("failed to add face recognized event")
		return
	}

	log.Info("face recognized")
}

func (rfh *RecognizedFaceHandler) addPersonRecognizeEvent(log *logrus.Entry, rf entity.RecognizedFace,
	pf entity.PersonFace, p entity.Person, photoID string, distance float64) {

	rpData, err := json.Marshal(entity.PersonRecognizeData{
		PhotoID:             photoID,
		PersonID:            pf.PersonID,
		PersonName:          p.Name,
		PersonPosition:      p.Position,
		PersonUnit:          p.Unit,
		DetectConfidence:    rf.DetectConfidence,
		FaceDescriptor:      rf.Descriptor,
		DescriptorsDistance: distance,
	})
	if err != nil {
		log.WithError(err).Error("failed to JSON marshal person recognized data")
		return
	}

	err = rfh.dbStorage.AddEvent(entity.Event{
		Time: time.Now(),
		Type: entity.PersonRecognize,
		Data: rpData,
	})
	if err != nil {
		log.WithError(err).Error("failed to add person recognized event")
		return
	}

	log.Info("person recognized")
}

func (rfh *RecognizedFaceHandler) openPassage(log *logrus.Entry, pf entity.PersonFace, p entity.Person) {
	err := rfh.passageOpener.OpenPassage()
	if err != nil {
		log.WithError(err).Error("failed to open passage")
		return
	}

	openTime := time.Now()

	log.WithField("passage_open_time", openTime).Info("passage opened")

	data, err := json.Marshal(entity.PassageOpenData{
		PersonID:       pf.PersonID,
		PersonName:     p.Name,
		PersonPosition: p.Position,
		PersonUnit:     p.Unit,
		PassageID:      rfh.passageID,
	})
	if err != nil {
		log.WithError(err).Error("failed to JSON marshal passage open data")
		return
	}

	err = rfh.dbStorage.AddEvent(entity.Event{
		Time: openTime,
		Type: entity.PassageOpen,
		Data: data,
	})
	if err != nil {
		log.WithError(err).Error("failed to add passage open event to DB storage")
		return
	}

	log.Info("passage opened")
}
