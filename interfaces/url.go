package interfaces

/*
Url defines functions every Url-model should implement.
*/
type Url interface {

	/*
		Use Validator interface.
	*/
	Validator

	/*
		Clean cleans up this Url struct.
	*/
	Clean()

	/*
		IsEmpty returns whether this Url is practically unusable.
	*/
	IsEmpty() bool

	/*
		Hostname returns this Url hostname.
	*/
	Hostname() string

	/*
		Path returns this Url path.
	*/
	Path() string

	/*
		HostPathConcat returns a concatenation of host, port and path.
	*/
	HostPathConcat() string

	/*
		HostPathConcat returns a concatenation of host, port and path with stricter output rules.
	*/
	HostPathConcatStrict() string

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
	MarshalText() ([]byte, error)
}
