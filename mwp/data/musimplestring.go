package data

import (
	"fmt"

	"github.com/Tanmay-R-Choudhary/murasaki/mwp/utils"
)

type MuSimpleString struct {
	data string
}

func (muSimpleString MuSimpleString) Serialize() string {
	return fmt.Sprintf("+%s%s", muSimpleString.data, utils.CRLF)
}

func (muSimpleString MuSimpleString) Type() DataType {
	return MuSimpleStringType
}

func (muSimpleString MuSimpleString) Equals(data Data) bool {
	if v, ok := data.(*MuSimpleString); !ok {
		return false
	} else {
		return muSimpleString.data == v.data
	}
}

func (muSimpleString MuSimpleString) Value() string {
	return muSimpleString.data
}

func NewMuSimpleString(data string) *MuSimpleString {
	return &MuSimpleString{data}
}

func (muSimpleString MuSimpleString) String() string {
	return muSimpleString.data
}
