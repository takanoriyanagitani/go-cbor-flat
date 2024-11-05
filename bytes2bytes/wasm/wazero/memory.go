package b2wazero2b

import (
	"errors"

	wa "github.com/tetratelabs/wazero/api"
)

var (
	ErrUnableToWrite error = errors.New("unable to write")
	ErrUnableToRead  error = errors.New("unable to read")
)

type Memory struct{ wa.Memory }

func (m Memory) WriteBytes(offset uint32, data []byte) error {
	var ok bool = m.Memory.Write(offset, data)
	var ng bool = !ok
	if ng {
		return ErrUnableToWrite
	}
	return nil
}

func (m Memory) GetView(offset uint32, size uint32) ([]byte, error) {
	view, ok := m.Memory.Read(offset, size)
	var ng bool = !ok
	if ng {
		return nil, ErrUnableToRead
	}
	return view, nil
}
