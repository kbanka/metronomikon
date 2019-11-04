BINARY=metronomikon

.phony: all build clean run test

all: build

build: $(BINARY)

clean:
	rm -f $(BINARY)

$(BINARY): $(shell find -name '*.go')
	go build -o $(BINARY)

run:
	go run $(BINARY).go

test:
	find -type f -name '*_test.go' | xargs -r dirname | sort -u | while read package; do \
		echo $$package; \
		go test -v $$package || exit 1; \
	done
