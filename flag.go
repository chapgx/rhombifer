package rhombifer

import (
	"errors"
	"fmt"
)

type Flag struct {
	// flag name
	Name string

	// Shot Desctiption
	Short string

	// Long Description
	Long string

	// Short Format for flag if any
	ShortFormat string

	// Is this a require flag
	// May remove this and add a require flags to the command itself
	Required bool

	// todo: not sure if to keep it
	// to define how the values would be parse
	ValuesDilimiter []string

	// Does this flag requires a value
	RequiresValue bool

	// Defines if the flag takes one value or multiple values
	SingleValue bool

	// Flag values parsed from the current command being run
	Values []string
}

func (f *Flag) AddValues(args ...string) error {
	if f.RequiresValue && len(args) <= 0 {
		return fmt.Errorf("flag requires values but got 0")
	}
	if f.SingleValue && len(args) > 1 {
		return fmt.Errorf("flag only accepts one value but got %d values", len(args))
	}
	if f.Values == nil {
		f.Values = make([]string, 0)
	}
	f.Values = append(f.Values, args...)
	return nil
}

// Returns the short and long format name for the flag
func (f *Flag) GetNames() (short, long string) {
	return f.ShortFormat, f.Name
}

// Grabs the first value of the flag
func (f Flag) GetSingleValue() (string, error) {
	if f.Values == nil {
		return "", ErrNilValues
	}

	if len(f.Values) < 0 {
		return "", ErrNoValues
	}

	return f.Values[0], nil
}

// TODO: add a validation function for type maybe
var (
	ErrNoValues  = errors.New("no values found on flag")
	ErrNilValues = errors.New("values is <nil>")
)
