# redis_spawn
Golang package to help spawn redis server for testing


```go
// file: ping.go

package ping

import (
  "github.com/garyburd/redigo/redis"
)

func Ping(host string) (string, error) {
  conn, err := redis.Dial("tcp", host)

  if err != nil {
    return "", err
  }

  return redis.String(conn.Do("PING"))

}

// file: ping_test.go

package ping

import (
  "github.com/nickdufresne/redis_spawn"
  "testing"
)

func TestPing(t *testing.T) {
  
  // spawn the server and check the INFO redis command

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


```
