package redis

import (
	"bufio"
	"bytes"
	"errors"
	"reflect"
	"testing"
)

func TestEncode(t *testing.T) {
	testCases := []struct {
		input    []interface{}
		expected []byte
		err      error
	}{
		{
			input:    []interface{}{"GET", "key"},
			expected: []byte("*2\r\n$3\r\nGET\r\n$3\r\nkey\r\n"),
			err:      nil,
		},
		{
			input:    []interface{}{"SET", "key", "value"},
			expected: []byte("*3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n"),
			err:      nil,
		},
		{
			input:    []interface{}{"INCR", "counter"},
			expected: []byte("*2\r\n$4\r\nINCR\r\n$7\r\ncounter\r\n"),
			err:      nil,
		},
		{
			input:    []interface{}{"LPUSH", "list", "value1", "value2"},
			expected: []byte("*4\r\n$5\r\nLPUSH\r\n$4\r\nlist\r\n$6\r\nvalue1\r\n$6\r\nvalue2\r\n"),
			err:      nil,
		},
	}

	for _, tc := range testCases {
		output, err := Encode(tc.input...)
		if !bytes.Equal(output, tc.expected) || !errors.Is(err, tc.err) {
			t.Errorf("Encode(%v) = (%v, %v), expected (%v, %v)", tc.input, output, err, tc.expected, tc.err)
		}
	}
}

func TestDecode(t *testing.T) {
	tests := []struct {
		input  string
		output interface{}
		err    error
	}{
		{
			input:  "+OK\r\n",
			output: "OK",
			err:    nil,
		},
		{
			input:  "-error message\r\n",
			output: nil,
			err:    errors.New("error message"),
		},
		{
			input:  ":0\r\n",
			output: int64(0),
			err:    nil,
		},
		{
			input:  ":123\r\n",
			output: int64(123),
			err:    nil,
		},
		{
			input:  ":-9223372036854775808\r\n",
			output: int64(-9223372036854775808),
			err:    nil,
		},
		{
			input:  ":9223372036854775807\r\n",
			output: int64(9223372036854775807),
			err:    nil,
		},
		{
			input:  "$11\r\nhello world\r\n",
			output: "hello world",
			err:    nil,
		},
		{
			input:  "$12\r\nhello\r\nworld\r\n",
			output: "hello\r\nworld",
			err:    nil,
		},
		{
			input:  "$0\r\n",
			output: "",
			err:    nil,
		},
		{
			input:  "$-1\r\n",
			output: nil,
			err:    nil,
		},
		{
			input:  "*9\r\n+OK\r\n:0\r\n:123\r\n:-9223372036854775808\r\n:9223372036854775807\r\n$11\r\nhello world\r\n$12\r\nhello\r\nworld\r\n$0\r\n$-1\r\n",
			output: []interface{}{"OK", int64(0), int64(123), int64(-9223372036854775808), int64(9223372036854775807), "hello world", "hello\r\nworld", "", nil},
			err:    nil,
		},
	}

	for _, test := range tests {
		reader := bufio.NewReader(bytes.NewBufferString(test.input))
		output, err := Decode(reader)
		if !reflect.DeepEqual(output, test.output) {
			t.Errorf("Decode(%q) output = %v, want %v", test.input, output, test.output)
		}
		if !reflect.DeepEqual(err, test.err) {
			t.Errorf("Decode(%q) err = %v, want %v", test.input, err, test.err)
		}
	}
}
