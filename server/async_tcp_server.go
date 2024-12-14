package server

import (
	"log"
	"net"
	"syscall"

	"github.com/Bihan001/MyDB/config"
	"github.com/Bihan001/MyDB/core"
)

type AsyncTCPServer struct {
}

func (tcpServer *AsyncTCPServer) RunServer() error {
	log.Println("starting a asynchronous TCP server on", config.Host, config.Port)

	var concurrent_clients int = 0

	// Create EPOLL Event Objects to hold events
	events := make([]syscall.EpollEvent, config.MaxClients)

	serverFd, err := syscall.Socket(syscall.AF_INET, syscall.O_NONBLOCK|syscall.SOCK_STREAM, 0)
	if err != nil {
		return err
	}

	defer syscall.Close(serverFd)

	// Set the server socket to operate in non-blocking mode
	if err = syscall.SetNonblock(serverFd, true); err != nil {
		return err
	}

	addr := syscall.SockaddrInet4{Port: config.Port}
	// IP addresses in Go are represented as 16-byte slices (net.IP), which can be used for both IPv4 and IPv6 addresses
	// IPv4 addresses are typically stored in a "IPv4-mapped IPv6 address" form, like ::ffff:192.0.2.1
	// To4() extracts the last 4 bytes of the 16-byte address,  making it a 4-byte IPv4 address like 192.0.2.1
	copy(addr.Addr[:], net.ParseIP(config.Host).To4())

	if err = syscall.Bind(serverFd, &addr); err != nil {
		return err
	}

	if err = syscall.Listen(serverFd, config.MaxClients); err != nil {
		return err
	}

	epfd, err := syscall.EpollCreate1(0)
	if err != nil {
		return err
	}
	defer syscall.Close(epfd)

	// Which events to get hints for and which file descriptor to watch for
	event := syscall.EpollEvent{
		Events: syscall.EPOLLIN,
		Fd:     int32(serverFd),
	}

	// Registering the event
	if err = syscall.EpollCtl(epfd, syscall.EPOLL_CTL_ADD, serverFd, &event); err != nil {
		return err
	}

	for {
		// Blocking call: Waiting until any file descriptor is ready for an I/O
		nevents, err := syscall.EpollWait(epfd, events[:], -1)
		if err != nil {
			log.Print("epoll_wait: ", err)
			continue
		}

		for i := 0; i < nevents; i++ {
			if events[i].Fd == int32(serverFd) {
				// The server is ready for I/O that means someone is trying to connect
				// Accept the new connection which creates a new fd
				// This fd is used to transfer data between the client and the server (1 fd per client)
				// Also register this fd to watch for changes

				fd, _, err := syscall.Accept(serverFd)
				if err != nil {
					log.Print(err)
					continue
				}

				concurrent_clients++

				if err = syscall.SetNonblock(fd, true); err != nil {
					log.Fatal(err)
					continue
				}

				event = syscall.EpollEvent{
					Events: syscall.EPOLLIN,
					Fd:     int32(fd),
				}

				if err = syscall.EpollCtl(epfd, syscall.EPOLL_CTL_ADD, fd, &event); err != nil {
					log.Fatal(err)
					continue
				}
			} else {
				// If events[i].Fd is not serverFd then it must be one of the client's Fd and is trying to send data
				// Read, evaluate and send response

				var clientServerFd int = int(events[i].Fd)
				buff := make([]byte, 512)
				var cmd *core.Cmd
				var err error
				var response []byte

				n, err := syscall.Read(clientServerFd, buff)

				if err == nil {
					cmd, err = readCommand(buff, n)
				}

				if err == nil {
					response, err = core.Evaluate(cmd)
				}

				if err == nil {
					_, err = syscall.Write(clientServerFd, response)
				}

				if err != nil {
					closeConnection(clientServerFd)
					concurrent_clients--
				}

			}
		}

	}

}

func closeConnection(fd int) {
	if err := syscall.Close(fd); err != nil {
		log.Print(err)
	}
}
