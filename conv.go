package g9p

import (
	"fmt"
	"encoding/binary"
	"unsafe"
	"bytes"
	"os"
)


// Convert uint32 to []byte -- little endian
func U32ToBytes(num uint32) []byte {
	width := unsafe.Sizeof(num)
	buf := bytes.NewBuffer(make([]byte, width))
	
	err := binary.Write(buf, binary.LittleEndian, num)
	if err != nil {
		// Maybe add something for this later
		fmt.Fprintln(os.Stderr, "Error, unable to set uint32: ", err)
		return []byte{}
	}
	
	bytesBuf := buf.Bytes()[width:]

	return bytesBuf
}

// Convert uint16 to []byte -- little endian
func U16ToBytes(num uint16) []byte {
	width := unsafe.Sizeof(num)
	buf := bytes.NewBuffer(make([]byte, width))
	
	err := binary.Write(buf, binary.LittleEndian, num)
	if err != nil {
		// Maybe add something for this later
		fmt.Fprintln(os.Stderr, "Error, unable to set uint16: ", err)
		return []byte{}
	}
	
	bytesBuf := buf.Bytes()[width:]

	return bytesBuf
}

// Convert byte to []byte -- little endian
func ByteToBytes(msgB byte) []byte {
	width := unsafe.Sizeof(msgB)
	buf := bytes.NewBuffer(make([]byte, width))
	
	err := binary.Write(buf, binary.LittleEndian, msgB)
	if err != nil {
		// Maybe add something for this later
		fmt.Fprintln(os.Stderr, "Error, unable to set byte: ", err)
		return []byte{}
	}
	
	bytesBuf := buf.Bytes()[width:]
	return bytesBuf
}

// Convert []byte to byte -- little endian
func BytesToByte(buf []byte) byte {
	var b byte
	err := binary.Read(bytes.NewReader(buf), binary.LittleEndian, &b)
	if err != nil {
		// Maybe add something for this later
		fmt.Fprintln(os.Stderr, "Error, unable to get byte: ", err)
		return ^byte(0)
	}
	return b
}
