package mwp

import (
	"errors"
	"strconv"
	"strings"

	"github.com/Tanmay-R-Choudhary/murasaki/mwp/data"
	"github.com/Tanmay-R-Choudhary/murasaki/mwp/utils"
)

func Deserialize(value string) (data.Data, int, error) {
	switch {
	case strings.HasPrefix(value, utils.SIMPLE_STRING_PREFIX):
		return deserializeSimpleString(value)
	case strings.HasPrefix(value, utils.STRING_PREFIX):
		return deserializeString(value)
	case strings.HasPrefix(value, utils.INT_PREFIX):
		return deserializeInt(value)
	case strings.HasPrefix(value, utils.ARRAY_PREFIX):
		return deserializeArray(value)
	default:
		return nil, 0, errors.New("Unexpected datatype received: " + value)
	}
}

func deserializeSimpleString(value string) (data.Data, int, error) {
	crlfIdx := strings.Index(value, utils.CRLF)
	if crlfIdx == -1 {
		return nil, 0, errors.New("Could not find CRLF terminator in stream")
	}
	return data.NewMuSimpleString(value[1:crlfIdx]), crlfIdx + 2, nil
}

func deserializeString(value string) (data.Data, int, error) {
	crlfIdx := strings.Index(value, utils.CRLF)
	if crlfIdx == -1 {
		return nil, 0, errors.New("Could not find CRLF terminator in stream")
	}
	unparsedStringLength := value[1:crlfIdx]
	parsedStringLength, err := strconv.Atoi(unparsedStringLength)
	if err != nil {
		return nil, 0, errors.New("Error occured while deserializing MuString length")
	}

	if parsedStringLength == -1 {
		return data.NewNullMuString(), 5, nil
	}
	return data.NewMuString(value[crlfIdx+2 : crlfIdx+2+parsedStringLength]), 5 + len(unparsedStringLength) + parsedStringLength, nil
}

func deserializeInt(value string) (data.Data, int, error) {
	crlfIdx := strings.Index(value, utils.CRLF)
	if crlfIdx == -1 {
		return nil, 0, errors.New("Could not find CRLF terminator in stream")
	}
	unparsedInt := value[1:crlfIdx]
	parsedInt, err := strconv.Atoi(unparsedInt)
	if err != nil {
		return nil, 0, errors.New("Error occured while deserializing MuInt value")
	}
	return data.NewMuInt(parsedInt), 3 + len(unparsedInt), nil
}

func deserializeArray(value string) (data.Data, int, error) {
	crlfIdx := strings.Index(value, utils.CRLF)
	if crlfIdx == -1 {
		return nil, 0, errors.New("Could not find CRLF terminator in stream")
	}
	unparsedArrayLength := value[1:crlfIdx]
	parsedArrayLength, err := strconv.Atoi(unparsedArrayLength)
	if err != nil {
		return nil, 0, errors.New("Error occured while deserializing MuArray length")
	}

	var elem []data.Data = []data.Data{}
	index := 3 + len(unparsedArrayLength)
	for range parsedArrayLength {
		arrayElement, offset, err := Deserialize(value[index:])
		if err != nil {
			return nil, 0, errors.New("Error occured while deserializing MuArray element -> " + err.Error())
		}

		elem = append(elem, arrayElement)
		index += offset
	}

	return data.NewMuArray(elem), index, nil
}
