package v0 //import github.com/f0o/turbo-pancake/spec/job/v0

import (
	"errors"

	metadata "github.com/f0o/turbo-pancake/spec/metadata/v0"
	"github.com/f0o/turbo-pancake/utils"
)

//Spec : Spec for creating a Job
type Spec struct {
	ID      string        `json:"id"`
	Input   []byte        `json:"input"`
	Command string        `json:"command"`
	Filter  metadata.Spec `json:"filters"`
}

var Queue chan Spec = make(chan Spec, 1024)

func (e *Spec) Validate() error {
	if e.ID != "" {
		return errors.New("Cannot override ID")
	}
	return e.Finalize()
}

func (e *Spec) Finalize() error {
	e.ID, _ = utils.UUIDv4()
	return nil
}
