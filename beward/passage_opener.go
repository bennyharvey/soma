package beward

import (
	"fmt"
	"net/http"
	"time"

	"github.com/bennyharvey/soma/entity"
)

// PassageOpener passge opener
type PassageOpener struct {
	BaseURI   string
	Direction entity.Direction
	lastOpen  time.Time
}

// NewPassageOpener passage opener
func NewPassageOpener(baseURI string, direction entity.Direction) *PassageOpener {
	return &PassageOpener{
		BaseURI:   baseURI,
		Direction: direction,
	}
}

const (
	openAPIPath = `/cgi-bin/intercom_cgi?user=admin&pwd=admin&action=maindoor`

	inDirectionNum  = "0"
	outDirectionNum = "1"
)

// OpenPassage opens passge
func (po *PassageOpener) OpenPassage() error {
	
	res, err := http.Get("http://" + po.BaseURI + openAPIPath)
	if err != nil {
		return fmt.Errorf("HTTP request: %w", err)
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("expected 200 status code but got %d",
			res.StatusCode)
	}

	po.lastOpen = time.Now()

	return nil
}

// LastOpenTime lot
func (po *PassageOpener) LastOpenTime() time.Time {
	return po.lastOpen
}
