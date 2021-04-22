package v0 //import github.com/f0o/turbo-pancake/spec/metadata/v0

//Spec : Spec for creating a Job
type Spec struct {
	GOARCH string `json:"arch"`
	GOOS   string `json:"os"`
}
