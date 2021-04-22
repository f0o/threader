module github.com/f0o/turbo-pancake/utils

go 1.12

require github.com/f0o/turbo-pancake/common v0.0.0

replace (
	github.com/f0o/turbo-pancake/common => ../common
	github.com/f0o/turbo-pancake/utils => ../utils
)
