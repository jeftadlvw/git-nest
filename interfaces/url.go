package interfaces

type Url interface {
	Clean()
	IsEmpty() bool
	HostPathConcat() string
	String() string
	UnmarshalText(text []byte) error
	MarshalText(text []byte) ([]byte, error)
}
