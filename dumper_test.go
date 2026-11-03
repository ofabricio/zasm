package zasm

import (
	"fmt"
	"os"
)

func ExampleDumpWriter() {

	src := `
		@dump
		abc 90
		    89 D8
			83 C0 01
		@end
		@dump
			89 D8
			FF C0
		@end
	`

	var prg Program
	oks := Parse(src, &prg)
	fmt.Println(oks)

	lb := prg.Body[0].(*DirDump).Body[0].(*LabelDef)
	i0 := prg.Body[0].(*DirDump).Body[1].(*DataDef)
	i1 := prg.Body[0].(*DirDump).Body[2].(*DataDef)
	i2 := prg.Body[0].(*DirDump).Body[3].(*DataDef)
	i3 := prg.Body[1].(*DirDump).Body[0].(*DataDef)
	i4 := prg.Body[1].(*DirDump).Body[1].(*DataDef)

	d := DumpWriter{out: os.Stdout}

	d.SetBytes(4)
	d.SetOffset(0)

	d.WriteLabel(lb)
	d.WriteInstr(i0, []byte{0x90})
	d.WriteInstr(i1, []byte{0x89, 0xD8})
	d.WriteInstr(i2, []byte{0x83, 0xC0, 0x01})

	d.SetBytes(1)
	d.WriteInstr(i3, []byte{0x89, 0xD8})
	d.WriteInstr(i4, []byte{0xFF, 0xC0})

	// Output:
	// true
	//     abc
	//     00000000  90             90
	//     00000001  89 D8          89 D8
	//     00000003  83 C0 01       83 C0 01
	//     00000006  89    89 D8
	//     00000007  D8
	//     00000008  FF    FF C0
	//     00000009  C0
}
