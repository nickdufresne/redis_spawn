package redis_spawn

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type Server struct {
	cmd   *exec.Cmd
	PID   int
	Port  int
	Ready <-chan bool
}

func (s *Server) Shutdown() error {
	return s.cmd.Process.Kill()
}

var AlreadyRunning = errors.New("Server already running")

func (s *Server) connect(address string) (net.Conn, error) {
	var conn net.Conn
	var err error
	for i := 1; i <= 5; i++ {
		conn, err = net.DialTimeout("tcp", address, time.Second)
		if err == nil {
			return conn, err
		}

		<-time.After(100 * time.Millisecond)
	}

	return conn, err
}

func (s *Server) checkRunning() (bool, error) {

	address := fmt.Sprintf("localhost:%d", s.Port)
	conn, err := s.connect(address)

	if err != nil {
		return false, err
	}

	defer conn.Close()

	fmt.Fprintf(conn, "INFO\n")

	reader := bufio.NewReader(conn)

	response, err := reader.ReadString('\n')

	if err != nil {
		return false, err
	}

	n, err := strconv.Atoi(response[1 : len(response)-2])

	if err != nil {
		return false, err
	}

	p := make([]byte, n)
	_, err = io.ReadFull(reader, p)

	if err != nil {
		return false, err
	}

	//fmt.Println(string(p))

	var pidCheck int
	buf := bytes.NewReader(p)
	scanner := bufio.NewScanner(buf)
	for scanner.Scan() {
		line := scanner.Text()
		strs := strings.Split(line, ":")
		if len(strs) == 2 && strs[0] == "process_id" {
			pidCheck, err = strconv.Atoi(strs[1])

			if err != nil {
				return false, err
			}

			break
		}
	}

	if err := scanner.Err(); err != nil {
		if err != nil {
			return false, err
		}
	}

	if pidCheck == s.PID {
		return true, nil
	} else {
		return false, AlreadyRunning
	}
}

func NewServer(port int) *Server {

	cmd := exec.Command("redis-server", "--port", fmt.Sprintf("%d", port), "--save", "\"\"")

	s := new(Server)
	s.Port = port
	s.cmd = cmd

	return s

}

func (s *Server) Start() error {
	err := s.cmd.Start()
	s.PID = s.cmd.Process.Pid
	if err != nil {
		return err
	}

	_, err = s.checkRunning()

	return err
}

func Spawn(port int) (*Server, error) {
	s := NewServer(port)
	err := s.Start()
	return s, err
}
