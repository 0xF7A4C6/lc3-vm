package vm

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"log"
)

type Vm struct {
	memory   []byte
	register register
}

func NewVM() *Vm {
	return &Vm{
		memory: make([]byte, 1<<16),
	}
}

func (vm *Vm) LoadRom(filePath string) error {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	baseAddr := binary.BigEndian.Uint16(content[:16]) // should be 3000
	log.Printf("Loading ROM base addr: %v\n", baseAddr)

	if len(content) > len(vm.memory)-int(baseAddr) {
		return fmt.Errorf("not enough space in memory to load ROM")
	}

	copy(vm.memory[baseAddr:], content[16:])
	vm.register.pc = baseAddr

	return nil
}

func (vm *Vm) getByte() byte {
	defer func() {
		vm.register.pc++
	}()

	return vm.memory[vm.register.pc]
}

func (vm *Vm) SignExtend(x byte, bits int32) uint32 {
	return uint32(int32(x<<bits) >> bits)
}

func (vm *Vm) updateFlag(val int) {
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
		switch op := vm.getByte(); op {
		default:
			panic(fmt.Sprintf("bad opcode '%d' at _pc %d", op, vm.register.pc))
		case OP_ADD:
			r0 := (op >> 9) & 0x7
			r1 := (op >> 6) & 0x7
			imm_flag := (op >> 5) & 0x1

			if r0 == 0 {
				continue
			}

			fmt.Println(r0, r1, imm_flag)

			if imm_flag == 1 {
				imm5 := vm.SignExtend(op&0x1F, 5)
				vm.register.general[r0] = vm.register.general[r1] + imm5
			}
		}
	}
}
