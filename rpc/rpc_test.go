package rpc_test

import (
	"simple-lsp/rpc"
	"testing"
)

type EncodingExample struct {
	Testing bool
}

func TestEncode(t *testing.T) {
	expected := "Content-Length: 16\r\n\r\n{\"Testing\":true}"
	actual := rpc.EncodeMessage(EncodingExample{Testing: true})

	if expected != actual {
		t.Fatalf("expected: %s, Actual: %s", expected, actual)
	}
}

func TestDecode(t *testing.T) {
	incomingMessage := "Content-Length: 15\r\n\r\n{\"method\":\"hi\"}"

	method, content, err := rpc.DecodeMessage([]byte(incomingMessage))
	if err != nil {
		t.Fatal(err)
	}

	if len(content) != 15 {
		t.Fatalf("Expected 16, Got: %d", content)
	}

	if method != "hi" {
		t.Fatalf("Expected: hi, Got: %s", method)
	}
}
