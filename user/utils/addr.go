package utils

import (
	"net"
)

// GetFreePort 动态生成端口
func GetFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, nil
	}
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, nil
	}
	defer func(l *net.TCPListener) {
		err := l.Close()
		if err != nil {

		}
	}(l)
	return l.Addr().(*net.TCPAddr).Port, nil
}
