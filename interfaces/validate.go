package interfaces

/*
Validator is an interface that provide a Validate() function.
*/
type Validator interface {

	/*
		Validate cleans and validates the underlying data structure and returns an error if anything is invalid.
	*/
	Validate() error
}
