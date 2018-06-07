package util

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"
)

// ParseInstruction parse modbus rtu instruction
func ParseInstruction(ins string) ([]byte, error) {
	instruction := []byte{}
	if ins == "" {
		return nil, fmt.Errorf("parse instruction error: instruction is empty")
	}

	for i := 2; i <= len(ins); i += 2 {
		t, err := parseByte(ins[i-2 : i])
		if err != nil {
			return nil, fmt.Errorf("parse instruction error: %v", err)
		}
		instruction = append(instruction, t)
	}
	return instruction, nil
}

func parseByte(byteStr string) (byte, error) {
	bufIns, err := strconv.ParseUint(byteStr, 16, 8)
	if err != nil {
		return 0, fmt.Errorf("parse byte %v error: %v", byteStr, err)
	}
	i := bytes.NewBuffer([]byte{})
	binary.Write(i, binary.BigEndian, &bufIns)
	return i.Bytes()[7], nil
}
