package core

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/Bihan001/MyDB/config"
)

func Evaluate(cmds Cmds) ([]byte, error) {
	buff := bytes.NewBuffer(make([]byte, 0))

	for _, cmd := range cmds {
		var response []byte
		var err error

		switch cmd.Cmd {
		case "PING":
			response, err = evaluatePing(cmd.Args)
		case "GET":
			response, err = evaluateGet(cmd.Args)
		case "SET":
			response, err = evaluateSet(cmd.Args)
		case "DEL":
			response, err = evaluateDelete(cmd.Args)
		case "TTL":
			response, err = evaluateTTL(cmd.Args)
		case "EXPIRE":
			response, err = evaluateExpire(cmd.Args)
		default:
			err = errors.New("invalid command")
		}
	
		if err != nil {
			response, err = Encode(fmt.Sprint(err), false)
		}

		if err == nil {
			buff.Write(response)
		}

		if err != nil {
			return make([]byte, 0), err
		}

	}
	
	return buff.Bytes(), nil
}

func evaluatePing(args []string) ([]byte, error) {
	switch len(args) {
	case 0:
		return Encode("PONG", true)
	case 1:
		return Encode(args[0], false)
	default:
		return nil, errors.New("invalid number of arguments")
	}
}

func evaluateGet(args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("(error) ERR wrong number of arguments for 'get' command")
	}

	key := args[0]
	val := store.Get(key)

	if val == nil {
		return RESP_NIL, nil
	}

	return Encode(val.value, false)
}

func evaluateSet(args []string) ([]byte, error) {
	if len(args) < 2 {
		return nil, errors.New("(error) ERR wrong number of arguments for 'set' command")
	}

	if (store.Length() >= config.MaxKeyLimit) {
		evict()
	}

	key, value := args[0], args[1]
	var expiryDurationMs int64 = -1

	for i := 2; i < len(args); i++ {
		switch args[i] {
		case "EX", "ex":
			i++
			if i == len(args) {
				return nil, errors.New("(error) ERR syntax error")
			}
			expiryDurationSec, err := strconv.ParseInt(args[i], 10, 64)
			if err != nil {
				return nil, errors.New("(error) ERR value is not an integer or out of range")
			}
			expiryDurationMs = expiryDurationSec * 1000
		default:
			return nil, errors.New("(error) ERR syntax error")
		}
	}

	store.Put(key, store.NewValue(value, expiryDurationMs))
	return Encode("OK", true)
}

func evaluateDelete(args []string) ([]byte, error) {
	deletedCount := 0

	for _, k := range args {
		if ok := store.Delete(k); ok {
			deletedCount++
		}
	}

	return Encode(deletedCount, false)
}

func evaluateTTL(args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("(error) ERR wrong number of arguments for 'ttl' command")
	}

	key := args[0]
	value := store.Get(key)

	if value == nil {
		return Encode(-2, false)
	}

	if value.expiresAt == -1 {
		return Encode(-1, false)
	}

	expiryDuration := (value.expiresAt - time.Now().UnixMilli()) / 1000

	return Encode(expiryDuration, false)
}

func evaluateExpire(args []string) ([]byte, error) {
	if len(args) < 2 {
		return nil, errors.New("(error) ERR wrong number of arguments for 'expire' command")
	}

	key := args[0]
	expiryDurationSec, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return nil, errors.New("(error) ERR value is not an integer or out of range")
	}

	value := store.Get(key)

	if value == nil {
		return Encode(0, false)
	}

	value.expiresAt = time.Now().UnixMilli() + expiryDurationSec * 1000;

	return Encode(1, false)
}
