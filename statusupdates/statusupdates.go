package statusupdates

import (
	"context"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/jadevelopmentgrp/Tickets-Utilities/utils"
)

const channel = "tickets:statusupdates"

func Publish(redis *redis.Client, botId uint64) {
	redis.Publish(utils.DefaultContext(), channel, botId)
}

func Listen(redis *redis.Client, ch chan uint64) {
	for payload := range redis.Subscribe(context.Background(), channel).Channel() {
		if id, err := strconv.ParseUint(payload.Payload, 10, 64); err == nil {
			ch <- id
		}
	}
}
