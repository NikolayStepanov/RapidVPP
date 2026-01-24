package domain

import (
	"net"
)

type IPWithPrefix struct {
	Address string
	Prefix  uint8
}

func (ip IPWithPrefix) ToNetIP() net.IP {
	return net.ParseIP(ip.Address)
}

type ACLInterfaceList struct {
	InterfaceID uint32
	Count       uint8
	InputACLs   []uint32
	OutputACLs  []uint32
}
