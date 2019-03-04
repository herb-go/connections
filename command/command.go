package command

import (
	"bytes"
)

// Command connection inbound command interface
type Command interface {
	// Encode encode command to byte slice.
	//Return encoded byte slice and any error raised.
	Encode() ([]byte, error)
	//Decode decode command by given byte slice.
	Decode([]byte) error
	// Type return command type
	Type() string
	//Data return command data
	Data() []byte
}

var defaultSeparator = []byte(" ")

// SeparatedCommand command which separated data with Separator
type SeparatedCommand struct {
	// CommandType command type
	CommandType string
	// CommandData command data
	CommandData []byte
	//Separator  data separator
	Separator []byte
}

// New create separated command with  default separator(space)
func New() *SeparatedCommand {
	return &SeparatedCommand{
		Separator: defaultSeparator,
	}
}

// Type return command type
func (c *SeparatedCommand) Type() string {
	return c.CommandType
}

//Data return command data
func (c *SeparatedCommand) Data() []byte {
	return c.CommandData
}

// Encode encode command to byte slice.
//Return encoded byte slice and any error raised.
func (c *SeparatedCommand) Encode() ([]byte, error) {
	return bytes.Join([][]byte{[]byte(c.CommandType), c.CommandData}, c.Separator), nil
}

//Decode decode command by given byte slice.
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
