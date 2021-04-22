PREFIX ?= /usr/local/bin
.PHONY: all clean install test
all: bin/turbo-pancake

bin/turbo-pancake : cmd/main.go
	@mkdir -p bin
	go build -o bin/turbo-pancake cmd/main.go

install: $(PREFIX)/turbo-pancake

$(PREFIX)/turbo-pancake: bin/turbo-pancake
	install -s -m 0755 bin/turbo-pancake $(PREFIX)/turbo-pancake

clean:
	rm -rf bin tmp

test: bin/turbo-pancake
	mkdir -p tmp
	@echo Running Test: Base64 Decoding
	base64 < bin/turbo-pancake | __THREADER_INTERNAL=base64 bin/turbo-pancake > tmp/base64
	test "$$(md5sum bin/turbo-pancake | cut -d\  -f1)" "==" "$$(md5sum tmp/base64 | cut -d\  -f1)"
	@rm tmp/base64
	@echo Running Test: Single-Thread
	time test "$$(echo 'Foo' | bin/turbo-pancake -delimiter ',' -threads 1 -command 'echo $$INPUT' -out-format '%s')" == Foo
	@echo Running Test: Multi-Thread \#1
	time test "$$(echo -n '1,2,3' | bin/turbo-pancake -delimiter ',' -threads 2 -command 'echo -n $$INPUT;sleep $$INPUT' -out-format '%s')" == 123
	@echo Running Test: Multi-Thread \#2
	time test "$$(echo -n '1,2,3' | bin/turbo-pancake -delimiter ',' -threads 2 -command 'echo -n $$INPUT;sleep $$((3-$$INPUT))' -out-format '%s')" == 231
