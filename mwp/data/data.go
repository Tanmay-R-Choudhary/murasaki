package data

type Data interface {
	Serialize() string
	Type() DataType
	Equals(Data) bool
	String() string
}
