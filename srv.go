package g9p

import (
	"net"
	"log"
	"os"
	"time"
)


// Represents a connection over 9p
type Conn9 struct {
	Conn	net.Conn
	Version	string
	Msize	uint32
}

// Represents a server for 9p
type Srv struct {
	Log		*log.Logger
	L		net.Listener
	// Supported versions
	Versions	[]string
	// Default msize
	Msize		uint32
	Debug		bool
}


// Start a new listener and process 9p connections ;; takes a list of supported 9p versions
func MkSrv(protocol string, port string, versions ...string) (Srv, error) {
	var srv Srv
	listener, err := net.Listen(protocol, ":" + port)
	if err != nil {
		log.Print("Error, unable to open listener: ", err)
		return srv, err
	}
	srv.Log = log.New(os.Stderr, "", log.Ldate | log.Ltime | log.Lshortfile)
	srv.L = listener
	srv.Versions = versions
	// Sensible default
	srv.Msize = 8216
	srv.Debug = true

	return srv, nil
}

// Initialize the server ;; splitting the two gives time to adjust defaults
func (s *Srv) Init() {
	go s.Listener()
}

// Listen for incoming 9p connections
func (s *Srv) Listener() {
	for {
		conn, err := s.L.Accept()
		if err != nil {
			s.Log.Print("Error, unable to accept conn: ", err)
			continue
		}

		var conn9 Conn9
		conn9.Conn = conn
		conn9.Msize = s.Msize

		go s.Handler(conn9)
	}
}

// Handle incoming 9p connections
func (s *Srv) Handler(c Conn9) {
	// Negotiate version
	buf := make([]byte, c.Msize)
	c.Conn.Read(buf)
	if msg, mt := Parse(buf); mt == Tversion {
		// Read Tversion call
		c.Msize, c.Version = msg.ReadTversion()
		
		// Find a way to ID client nicely
		if s.Debug {
			//msg.Print()
			s.Log.Printf("← Tversion tag=%d msize=%d version=\"%s\"", msg.Tag, c.Msize, c.Version)
		}
		
		rmsg, err := MkRversion(msg)
		if err != nil {
			s.Log.Print("Error, failure creating Rversion: ", err)
			return
		}
		
		if s.Debug {
			//rmsg.Print()
			s.Log.Printf("→ Rversion tag=%d msize=%d version=\"%s\"", msg.Tag, c.Msize, c.Version)
		}
		
		_, err = c.Conn.Write(rmsg.Full)
		if err != nil {
			s.Log.Print("Error, sending Rversion: ", err)
			return
		}
	} else {
		c.Rerror("Expected Tversion")
		return
	}

	for {
		// Read from client
		
		// Determine function to perform (if any)
		
		// Respond to client
	
		time.Sleep(5 * time.Millisecond)
	}
}

// The only one allowed to throw errors is me, Dio
func (c *Conn9) Rerror(msg string) error {
	// TODO
	c.Conn.Write([]byte(msg))

	return nil
}
