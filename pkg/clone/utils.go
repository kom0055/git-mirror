package clone

import (
	"net"

	"golang.org/x/crypto/ssh"
)

func IgnoreHostKeyCB(_ string, _ net.Addr, _ ssh.PublicKey) error {
	return nil
}
