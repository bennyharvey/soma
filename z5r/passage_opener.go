package z5r

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/bennyharvey/soma/entity"
)

type PassageOpener struct {
	BaseURI   string
	Direction entity.Direction
	lastOpen  time.Time
}

func NewPassageOpener(baseURI string, direction entity.Direction) *PassageOpener {
	return &PassageOpener{
		BaseURI:   baseURI,
		Direction: direction,
	}
}

const (
	openAPIPath = `/cgi-bin/command`

	openBodyPrefix = `DIR=`

	inDirectionNum  = "0"
	outDirectionNum = "1"
)

func (po *PassageOpener) OpenPassage() error {
	var directionNum string

	if po.Direction == entity.In {
		directionNum = inDirectionNum
	} else {
		directionNum = outDirectionNum
	}

	res, err := http.Post(po.BaseURI+openAPIPath, "text/plain",
		strings.NewReader(openBodyPrefix+directionNum))
	if err != nil {
		return fmt.Errorf("HTTP post: %w", err)
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("expected 200 status code but got %d",
			res.StatusCode)
	}

	po.lastOpen = time.Now()

	return nil
}

func (po *PassageOpener) LastOpenTime() time.Time {
	return po.lastOpen
}
