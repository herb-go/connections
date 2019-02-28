package command

import (
	"bytes"
)

type Command interface {
	Encode() ([]byte, error)
	Decode([]byte) error
	Type() string
	Data() []byte
}

var defaultSeparator = []byte(" ")

type SeparatedCommand struct {
	CommandType string
	CommandData []byte
	Separator   []byte
}

func New() *SeparatedCommand {
	return &SeparatedCommand{
		Separator: defaultSeparator,
	}
}
func (c *SeparatedCommand) Type() string {
	return c.CommandType
}
func (c *SeparatedCommand) Data() []byte {
	return c.CommandData
}
func (c *SeparatedCommand) Encode() ([]byte, error) {
	return bytes.Join([][]byte{[]byte(c.CommandType), c.CommandData}, c.Separator), nil
}
func (c *SeparatedCommand) Decode(bs []byte) error {
	cmds := bytes.SplitN(bs, c.Separator, 2)
	if len(cmds) > 0 {
		c.CommandType = string(cmds[0])
	}
	if len(cmds) > 1 {
		c.CommandData = cmds[1]
	}
	return nil
}
