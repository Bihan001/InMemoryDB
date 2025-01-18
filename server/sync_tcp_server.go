package server

import (
	"fmt"
	"log"
	"net"
	"strconv"

	"github.com/Bihan001/MyDB/config"
	"github.com/Bihan001/MyDB/core"
)

type SyncTCPServer struct {
}

func (tcpServer *SyncTCPServer) RunServer() error {
	log.Println("starting a synchronous TCP server on", config.Host, config.Port)

	var concurrent_clients int = 0

	lsnr, err := net.Listen("tcp", config.Host+":"+strconv.Itoa(config.Port))
	if err != nil {
		return err
	}

	for {
		// blocking call: waiting for the new client to connect
		connection, err := lsnr.Accept()
		if err != nil {
			continue
		}

		concurrent_clients++

		fmt.Println("A client connected")

		for {
			buff := make([]byte, 512)
			var cmds core.Cmds
			var response []byte

			n, err := connection.Read(buff)
			if err == nil {
				cmds, err = readCommands(buff, n)
			}

			if err == nil {
				response, err = core.Evaluate(cmds)
			}

			if err == nil {
				_, err = connection.Write(response)
			}

			if err != nil {
				fmt.Println(err)
				connection.Close()
				concurrent_clients--
				break
			}
		}
	}
}
