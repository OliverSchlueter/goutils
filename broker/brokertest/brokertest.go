package brokertest

import (
	"github.com/OliverSchlueter/goutils/broker"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
	"time"
)

func TestPublish(t *testing.T, b broker.Broker) {
	subject := "test.publish"
	testData := []byte("test publish data")
	receivedCh := make(chan []byte, 1)

	// Set up a subscriber to verify the publish worked
	err := b.Subscribe(subject, func(msg *nats.Msg) {
		receivedCh <- msg.Data
	})
	require.NoError(t, err, "Failed to subscribe")

	// Publish the message
	err = b.Publish(subject, testData)
	require.NoError(t, err, "Failed to publish message")

	// Wait for the message with timeout
	select {
	case received := <-receivedCh:
		assert.Equal(t, testData, received, "Received data doesn't match sent data")
	case <-time.After(500 * time.Millisecond):
		assert.Fail(t, "Timed out waiting for published message")
	}
}

func TestRequest(t *testing.T, b broker.Broker) {
	subject := "test.request"
	requestData := []byte("test request data")
	responseData := []byte("test response data")

	// For NatsBroker, set up a responder
	if _, ok := b.(*broker.NatsBroker); ok {
		err := b.Subscribe(subject, func(msg *nats.Msg) {
			b.Publish(msg.Reply, responseData)
		})
		require.NoError(t, err, "Failed to set up responder")
	}

	// For FakeBroker, we need to set the request handler
	if fakeBroker, ok := b.(*broker.FakeBroker); ok {
		fakeBroker.SetRequestHandler(func(msg *nats.Msg) (*nats.Msg, error) {
			return &nats.Msg{
				Subject: "response",
				Data:    responseData,
			}, nil
		})
	}

	// Send the request
	response, err := b.Request(subject, requestData)
	require.NoError(t, err, "Request failed")
	require.NotNil(t, response, "Response should not be nil")
	assert.Equal(t, responseData, response.Data, "Response data doesn't match expected")
}

func TestSubscribe(t *testing.T, b broker.Broker) {
	subject := "test.subscribe"
	testData := []byte("test subscribe data")
	receivedCount := 0
	expectedCount := 3

	// Create a channel to synchronize test completion
	doneCh := make(chan struct{})

	// Subscribe to the subject
	err := b.Subscribe(subject, func(msg *nats.Msg) {
		assert.Equal(t, subject, msg.Subject, "Received message has wrong subject")
		assert.Equal(t, testData, msg.Data, "Received message has wrong data")
		receivedCount++
		if receivedCount == expectedCount {
			close(doneCh)
		}
	})
	require.NoError(t, err, "Failed to subscribe")

	// Publish multiple messages
	for i := 0; i < expectedCount; i++ {
		err = b.Publish(subject, testData)
		require.NoError(t, err, "Failed to publish message")
	}

	// Wait for all messages or timeout
	select {
	case <-doneCh:
		assert.Equal(t, expectedCount, receivedCount, "Incorrect number of messages received")
	case <-time.After(1 * time.Second):
		assert.Fail(t, "Timed out waiting for all messages")
	}
}

func TestSubscribeQueue(t *testing.T, b broker.Broker) {
	subject := "test.queue"
	queue := "test-group"
	testData := []byte("test queue data")

	// For testing queue subscription, we'll create multiple subscribers
	// and ensure messages are delivered appropriately

	receiverCount := 3
	messageCount := 5

	// Track received messages per subscriber
	received := make([]int, receiverCount)
	mutex := &sync.Mutex{}

	// Track total messages received
	totalReceived := 0
	wg := sync.WaitGroup{}
	wg.Add(messageCount) // We expect to receive exactly messageCount messages in total

	// Create multiple queue subscribers
	for i := 0; i < receiverCount; i++ {
		subscriberID := i
		err := b.SubscribeQueue(subject, queue, func(msg *nats.Msg) {
			assert.Equal(t, testData, msg.Data, "Message data doesn't match")

			mutex.Lock()
			received[subscriberID]++
			totalReceived++
			// Only call Done() for the first messageCount messages received
			if totalReceived <= messageCount {
				wg.Done()
			}
			mutex.Unlock()
		})
		require.NoError(t, err, "Failed to create queue subscriber")
	}

	// Publish messages
	for i := 0; i < messageCount; i++ {
		err := b.Publish(subject, testData)
		require.NoError(t, err, "Failed to publish message")
	}

	// Wait for all messages to be processed with timeout
	waitCh := make(chan struct{})
	go func() {
		wg.Wait()
		close(waitCh)
	}()

	select {
	case <-waitCh:
		// Success - continue with validation
	case <-time.After(1 * time.Second):
		assert.Fail(t, "Timed out waiting for queue messages")
	}

	// Verify that at least messageCount messages were received
	mutex.Lock()
	finalTotal := 0
	for _, count := range received {
		finalTotal += count
	}
	mutex.Unlock()

	assert.GreaterOrEqual(t, finalTotal, messageCount,
		"At least messageCount messages should be received across all subscribers")

	if _, ok := b.(*broker.NatsBroker); ok {
		// For real NATS, each message should go to exactly one subscriber
		assert.Equal(t, messageCount, finalTotal,
			"Each message should be delivered to exactly one subscriber")
	}
}
