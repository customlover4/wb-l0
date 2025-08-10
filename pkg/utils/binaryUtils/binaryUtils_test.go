package binaryutils

import (
	"bytes"
	"testing"
)

type TestCase struct {
	Data string
}

func TestWriteReadString(t *testing.T) {
	test := []TestCase{
		{
			"test",
		},
		{
			"sdfgsdh",
		},
		{
			"2335235",
		},
		{
			"qwrfasd",
		},
		{
			"a",
		},
		{
			"teaerhibj;aejaeoibroiajeoibjoiajboiaejoibst",
		},
		{
			"sdfkj",
		},
	}

	for _, v := range test {
		buf := new(bytes.Buffer)
		err := WriteString(buf, v.Data)
		if err != nil {
			t.Error(err.Error())
		}

		r := bytes.NewReader(buf.Bytes())
		if data, err := ReadString(r); err != nil || data != v.Data {
			t.Errorf(
				"Error: %s\nGiven data:%s\nTaken data: %s",
				err.Error(), v.Data, data,
			)
		}
	}
}

type TestCase2 struct {
	Data []byte
}

func TestWriteReadWithBytesLength(t *testing.T) {
	test := []TestCase2{
		{
			[]byte("test"),
		},
		{
			[]byte("test"),
		},
		{
			[]byte("test"),
		},
		{
			[]byte("test"),
		},
		{
			[]byte("test"),
		},
		{
			[]byte("test"),
		},
		{
			[]byte("test"),
		},
	}

	for _, v := range test {
		buf := new(bytes.Buffer)
		err := WriteBytesWithLength(buf, v.Data)
		if err != nil {
			t.Error(err.Error())
		}

		r := bytes.NewReader(buf.Bytes())
		data, err := ReadBytesWithLength(r)
		if err != nil {
			t.Errorf(
				"Error: %s", err.Error(),
			)
		}

		if len(data) != len(v.Data) {
			t.Errorf(
				"Len is not compared\nGiven data len: %d\nTaken data len: %d",
				len(v.Data), len(data),
			)
		}

		for i := 0; i < len(data); i++ {
			if data[i] != v.Data[i] {
				t.Errorf(
					"Value compared\nGiven data bt: %v\nTaken data bt: %v",
					v.Data[i], data[i],
				)
			}
		}
	}
}
