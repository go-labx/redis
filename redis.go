package redis

import (
	"bufio"
	"fmt"
	"net"
)

type Client struct {
	conn   net.Conn
	reader *bufio.Reader
}

type ClientOptions struct {
	Addr string
}

func NewClient(addr string) (*Client, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &Client{
		conn:   conn,
		reader: bufio.NewReader(conn),
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) SendCommand(args []interface{}) (interface{}, error) {
	data, err := Encode(args...)
	if err != nil {
		return nil, err
	}
	fmt.Println("write data ---->")
	fmt.Println(data)
	fmt.Println(string(data))

	if _, err = c.conn.Write(data); err != nil {
		return nil, err
	}

	resp, err := Decode(c.reader)
	if err != nil {
		return nil, err
	}
	fmt.Println("resp ------>", resp)

	return resp, nil
}

func (c *Client) Ping() (string, error) {
	resp, err := c.SendCommand([]interface{}{"Ping"})
	if err != nil {
		return "", err
	}
	return resp.(string), nil
}

func (c *Client) Append(key string, value string) (int64, error) {
	resp, err := c.SendCommand([]interface{}{"Append", key, value})
	if err != nil {
		return 0, err
	}
	return resp.(int64), nil
}

func (c *Client) Decr(key string) (int64, error) {
	resp, err := c.SendCommand([]interface{}{"Decr", key})
	if err != nil {
		return 0, err
	}
	return resp.(int64), nil
}

func (c *Client) Decrby(key string, decrement int64) (int64, error) {
	resp, err := c.SendCommand([]interface{}{"Decrby", key, decrement})
	if err != nil {
		return 0, err
	}
	return resp.(int64), nil
}

func (c *Client) Del(keys ...string) (int, error) {
	resp, err := c.SendCommand(append([]interface{}{"DEL"}, keys...))
	if err != nil {
		return 0, err
	}
	return resp.(int), nil
}

func (c *Client) Get(key string) (string, error) {
	resp, err := c.SendCommand([]interface{}{"GET", key})
	if err != nil {
		return "", err
	}
	if resp == nil {
		return "", nil
	}
	return resp.(string), nil
}
