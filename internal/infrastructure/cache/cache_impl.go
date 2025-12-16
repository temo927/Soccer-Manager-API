package cache

import (
	"context"
	"encoding/json"
	"fmt"

	cachePort "soccer-manager-api/internal/ports/cache"
)


func CacheKey(prefix string, id string) string {
	return fmt.Sprintf("%s:%s", prefix, id)
}


type CacheHelper struct {
	cache cachePort.Cache
}


func NewCacheHelper(c cachePort.Cache) *CacheHelper {
	return &CacheHelper{cache: c}
}


func (h *CacheHelper) Get(ctx context.Context, key string, dest interface{}) error {
	data, err := h.cache.Get(ctx, key)
	if err != nil {
		return err
	}
	if data == nil {
		return fmt.Errorf("cache miss")
	}
	return json.Unmarshal(data, dest)
}


func (h *CacheHelper) Set(ctx context.Context, key string, value interface{}, ttl int) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return h.cache.Set(ctx, key, data, ttl)
}


func (h *CacheHelper) InvalidateTeamCache(ctx context.Context, teamID string) error {
	patterns := []string{
		CacheKey("team", teamID),
		CacheKey("team:players", teamID),
		CacheKey("team:value", teamID),
	}
	for _, pattern := range patterns {
		if err := h.cache.Delete(ctx, pattern); err != nil {
			return err
		}
	}
	return nil
}


func (h *CacheHelper) InvalidateTransferListCache(ctx context.Context) error {
	return h.cache.DeleteByPattern(ctx, "transfer_list:*")
}

