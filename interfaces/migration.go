package interfaces

/*
Migration defines an interface for performing migration on
nested modules.
*/
type Migration interface {

	/*
		Migrate migrates the module.
	*/
	Migrate() error
}
