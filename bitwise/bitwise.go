package bitwise

// Sets the bit at pos in the byte n.
func SetBit(n byte, pos uint64) byte {
	n |= (1 << pos)
	return n
}

// Clears the bit at pos in the byte n.
func ClearBit(n byte, pos uint64) byte {
	n &^= (1 << pos)
	return n
}

// Checks if the bit at pos in the byte n is set
func HasBit(n byte, pos uint64) bool {
	val := n & (1 << pos)
	return (val > 0)
}
