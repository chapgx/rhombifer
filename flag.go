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

	// Does this flag requires a value
	RequiresValue bool

	// Defines if the flag takes one value or multiple values
	SingleValue bool

	// Flag values parsed from the current command being run
	Values []string

	Command *Command

	Run Run
}

// AddValues add values to the flag
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

// GetNames returns the short and long format name for the flag
func (f *Flag) GetNames() (short, long string) {
	return f.ShortFormat, f.Name
}

// GetSingleValue returns the first value of the flag
func (f Flag) GetSingleValue() (string, error) {
	if f.Values == nil {
		return "", ErrNilValues
	}

	if len(f.Values) == 0 {
		return "", ErrNoValues
	}

	return f.Values[0], nil
}

// SetRequired sets the flag to require and returns a pointer to instance f
func (f *Flag) SetRequired() *Flag {
	f.Required = true
	return f
}

// Exec sets a run function if any and returns a pointer to instane f
func (f *Flag) Exec(fn Run) *Flag {
	f.Run = fn
	return f
}

// SetShortFormat sets the short format for the flag
func (f *Flag) SetShortFormat(sf string) *Flag {
	f.ShortFormat = sf
	return f
}

// SetValuesRequired sets the values to required for the flag
func (f *Flag) SetValuesRequired() *Flag {
	f.RequiresValue = true
	return f
}

// NewFlag is a short hand for when you really only need flag present and
// the description
func NewFlag(name, shortdesc string) *Flag {
	return &Flag{Name: name, Short: shortdesc}
}

var (
	ErrNoValues  = errors.New("no values found on flag")
	ErrNilValues = errors.New("values is <nil>")
)
