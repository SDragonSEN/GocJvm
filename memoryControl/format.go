package memCtrl

/******************************************************************
    []byte转int64型
******************************************************************/
func BytesToUint32(b []byte) uint32 {
	return uint32(b[0])<<24 | uint32(b[1])<<16 | uint32(b[2])<<8 | uint32(b[3])
}

/******************************************************************
    []byte转int64型
******************************************************************/
func Uint32ToBytes(n uint32, b []byte) {
	b[0] = byte(n >> 24)
	b[1] = byte(n << 8 >> 24)
	b[2] = byte(n << 16 >> 24)
	b[3] = byte(n << 24 >> 24)
}

/******************************************************************
    []byte转NodeHeader型
******************************************************************/
func FormatHeader(b []byte) NodeHeader {
	return NodeHeader{Size: BytesToUint32(b[0:4]), PreNode: BytesToUint32(b[4:8]), NextNode: BytesToUint32(b[8:12]), Type: b[12]}
}

/******************************************************************
    NodeHeader转[]byte型
******************************************************************/
func WriteHeader(nodeHeadr NodeHeader, b []byte) {
	Uint32ToBytes(nodeHeadr.Size, b[0:4])
	Uint32ToBytes(nodeHeadr.PreNode, b[4:8])
	Uint32ToBytes(nodeHeadr.NextNode, b[8:12])
	nodeHeadr.Type = b[12]
}
