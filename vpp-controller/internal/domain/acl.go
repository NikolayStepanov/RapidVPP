package domain

type ACLAction uint8
type AclID uint32

const (
	ACLDeny   ACLAction = 0
	ACLPermit ACLAction = 1
)

type ACLRule struct {
	Action        ACLAction
	Proto         uint8
	Src           IPWithPrefix
	Dst           IPWithPrefix
	SrcPortLow    uint16
	SrcPortHigh   uint16
	DstPortLow    uint16
	DstPortHigh   uint16
	TCPFlagsMask  uint8
	TCPFlagsValue uint8
}

type ACLInfo struct {
	ID    AclID
	Name  string
	Rules []ACLRule
}
