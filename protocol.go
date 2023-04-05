package redis

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

const (
	SimpleString = '+'
	ErrorString  = '-'
	Integer      = ':'
	BulkString   = '$'
	Array        = '*'
	LF           = '\n'
	CRLF         = "\r\n"
)

var (
	ErrInvalidPrefix = errors.New("invalid prefix")
)

func Encode(args ...interface{}) ([]byte, error) {
	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)

	if _, err := fmt.Fprintf(writer, "*%d%s", len(args), CRLF); err != nil {
		return nil, err
	}

	for _, arg := range args {
		switch arg := arg.(type) {
		case int, int8, int16, int32, uint, uint8, uint16, uint32, uint64:
			if _, err := fmt.Fprintf(writer, ":%d%s", arg, CRLF); err != nil {
				return nil, err
			}
		case float32, float64, bool:
			if _, err := fmt.Fprintf(writer, "+%v%s", arg, CRLF); err != nil {
				return nil, err
			}
		case int64:
			str := strconv.Itoa(int(arg))
			if _, err := fmt.Fprintf(writer, "$%d%s%s%s", len(str), CRLF, str, CRLF); err != nil {
				return nil, err
			}
		case string:
			if _, err := fmt.Fprintf(writer, "$%d%s%s%s", len(arg), CRLF, arg, CRLF); err != nil {
				return nil, err
			}
		case []byte:
			if _, err := fmt.Fprintf(writer, "$%d%s%s%s", len(arg), CRLF, arg, CRLF); err != nil {
				return nil, err
			}
		case nil:
			if _, err := fmt.Fprintf(writer, "$-1%s", CRLF); err != nil {
				return nil, err
			}
		case []interface{}:
			if _, err := fmt.Fprintf(writer, "*%d%s", len(arg), CRLF); err != nil {
				return nil, err
			}
			for _, elem := range arg {
				b, err := Encode(elem)
				if err != nil {
					return nil, err
				}
				if _, err := writer.Write(b); err != nil {
					return nil, err
				}
			}
		default:
			return nil, fmt.Errorf("unsupported type: %T", arg)
		}
	}

	if err := writer.Flush(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func Decode(reader *bufio.Reader) (interface{}, error) {
	prefix, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}

	switch prefix {
	case SimpleString:
		if line, err := reader.ReadString(LF); err != nil {
			return nil, err
		} else {
			return strings.TrimRight(line, CRLF), nil
		}
	case ErrorString:
		line, err := reader.ReadString(LF)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(strings.TrimRight(line, CRLF))
	case Integer:
		line, err := reader.ReadString(LF)
		if err != nil {
			return nil, err
		}
		return strconv.ParseInt(line[:len(line)-2], 10, 64)
	case BulkString:
		lengthLine, err := reader.ReadString(LF)
		if err != nil {
			return nil, err
		}
		length, err := strconv.Atoi(lengthLine[:len(lengthLine)-2])
		if err != nil {
			return nil, err
		}
		if length == 0 {
			return "", nil
		}
		if length == -1 {
			return nil, nil
		}
		data := make([]byte, length+2)
		if _, err := io.ReadFull(reader, data); err != nil {
			return nil, err
		}
		return string(data[:len(data)-2]), nil
	case Array:
		lengthLine, err := reader.ReadString(LF)
		if err != nil {
			return nil, err
		}
		length, err := strconv.Atoi(lengthLine[:len(lengthLine)-2])
		if err != nil {
			return nil, err
		}
		if length == -1 {
			return nil, nil
		}
		array := make([]interface{}, length)
		for i := 0; i < length; i++ {
			value, err := Decode(reader)
			if err != nil {
				return nil, err
			}
			array[i] = value
		}
		return array, nil
	default:
		return nil, ErrInvalidPrefix
	}

}
