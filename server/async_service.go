package server

import (
	"log"
	"net"
	"syscall"

	"github.com/Bihan001/MyDB/config"
	"github.com/Bihan001/MyDB/core"
)

type asyncService struct {
    context *core.Context
    evaluator core.Evaluator
}

func NewAsyncService(ctx *core.Context, evaluator core.Evaluator) *asyncService {
    return &asyncService{
        context: ctx,
        evaluator: evaluator,
    }
}

func (srv *asyncService) RunService() error {
    preRun(srv.context, srv.evaluator)

    log.Println("starting an asynchronous TCP server on", config.ServerHost, config.ServerPort)

    var activeConnections int = 0
    events := make([]syscall.EpollEvent, config.ConnectionLimit)

    serverFd, err := syscall.Socket(syscall.AF_INET, syscall.O_NONBLOCK|syscall.SOCK_STREAM, 0)
    if err != nil {
        return err
    }
    defer syscall.Close(serverFd)

    if err = syscall.SetNonblock(serverFd, true); err != nil {
        return err
    }

    addr := syscall.SockaddrInet4{Port: config.ServerPort}
    copy(addr.Addr[:], net.ParseIP(config.ServerHost).To4())

    if err = syscall.Bind(serverFd, &addr); err != nil {
        return err
    }

    if err = syscall.Listen(serverFd, config.ConnectionLimit); err != nil {
        return err
    }

    epfd, err := syscall.EpollCreate1(0)
    if err != nil {
        return err
    }
    defer syscall.Close(epfd)

    event := syscall.EpollEvent{
        Events: syscall.EPOLLIN,
        Fd:     int32(serverFd),
    }

    if err = syscall.EpollCtl(epfd, syscall.EPOLL_CTL_ADD, serverFd, &event); err != nil {
        return err
    }

    for {
        srv.context.ExpiryManager.PurgeExpiredEntries()
        nevents, err := syscall.EpollWait(epfd, events[:], -1)
        if err != nil {
            log.Println("epoll_wait: ", err)
            continue
        }

        for i := 0; i < nevents; i++ {
            if events[i].Fd == int32(serverFd) {
                fd, _, err := syscall.Accept(serverFd)
                if err != nil {
                    log.Println(err)
                    continue
                }

                activeConnections++
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
                var clientFd int = int(events[i].Fd)
                buffer := make([]byte, 512)
                var ops core.OperationList
                var reply []byte

                n, err := syscall.Read(clientFd, buffer)
                if err == nil {
                    ops, err = parseCommands(buffer, n, srv.context.Decoder, srv.context.WAL)
                }

                if err == nil {
                    reply, err = srv.evaluator.Evaluate(ops)
                }

                if err == nil {
                    _, err = syscall.Write(clientFd, reply)
                }

                if err != nil {
                    closeClientConnection(clientFd)
                    activeConnections--
                }
            }
        }
    }
}

func closeClientConnection(fd int) {
    if err := syscall.Close(fd); err != nil {
        log.Println(err)
    }
}
