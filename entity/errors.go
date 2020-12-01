package entity

import "errors"

var (
	ErrUserNotFound   = errors.New("user not found")
	ErrPersonNotFound = errors.New("person not found")
)

type InvalidParamErr struct {
	Param string
}

func (e InvalidParamErr) Error() string {
	return e.Param + "is invalid"
}
