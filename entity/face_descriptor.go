package entity

import (
	"bytes"
	"database/sql/driver"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/lib/pq"

	"git.tattelecom.ru/dimuls/face"
)

const FaceDescriptorSize = face.DescriptorSize

type FaceDescriptor face.Descriptor

func (fd *FaceDescriptor) Scan(src interface{}) error {
	switch src := src.(type) {
	case []byte:
		return fd.scanBytes(src)
	case string:
		return fd.scanBytes([]byte(src))
	case nil:
		return nil
	}

	return fmt.Errorf("cannot convert %T to FaceDescriptor", src)
}

func (fd FaceDescriptor) Value() (driver.Value, error) {
	return pq.Array(fd).Value()
}

func (fd *FaceDescriptor) scanBytes(src []byte) error {
	elems, err := scanLinearArray(src, []byte{','}, "Float32Array")
	if err != nil {
		return err
	}
	for i, vs := range elems {
		v, err := strconv.ParseFloat(string(vs), 32)
		if err != nil {
			return fmt.Errorf("parsing array element index %d: %w",
				i, err)
		}
		fd[i] = float32(v)
	}
	return nil
}

func scanLinearArray(src, del []byte, typ string) (elems [][]byte, err error) {
	dims, elems, err := parseArray(src, del)
	if err != nil {
		return nil, err
	}
	if len(dims) > 1 {
		return nil, fmt.Errorf("cannot convert ARRAY%s to %s", strings.Replace(fmt.Sprint(dims), " ", "][", -1), typ)
	}
	return elems, err
}

func parseArray(src, del []byte) (dims []int, elems [][]byte, err error) {
	var depth, i int

	if len(src) < 1 || src[0] != '{' {
		return nil, nil, fmt.Errorf("pq: unable to parse array; expected %q at offset %d", '{', 0)
	}

Open:
	for i < len(src) {
		switch src[i] {
		case '{':
			depth++
			i++
		case '}':
			elems = make([][]byte, 0)
			goto Close
		default:
			break Open
		}
	}
	dims = make([]int, i)

Element:
	for i < len(src) {
		switch src[i] {
		case '{':
			if depth == len(dims) {
				break Element
			}
			depth++
			dims[depth-1] = 0
			i++
		case '"':
			var elem = []byte{}
			var escape bool
			for i++; i < len(src); i++ {
				if escape {
					elem = append(elem, src[i])
					escape = false
				} else {
					switch src[i] {
					default:
						elem = append(elem, src[i])
					case '\\':
						escape = true
					case '"':
						elems = append(elems, elem)
						i++
						break Element
					}
				}
			}
		default:
			for start := i; i < len(src); i++ {
				if bytes.HasPrefix(src[i:], del) || src[i] == '}' {
					elem := src[start:i]
					if len(elem) == 0 {
						return nil, nil, fmt.Errorf("pq: unable to parse array; unexpected %q at offset %d", src[i], i)
					}
					if bytes.Equal(elem, []byte("NULL")) {
						elem = nil
					}
					elems = append(elems, elem)
					break Element
				}
			}
		}
	}

	for i < len(src) {
		if bytes.HasPrefix(src[i:], del) && depth > 0 {
			dims[depth-1]++
			i += len(del)
			goto Element
		} else if src[i] == '}' && depth > 0 {
			dims[depth-1]++
			depth--
			i++
		} else {
			return nil, nil, fmt.Errorf("pq: unable to parse array; unexpected %q at offset %d", src[i], i)
		}
	}

Close:
	for i < len(src) {
		if src[i] == '}' && depth > 0 {
			depth--
			i++
		} else {
			return nil, nil, fmt.Errorf("pq: unable to parse array; unexpected %q at offset %d", src[i], i)
		}
	}
	if depth > 0 {
		err = fmt.Errorf("pq: unable to parse array; expected %q at offset %d", '}', i)
	}
	if err == nil {
		for _, d := range dims {
			if (len(elems) % d) != 0 {
				err = fmt.Errorf("pq: multidimensional arrays must have elements with matching dimensions")
			}
		}
	}
	return
}

func FaceDescriptorDistance(a, b FaceDescriptor) float64 {
	var sum float64
	for i := range a {
		sum += math.Pow(float64(a[i]-b[i]), 2)
	}
	return math.Sqrt(sum)
}
