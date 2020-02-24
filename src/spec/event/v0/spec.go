package v0 //import event "github.com/f0o/turbo-pancake/spec/event/v0"

import (
	"errors"
	"fmt"

	"github.com/f0o/turbo-pancake/utils"
)

type Spec struct {
	ID      string
	Version version
	Type    Type
	Verb    Verb
	Payload interface{}
}

type Type string

const (
	TypeWorker Type = "worker"
	TypeJob    Type = "job"
	TypeReturn Type = "return"
)

type Verb string

const (
	VerbCreate Verb = "create"
	VerbUpdate Verb = "update"
	VerbDelete Verb = "delete"
	VerbReturn Verb = "return"
)

type version int

const (
	Version version = 0
)

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

//Publish : Publish Event
func (e *Spec) Publish() {
	if e.Validate() == nil {
		fmt.Printf("%+v\n", *e)
	} else {
		//something
	}
}
