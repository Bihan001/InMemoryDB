package ioprotocol

type Encoder interface {
    Encode(data interface{}, useSimpleString bool) ([]byte, error)
}

type Decoder interface {
    Decode(data []byte) ([]interface{}, error)
}
