package core

import (
	"errors"
	"strconv"
)
func getType(te uint8) uint8 {
	return te & 0b11110000
}

func getEncoding(te uint8) uint8 {
	return te & 0b00001111
}

func assertType(te uint8, t uint8) error {
	if getType(te) != t {
		return errors.New("the operation is not permitted on this type")
	}
	return nil
}

func assertEncoding(te uint8, e uint8) error {
	if getEncoding(te) != e {
		return errors.New("the operation is not permitted on this encoding")
	}
	return nil
}

func tryObjectEncoding(v string) (uint8, uint8) {
	objectType := OBJECT_TYPE_STRING

	if _, err := strconv.ParseInt(v, 10, 64); err == nil {
		return objectType, OBJ_ENCODING_INT
	}
	
	if len(v) <= 44 {
		return objectType, OBJ_ENCODING_EMBEDDED_STR
	}

	return objectType, OBJ_ENCODING_RAW
}