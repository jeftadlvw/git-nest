package models

/*
Checksums contains checksums for required files at startup
*/
type Checksums struct {
	/*
		ConfigurationFile contains the checksum for the `nestmodules.toml` configuration file.
	*/
	ConfigurationFile string
}

/*
Validate performs validation on this Config.
*/
func (c Checksums) Validate() error {
	return nil
}
