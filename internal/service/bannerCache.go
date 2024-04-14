package service

import (
	"avito/internal/model"
	"fmt"
	"strconv"
	"time"

	"github.com/dgraph-io/ristretto"
)

type BannerCache struct {
	cache *ristretto.Cache
	ttl   time.Duration
}

type ConfigCache struct {
	TTL string `yaml:"ttl" env:"TTL"`
}

func NewBannerCache(cfg *ConfigCache) (*BannerCache, error) {
	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create new cache [banner cache service ~ NewBannerCache]: %w", err)
	}

	ttl, err := time.ParseDuration(cfg.TTL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse duration [banner cache service ~ NewBannerCache]: %w", err)
	}

	return &BannerCache{
		cache: cache,
		ttl:   ttl,
	}, nil
}

func (c *BannerCache) Set(tagID, featureID int, banner *model.Banner) {
	key := tagAndFeatureKey(tagID, featureID)
	c.cache.SetWithTTL(key, banner, 0, c.ttl)
	c.cache.Wait()
}

func (c *BannerCache) Get(tagID, featureID int) (banner *model.Banner, ok bool) {
	key := tagAndFeatureKey(tagID, featureID)
	if x, found := c.cache.Get(key); found {
		b, ok := x.(*model.Banner)
		return b, ok
	}
	return nil, false
}

func tagAndFeatureKey(tagID, featureID int) string {
	return strconv.Itoa(tagID) + strconv.Itoa(featureID)
}
