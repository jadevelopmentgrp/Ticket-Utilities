package bloxlink

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/jadevelopmentgrp/Tickets-Utilities/webproxy"
)

type (
	BloxlinkIntegration struct {
		redis  *redis.Client
		proxy  *webproxy.WebProxy
		apiKey string
	}

	cachedUser struct {
		User *User `json:"user"` // If the user does not exist, this will be nil, making a separate bool redundant
	}
)

func NewBloxlinkIntegration(redis *redis.Client, proxy *webproxy.WebProxy, apiKey string) *BloxlinkIntegration {
	return &BloxlinkIntegration{
		redis:  redis,
		proxy:  proxy,
		apiKey: apiKey,
	}
}

func newCachedUser(user User) cachedUser {
	return cachedUser{
		User: &user,
	}
}

func newNullUser() cachedUser {
	return cachedUser{
		User: nil,
	}
}

const cacheLength = time.Hour * 24

func (i *BloxlinkIntegration) GetRobloxUser(ctx context.Context, discordUserId uint64) (User, error) {
	redisKey := fmt.Sprintf("bloxlink:%d", discordUserId)

	// See if we have a cached value
	cached, err := i.redis.Get(ctx, redisKey).Result()
	if err == nil {
		var user cachedUser
		if err := json.Unmarshal([]byte(cached), &user); err != nil {
			return User{}, err
		}

		if user.User == nil {
			return User{}, ErrUserNotFound
		} else {
			return *user.User, nil
		}
	} else if err != redis.Nil { // If the error is redis.Nil, this means that the key does not exist, and we should continue
		return User{}, err
	}

	// Fetch user ID from Bloxlink
	robloxId, err := RequestUserId(ctx, i.proxy, i.apiKey, discordUserId)
	if err != nil {
		if err == ErrUserNotFound { // If user not found, we should still cache this
			encoded, err := json.Marshal(newNullUser())
			if err != nil {
				return User{}, err
			}

			i.redis.SetEX(context.Background(), redisKey, encoded, cacheLength)
		}

		return User{}, err
	}

	// Fetch user object
	user, err := RequestUserData(ctx, i.proxy, robloxId)
	if err != nil {
		return User{}, err
	}

	// Cache response
	encoded, err := json.Marshal(newCachedUser(user))
	if err != nil {
		return User{}, err
	}

	i.redis.SetEX(ctx, redisKey, string(encoded), cacheLength)

	return user, nil
}
