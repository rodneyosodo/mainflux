// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package nats_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/mainflux/mainflux/pkg/messaging"
	"github.com/mainflux/mainflux/pkg/messaging/nats"
	"github.com/stretchr/testify/assert"
)

const (
	topic       = "topic"
	chansPrefix = "channels"
	channel     = "9b7b1b3f-b1b0-46a8-a717-b8213f9eda3b"
	subtopic    = "engine"
	clientID    = "9b7b1b3f-b1b0-46a8-a717-b8213f9eda3b"
)

var (
	errFailed = fmt.Errorf("failed")
	msgChan   = make(chan *messaging.Message)
	message   = &messaging.Message{
		Channel:   channel,
		Subtopic:  subtopic,
		Publisher: "9b7b1b3f-b1b0-46a8-a717-b8213f9eda3b",
		Protocol:  "mqtt",
		Payload:   []byte("payload"),
		Created:   time.Now().UnixNano(),
	}
)

func TestPublisher(t *testing.T) {
	err := pubsub.Subscribe(context.TODO(), clientID, fmt.Sprintf("%s.>", chansPrefix), handler{})
	assert.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))
	cases := []struct {
		desc     string
		topic    string
		subtopic string
		message  *messaging.Message
		error    error
	}{
		{
			desc:     "publish message with empty message",
			topic:    channel,
			subtopic: subtopic,
			message:  &messaging.Message{},
			error:    nil,
		},
		{
			desc:     "publish message with message",
			topic:    channel,
			subtopic: subtopic,
			message:  message,
			error:    nil,
		},
		{
			desc:     "publish message with topic and empty subtopic",
			topic:    channel,
			subtopic: "",
			message:  message,
			error:    nil,
		},
		{
			desc:     "publish message with subtopic and empty topic",
			topic:    "",
			subtopic: subtopic,
			message:  message,
			error:    nats.ErrEmptyTopic,
		},
		{
			desc:     "publish message with topic and subtopic",
			topic:    channel,
			subtopic: subtopic,
			message:  message,
			error:    nil,
		},
	}

	for _, tc := range cases {
		tc.message.Subtopic = tc.subtopic
		err := pubsub.Publish(context.TODO(), tc.topic, tc.message)
		assert.Equal(t, tc.error, err, fmt.Sprintf("%s: expected %+v got %+v\n", tc.desc, tc.error, err))

		if err == nil {
			receivedMsg := <-msgChan
			assert.Equal(t, tc.message.Payload, receivedMsg.Payload, fmt.Sprintf("%s: expected %+v got %+v\n", tc.desc, tc.message.Payload, receivedMsg))
			assert.Equal(t, tc.message.Channel, receivedMsg.Channel, fmt.Sprintf("%s: expected %+v got %+v\n", tc.desc, &tc.message, receivedMsg))
			assert.Equal(t, tc.message.Created, receivedMsg.Created, fmt.Sprintf("%s: expected %+v got %+v\n", tc.desc, &tc.message, receivedMsg))
			assert.Equal(t, tc.message.Protocol, receivedMsg.Protocol, fmt.Sprintf("%s: expected %+v got %+v\n", tc.desc, &tc.message, receivedMsg))
			assert.Equal(t, tc.message.Publisher, receivedMsg.Publisher, fmt.Sprintf("%s: expected %+v got %+v\n", tc.desc, &tc.message, receivedMsg))
			assert.Equal(t, tc.message.Subtopic, receivedMsg.Subtopic, fmt.Sprintf("%s: expected %+v got %+v\n", tc.desc, &tc.message, receivedMsg))
			assert.Equal(t, tc.message.Payload, receivedMsg.Payload, fmt.Sprintf("%s: expected %+v got %+v\n", tc.desc, &tc.message, receivedMsg))
		}
	}
}

