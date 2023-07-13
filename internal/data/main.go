package data

import (
	"github.com/adjust/rmq/v5"
	"github.com/go-redis/redis/v8"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type Rediser interface {
	Client() *redis.Client
	CleanQueues() (int64, error)
	OpenRelayQueue() rmq.Queue
}

type rediser struct {
	client     *redis.Client
	connection rmq.Connection
	cleaner    *rmq.Cleaner

	relayQueueOnce comfig.Once
}

func (r *rediser) Client() *redis.Client {
	return r.client
}

func (r *rediser) CleanQueues() (int64, error) {
	return r.cleaner.Clean()
}

func (r *rediser) OpenRelayQueue() rmq.Queue {
	return r.relayQueueOnce.Do(func() interface{} {
		taskQueue, err := r.connection.OpenQueue("relay")
		if err != nil {
			panic(errors.Wrap(err, "failed to open a task queue"))
		}

		return taskQueue
	}).(rmq.Queue)
}

func NewRediser(cfg config, log *logan.Entry) Rediser {
	client := redis.NewClient(&redis.Options{Addr: cfg.Addr, Password: cfg.Password})
	errChan := make(chan error)
	connection, err := rmq.OpenConnectionWithRedisClient("redis-queue", client, errChan)
	if err != nil {
		panic(errors.Wrap(err, "failed to set up redis queue"))
	}
	go func() {
		for {
			err, more := <-errChan
			if !more {
				return
			}
			log.WithError(err).Error("a background redis queue error happened")
		}
	}()

	return &rediser{
		client:     client,
		connection: connection,
		cleaner:    rmq.NewCleaner(connection),
	}
}
