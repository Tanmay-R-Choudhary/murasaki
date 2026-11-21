package data

import (
	"fmt"
	"strconv"

	"github.com/Tanmay-R-Choudhary/murasaki/mwp/utils"
)

type MuInt struct {
	value int
}

func (muInt MuInt) Serialize() string {
	return fmt.Sprintf(":%d%s", muInt.value, utils.CRLF)
}

func (muInt MuInt) Type() DataType {
	return MuIntType
}

func (muInt MuInt) Value() int {
	return muInt.value
}

func (muInt MuInt) Equals(data Data) bool {
	if v, ok := data.(*MuInt); !ok {
		return false
	} else {
		return muInt.value == v.value
	}
}

func NewMuInt(value int) *MuInt {
	return &MuInt{value}
}

func (muInt MuInt) String() string {
	return strconv.Itoa(muInt.value)
}
