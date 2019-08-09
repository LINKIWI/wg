package cli

import (
	"flag"
	"fmt"
)

// ChoicesFlag decribes a flag.Value that only allows one of several candidate choices.
type ChoicesFlag struct {
	flag.Value

	candidates    []string
	defaultChoice string
	choice        string
}

// NewChoicesFlag creates a new ChoicesFlag with candidate choices and a default choice.
func NewChoicesFlag(choices []string, defaultChoice string) *ChoicesFlag {
	return &ChoicesFlag{candidates: choices, defaultChoice: defaultChoice}
}

// Set is called when the associated flag is specified on the CLI.
func (cf *ChoicesFlag) Set(value string) error {
	for _, candidate := range cf.candidates {
		if value == candidate {
			cf.choice = value
			return nil
		}
	}

	return fmt.Errorf("supplied value is not a valid choice; candidates=%v", cf.candidates)
}

// Choice returns the current choice set by the command-line flag.
func (cf *ChoicesFlag) Choice() string {
	if cf.choice == "" {
		return cf.defaultChoice
	}

	return cf.choice
}

// String acts as an alias for Choice.
func (cf *ChoicesFlag) String() string {
	return cf.Choice()
}
