package engine

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/Bihan001/MyDB/internal/config"
	"github.com/Bihan001/MyDB/internal/interfaces"
)

type Evaluator interface {
    Evaluate(ops OperationList) ([]byte, error)
}

type defaultEvaluator struct {
    context *Context
    store  interfaces.DataStore
    expirableStore interfaces.Expirable
}

func GetNewEvaluator(context *Context) Evaluator {
    return &defaultEvaluator{
        context: context,
        store: context.Store,
        expirableStore: context.Store.(interfaces.Expirable),
    }
}

func (d *defaultEvaluator) Evaluate(ops OperationList) ([]byte, error) {
    buff := bytes.NewBuffer(make([]byte, 0))

    for _, op := range ops {
        var output []byte
        var err error

        switch op.Name {
        case "PING":
            output, err = d.evaluatePing(op.Args)
        case "GET":
            output, err = d.evaluateGet(op.Args)
        case "SET":
            output, err = d.evaluateSet(op.Args)
        case "DEL":
            output, err = d.evaluateDelete(op.Args)
        case "TTL":
            output, err = d.evaluateTTL(op.Args)
        case "EXPIRE":
            output, err = d.evaluateExpire(op.Args)
        case "INCR":
            output, err = d.evaluateIncrement(op.Args)
        case "INFO":
            output, err = d.evaluateInfo(op.Args)
        case "CLIENT":
            output, err = d.evaluateClient(op.Args)
        case "LATENCY":
            output, err = d.evaluateLatency(op.Args)
        default:
            err = errors.New("invalid command")
        }

        if err != nil {
            output, _ = d.context.Encoder.Encode(fmt.Sprint(err), false)
            buff.Write(output)
            continue
        }

        buff.Write(output)
    }

    return buff.Bytes(), nil
}

func (d *defaultEvaluator) evaluatePing(args []string) ([]byte, error) {
    switch len(args) {
    case 0:
        return d.context.Encoder.Encode("PONG", true)
    case 1:
        return d.context.Encoder.Encode(args[0], false)
    default:
        return nil, errors.New("invalid number of arguments for PING")
    }
}

func (d *defaultEvaluator) evaluateGet(args []string) ([]byte, error) {
    if len(args) != 1 {
        return nil, errors.New("(error) ERR wrong number of arguments for 'get' command")
    }
    val := d.context.Store.Get(args[0])
    if val == nil {
        return []byte("$-1\r\n"), nil
    }
    return d.context.Encoder.Encode(val.GetValue(), false)
}

func (d *defaultEvaluator) evaluateSet(args []string) ([]byte, error) {
    if len(args) < 2 {
        return nil, errors.New("(error) ERR wrong number of arguments for 'set' command")
    }
    if d.context.Store.Size() >= config.KeyCountLimit {
        d.context.EvictionPolicy.Evict()
    }

    key, value := args[0], args[1]
    objType, objEncoding := d.context.ObjectEncoder.EvaluateObjectEncoding(value)
    var expireMs int64 = -1

    for i := 2; i < len(args); i++ {
        switch args[i] {
        case "EX", "ex":
            i++
            if i == len(args) {
                return nil, errors.New("(error) ERR syntax error")
            }
            expireSec, parseErr := strconv.ParseInt(args[i], 10, 64)
            if parseErr != nil {
                return nil, errors.New("(error) ERR value is not an integer or out of range")
            }
            expireMs = expireSec * 1000
        default:
            return nil, errors.New("(error) ERR syntax error")
        }
    }

    entry := d.context.Store.CreateEntry(value, expireMs, objType, objEncoding)
    d.context.Store.Set(key, entry)
    return d.context.Encoder.Encode("OK", true)
}

func (d *defaultEvaluator) evaluateDelete(args []string) ([]byte, error) {
    deletedCount := 0
    for _, k := range args {
        if d.context.Store.Del(k) {
            deletedCount++
        }
    }
    return d.context.Encoder.Encode(deletedCount, false)
}

func (d *defaultEvaluator) evaluateTTL(args []string) ([]byte, error) {
    if len(args) != 1 {
        return nil, errors.New("(error) ERR wrong number of arguments for 'ttl' command")
    }
    entry := d.context.Store.Get(args[0])
    if entry == nil {
        return d.context.Encoder.Encode(int64(-2), false)
    }

    exp, exists := d.expirableStore.GetExpiryMs(entry)

    if !exists {
        return d.context.Encoder.Encode(int64(-1), false)
    }

    if exp <= uint64(time.Now().UnixMilli()) {
        return d.context.Encoder.Encode(int64(-2), false)
    }

    secondsLeft := (exp - (uint64(time.Now().UnixMilli()))) / 1000
    
    return d.context.Encoder.Encode(secondsLeft, false)
}

func (d *defaultEvaluator) evaluateExpire(args []string) ([]byte, error) {
    if len(args) < 2 {
        return nil, errors.New("(error) ERR wrong number of arguments for 'expire' command")
    }
    key := args[0]
    expireSec, err := strconv.ParseInt(args[1], 10, 64)
    if err != nil {
        return nil, errors.New("(error) ERR value is not an integer or out of range")
    }
    entry := d.context.Store.Get(key)
    if entry == nil {
        return d.context.Encoder.Encode(0, false)
    }
    d.expirableStore.SetExpiry(entry, time.Now().UnixMilli() + expireSec*1000)
    return d.context.Encoder.Encode(1, false)
}

func (d *defaultEvaluator) evaluateIncrement(args []string) ([]byte, error) {
    if len(args) != 1 {
        return d.context.Encoder.Encode(errors.New("ERR wrong number of arguments for 'incr' command"), false)
    }
    key := args[0]
    entry := d.context.Store.Get(key)
    if entry == nil {
        entry = d.context.Store.CreateEntry("0", -1, TypeString, EncodingInt)
        d.context.Store.Set(key, entry)
    }
    if err := checkTypeMask(entry.GetTypeEncoding(), TypeString); err != nil {
        return d.context.Encoder.Encode(err, false)
    }
    if err := checkEncMask(entry.GetTypeEncoding(), EncodingInt); err != nil {
        return d.context.Encoder.Encode(err, false)
    }

    numericVal, _ := strconv.ParseInt(entry.GetValue().(string), 10, 64)
    numericVal++
    entry.SetValue(strconv.FormatInt(numericVal, 10))
    return d.context.Encoder.Encode(numericVal, false)
}

func (d *defaultEvaluator) evaluateInfo(args []string) ([]byte, error) {
    var info []byte
    buf := bytes.NewBuffer(info)
    buf.WriteString("# Keyspace\r\n")
    dbStats := d.context.StatsManager.GetDbStats()
    for i := range dbStats {
        buf.WriteString(fmt.Sprintf("db%d:keys=%d,expires=0,avg_ttl=0\r\n", i, dbStats[i]["keys"]))
    }
    return d.context.Encoder.Encode(buf.String(), false)
}

func (d *defaultEvaluator) evaluateClient(args []string) ([]byte, error) {
    return d.context.Encoder.Encode("OK", true)
}

func (d *defaultEvaluator) evaluateLatency(args []string) ([]byte, error) {
    return d.context.Encoder.Encode([]string{}, false)
}
