package ioprotocol

import (
	"errors"
	"fmt"
)

type respEncoder struct {
}

type respDecoder struct {
}

func GetNewRespEncoder() Encoder {
    return &respEncoder{}
}

func GetNewRespDecoder() Decoder {
    return &respDecoder{}
}

func (rd *respDecoder) Decode(data []byte) ([]interface{}, error) {
    if len(data) == 0 {
        return nil, errors.New("no data")
    }
    handler := &respHandler{
        payload: data,
        index:   0,
    }

    var results []interface{}
    for handler.index < len(handler.payload) {
        val, err := handler.decodeSingle()
        if err != nil {
            return results, err
        }
        results = append(results, val)
    }
    return results, nil
}

func (rd *respEncoder) Encode(data interface{}, useSimpleString bool) ([]byte, error) {
    switch val := data.(type) {
    case string:
        if useSimpleString {
            return []byte(fmt.Sprintf("+%s\r\n", val)), nil
        }
        return []byte(fmt.Sprintf("$%d\r\n%s\r\n", len(val), val)), nil
    case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
        return []byte(fmt.Sprintf(":%d\r\n", val)), nil
    default:
        return []byte("$-1\r\n"), nil
    }
}

type respHandler struct {
    payload []byte
    index   int
}

func (h *respHandler) decodeSingle() (interface{}, error) {
    if len(h.payload) == 0 {
        return nil, errors.New("no data")
    }
    start := h.index
    h.index++

    switch h.payload[start] {
    case '+':
        return h.readSimpleString()
    case '-':
        return h.readErrorString()
    case ':':
        return h.readInteger()
    case '$':
        return h.readBulkString()
    case '*':
        return h.readArray()
    }
    return nil, errors.New("invalid command prefix")
}

func (h *respHandler) readInteger() (int64, error) {
    var val int64
    for ; h.index < len(h.payload); h.index++ {
        if h.payload[h.index] == '\r' {
            h.index += 2
            return val, nil
        }
        val = val*10 + int64(h.payload[h.index]-'0')
    }
    return 0, errors.New("invalid integer format")
}

func (h *respHandler) readSimpleString() (string, error) {
    start := h.index
    var end int
    for ; h.index < len(h.payload); h.index++ {
        if h.payload[h.index] == '\r' {
            end = h.index
            break
        }
    }
    str := string(h.payload[start:end])
    h.index += 2
    return str, nil
}

func (h *respHandler) readErrorString() (string, error) {
    return h.readSimpleString()
}

func (h *respHandler) readBulkString() (string, error) {
    length, err := h.readInteger()
    if err != nil {
        return "", err
    }
    intLen := int(length)
    if intLen < 0 || (h.index+intLen+2) > len(h.payload) {
        return "", errors.New("invalid bulk string length")
    }
    str := string(h.payload[h.index : h.index+intLen])
    h.index += intLen + 2
    return str, nil
}

func (h *respHandler) readArray() ([]interface{}, error) {
    count, err := h.readInteger()
    if err != nil {
        return nil, err
    }
    arr := make([]interface{}, count)
    for i := range arr {
        elem, err := h.decodeSingle()
        if err != nil {
            return nil, err
        }
        arr[i] = elem
    }
    return arr, nil
}
