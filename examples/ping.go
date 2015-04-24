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
