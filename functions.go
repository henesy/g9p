package g9p

import (
	"encoding/binary"
	"bytes"
	// Heresy
	"fmt"
	"os"
)


const (
	// See: /sys/include/fcall.h
	Tversion byte	= iota + 100
	Rversion
	Tauth 		
	Rauth
	Tattach 	
	Rattach
	// Illegal
	Terror 		
	Rerror
	Tflush 	
	Rflush
	Twalk 		
	Rwalk
	Topen 		
	Ropen
	Tcreate 	
	Rcreate
	Tread 		
	Rread
	Twrite 	
	Rwrite
	Tclunk 	
	Rclunk
	Tremove 	
	Rremove
	Tstat 		
	Rstat
	Twstat 	
	Rwstat
	Tmax
	Invalid
)

// Lengths of various 9p elements in bytes -- see intro(5)
const SizeLen = 4
const TypeLen = 1
const TagLen = 2
const MsizeLen = 4

/* Snagged from droyo's styx */
// QidLen is the length of a Qid in bytes.
const QidLen = 13

// NoTag is the tag for Tversion and Rversion requests.
const NoTag = ^uint16(0)

// NoFid is a reserved fid used in a Tattach request for the
// afid field, that indicates that the client does not wish
// to authenticate his session.
const NoFid = ^uint32(0)

// Flags for the mode field in Topen and Tcreate messages
const (
	OREAD   = 0  // open read-only
	OWRITE  = 1  // open write-only
	ORDWR   = 2  // open read-write
	OEXEC   = 3  // execute (== read but check execute permission)
	OTRUNC  = 16 // or'ed in (except for exec), truncate file first
	OCEXEC  = 32 // or'ed in, close on exec
	ORCLOSE = 64 // or'ed in, remove on close
)

// File modes
const (
	DMDIR    = 0x80000000 // mode bit for directories
	DMAPPEND = 0x40000000 // mode bit for append only files
	DMEXCL   = 0x20000000 // mode bit for exclusive use files
	DMMOUNT  = 0x10000000 // mode bit for mounted channel
	DMAUTH   = 0x08000000 // mode bit for authentication file
	DMTMP    = 0x04000000 // mode bit for non-backed-up file
	DMREAD   = 0x4        // mode bit for read permission
	DMWRITE  = 0x2        // mode bit for write permission
	DMEXEC   = 0x1        // mode bit for execute permission

	// Mask for the type bits
	DMTYPE = DMDIR | DMAPPEND | DMEXCL | DMMOUNT | DMTMP

	// Mask for the permissions bits
	DMPERM = DMREAD | DMWRITE | DMEXEC
)


type Msg struct {
	Size	uint32
	T		uint8
	Tag		uint16
	Payload	[]byte
	Orig	[]byte
}


// For debugging, mostly
func (m Msg) Print() {
	fmt.Println("Size:", m.Size)
	fmt.Println("Type:", m.T)
	fmt.Println("Tag:", m.Tag)
	fmt.Println("Payload size:", len(m.Payload))
}


// Identify a message's type and operate accordingly -- both srv and client use
func Parse(buf []byte) (Msg, byte) {
	var msg Msg
	
	// Size
	msg.Size = binary.LittleEndian.Uint32(buf[:SizeLen])
	
	msg.Orig = buf[:msg.Size]
	
	// Type
	err := binary.Read(bytes.NewReader(buf[SizeLen:SizeLen+TypeLen]), binary.LittleEndian, &msg.T)
	if err != nil {
		// Maybe add something for this later
		fmt.Fprintln(os.Stderr, "Error, unable to get type: ", err)
		return msg, Invalid
	}
	
	// Tag
	msg.Tag = binary.LittleEndian.Uint16(buf[SizeLen+TypeLen:SizeLen+TypeLen+TagLen])
	
	// Payload
	msg.Payload = buf[SizeLen+TypeLen+TagLen:msg.Size]

	return msg, msg.T
}

// Get extra fields for Tversion
func ReadTversion(msg Msg) (msize uint32, version string) {
	msize = binary.LittleEndian.Uint32(msg.Payload[:4])

	version = string(msg.Payload[MsizeLen+1:])

	return
}

// Write an Rversion -- Call after reading a Tversion ;; maybe move to Msg.Send() or otherwise break things up?
func (c *Conn9) Rversion(msg Msg) error {
	var buf []byte

	// Append type
	t := bytes.NewBuffer(make([]byte, TypeLen))
	
	err := binary.Write(t, binary.LittleEndian, Rversion)
	if err != nil {
		// Maybe add something for this later
		fmt.Fprintln(os.Stderr, "Error, unable to set type: ", err)
		return err
	}
	
	buf = append(buf, t.Bytes()...)

	// Append tag
	tag := bytes.NewBuffer(make([]byte, TagLen))
	
	err = binary.Write(tag, binary.LittleEndian, msg.Tag)
	if err != nil {
		// Maybe add something for this later
		fmt.Fprintln(os.Stderr, "Error, unable to set tag: ", err)
		return err
	}
	
	buf = append(buf, tag.Bytes()...)
	
	// Append msize
	msize := bytes.NewBuffer(make([]byte, MsizeLen))
	
	err = binary.Write(msize, binary.LittleEndian, c.Msize)
	if err != nil {
		// Maybe add something for this later
		fmt.Fprintln(os.Stderr, "Error, unable to set msize: ", err)
		return err
	}
	
	buf = append(buf, msize.Bytes()...)

	// Append version
	buf = append(buf, []byte(c.Version)...)
	
	// Prepend size
	// 4 bytes for size
	sizeBytes := uint32(len(buf) + SizeLen)
	size := bytes.NewBuffer(make([]byte, SizeLen))
	
	err = binary.Write(size, binary.LittleEndian, sizeBytes)
	if err != nil {
		// Maybe add something for this later
		fmt.Fprintln(os.Stderr, "Error, unable to set size: ", err)
		return err
	}

	buf = append(size.Bytes(), buf...)
	
	fmt.Println(msg.Orig)
	fmt.Println(sizeBytes)
	fmt.Println(buf)
	
	_, err = c.Conn.Write(buf)
	return err
}
