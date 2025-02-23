package ioprotocol_test

import (
    "fmt"
    "testing"

    "github.com/Bihan001/MyDB/internal/ioprotocol"
)

func TestSimpleStringDecode(t *testing.T) {
    decoder := ioprotocol.GetNewRespDecoder()
    cases := map[string]string{
        "+OK\r\n": "OK",
    }
    for input, expected := range cases {
        values, err := decoder.Decode([]byte(input))
        if err != nil {
            t.Fail()
        }
        actual := values[0]
        if expected != actual {
            fmt.Println(actual, err)
            t.Fail()
        }
    }
}

func TestError(t *testing.T) {
    decoder := ioprotocol.GetNewRespDecoder()
    cases := map[string]string{
        "-Error message\r\n": "Error message",
    }
    for input, expected := range cases {
        values, err := decoder.Decode([]byte(input))
        actual := values[0]
        if expected != actual {
            fmt.Println(actual, err)
            t.Fail()
        }
    }
}

func TestInt64(t *testing.T) {
    decoder := ioprotocol.GetNewRespDecoder()
    cases := map[string]int64{
        ":0\r\n":    0,
        ":1000\r\n": 1000,
    }
    for input, expected := range cases {
        values, err := decoder.Decode([]byte(input))
        if err != nil {
            t.Fail()
        }
        actual := values[0]
        if expected != actual {
            fmt.Println(actual, err)
            t.Fail()
        }
    }
}

func TestBulkStringDecode(t *testing.T) {
    decoder := ioprotocol.GetNewRespDecoder()
    cases := map[string]string{
        "$5\r\nhello\r\n": "hello",
        "$0\r\n\r\n":      "",
    }
    for input, expected := range cases {
        values, err := decoder.Decode([]byte(input))
        if err != nil {
            t.Fail()
        }
        actual := values[0]
        if expected != actual {
            fmt.Println(actual, err)
            t.Fail()
        }
    }
}

func TestArrayDecode(t *testing.T) {
    decoder := ioprotocol.GetNewRespDecoder()
    cases := map[string][]interface{}{
        "*0\r\n": {},
        "*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n": {"hello", "world"},
        "*3\r\n:1\r\n:2\r\n:3\r\n": {int64(1), int64(2), int64(3)},
        "*5\r\n:1\r\n:2\r\n:3\r\n:4\r\n$5\r\nhello\r\n": {
            int64(1), int64(2), int64(3), int64(4), "hello"},
        "*2\r\n*3\r\n:1\r\n:2\r\n:3\r\n*2\r\n+Hello\r\n-World\r\n": {
            []int64{int64(1), int64(2), int64(3)},
            []interface{}{"Hello", "World"},
        },
    }
    for input, expected := range cases {
        values, _ := decoder.Decode([]byte(input))
        actual := values[0]
        arr := actual.([]interface{})
        if len(arr) != len(expected) {
            t.Fail()
        }
        for i := range arr {
            if fmt.Sprintf("%v", expected[i]) != fmt.Sprintf("%v", arr[i]) {
                t.Fail()
            }
        }
    }
}
