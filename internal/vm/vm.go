package vm

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"log"
)

type Vm struct {
	memory   []uint16
	register register
}

func NewVM() *Vm {
	return &Vm{
		memory: make([]uint16, 1<<16),
	}
}

func (vm *Vm) LoadRom(filePath string) error {
	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	content := fileContent[2:]

	var origin uint16
	origin = binary.BigEndian.Uint16(fileContent[:2])

	log.Printf("Origin memory location: 0x%04X", origin)
	vm.register.pc = origin

	for count := 0; count < len(content); count += 2 {
		value := binary.BigEndian.Uint16(content[count : count+2])

		vm.memory[vm.register.pc] = value
		vm.register.pc++
	}

	return nil

}

func (vm *Vm) getByte() uint16 {
	defer func() {
		vm.register.pc++
	}()

	return vm.memory[vm.register.pc]
}

func (vm *Vm) signExtend(x, bitCount uint16) uint16 {
	if (x>>(bitCount-1))&1 == 1 {
		x |= 0xFFFF << bitCount
	}
	return x & 0xFFFF
}

func (vm *Vm) updateFlag(val uint16) {
	switch val {
	case val >> 15:
		vm.register.r_cond = FL_NEG
	case 0:
		vm.register.r_cond = FL_ZRO
	default:
		vm.register.r_cond = FL_POS
	}
}

func (vm *Vm) Run() {
	for {
		op := vm.getByte()

		switch op {
		default:
			panic(fmt.Sprintf("bad opcode '%d' at _pc %d", op, vm.register.pc))
		case OP_RES:
			panic(fmt.Sprintf("OP_RES at _pc %d", vm.register.pc))
		case OP_TRAP:
			vm.execTrap(op)
		case OP_BR:
			pc_offset := vm.signExtend(op&0x1FF, 9)
			cond := (op >> 9) & 0x7

			if cond == vm.register.r_cond {
				vm.register.pc += pc_offset
			}
		case OP_ADD:
			r0 := (op >> 9) & 0x7
			r1 := (op >> 6) & 0x7
			imm_flag := (op >> 5) & 0x1

			fmt.Println(r0, r1, imm_flag)

			if imm_flag == MODE_IMMEDIATE {
				imm5 := vm.signExtend(op&0x1F, 5)
				vm.register.general[r0] = vm.register.general[r1] + imm5
			} else {
				r2 := op & 0x7
				vm.register.general[r0] = vm.register.general[r1] + vm.register.general[r2]
			}

			vm.updateFlag(r0)
		case OP_LDI:
			r0 := (op >> 9) & 0x7
			pc_offset := vm.signExtend(op&0x1FF, 9)

			vm.register.general[r0] = vm.memRead(vm.memRead(vm.register.pc + pc_offset))
			vm.updateFlag(r0)
		case OP_AND:
			r0 := (op >> 9) & 0x7
			r1 := (op >> 6) & 0x7
			imm_flag := (op >> 5) & 0x1

			if imm_flag == MODE_IMMEDIATE {
				imm5 := vm.signExtend(op&0x1F, 5)
				vm.register.general[r0] = vm.register.general[r1] & imm5
			} else {
				r2 := op & 0x7
				vm.register.general[r0] = vm.register.general[r1] & vm.register.general[r2]
			}

			vm.updateFlag(r0)
		case OP_NOT:
			r0 := (op >> 9) & 0x7
			r1 := (op >> 6) & 0x7

			vm.register.general[r0] = ^vm.register.general[r1]
			vm.updateFlag(r0)
		case OP_JMP:
			r1 := (op >> 6) & 0x7
			vm.register.pc = vm.register.general[r1]
		case OP_JSR:
			long_flag := (op >> 11) & 1
			vm.register.general[7] = vm.register.pc

			if long_flag == MODE_IMMEDIATE {
				long_pc_offset := vm.signExtend(op>>0x7FF, 11)
				vm.register.pc += long_pc_offset // JSR
			} else {
				r1 := (op >> 6) & 0x7
				vm.register.pc = vm.register.general[r1] // JSRR
			}
		case OP_LD:
			r0 := (op >> 9) & 0x7
			pc_offset := vm.signExtend(op&0x1FF, 9)

			vm.register.general[r0] = vm.memRead(vm.register.pc + pc_offset)
			vm.updateFlag(r0)
		case OP_LDR:
			r0 := (op >> 9) & 0x7
			r1 := (op >> 6) & 0x7
			offset := vm.signExtend(op&0x3F, 6)

			vm.register.general[r0] = vm.memRead(vm.register.general[r1] + offset)
			vm.updateFlag(r0)
		case OP_LEA:
			r0 := (op >> 9) & 0x7
			pc_offset := vm.signExtend(op&0x1FF, 9)

			vm.register.general[r0] = vm.register.pc + pc_offset
			vm.updateFlag(r0)
		case OP_ST:
			r0 := (op >> 9) & 0x7
			pc_offset := vm.signExtend(op&0x1FF, 9)

			vm.memWrite(vm.register.pc+pc_offset, vm.register.general[r0])
		case OP_STI:
			r0 := (op >> 9) & 0x7
			pc_offset := vm.signExtend(op&0x1FF, 9)

			vm.memWrite(vm.memRead(vm.register.pc+pc_offset), vm.register.general[r0])
		case OP_STR:
			r0 := (op >> 9) & 0x7
			r1 := (op >> 6) & 0x7
			offset := vm.signExtend(op&0x3F, 6)

			vm.memWrite(vm.register.general[r1]+offset, vm.register.general[r0])
		}
	}
}

/*
 */
