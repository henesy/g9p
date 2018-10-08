package g9p

import (
	"encoding/binary"
	// Heresy
	"fmt"
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
const PrefixLen = SizeLen+TypeLen+TagLen

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
	Full	[]byte
	Extra	map[string]interface{}
}


// For debugging, mostly
func (m Msg) Print() {
	fmt.Println("---")
	fmt.Println("Size:", m.Size)
	fmt.Println("Type:", m.T)
	fmt.Println("Tag:", m.Tag)
	fmt.Println("Payload size:", len(m.Payload()))
	fmt.Println("---")
}

// Returns the payload of a Msg
func (m Msg) Payload() []byte {
	return m.Full[PrefixLen:]
}

// Identify a message's type and operate accordingly -- both srv and client use
func Parse(buf []byte) (Msg, byte) {
	var msg Msg
	
	// Size
	msg.Size = binary.LittleEndian.Uint32(buf[:SizeLen])
	
	msg.Full = buf[:msg.Size]
	
	// Type
	msg.T = BytesToByte(buf[SizeLen:SizeLen+TypeLen])
	
	// Tag
	msg.Tag = binary.LittleEndian.Uint16(buf[SizeLen+TypeLen:PrefixLen])
	
	// Payload
	//msg.Payload = buf[PrefixLen:msg.Size]
	
	// Extra
	msg.Extra = make(map[string]interface{})

	return msg, msg.T
}

// Get extra fields for Tversion
func (msg *Msg) ReadTversion() (msize uint32, version string) {
	msize = binary.LittleEndian.Uint32(msg.Payload()[:4])
	msg.Extra["msize"] = msize

	version = string(msg.Payload()[MsizeLen+1:])
	msg.Extra["version"] = version

	return
}

// Create prefix for a Msg
func MkMsg(msgT byte, msgTag uint16) Msg {
	var msg Msg
	var buf []byte
	
	// Append type
	typeBuf := ByteToBytes(msgT)
	fmt.Println(typeBuf)
	buf = append(buf, typeBuf...)

	// Append tag
	tagBuf := U16ToBytes(msgTag)
	buf = append(buf, tagBuf...)
	
	msg.Full = buf
	msg.Tag = msgTag
	msg.T = msgT
	
	return msg
}

// Write an Rversion -- Call after reading a Tversion
func MkRversion(msg Msg) (Msg) {
	rmsg := MkMsg(Rversion, msg.Tag)

	// Append msize
	msizeBuf := U32ToBytes((msg.Extra["msize"]).(uint32))
	rmsg.Full = append(rmsg.Full, msizeBuf...)

	// Append version
	vStr := (msg.Extra["version"]).(string)
	vBuf := make([]byte, 1)

	// This is magic discovered through trial and error, why does this exist?
	vBuf[0] = byte(0x6)
	vBuf = append(vBuf, []byte(vStr)[:len(vStr)]...)
	rmsg.Full = append(rmsg.Full, vBuf...)

	// Prepend size
	// 4 bytes for size
	sizeBytes := uint32(len(rmsg.Full) + SizeLen)
	sizeBuf := U32ToBytes(sizeBytes)
	rmsg.Full = append(sizeBuf, rmsg.Full...)

	// Load rmsg
	rmsg.Size = sizeBytes
	rmsg.Tag = msg.Tag
	rmsg.T = Rversion
	
	return rmsg
}
