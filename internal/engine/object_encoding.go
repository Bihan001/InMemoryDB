package engine

import (
    "errors"
    "strconv"
)

type ObjectEncoder interface {
    EvaluateObjectEncoding(val string) (uint8, uint8)
}

type simpleObjectEncoder struct {
}

func GetNewObjectEncoder() ObjectEncoder {
    return &simpleObjectEncoder{}
}

func (soe *simpleObjectEncoder) EvaluateObjectEncoding(val string) (uint8, uint8) {
    objType := TypeString
    if _, err := strconv.ParseInt(val, 10, 64); err == nil {
        return objType, EncodingInt
    }
    if len(val) <= 44 {
        return objType, EncodingEmbeddedString
    }
    return objType, EncodingRaw
}

// Reuse type-checking from original code
func checkTypeMask(te uint8, t uint8) error {
    if (te & 0b11110000) != t {
        return errors.New("the operation is not permitted on this type")
    }
    return nil
}

func checkEncMask(te uint8, e uint8) error {
    if (te & 0b00001111) != e {
        return errors.New("the operation is not permitted on this encoding")
    }
    return nil
}

var (
    TypeString             uint8 = 0
    EncodingInt            uint8 = 1
    EncodingRaw            uint8 = 0
    EncodingEmbeddedString uint8 = 8
)
