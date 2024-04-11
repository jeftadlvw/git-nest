package models

type Config struct {
	IgnoreDuplicateOrigins    bool
	IgnoreDuplicateOriginRefs bool
}

func (c Config) Validate() error {
	return nil
}
