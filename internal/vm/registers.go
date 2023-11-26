package vm

type register struct {
	// program counter
	pc uint16

	// condition flags
	r_cond  uint16
	r_count uint16

	general [8]uint16
}
