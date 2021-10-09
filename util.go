package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

func uinttobyte(source uint32) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, source)
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}
	//	fmt.Printf("Encoded: % x\n", buf.Bytes())
	return buf.Bytes()
}
