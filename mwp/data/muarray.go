package data

import (
	"fmt"
	"strings"

	"github.com/Tanmay-R-Choudhary/murasaki/mwp/utils"
)

type MuArray struct {
	size uint
	elem []Data
}

func (muArr MuArray) Serialize() string {
	var elements strings.Builder
	for _, e := range muArr.elem {
		elements.WriteString(e.Serialize())
	}

	return fmt.Sprintf("*%d%s%s", muArr.size, utils.CRLF, elements.String())
}

func (muArr MuArray) Type() DataType {
	return MuArrayType
}

func (muArr MuArray) Len() (uint, error) {
	return muArr.size, nil
}

func (muArr MuArray) Value() []Data {
	return muArr.elem
}

func (muArr MuArray) Equals(data Data) bool {
	if v, ok := data.(*MuArray); !ok {
		return false
	} else {
		if len(v.elem) != len(muArr.elem) {
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

func (muArr MuArray) String() string {
	res := strings.Builder{}
	for _, v := range muArr.elem {
		res.WriteString(v.String())
		res.WriteString("\n")
	}

	return res.String()
}
