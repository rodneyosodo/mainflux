// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package producer

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/mainflux/mainflux/bootstrap"
	"github.com/mainflux/mainflux/pkg/errors"
)

const (
	streamID                       = "mainflux.bootstrap"
	streamLen                      = 1000
	checkUnpublishedEventsInterval = 1 * time.Minute
)

var _ bootstrap.Service = (*eventStore)(nil)

type eventStore struct {
	svc               bootstrap.Service
	client            *redis.Client
	unpublishedEvents []*redis.XAddArgs
}

// NewEventStoreMiddleware returns wrapper around bootstrap service that sends
// events to event store.
func NewEventStoreMiddleware(ctx context.Context, svc bootstrap.Service, client *redis.Client) bootstrap.Service {
	es := &eventStore{
		svc:    svc,
		client: client,
	}

	go es.startPublishingRoutine(ctx)

	return es
}

func (es *eventStore) Add(ctx context.Context, token string, cfg bootstrap.Config) (bootstrap.Config, error) {
	saved, err := es.svc.Add(ctx, token, cfg)
	if err != nil {
		return saved, err
	}

	ev := configEvent{
		saved, configCreate,
	}

	if err1 := es.publish(ctx, ev); err1 != nil {
		return saved, errors.Wrap(err, err1)
	}

	return saved, err
}

func (es *eventStore) View(ctx context.Context, token, id string) (bootstrap.Config, error) {
	cfg, err := es.svc.View(ctx, token, id)
	if err != nil {
		return cfg, err
	}
	ev := configEvent{
		cfg, configList,
	}

	if err1 := es.publish(ctx, ev); err1 != nil {
		return cfg, errors.Wrap(err, err1)
	}

	return cfg, err
}

func (es *eventStore) Update(ctx context.Context, token string, cfg bootstrap.Config) error {
	if err := es.svc.Update(ctx, token, cfg); err != nil {
		return err
	}

	ev := configEvent{
		cfg, configUpdate,
	}

	return es.publish(ctx, ev)
}

func (es eventStore) UpdateCert(ctx context.Context, token, thingKey, clientCert, clientKey, caCert string) (bootstrap.Config, error) {
	cfg, err := es.svc.UpdateCert(ctx, token, thingKey, clientCert, clientKey, caCert)
	if err != nil {
		return cfg, err
	}

	ev := updateCertEvent{
		thingKey:   thingKey,
		clientCert: clientCert,
		clientKey:  clientKey,
		caCert:     caCert,
	}

	return es.publish(ctx, ev)
}

func (es *eventStore) UpdateConnections(ctx context.Context, token, id string, connections []string) error {
	if err := es.svc.UpdateConnections(ctx, token, id, connections); err != nil {
		return err
	}

	ev := updateConnectionsEvent{
		mfThing:    id,
		mfChannels: connections,
	}

	return es.publish(ctx, ev)
}

func (es *eventStore) List(ctx context.Context, token string, filter bootstrap.Filter, offset, limit uint64) (bootstrap.ConfigsPage, error) {
	bp, err := es.svc.List(ctx, token, filter, offset, limit)
	if err != nil {
		return bp, err
	}

	ev := listConfigsEvent{
		offset:       offset,
		limit:        limit,
		fullMatch:    filter.FullMatch,
		partialMatch: filter.PartialMatch,
	}

	if err1 := es.publish(ctx, ev); err1 != nil {
		return bp, errors.Wrap(err, err1)
	}

	return bp, nil
}

func (es *eventStore) Remove(ctx context.Context, token, id string) error {
	if err := es.svc.Remove(ctx, token, id); err != nil {
		return err
	}

	ev := removeConfigEvent{
		mfThing: id,
	}

	return es.publish(ctx, ev)
}

func (es *eventStore) Bootstrap(ctx context.Context, externalKey, externalID string, secure bool) (bootstrap.Config, error) {
	cfg, err := es.svc.Bootstrap(ctx, externalKey, externalID, secure)

	ev := bootstrapEvent{
		cfg,
		externalID,
		true,
	}

	if err != nil {
		ev.success = false
	}

	if err1 := es.publish(ctx, ev); err1 != nil {
		return cfg, err1
	}

	return cfg, err
}

func (es *eventStore) ChangeState(ctx context.Context, token, id string, state bootstrap.State) error {
	if err := es.svc.ChangeState(ctx, token, id, state); err != nil {
		return err
	}

	ev := changeStateEvent{
		mfThing: id,
		state:   state,
	}

	return es.publish(ctx, ev)
}

func (es *eventStore) RemoveConfigHandler(ctx context.Context, id string) error {
	if err := es.svc.RemoveConfigHandler(ctx, id); err != nil {
		return err
	}

	ev := removeHandlerEvent{
		id:        id,
		operation: configHandlerRemove,
	}

	return es.publish(ctx, ev)
}

func (es *eventStore) RemoveChannelHandler(ctx context.Context, id string) error {
	if err := es.svc.RemoveChannelHandler(ctx, id); err != nil {
		return err
	}

	ev := removeHandlerEvent{
		id:        id,
		operation: channelHandlerRemove,
	}

	return es.publish(ctx, ev)
}

func (es *eventStore) UpdateChannelHandler(ctx context.Context, channel bootstrap.Channel) error {
	if err := es.svc.UpdateChannelHandler(ctx, channel); err != nil {
		return err
	}

	ev := updateChannelHandlerEvent{
		channel,
	}

	return es.publish(ctx, ev)
}

func (es *eventStore) DisconnectThingHandler(ctx context.Context, channelID, thingID string) error {
	if err := es.svc.DisconnectThingHandler(ctx, channelID, thingID); err != nil {
		return err
	}

	ev := disconnectThingEvent{
		channelID,
		thingID,
	}

	return es.publish(ctx, ev)
}

func (es *eventStore) checkRedisConnection(ctx context.Context) error {
	// A timeout is used to avoid blocking the main thread
	ctx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	return es.client.Ping(ctx).Err()
}

func (es *eventStore) publish(ctx context.Context, ev event) error {
	values, err := ev.encode()
	if err != nil {
		return err
	}

	record := &redis.XAddArgs{
		Stream:       streamID,
		MaxLenApprox: streamLen,
		Values:       values,
	}

	if err := es.checkRedisConnection(ctx); err != nil {
		es.unpublishedEvents = append(es.unpublishedEvents, record)
		return nil
	}

	return es.client.XAdd(ctx, record).Err()
}

func (es *eventStore) startPublishingRoutine(ctx context.Context) {
	ticker := time.NewTicker(checkUnpublishedEventsInterval)
	for {
		select {
		case <-ticker.C:
			if err := es.checkRedisConnection(ctx); err == nil {
				for i := len(es.unpublishedEvents) - 1; i >= 0; i-- {
					if err := es.client.XAdd(ctx, es.unpublishedEvents[i]).Err(); err == nil {
						es.unpublishedEvents = append(es.unpublishedEvents[:i], es.unpublishedEvents[i+1:]...)
					}
				}
			}
		case <-ctx.Done():
			return
		}
	}
}
