package book

type Meta interface {
	Value() string
	String() string
	URL() string
	IsNull() bool
}
