package server

import (
	"log"
	"strings"

	"github.com/Bihan001/MyDB/core"
)

type Server interface {
	RunServer() error
}

func init() {
	bytes := core.ReadWAL()
	if len(bytes) == 0 {
		return
	}
	cmds, err := readCommands(bytes, len(bytes))
	if err != nil {
		log.Fatal("error reading commands from WAL file: ", err)
	}
	core.Evaluate(cmds)
}

func toStringArray(arr []interface{}) ([]string, error) {
	res := make([]string, len(arr))
	for i := range arr {
		res[i] = arr[i].(string)
	}
	return res, nil
}

func readCommands(buff []byte, buffLength int) (core.Cmds, error) {
	core.WriteOperationInWAL(buff)

	resp := core.GetNewResp(buff[:buffLength])

	values, err := resp.Decode()

	if err != nil {
		return nil, err
	}

	var cmds []*core.Cmd = make([]*core.Cmd, 0)
	
	for _, value := range values {
		tokens, err := toStringArray(value.([]interface{}))
		if err != nil {
			return nil, err
		}
		cmds = append(cmds, &core.Cmd{
			Cmd:  strings.ToUpper(tokens[0]),
			Args: tokens[1:],
		})
	}
	return cmds, nil
}
