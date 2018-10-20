package g9p

import (
	"net"
	"log"
	"os"
	"time"
)

// Globals suck, extra arguments to everything suck more
var Log		*log.Logger
var Debug	bool

// Represents a connection over 9p
type Conn9 struct {
	Conn	net.Conn
	Version	string
	Msize	uint32
}

// Represents a server for 9p
type Srv struct {
	L			net.Listener
	// Supported versions
	Versions	[]string
	// Default msize
	Msize		uint32
	Fs			FS
}


// Send Msg over connection (and log)
func (c *Conn9) Send(msg Msg) error {
	if Debug {
		Chatty(msg, c)
	}

	_, err := c.Conn.Write(msg.Full)
	return err
}

// Read and parse a message
func (c *Conn9) Read() (msg Msg, err error) {
	err = nil
	buf := make([]byte, c.Msize)
	_, err = c.Conn.Read(buf)
	msg = Parse(buf)
	return msg, err
}

// Start a new listener and process 9p connections ;; takes a list of supported 9p versions
func MkSrv(protocol string, port string, versions ...string) (Srv, error) {
	var srv Srv
	listener, err := net.Listen(protocol, ":" + port)
	if err != nil {
		log.Print("Error, unable to open listener: ", err)
		return srv, err
	}
	Log = log.New(os.Stderr, "", log.Ldate | log.Ltime | log.Lshortfile)
	Debug = true
	srv.L = listener
	srv.Versions = versions
	// Sensible default
	srv.Msize = 8216
	//srv.Fs = MkFs()

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
			Log.Print("Error, unable to accept conn: ", err)
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
	msg, err := c.Read()
	if err != nil {
		Log.Println("Error, bad read from client: ", err)
		return
	}
	if msg.T == Tversion {
		// Accept versions -- should check if supported
		c.Msize, c.Version = msg.Extra["msize"].(uint32), msg.Extra["version"].(string)
		
		// TODO - Find a way to ID client nicely
		rmsg := MkRversion(msg)
		//Log.Println(s.Fs.Get("/").Name())
		
		err := c.Send(rmsg)
		if err != nil {
			Log.Print("Error, sending Rversion: ", err)
			return
		}
	} else {
		c.Send(MkRerror("Expected Tversion", &msg))
		return
	}

	for {
		// Read from client
		
		// Determine function to perform (if any)
		
		// Respond to client
	
		time.Sleep(5 * time.Millisecond)
	}
}
