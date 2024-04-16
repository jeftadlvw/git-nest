package models

type Config struct {
	AllowDuplicateOrigins bool `toml:"allow_duplicate_origins"`
	AllowUnequalRoots     bool `toml:"allow_unequal_roots"`
}

func (c Config) Validate() error {
	return nil
}