func TestPubsub(t *testing.T) {
	// Test Subscribe and Unsubscribe.
	subcases := []struct {
		desc         string
		topic        string
		clientID     string
		errorMessage error
		pubsub       bool // true for subscribe and false for unsubscribe.
		handler      messaging.MessageHandler
	}{
		{
			desc:         "Subscribe to a topic with an ID",
			topic:        fmt.Sprintf("%s.%s", chansPrefix, topic),
			clientID:     "clientid1",
			errorMessage: nil,
			pubsub:       true,
			handler:      handler{false},
		},
		{
			desc:         "Subscribe to the same topic with a different ID",
			topic:        fmt.Sprintf("%s.%s", chansPrefix, topic),
			clientID:     "clientid2",
			errorMessage: nil,
			pubsub:       true,
			handler:      handler{false},
		},
		{
			desc:         "Subscribe to an already subscribed topic with an ID",
			topic:        fmt.Sprintf("%s.%s", chansPrefix, topic),
			clientID:     "clientid1",
			errorMessage: nil,
			pubsub:       true,
			handler:      handler{false},
		},
		{
			desc:         "Unsubscribe from a topic with an ID",
			topic:        fmt.Sprintf("%s.%s", chansPrefix, topic),
			clientID:     "clientid1",
			errorMessage: nil,
			pubsub:       false,
			handler:      handler{false},
		},
		{
			desc:         "Unsubscribe from a non-existent topic with an ID",
			topic:        "h",
			clientID:     "clientid1",
			errorMessage: nats.ErrNotSubscribed,
			pubsub:       false,
			handler:      handler{false},
		},
		{
			desc:         "Unsubscribe from the same topic with a different ID",
			topic:        fmt.Sprintf("%s.%s", chansPrefix, topic),
			clientID:     "clientidd2",
			errorMessage: nats.ErrNotSubscribed,
			pubsub:       false,
			handler:      handler{false},
		},
		{
			desc:         "Unsubscribe from the same topic with a different ID not subscribed",
			topic:        fmt.Sprintf("%s.%s", chansPrefix, topic),
			clientID:     "clientidd3",
			errorMessage: nats.ErrNotSubscribed,
			pubsub:       false,
			handler:      handler{false},
		},
		{
			desc:         "Unsubscribe from an already unsubscribed topic with an ID",
			topic:        fmt.Sprintf("%s.%s", chansPrefix, topic),
			clientID:     "clientid1",
			errorMessage: nats.ErrNotSubscribed,
			pubsub:       false,
			handler:      handler{false},
		},
		{
			desc:         "Subscribe to a topic with a subtopic with an ID",
			topic:        fmt.Sprintf("%s.%s.%s", chansPrefix, topic, subtopic),
			clientID:     "clientidd1",
			errorMessage: nil,
			pubsub:       true,
			handler:      handler{false},
		},
		{
			desc:         "Subscribe to an already subscribed topic with a subtopic with an ID",
			topic:        fmt.Sprintf("%s.%s.%s", chansPrefix, topic, subtopic),
			clientID:     "clientidd1",
			errorMessage: nil,
			pubsub:       true,
			handler:      handler{false},
		},
		{
			desc:         "Unsubscribe from a topic with a subtopic with an ID",
			topic:        fmt.Sprintf("%s.%s.%s", chansPrefix, topic, subtopic),
			clientID:     "clientidd1",
			errorMessage: nil,
			pubsub:       false,
			handler:      handler{false},
		},
		{
			desc:         "Unsubscribe from an already unsubscribed topic with a subtopic with an ID",
			topic:        fmt.Sprintf("%s.%s.%s", chansPrefix, topic, subtopic),
			clientID:     "clientid1",
			errorMessage: nats.ErrNotSubscribed,
			pubsub:       false,
			handler:      handler{false},
		},
		{
			desc:         "Subscribe to an empty topic with an ID",
			topic:        "",
			clientID:     "clientid1",
			errorMessage: nats.ErrEmptyTopic,
			pubsub:       true,
			handler:      handler{false},
		},
		{
			desc:         "Unsubscribe from an empty topic with an ID",
			topic:        "",
			clientID:     "clientid1",
			errorMessage: nats.ErrEmptyTopic,
			pubsub:       false,
			handler:      handler{false},
		},
		{
			desc:         "Subscribe to a topic with empty id",
			topic:        fmt.Sprintf("%s.%s", chansPrefix, topic),
			clientID:     "",
			errorMessage: nats.ErrEmptyID,
			pubsub:       true,
			handler:      handler{false},
		},
		{
			desc:         "Unsubscribe from a topic with empty id",
			topic:        fmt.Sprintf("%s.%s", chansPrefix, topic),
			clientID:     "",
			errorMessage: nats.ErrEmptyID,
			pubsub:       false,
			handler:      handler{false},
		},
		{
			desc:         "Subscribe to another topic with an ID",
			topic:        fmt.Sprintf("%s.%s", chansPrefix, topic+"1"),
			clientID:     "clientid3",
			errorMessage: nil,
			pubsub:       true,
			handler:      handler{true},
		},
		{
			desc:         "Subscribe to another already subscribed topic with an ID with Unsubscribe failing",
			topic:        fmt.Sprintf("%s.%s", chansPrefix, topic+"1"),
			clientID:     "clientid3",
			errorMessage: errFailed,
			pubsub:       true,
			handler:      handler{true},
		},
		{
			desc:         "Subscribe to a new topic with an ID",
			topic:        fmt.Sprintf("%s.%s", chansPrefix, topic+"2"),
			clientID:     "clientid4",
			errorMessage: nil,
			pubsub:       true,
			handler:      handler{true},
		},
		{
			desc:         "Unsubscribe from a topic with an ID with failing handler",
			topic:        fmt.Sprintf("%s.%s", chansPrefix, topic+"2"),
			clientID:     "clientid4",
			errorMessage: errFailed,
			pubsub:       false,
			handler:      handler{true},
		},
	}

	for _, pc := range subcases {
		if pc.pubsub == true {
			err := pubsub.Subscribe(context.TODO(), pc.clientID, pc.topic, pc.handler)
			if pc.errorMessage == nil {
				assert.Nil(t, err, fmt.Sprintf("%s got unexpected error: %s", pc.desc, err))
			} else {
				assert.Equal(t, err, pc.errorMessage)
			}
		} else {
			err := pubsub.Unsubscribe(context.TODO(), pc.clientID, pc.topic)
			if pc.errorMessage == nil {
				assert.Nil(t, err, fmt.Sprintf("%s got unexpected error: %s", pc.desc, err))
			} else {
				assert.Equal(t, err, pc.errorMessage)
			}
		}
	}
}

type handler struct {
	fail bool
}

func (h handler) Handle(msg *messaging.Message) error {
	msgChan <- msg

	return nil
}

func (h handler) Cancel() error {
	if h.fail {
		return errFailed
	}

	return nil
}
