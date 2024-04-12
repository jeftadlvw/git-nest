package models

type Config struct {
	AllowDuplicateOrigins    bool
	AllowDuplicateOriginRefs bool
	AllowUnequalRoots        bool
}

func (c Config) Validate() error {
	return nil
}
