package vm

type register struct {
	// program counter
	pc uint16

	// condition flags
	r_cond  byte
	r_count byte

	general [8]uint32
}
