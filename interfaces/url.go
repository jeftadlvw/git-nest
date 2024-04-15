package interfaces

type Url interface {
	String() string
	UnMarshalText(text []byte) error
	Marshaltext(text []byte) ([]byte, error)
}
