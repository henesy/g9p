package g9p

import (
	"net"
	"log"
	"os"
)


// Represents a connection over 9p
type Conn9 struct {
	conn	net.Conn
	version	string
}

// Represents a server for 9p
type Srv struct {
	LogFile	*os.File
	L		net.Listener
}


// Start a new listener and process 9p connections
func MkSrv(protocol string, port string) (Srv, error) {
	var srv Srv
	listener, err := net.Listen(protocol, port)
	if err != nil {
		log.Print("Error, unable to open listener, ", err)
		return srv, err
	}
	srv.LogFile = os.Stderr
	srv.L = listener
	return srv, nil
}

