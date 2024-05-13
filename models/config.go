package models

/*
Config represents all git-nest configurable options.
*/
type Config struct {
	/*
		AllowDuplicateOrigins defines whether duplicate remote origins are allowed.
	*/
	AllowDuplicateOrigins bool `toml:"allow_duplicate_origins"`

	/*
		AllowUnequalRoots defines whether the project's git root and git-nest root are allowed to not be aligned.
	*/
	AllowUnequalRoots bool `toml:"allow_unequal_roots"`
}

/*
Validate performs validation on this Config.
*/
func (c Config) Validate() error {
	return nil
}
