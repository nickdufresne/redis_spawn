package ping

import (
	"github.com/nickdufresne/redis_spawn"
	"testing"
)

func TestPing(t *testing.T) {
	server, err := redis_spawn.Spawn(9999)

	if err != nil {
		t.Errorf("Could not start redis: %s", err.Error())
	}

	pong, err := Ping("localhost:9999")

	if err != nil {
		t.Errorf("Could not connect to redis: %s", err.Error())
	}

	if pong != "PONG" {
		t.Errorf("Unexpected result from redis: %s != 'PONG'", pong)
	}

	server.Shutdown()
}
