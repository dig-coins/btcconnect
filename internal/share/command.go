package share

import (
	"encoding/hex"
	"encoding/json"

	"github.com/sgostarter/i/commerr"
)

type CommandType int

const (
	CommandTypeUnknown CommandType = iota
	CommandTypeGenTx
	CommandTypeUpdateTx
	CommandTypeMax
)

type Command struct {
	CommandType CommandType `json:"command_type" yaml:"command_type"`
	Input       string      `json:"input" yaml:"input"`
}

func (command *Command) Valid() bool {
	if command.CommandType <= CommandTypeUnknown || command.CommandType >= CommandTypeMax {
		return false
	}

	return command.Input != ""
}

func UnmarshalCommand(s string) (command Command, err error) {
	d, err := hex.DecodeString(s)
	if err != nil {
		return
	}

	err = json.Unmarshal(d, &command)
	if err != nil {
		return
	}

	if !command.Valid() {
		err = commerr.ErrInvalidArgument

		return
	}

	return
}

func MarshalCommand(command Command) (s string, err error) {
	if !command.Valid() {
		err = commerr.ErrInvalidArgument

		return
	}

	d, err := json.Marshal(&command)
	if err != nil {
		return
	}

	s = hex.EncodeToString(d)

	return
}

func MarshalCommandToJSON(command Command) (s string, err error) {
	if !command.Valid() {
		err = commerr.ErrInvalidArgument

		return
	}

	d, err := json.Marshal(&command)
	if err != nil {
		return
	}

	s = string(d)

	return
}
