package interfaces

/*
Url defines functions every Url-model should implement.
*/
type Url interface {
	/*
		Clean cleans up this Url struct.
	*/
	Clean()

	/*
		IsEmpty returns whether this Url is practically unusable.
	*/
	IsEmpty() bool

	/*
		HostPathConcat returns a concatenation of host, port and path.
	*/
	HostPathConcat() string

	/*
		String ensures the implementation of the Stringer interface.
		Provides a string representation of this Url.
	*/
	String() string

	/*
		UnmarshalText ensures the implementation of encoding.TextUnmarshaler interface.
		Unmarshals an Url struct from a slice of bytes.
	*/
	UnmarshalText(text []byte) error

	/*
		MarshalText ensures the implementation of encoding.TextMarshaler interface.
		Marshals this Url struct back into a slice of bytes.
	*/
	MarshalText(text []byte) ([]byte, error)
}
