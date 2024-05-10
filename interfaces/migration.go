package interfaces

/*
ModuleMigration defines an interface for performing migration on
nested modules.
*/
type ModuleMigration interface {

	/*
		Migrate migrates the module.
	*/
	Migrate() error
}
