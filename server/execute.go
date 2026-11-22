package main

import (
	"strings"

	"github.com/Tanmay-R-Choudhary/murasaki/mwp/data"
)

func execute(cmd data.Data) data.Data {
	arr, ok := cmd.(data.MuArray)
	if !ok || len(arr.Value()) == 0 {
		return data.NewMuString("ERR invalid command")
	}

	args := arr.Value()
	name, ok := args[0].(data.MuString)
	if !ok {
		return data.NewMuString("ERR invalid command")
	}

	switch strings.ToUpper(name.Value()) {
	case "PING":
		return data.NewMuString("PONG")
	case "ECHO":
		if len(args) != 2 {
			return data.NewMuString("ERR wrong number of arguments for 'echo' command")
		}
		return args[1]
	case "GET":
		if len(args) != 2 {
			return data.NewMuString("ERR wrong number of arguments for 'get' command")
		}
		return getStore().Get(args[1])
	case "SET":
		if len(args) != 3 {
			return data.NewMuString("ERR wrong number of arguments for 'set' command")
		}

		getStore().Set(args[1], args[2])
		return data.NewMuSimpleString("OK")
	default:
		return data.NewMuArray([]data.Data{data.NewMuString("Unrecognized")})
	}
}
