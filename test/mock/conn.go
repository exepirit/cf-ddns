package mock

import (
	"github.com/stretchr/testify/mock"
	"net"
	"time"
)

type Conn struct {
	mock.Mock
}

func (c *Conn) Read(b []byte) (n int, err error) {
	args := c.Called(b)
	return args.Int(0), args.Error(1)
}

func (c *Conn) Write(b []byte) (n int, err error) {
	args := c.Called(b)
	return args.Int(0), args.Error(1)
}

func (c *Conn) Close() error {
	return c.Called().Error(0)
}

func (c *Conn) LocalAddr() net.Addr {
	return c.Called().Get(0).(net.Addr)
}

func (c *Conn) RemoteAddr() net.Addr {
	return c.Called().Get(0).(net.Addr)
}

func (c *Conn) SetDeadline(t time.Time) error {
	return c.Called(t).Error(0)
}

func (c *Conn) SetReadDeadline(t time.Time) error {
	return c.Called(t).Error(0)
}

func (c *Conn) SetWriteDeadline(t time.Time) error {
	return c.Called(t).Error(0)
}
