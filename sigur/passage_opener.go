package sigur

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/bennyharvey/soma/entity"
)

type PassageOpener struct {
	Address   string
	Direction entity.Direction
	lastOpen  time.Time
}

func NewPassageOpener(address string, direction entity.Direction) *PassageOpener {
	return &PassageOpener{
		Address:   address,
		Direction: direction,
	}
}

const (
	loginMsg   = "LOGIN 1.8 Administrator password\n"
	openMsgFmt = "ALLOWPASS 1 ANONYMOUS %s\n"
	exitMsg    = "EXIT\n"
)

func (po *PassageOpener) OpenPassage() error {
	conn, err := net.Dial("tcp", po.Address)
	if err != nil {
		return fmt.Errorf("dial sigur controller: %w", err)
	}

	defer func() {
		err = conn.Close()
		if err != nil {
			logrus.WithField("subsystem", "sigur_passage_opener").
				WithError(err).Error("failed to close connection")
		}
	}()

	buf := make([]byte, 1024)

	err = conn.SetDeadline(time.Now().Add(100 * time.Millisecond))
	if err != nil {
		return fmt.Errorf("set connection deadline: %w", err)
	}

	_, err = conn.Write([]byte(loginMsg))
	if err != nil {
		return fmt.Errorf("write login message: %w", err)
	}

	n, err := conn.Read(buf)
	if err != nil {
		return fmt.Errorf("read login response: %w", err)
	}

	if n != 4 && string(buf[:4]) != "OK\r\n" {
		return fmt.Errorf("unexpected login reponse: %s",
			string(buf[:n]))
	}

	_, err = conn.Write([]byte(fmt.Sprintf(openMsgFmt,
		strings.ToUpper(string(po.Direction)))))
	if err != nil {
		return fmt.Errorf("write open message: %w", err)
	}

	n, err = conn.Read(buf)
	if err != nil {
		return fmt.Errorf("read open response: %w", err)
	}

	if n != 4 && string(buf[:4]) != "OK\r\n" {
		return fmt.Errorf("unexpected open reponse: %s",
			string(buf[:n]))
	}

	_, err = conn.Write([]byte(exitMsg))
	if err != nil {
		return fmt.Errorf("write exit message: %w", err)
	}

	po.lastOpen = time.Now()

	return nil
}

func (po *PassageOpener) LastOpenTime() time.Time {
	return po.lastOpen
}
