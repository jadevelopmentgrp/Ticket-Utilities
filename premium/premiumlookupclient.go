package premium

import (
	"context"

	"github.com/go-redis/redis/v8"
	database "github.com/jadevelopmentgrp/Tickets-Database"
	"github.com/jadevelopmentgrp/Tickets-Utilities/model"
	"github.com/rxdn/gdl/cache"
	"github.com/rxdn/gdl/objects/guild"
	"github.com/rxdn/gdl/rest/ratelimit"
)

type IPremiumLookupClient interface {
	GetCachedTier(ctx context.Context, guildId uint64) (CachedTier, error)
	SetCachedTier(ctx context.Context, guildId uint64, tier CachedTier) error
	DeleteCachedTier(ctx context.Context, guildId uint64) error

	GetTierByGuild(ctx context.Context, guild guild.Guild) (PremiumTier, model.EntitlementSource, error)
	GetTierByGuildId(ctx context.Context, guildId uint64, includeVoting bool, botToken string, rateLimiter *ratelimit.Ratelimiter) (PremiumTier, error)
	GetTierByGuildIdWithSource(ctx context.Context, guildId uint64, botToken string, rateLimiter *ratelimit.Ratelimiter) (PremiumTier, model.EntitlementSource, error)

	GetTierByUser(ctx context.Context, userId uint64, includeVoting bool) (PremiumTier, error)
	GetTierByUserWithSource(ctx context.Context, userId uint64) (PremiumTier, model.EntitlementSource, error)
}

type PremiumLookupClient struct {
	redis    *redis.Client
	cache    *cache.PgCache
	database *database.Database
}

var _ IPremiumLookupClient = (*PremiumLookupClient)(nil)

func NewPremiumLookupClient(redisClient *redis.Client, cache *cache.PgCache, database *database.Database) *PremiumLookupClient {
	return &PremiumLookupClient{
		redis:    redisClient,
		cache:    cache,
		database: database,
	}
}
