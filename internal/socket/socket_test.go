package socket

import (
	"testing"

	"github.com/mantzoun/yannena/internal/config"
)

// Interface defined for the test, to accomodate the helper functions
type TestMessenger interface {
	MessagePop() (string, error)
	MessageSend(string) (int, error)
}

var (
	test_message   = "Hello"
	test_message_2 = "Goodbye"
	test_config    = config.Config{
		ServerListenAddress:  "127.0.0.1:34001",
		ClientConnectAddress: "127.0.0.1:34001",
	}
)

func TestClientSendMessageServerDown(t *testing.T) {
	var err error

	client := NewClient(&test_config)

	client.Connect()
	_, err = client.MessageSend(test_message)
	if err != ErrNoConnection {
		t.Fatal("Unexpected error")
	}

	client.Disconnect()
	client.Stop()

	waitForClientStopped(client)
}

func TestClientSendMessageOK(t *testing.T) {
	client := NewClient(&test_config)
	server := NewServer(&test_config)

	server.Start()
	client.Connect()

	t.Cleanup(func() {
		server.Stop()

		client.Stop()
		client.Disconnect()
		waitForServerStopped(server)
		waitForClientStopped(client)
	})

	waitForConnection(client)
	sendMessage(client, test_message)
	received := receiveMessage(server)

	verifyMessage(t, received, test_message)
}

func TestClientConnectDisconnect(t *testing.T) {
	client := NewClient(&test_config)
	server := NewServer(&test_config)

	server.Start()
	client.Connect()

	t.Cleanup(func() {
		server.Stop()

		client.Stop()
		client.Disconnect()
		waitForClientStopped(client)
		waitForServerStopped(server)
	})

	waitForConnection(client)

	sendMessage(client, test_message)
	sendMessage(client, test_message_2)

	received := receiveMessage(server)
	verifyMessage(t, received, test_message)

	received = receiveMessage(server)
	verifyMessage(t, received, test_message_2)

	client.Stop()
	client.Disconnect()
	waitForClientStopped(client)

	client.Connect()
	waitForConnection(client)

	sendMessage(client, test_message)
	received = receiveMessage(server)
	verifyMessage(t, received, test_message)

	sendMessage(client, test_message_2)
	received = receiveMessage(server)
	verifyMessage(t, received, test_message_2)
}
