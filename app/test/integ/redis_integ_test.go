package integ

import (
	"testing"
)

func TestCanSetValue(t *testing.T) {
	if _, err := rdb.Set(ctx, key, value, ttl).Result(); err != nil {
		t.Error(err)
	}
}

func TestCanGetValue(t *testing.T) {
	res, err := rdb.Get(ctx, key).Result()
	if err != nil {
		t.Error(err)
	}

	if res != value {
		t.Errorf("Wrong value. Expected %s but got %s", value, res)
	}
}

func TestCanTTLValue(t *testing.T) {
	res, err := rdb.TTL(ctx, key).Result()
	if err != nil {
		t.Error(err)
	}

	if res == 0 || res > ttl {
		t.Errorf("Wrong value. Expected less than %s and not 0 but got %s", ttl, res)
	}
}

func TestCanDeleteValue(t *testing.T) {
	if _, err := rdb.Del(ctx, key).Result(); err != nil {
		t.Error(err)
	}

	exists, err := rdb.Exists(ctx, key).Result()
	if err != nil {
		t.Error(err)
	}

	if exists != 0 {
		t.Errorf("Wrong value. Expected %v but got %v", 0, exists)
	}
}
