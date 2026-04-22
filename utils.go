package rhombifer

import (
	"fmt"
	"strings"
)

// IsFirstArgFlag if the argument being passed in is a flag. Returns true if it is and false if it isn't
func IsFirstArgFlag(arg string) bool {
	if strings.HasPrefix(arg, "-") || strings.HasPrefix(arg, "--") {
		return true
	}
	return false
}

// ExtractFlagValues expecified quantity of values from flag
func ExtractFlagValues(flag *Flag, quantity int) ([]string, error) {
	vals := make([]string, 0)
	if flag == nil {
		return vals, ErroFlagUndefined
	}

	if len(flag.Values) == 0 {
		ErroFlagHasNoValues.AppendMessage("flag name: " + flag.Name)
		return vals, ErroFlagHasNoValues
	}

	if quantity == 0 {
		return vals, fmt.Errorf("ExtractFlagValues func must have a quatity parameter equal or greater than 1")
	}

	for count := 0; count < quantity; count++ {
		vals = append(vals, flag.Values[count])
	}

	return vals, nil
}

// GetFlags found flags for the current executed command
func GetFlags() ([]*Flag, error) {
	if foundflags == nil {
		return nil, ErroFoundFlagsIsNil
	}
	return foundflags, nil
}

// FindFlags found flags in the current executed command and returns all flags specified in the aliases
func FindFlags(aliases ...string) ([]*Flag, error) {
	var flags []*Flag
	if foundflags == nil {
		return flags, ErroFoundFlagsIsNil
	}

	if len(aliases) == 0 {
		return flags, fmt.Errorf("no aliases provided")
	}

	flags = make([]*Flag, 0)

	for _, alias := range aliases {
		alias = strings.TrimPrefix(alias, "--")
		alias = strings.TrimPrefix(alias, "-")
		for _, f := range foundflags {
			if f.Name == alias || f.ShortFormat == alias {
				flags = append(flags, f)
			}
		}
	}
	return flags, nil
}

// FindFlag found flag and returns the first flag that matches any of the aliases provided. Once a flag has been
// matched the rest of the aliases are no searched, use [FindFlags] if that is your intent
func FindFlag(aliases ...string) (*Flag, error) {
	if foundflags == nil {
		return nil, ErroFoundFlagsIsNil
	}
	var flag *Flag
floop:
	for _, f := range foundflags {
		for _, alias := range aliases {
			alias = strings.TrimPrefix(alias, "--")
			alias = strings.TrimPrefix(alias, "-")
			if alias == f.Name || alias == f.ShortFormat {
				flag = f
				break floop
			}
		}
	}
	if flag == nil {
		return nil, ErroFlagNotFound
	}
	return flag, nil
}
