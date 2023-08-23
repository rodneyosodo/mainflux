//go:build rabbitmq
// +build rabbitmq

// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package brokers

import (
	"context"
	"log"

	mflog "github.com/mainflux/mainflux/logger"
	"github.com/mainflux/mainflux/pkg/messaging"
	"github.com/mainflux/mainflux/pkg/messaging/rabbitmq"
)

// SubjectAllChannels represents subject to subscribe for all the channels.
const SubjectAllChannels = "channels.#"

func init() {
	log.Println("The binary was build using RabbitMQ as the message broker")
}

func NewPublisher(_ context.Context, url string) (messaging.Publisher, error) {
	pb, err := rabbitmq.NewPublisher(url)
	if err != nil {
		return nil, err
	}

	return pb, nil
}

func NewPubSub(_ context.Context, url string, logger mflog.Logger) (messaging.PubSub, error) {
	pb, err := rabbitmq.NewPubSub(url, logger)
	if err != nil {
		return nil, err
	}

	return pb, nil
}
