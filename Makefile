## test/all runs all related tests.
test/all: test/lint test

## test runs `go test`
test:
	go test -v -p 2 -count 1 -timeout 240s -race ./...

## test runs `go test -run $(RUN)`
test/run:
	go test -v -p 2 -count 1 -timeout 240s -race ./... -run $(RUN)

## test/lint runs linter
test/lint:
	# checks the coding style.
	(! gofmt -s -d `find . -name vendor -prune -type f -o -name '*.go'` | grep '^')
	# checks the import format.
	(! goimports -l `find . -name vendor -prune -type f -o -name '*.go'` | grep -v 'pb.go' | grep 'go')
	# checks the error the compiler can't find.
	go vet ./...
	# checks shadowed variables.
	go vet -vettool=$(which shadow) ./...
	# checks no used assigned value.
	ineffassign ./...
	# checks not to ignore the error.
	errcheck ./...
	# checks unused global variables and constants.
	varcheck ./...
	# checks dispensable type conversions.
	unconvert -v ./...

## dev/install/dep installs depenencies required for development.
dev/install/dep:
	sh ./.github/install_dep.sh

## dev/build/proto builds proto files under the _proto directory.
dev/build/proto:
	protoc -I _proto _proto/*.proto --go_out=plugins=grpc:internal/addon/plugin/proto

## ARG is command arguments.
ARG=lint _example/proto

## run/cmd/protolint runs protolint with ARG
run/cmd/protolint:
	go run cmd/protolint/main.go $(ARG)

## run/cmd/protolint/exampleconfig runs protolint with ARG under _example/config
run/cmd/protolint/exampleconfig:
	cd _example/config && go run ../../cmd/protolint/main.go $(ARG)

## build/cmd/protolint builds protolint
build/cmd/protolint:
	go build \
		-ldflags "-X github.com/maramkhaledn/protolint/internal/cmd.version=`git describe --tags --abbrev=0` -X github.com/maramkhaledn/protolint/internal/cmd.revision=`git rev-parse --short HEAD`" \
		-o protolint \
		cmd/protolint/main.go

## build/example/plugin builds a plugin
build/example/plugin:
	go build -o plugin_example _example/plugin/main.go

## build/cmd/protoc-gen-protolint builds protoc-gen-protolint
build/cmd/protoc-gen-protolint:
	go build \
		-ldflags "-X github.com/maramkhaledn/protolint/internal/cmd/protocgenprotolint.version=`git describe --tags --abbrev=0` -X github.com/maramkhaledn/protolint/internal/cmd/protocgenprotolint.revision=`git rev-parse --short HEAD`" \
		-o protoc-gen-protolint \
		cmd/protoc-gen-protolint/main.go
