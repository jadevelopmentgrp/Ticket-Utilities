package whitelabeldelete

import (
	"context"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/jadevelopmentgrp/Tickets-Utilities/utils"
)

const channel = "tickets:whitelabeldelete"

func Publish(redis *redis.Client, botId uint64) {
	redis.Publish(utils.DefaultContext(), channel, botId)
}

// bot id
func Listen(redis *redis.Client, ch chan uint64) {
	for payload := range redis.Subscribe(context.Background(), channel).Channel() {
		if id, err := strconv.ParseUint(payload.Payload, 10, 64); err == nil {
			ch <- id
		}
	}
}
