package redis

import (
	"testing"
)

func TestPing(t *testing.T) {
	client, err := NewClient("localhost:6379")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer func(client *Client) {
		err := client.Close()
		if err != nil {
		}
	}(client)

	resp, err := client.Ping()
	if err != nil {
		t.Fatalf("Failed to execute Ping command: %v", err)
	}

	expected := "PONG"
	if resp != expected {
		t.Fatalf("Expected response to be %q, but got %q", expected, resp)
	}
}

func TestAppend(t *testing.T) {
	client, err := NewClient("localhost:6379")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer func(client *Client) {
		err := client.Close()
		if err != nil {
		}
	}(client)

	key := "test_key"
	value := "test_value"
	resp, err := client.Append(key, value)
	if err != nil {
		t.Fatalf("Failed to execute Append command: %v", err)
	}

	expected := int64(len(value))
	if resp != expected {
		t.Fatalf("Expected response to be %d, but got %d", expected, resp)
	}
}

func TestDecrby(t *testing.T) {
	client, err := NewClient("localhost:6379")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer func(client *Client) {
		err := client.Close()
		if err != nil {
		}
	}(client)

	resp, err := client.Decrby("num", int64(100))
	if err != nil {
		t.Fatalf("Failed to execute Decrby command: %v", err)
	}

	expected := int64(-100)
	if resp != expected {
		t.Fatalf("Expected response to be %d, but got %d", expected, resp)
	}
}
