package wal

import (
    "log"
    "os"

    "github.com/Bihan001/MyDB/internal/config"
)

type WAL interface {
    WriteToWAL(input []byte)
    ReadWALFile() []byte
}

type fileWAL struct {
    fileHandle *os.File
    filePath   string
}

var defaultWAL WAL = func() WAL {
    handle, err := os.OpenFile(config.LogFilePath, os.O_CREATE|os.O_RDWR, 0600)
    if err != nil {
        log.Fatal("error while opening WAL file: ", err)
    }
    return &fileWAL{
        fileHandle: handle,
        filePath:   config.LogFilePath,
    }
}()

func GetWAL() WAL {
    return defaultWAL
}

func (fw *fileWAL) WriteToWAL(input []byte) {
    size := fw.indexOfNullTerminator(input)
    data := input[:size]
    _, err := fw.fileHandle.Write(data)
    if err != nil {
        log.Fatal("error while writing WAL: ", err)
    }
}

func (fw *fileWAL) ReadWALFile() []byte {
    bytes, err := os.ReadFile(fw.filePath)
    if err != nil {
        log.Fatal("error while reading WAL file: ", err)
    }
    return bytes
}

func (fw *fileWAL) indexOfNullTerminator(data []byte) int {
    for i := 0; i < len(data); i++ {
        if data[i] == 0 {
            return i
        }
    }
    return len(data)
}
