package entity

import (
	"encoding/json"
	"time"
)

type EventType string

const (
	PassageOpen     EventType = "passage_open"
	FaceRecognize   EventType = "face_recognize"
	PersonRecognize EventType = "person_recognize"
)

type PassageOpenData struct {
	PersonID       int64  `json:"person_id"`
	PersonName     string `json:"person_name"`
	PersonPosition string `json:"person_position"`
	PersonUnit     string `json:"person_unit"`
	PassageID      string `json:"passage_id"`
}

type FaceRecognizedData struct {
	PhotoID          string         `json:"photo_id"`
	FaceDescriptor   FaceDescriptor `json:"face_descriptor"`
	DetectConfidence float64        `json:"detect_confidence"`
}

type PersonRecognizeData struct {
	PhotoID             string         `json:"photo_id"`
	PersonID            int64          `json:"person_id"`
	PersonName          string         `json:"person_name"`
	PersonPosition      string         `json:"person_position"`
	PersonUnit          string         `json:"person_unit"`
	DetectConfidence    float64        `json:"detect_confidence"`
	FaceDescriptor      FaceDescriptor `json:"face_descriptor"`
	DescriptorsDistance float64        `json:"descriptors_distance"`
}

type Event struct {
	ID   int64           `json:"id" db:"id"`
	Time time.Time       `json:"time" db:"time"`
	PassageID string `json:"passageID" db:"passage_id"`
	Type EventType       `json:"type" db:"type"`
	Data json.RawMessage `json:"data" db:"data"`
}

type EventsFilters struct {
	From time.Time
	To   time.Time
	Type EventType

	PassageID string

	PersonID       int64
	PersonName     string
	PersonPosition string
	PersonUnit     string

	OrderBy        string
	OrderDirection string

	Offset int
	Limit  int

	Set struct {
		From, To, Type, PassageID, PersonID, PersonName, PersonPosition, PersonUnit, OrderBy, OrderDirection,
		OffSet, Limit bool
	}
}

type EventsFilter func(p *EventsFilters)

func EventsFrom(from time.Time) EventsFilter {
	return func(p *EventsFilters) {
		p.From = from
		p.Set.From = true
	}
}

func EventsTo(to time.Time) EventsFilter {
	return func(p *EventsFilters) {
		p.To = to
		p.Set.To = true
	}
}

func EventsType(et EventType) EventsFilter {
	return func(p *EventsFilters) {
		p.Type = et
		p.Set.Type = true
	}
}

func EventsPassageID(id string) EventsFilter {
	return func(p *EventsFilters) {
		p.PassageID = id
		p.Set.PassageID = true
	}
}

func EventsPersonID(id int64) EventsFilter {
	return func(p *EventsFilters) {
		p.PersonID = id
		p.Set.PersonID = true
	}
}

func EventsPersonName(name string) EventsFilter {
	return func(p *EventsFilters) {
		p.PersonName = name
		p.Set.PersonName = true
	}
}

func EventsOrderBy(orderBy string) EventsFilter {
	return func(p *EventsFilters) {
		p.OrderBy = orderBy
		p.Set.OrderBy = true
	}
}

func EventsOrderDirection(orderDirection string) EventsFilter {
	return func(p *EventsFilters) {
		p.OrderDirection = orderDirection
		p.Set.OrderDirection = true
	}
}

func EventsLimit(limit int) EventsFilter {
	return func(p *EventsFilters) {
		p.Limit = limit
		p.Set.Limit = true
	}
}

func EventsOffset(offset int) EventsFilter {
	return func(p *EventsFilters) {
		p.Offset = offset
		p.Set.OffSet = true
	}
}
