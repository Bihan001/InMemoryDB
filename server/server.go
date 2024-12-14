package server

import (
	"fmt"
	"strings"

	"github.com/Bihan001/MyDB/core"
)

type Server interface {
	RunServer() error
}

func readCommand(buff []byte, buffLength int) (*core.Cmd, error) {
	fmt.Println(string(buff[:buffLength]))

	resp := core.GetNewResp(buff[:buffLength])

	tokens, err := resp.DecodeArrayString()

	if err != nil {
		return nil, err
	}

	return &core.Cmd{
		Cmd:  strings.ToUpper(tokens[0]),
		Args: tokens[1:],
	}, nil
}
