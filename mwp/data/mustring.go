package data

import (
	"fmt"

	"github.com/Tanmay-R-Choudhary/murasaki/mwp/utils"
)

type MuString struct {
	size int
	data string
}

func (muString MuString) Serialize() string {
	return fmt.Sprintf("$%d%s%s%s", muString.size, utils.Crlf, muString.data, utils.Crlf)
}

func (muString MuString) Type() DataType {
	return MuStringType
}

func (muString MuString) Equals(data Data) bool {
	if v, ok := data.(MuString); !ok {
		return false
	} else {
		return muString.data == v.data && muString.size == v.size
	}
}

func (muString MuString) Value() string {
	return muString.data
}

func NewMuString(data string) MuString {
	size := len(data)
	return MuString{size, data}
}

func NewNullMuString() MuString {
	return MuString{size: -1, data: ""}
}

func (muString MuString) String() string {
	if muString.size == -1 {
		return "MuString: NULL"
	} else {
		return "MuString: " + muString.data
	}
}
