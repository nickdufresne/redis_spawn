package redis_spawn_test

import (
	"github.com/garyburd/redigo/redis"
	redspawn "github.com/nickdufresne/redis_spawn"
	"net"
	"testing"
	"time"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type MySuite struct{}

var _ = Suite(&MySuite{})

func (s *MySuite) TestSpawnLaunches(c *C) {
	server, err := redspawn.Spawn(9999)

	c.Check(server.PID, Not(Equals), 0)
	c.Check(err, Equals, nil)

	conn, err := redis.Dial("tcp", "localhost:9999")
	c.Check(err, Equals, nil)

	pong, err := redis.String(conn.Do("PING"))
	c.Check(err, Equals, nil)

	c.Check(pong, Equals, "PONG")

	err = server.Shutdown()

	c.Check(err, Equals, nil)

	<-time.After(100 * time.Millisecond)
	_, err = net.Dial("tcp", "localhost:9999")
	c.Check(err, Not(Equals), nil)

}

func (s *MySuite) TestSpawnDetectsAlreadyRunning(c *C) {
	server, err := redspawn.Spawn(9999)

	c.Check(server.PID, Not(Equals), 0)
	c.Check(err, Equals, nil)

	conn, err := redis.Dial("tcp", "localhost:9999")
	c.Check(err, Equals, nil)

	pong, err := redis.String(conn.Do("PING"))
	c.Check(err, Equals, nil)

	c.Check(pong, Equals, "PONG")

	server2, err := redspawn.Spawn(9999)
	c.Check(err, Equals, redspawn.AlreadyRunning)

	server.Shutdown()
	server2.Shutdown()

	<-time.After(100 * time.Millisecond)
}
