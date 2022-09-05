.PHONY: format
format:
	go mod download; \
	go mod tidy; \
	build/format.sh

.PHONY: clean
clean:
	build/clean.sh

.PHONY: all
all: format
	build/build.sh
