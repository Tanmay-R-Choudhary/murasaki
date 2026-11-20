package mwp

import (
	"errors"
	"strconv"
	"strings"
)

func Deserialize(value string) (Data, int, error) {
	switch {
	case strings.HasPrefix(value, STRING_PREFIX):
		return deserializeString(value)
	case strings.HasPrefix(value, INT_PREFIX):
		return deserializeInt(value)
	case strings.HasPrefix(value, ARRAY_PREFIX):
		return deserializeArray(value)
	default:
		return nil, 0, errors.New("Unexpected datatype received: " + value)
	}
}

func deserializeString(value string) (Data, int, error) {
	crlfIdx := strings.Index(value, CRLF)
	unparsedStringLength := value[1:crlfIdx]
	parsedStringLength, err := strconv.Atoi(unparsedStringLength)
	if err != nil {
		return nil, 0, errors.New("Error occured while deserializing MuString length")
	}
	return NewMuString(value[crlfIdx+2 : crlfIdx+2+parsedStringLength]), 5 + len(unparsedStringLength) + parsedStringLength, nil
}

func deserializeInt(value string) (Data, int, error) {
	crlfIdx := strings.Index(value, CRLF)
	unparsedInt := value[1:crlfIdx]
	parsedInt, err := strconv.Atoi(unparsedInt)
	if err != nil {
		return nil, 0, errors.New("Error occured while deserializing MuInt value")
	}
	return NewMuInt(parsedInt), 3 + len(unparsedInt), nil
}

func deserializeArray(value string) (Data, int, error) {
	crlfIdx := strings.Index(value, CRLF)
	unparsedArrayLength := value[1:crlfIdx]
	parsedArrayLength, err := strconv.Atoi(unparsedArrayLength)
	if err != nil {
		return nil, 0, errors.New("Error occured while deserializing MuArray length")
	}

	var elem []Data = []Data{}
	index := 3 + len(unparsedArrayLength)
	for range parsedArrayLength {
		arrayElement, offset, err := Deserialize(value[index:])
		if err != nil {
			return nil, 0, errors.New("Error occured while deserializing MuArray element -> " + err.Error())
		}

		elem = append(elem, arrayElement)
		index += offset
	}

	return NewMuArray(elem), index, nil
}
