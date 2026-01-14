package socket

import (
	"net"
	"strings"
	"testing"
	"time"
)

var (
	test_message = "Hello\n"
	test_address = "127.0.0.1:34001"
)

func TestServerReceiveMEssage(t *testing.T) {
	Start(test_address)

	conn, err := net.Dial("tcp", test_address)

	if err != nil {
		t.Fatal("Error connecting to socket")
	}
	defer conn.Close()

	_, err = conn.Write([]byte(test_message))
	if err != nil {
		t.Fatal("Error writing to socket")
	}

	retries := 0

	for {
		var received string
		received, err = MessagePop()
		if err != nil {
			retries += 1
			if retries > 3 {
				t.Error("Error retrieving message from queue")
			} else {
				time.Sleep(200 * time.Millisecond)
			}
		} else {
			if received != strings.TrimSuffix(test_message, "\n") {
				t.Error(received + " != " + test_message)
				t.Fatal("Message verification failed")
			} else {
				break
			}
		}
	}
}
