module github.com/f0o/turbo-pancake/cli

go 1.12

require (
	github.com/f0o/turbo-pancake/common v0.0.0
	github.com/f0o/turbo-pancake/utils v0.0.0
	github.com/f0o/turbo-pancake/worker v0.0.0
)

replace (
	github.com/f0o/turbo-pancake/cli => ../cli
	github.com/f0o/turbo-pancake/common => ../common
	github.com/f0o/turbo-pancake/utils => ../utils
	github.com/f0o/turbo-pancake/worker => ../worker
)
