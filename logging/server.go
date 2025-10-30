package logging

import (
	"fmt"
	"io"
	"net"
	"sync"
)

const (
	socketPath = "./ctop.sock"
	socketAddr = "0.0.0.0:9000"
)

var server struct {
	wg sync.WaitGroup
	ln net.Listener
}

func getListener() net.Listener {
	var ln net.Listener
	var err error
	if debugModeTCP() {
		ln, err = net.Listen("tcp", socketAddr)
	} else {
		ln, err = net.Listen("unix", socketPath)
	}
	if err != nil {
		panic(err)
	}
	return ln
}

func StartServer() {
	server.ln = getListener()

	go func() {
		for {
			conn, err := server.ln.Accept()
			if err != nil {
				if ne, ok := err.(net.Error); ok && ne.Timeout() {
					continue
				}
				return
			}
			go handler(conn)
		}
	}()

	Log.Notice("logging server started")
}

func StopServer() {
	server.wg.Wait()
	if server.ln != nil {
		if err := server.ln.Close(); err != nil {
			Log.Errorf("failed to close log server listener: %s", err)
		}
	}
}

func handler(wc io.WriteCloser) {
	server.wg.Add(1)
	defer server.wg.Done()
	defer func() {
		if err := wc.Close(); err != nil {
			Log.Errorf("failed to close log handler: %s", err)
		}
	}()
	for msg := range Log.tail() {
		msg = fmt.Sprintf("%s\n", msg)
		if _, err := wc.Write([]byte(msg)); err != nil {
			Log.Errorf("failed to write to log handler: %s", err)
		}
	}
	if _, err := wc.Write([]byte("bye\n")); err != nil {
		Log.Errorf("failed to write to log handler: %s", err)
	}
}
