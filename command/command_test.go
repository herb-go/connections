package command

import (
	"bytes"
	"testing"
)

func TestCmds(t *testing.T) {
	var type1 = "cmdtype1"
	var type2 = "cmdtype2"
	var data1 = "cmd1"
	var dataempty = ""
	cmd1 := New()
	err := cmd1.Decode([]byte(type1 + " " + data1))
	if err != nil {
		t.Fatal(err)
	}
	if cmd1.Type() != type1 {
		t.Fatal(cmd1.Type())
	}
	if string(cmd1.Data()) != data1 {
		t.Fatal(cmd1.Type())
	}
	bs, err := cmd1.Encode()
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Compare(bs, []byte(type1+" "+data1)) != 0 {
		t.Fatal(bs)
	}
	cmd2 := New()
	err = cmd2.Decode([]byte(type2 + " "))
	if err != nil {
		t.Fatal(err)
	}
	if cmd2.Type() != type2 {
		t.Fatal(cmd2.Type())
	}
	if string(cmd2.Data()) != "" {
		t.Fatal(cmd2.Data())
	}
	cmdempty := New()
	err = cmdempty.Decode([]byte(type2 + dataempty))
	if err != nil {
		t.Fatal(err)
	}
	if cmdempty.Type() != type2 {
		t.Fatal(cmd1.Type())
	}
	if string(cmdempty.Data()) != "" {
		t.Fatal(cmdempty.Type())
	}
}
