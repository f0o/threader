module github.com/f0o/turbo-pancake

go 1.12

require github.com/f0o/turbo-pancake/cli v0.0.0

replace (
	github.com/f0o/turbo-pancake/cli => ./src/cli
	github.com/f0o/turbo-pancake/common => ./src/common
	github.com/f0o/turbo-pancake/spec => ./src/spec
	github.com/f0o/turbo-pancake/utils => ./src/utils
	github.com/f0o/turbo-pancake/worker => ./src/worker
)
