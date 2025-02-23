package server

import (
	"log"
	"strings"

	"github.com/Bihan001/MyDB/internal/engine"
	"github.com/Bihan001/MyDB/internal/ioprotocol"
	"github.com/Bihan001/MyDB/internal/wal"
)

type ServiceRunner interface {
    RunService() error
}

func preRun(context *engine.Context, evaluator engine.Evaluator) {
    // Load WAL if any content is present, replay it
    walData := context.WAL.ReadWALFile()
    if len(walData) == 0 {
        return
    }
    ops, err := parseCommands(walData, len(walData), context.Decoder, context.WAL)
    if err != nil {
        log.Fatal("error reading commands from WAL file: ", err)
    }
    _, _ = evaluator.Evaluate(ops)
}

func parseCommands(buffer []byte, bufferLen int, decoder ioprotocol.Decoder, wal wal.WAL) (engine.OperationList, error) {
    wal.WriteToWAL(buffer)
    parsed, err := decoder.Decode(buffer[:bufferLen])
    if err != nil {
        return nil, err
    }

    var ops engine.OperationList
    for _, val := range parsed {
        tokens, err := toStringArray(val.([]interface{}))
        if err != nil {
            return nil, err
        }
        ops = append(ops, &engine.Operation{
            Name: strings.ToUpper(tokens[0]),
            Args: tokens[1:],
        })
    }
    return ops, nil
}

func toStringArray(arr []interface{}) ([]string, error) {
    res := make([]string, len(arr))
    for i := range arr {
        res[i] = arr[i].(string)
    }
    return res, nil
}
