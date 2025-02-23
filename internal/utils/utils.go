package utils

import (
	"errors"
	"time"
)

func GetCurrentTime() uint64 {
	return uint64(time.Now().UnixMilli())
}

func CheckTypeMask(te uint8, t uint8) error {
    if (te & 0b11110000) != t {
        return errors.New("the operation is not permitted on this type")
    }
    return nil
}

func CheckEncMask(te uint8, e uint8) error {
    if (te & 0b00001111) != e {
        return errors.New("the operation is not permitted on this encoding")
    }
    return nil
}