package internal

import (
	"fmt"
)

func ArgExactN(n int) func(c int) error {
	return func(c int) error {
		if c != n {
			return fmt.Errorf("exactly %d argument%s allowed", n, pluralS(n))
		}

		return nil
	}
}

func ArgMinN(n int) func(c int) error {
	return func(c int) error {
		if c < n {
			return fmt.Errorf("minimal %d argument%s required", n, pluralS(n))
		}

		return nil
	}
}

func ArgMaxN(n int) func(c int) error {
	return func(c int) error {
		if c > n {
			return fmt.Errorf("maximal %d argument%s allowed", n, pluralS(n))
		}

		return nil
	}
}

func ArgRangeNM(n int, m int) func(c int) error {
	return func(c int) error {
		if c < n || c > m {
			return fmt.Errorf("argument count must be between %d and %d", n, m)
		}

		return nil
	}
}

func pluralS(c int) string {
	plural := ""
	if c != 1 {
		plural = "s"
	}
	return plural
}
