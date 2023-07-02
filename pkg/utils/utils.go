package utils

import (
	"golang.org/x/crypto/ssh"
	"net"
)

func IgnoreHostKeyCB(_ string, _ net.Addr, _ ssh.PublicKey) error {
	return nil
}
