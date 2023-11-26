package vm

import (
	"fmt"
	"os"
)

const (
	TRAP_GETC  = 0x20 /* get character from keyboard, not echoed onto the terminal */
	TRAP_OUT   = 0x21 /* output a character */
	TRAP_PUTS  = 0x22 /* output a word string */
	TRAP_IN    = 0x23 /* get character from keyboard, echoed onto the terminal */
	TRAP_PUTSP = 0x24 /* output a byte string */
	TRAP_HALT  = 0x25 /* halt the program */
)

func (vm *Vm) execTrap(op uint16) {
	vm.register.general[7] = vm.register.pc

	switch op & 0xFF {
	case TRAP_GETC:
		for {
			var input byte
			_, err := fmt.Scanf("%c", &input)
			if err != nil {
				continue
			}

			vm.register.general[0] = uint16(input)
			break
		}
	case TRAP_OUT:
		character := byte(vm.register.general[0] & 0xFF)

		fmt.Printf("%c", character)
		if err := os.Stdout.Sync(); err != nil {
			panic(err)
		}
	case TRAP_PUTS:
		c := &vm.memory[vm.register.general[0]]

		for *c != 0 {
			char := byte(*c & 0xFF)
			fmt.Printf("%c", char)
			c = &vm.memory[*c]
		}

		err := os.Stdout.Sync()
		if err != nil {
			panic(err)
		}
	case TRAP_IN:
		for {
			var input byte
			_, err := fmt.Scanf("%c", &input)
			if err != nil {
				continue
			}

			fmt.Fprintf(os.Stdout, "%c", input)

			vm.register.general[0] = uint16(input)
			vm.updateFlag(0)
			break
		}
	case TRAP_PUTSP:
		c := &vm.memory[0]

		for *c != 0 {
			char1 := byte(*c & 0xFF)
			fmt.Printf("%c", char1)

			char2 := byte((*c >> 8) & 0xFF)
			if char2 != 0 {
				fmt.Printf("%c", char2)
			}

			c = &vm.memory[*c]
		}

		if err := os.Stdout.Sync(); err != nil {
			panic(err)
		}
	case TRAP_HALT:
		panic("halt")
	}
}

/*

 */
