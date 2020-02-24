package v0 //import github.com/f0o/turbo-pancake/spec/return/v0

import (
	"errors"

	metadata "github.com/f0o/turbo-pancake/spec/metadata/v0"
	"github.com/f0o/turbo-pancake/utils"
)

//Spec : Spec for creating a Job
type Spec struct {
	ID       string        `json:"id"`
	Version  version       `json:"version"`
	Worker   string        `json:"worker"`
	Return   []byte        `json:"return"`
	Exit     int           `json:"exit"`
	Metadata metadata.Spec `json:"metadata"`
}

type version int

const (
	Version version = 0
)

//Queue : Object Queue
var Queue chan Spec = make(chan Spec, 1024)

//Validate : Validate Object
func (e *Spec) Validate() error {
	if e.ID != "" {
		return errors.New("Cannot override ID")
	}
	return e.Finalize()
}

//Finalize : Finalize Object
func (e *Spec) Finalize() error {
	e.ID, _ = utils.UUIDv4()
	e.Version = Version
	return nil
}

//Publish : Publish Object
func (e *Spec) Publish() error {
	err := e.Validate()
	if err == nil {
		Queue <- *e
		return nil
	}
	return err
}
