package cache_test

import (
	"testing"
	"time"

	"github.com/Oleexo/config-go"
	"github.com/Oleexo/config-go/mem"

	"github.com/mistribe/subtracker/internal/adapters/cache"
	"github.com/mistribe/subtracker/internal/ports"
)

func TestSetAndGet(t *testing.T) {
	cfg := config.NewConfiguration()
	local := cache.NewLocal(cfg)

	t.Run("Set and Get without expiration", func(t *testing.T) {
		local.Set("key1", "value1")
		if got := local.Get("key1"); got != "value1" {
			t.Errorf("Get() = %v, want %v", got, "value1")
		}
	})

	t.Run("Set and Get with expiration", func(t *testing.T) {
		local.Set("key2", "value2", ports.WithDuration(100*time.Millisecond))
		if got := local.Get("key2"); got != "value2" {
			t.Errorf("Get() = %v, want %v", got, "value2")
		}
	})

	t.Run("Get expired item", func(t *testing.T) {
		local.Set("key3", "value3", ports.WithDuration(1*time.Millisecond))
		time.Sleep(2 * time.Millisecond)
		if got := local.Get("key3"); got != nil {
			t.Errorf("Get() = %v, want nil", got)
		}
	})

	t.Run("Get non-existent key", func(t *testing.T) {
		if got := local.Get("non-existent"); got != nil {
			t.Errorf("Get() = %v, want nil", got)
		}
	})
}

func TestSetAndGet_WithExpiration(t *testing.T) {
	memConfig := make(map[string]config.Entry)
	memConfig["CACHE_DEFAULT_TTL"] = config.NewEntryInt(int64(100 * time.Millisecond))
	cfg := config.NewConfiguration(mem.WithMemory(memConfig))
	local := cache.NewLocal(cfg)

	t.Run("Set with default TTL", func(t *testing.T) {
		local.Set("key1", "value1")
		if got := local.Get("key1"); got != "value1" {
			t.Errorf("Get() = %v, want %v", got, "value1")
		}
		time.Sleep(150 * time.Millisecond)
		if got := local.Get("key1"); got != nil {
			t.Errorf("Get() = %v, want nil", got)
		}
	})

	t.Run("Set with custom TTL overriding default", func(t *testing.T) {
		local.Set("key2", "value2", ports.WithDuration(200*time.Millisecond))
		if got := local.Get("key2"); got != "value2" {
			t.Errorf("Get() = %v, want %v", got, "value2")
		}
		time.Sleep(150 * time.Millisecond)
		if got := local.Get("key2"); got != "value2" {
			t.Errorf("Get() = %v, want %v", got, "value2")
		}
		time.Sleep(100 * time.Millisecond)
		if got := local.Get("key2"); got != nil {
			t.Errorf("Get() = %v, want nil", got)
		}
	})

	t.Run("Multiple items with different expirations", func(t *testing.T) {
		local.Set("key3", "value3", ports.WithDuration(300*time.Millisecond))
		local.Set("key4", "value4", ports.WithDuration(50*time.Millisecond))
		time.Sleep(100 * time.Millisecond)
		if got := local.Get("key3"); got != "value3" {
			t.Errorf("Get() = %v, want %v", got, "value3")
		}
		if got := local.Get("key4"); got != nil {
			t.Errorf("Get() = %v, want nil", got)
		}
	})
}
