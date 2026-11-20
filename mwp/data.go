package mwp

import (
	"fmt"
	"strings"
)

type Data interface {
	Serialize() string
	Type() DataType
	Equals(Data) bool
}

type MuString struct {
	data string
}

func (muString *MuString) Serialize() string {
	return fmt.Sprintf("$%d%s%s%s", len(muString.data), CRLF, muString.data, CRLF)
}

func (*MuString) Type() DataType {
	return MuStringType
}

func (muString *MuString) Equals(data Data) bool {
	if v, ok := data.(*MuString); !ok {
		return false
	} else {
		return muString.data == v.Value()
	}
}

func (muString *MuString) Value() string {
	return muString.data
}

func NewMuString(data string) *MuString {
	return &MuString{data}
}

type MuInt struct {
	value int
}

func (muInt *MuInt) Serialize() string {
	return fmt.Sprintf(":%d%s", muInt.value, CRLF)
}

func (*MuInt) Type() DataType {
	return MuIntType
}

func (muInt *MuInt) Value() int {
	return muInt.value
}

func (muInt *MuInt) Equals(data Data) bool {
	if v, ok := data.(*MuInt); !ok {
		return false
	} else {
		return muInt.value == v.Value()
	}
}

func NewMuInt(value int) *MuInt {
	return &MuInt{value}
}

type MuArray struct {
	size uint
	elem []Data
}

func (muArr *MuArray) Serialize() string {
	var elements strings.Builder
	for _, e := range muArr.elem {
		elements.WriteString(e.Serialize())
	}

	return fmt.Sprintf("*%d%s%s", muArr.size, CRLF, elements.String())
}

func (*MuArray) Type() DataType {
	return MuArrayType
}

func (muArr *MuArray) Len() (uint, error) {
	return muArr.size, nil
}

func (muArr *MuArray) Value() []Data {
	return muArr.elem
}

func (muArr *MuArray) Equals(data Data) bool {
	if v, ok := data.(*MuArray); !ok {
		return false
	} else {
		if len(v.Value()) != len(muArr.elem) {
			return false
		}

		for idx, e := range v.Value() {
			if !e.Equals(muArr.elem[idx]) {
				return false
			}
		}

		return true
	}
}

func NewMuArray(elem []Data) *MuArray {
	return &MuArray{uint(len(elem)), elem}
}
