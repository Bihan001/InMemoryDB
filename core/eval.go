package core

import (
	"errors"
	"fmt"
)

func Evaluate(cmd *Cmd) ([]byte, error) {
	var response []byte
	var err error

	switch cmd.Cmd {
	case "PING":
		response, err = evaluatePing(cmd.Args)
	default:
		err = errors.New("invalid command")
	}

	if err != nil {
		response, err = Encode(fmt.Sprint(err), false)
	}

	return response, err
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
