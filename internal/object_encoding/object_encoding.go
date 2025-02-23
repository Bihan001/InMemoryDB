package object_encoding

import (
	"strconv"

	"github.com/Bihan001/MyDB/internal/config"
	"github.com/Bihan001/MyDB/internal/interfaces"
)

type simpleObjectEncoder struct {}

func GetNewObjectEncoder() interfaces.ObjectEncoder {
    return &simpleObjectEncoder{}
}

func (soe *simpleObjectEncoder) EvaluateObjectEncoding(val string) (uint8, uint8) {
    objType := config.TypeString
    if _, err := strconv.ParseInt(val, 10, 64); err == nil {
        return objType, config.EncodingInt
    }
    if len(val) <= 44 {
        return objType, config.EncodingEmbeddedString
    }
    return objType, config.EncodingRaw
}

