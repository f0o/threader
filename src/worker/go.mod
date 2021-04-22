module github.com/f0o/turbo-pancake/worker

go 1.12

require (
	github.com/f0o/turbo-pancake/common v0.0.0
	github.com/f0o/turbo-pancake/spec v0.0.0
	github.com/f0o/turbo-pancake/utils v0.0.0
)

replace (
	github.com/f0o/turbo-pancake/common => ../common
	github.com/f0o/turbo-pancake/spec => ../spec
	github.com/f0o/turbo-pancake/utils => ../utils
	github.com/f0o/turbo-pancake/worker => ../worker
)
