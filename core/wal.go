package core

import (
	"log"
	"os"

	"github.com/Bihan001/MyDB/config"
)

var walFile *os.File

func init() {
	var err error
	walFile, err = os.OpenFile(config.WALFilePath, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		log.Fatal("error while opening WAL file: ", err)
	}
}

func WriteOperationInWAL(input []byte) {
	// The input slice has a hardcoded size.
	// If the actual string is smaller, the rest of the characters are null terminators which break the WAL
	inputLength := getNullTerminatorIndex(input)
	data := input[:inputLength]
	_, err := walFile.Write(data)
	if err != nil {
		log.Fatal("error while writing WAL: ", err)
	}
}

func ReadWAL() []byte {
	bytes, err := os.ReadFile(config.WALFilePath)
	if err != nil {
		log.Fatal("error while reading WAL file: ", err)
	}
	return bytes
}

func getNullTerminatorIndex(data []byte) int {
	for i := 0; i < len(data); i++ {
		if data[i] == 0 {
			return i
		}
	}
	return len(data)
}