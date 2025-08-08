package binaryutils

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

func WriteString(buf *bytes.Buffer, s string) error {
	bt := []byte(s)

	if err := binary.Write(buf, binary.LittleEndian, uint32(len(bt))); err != nil {
		return fmt.Errorf("failed to write string length: %w", err)
	}
	if err := binary.Write(buf, binary.LittleEndian, bt); err != nil {
		return fmt.Errorf("failed to write string data: %w", err)
	}
	return nil
}

func ReadString(r *bytes.Reader) (string, error) {
	var length uint32
	if err := binary.Read(r, binary.LittleEndian, &length); err != nil {
		return "", fmt.Errorf("failed to read string length: %w", err)
	}

	data := make([]byte, length)
	if err := binary.Read(r, binary.LittleEndian, &data); err != nil {
		return "", fmt.Errorf("failed to read string data: %w", err)
	}

	return string(data), nil
}

func WriteBytesWithLength(buf *bytes.Buffer, data []byte) error {
	if err := binary.Write(buf, binary.LittleEndian, uint32(len(data))); err != nil {
		return err
	}
	return binary.Write(buf, binary.LittleEndian, data)
}

func ReadBytesWithLength(r *bytes.Reader) ([]byte, error) {
	var length uint32
	if err := binary.Read(r, binary.LittleEndian, &length); err != nil {
		return nil, err
	}
	data := make([]byte, length)
	if err := binary.Read(r, binary.LittleEndian, &data); err != nil {
		return nil, err
	}
	return data, nil
}