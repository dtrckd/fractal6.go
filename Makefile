.ONESHELL:
GOFLAGS ?= $(GOFLAGS:) -v
GOFLAGS_PROD ?= $(GOFLAGS:) -mod=vendor
GOBIN := $(PWD)/bin
RELEASE := "fractal6"
MOD := "zerogov/fractal6.go"

# TODO: versioning
# LDFLAGS see versioning, hash etc...

.PHONY: build prod vendor
default: build gen

#
# Build commands
#

run:
	go run main.go run

build:
	go build $(GOFLAGS) -o $(GOBIN)/$(RELEASE) main.go

prod:
	go build -trimpath $(GOFLAGS_PROD) -ldflags "-X $(MOD)/cmd.buildMode=PROD"  -o $(GOBIN)/$(RELEASE) main.go

vendor:
	go mod vendor

#
# Generate Graphql code and schema
#

generate: schema gen
genall: schema_all gen

schema:
	cd ../schema
	make schema # Do not alter Dgraph
	cd -

schema_all:
	cd ../schema
	make all # Do alter Dgraph
	cd -

gen: _gen _named_returns_resolver _add_omitempty

_gen:
	go run ./scripts/gqlgen.go
	# Or @DEBUG: why it doesnt work anymore ?
	#go generate ./...
	# go run github.com/99designs/gqlgen generate

_named_returns_resolver:
	sed -i "s/\(func.*\)(\([^,]*\),\([^,]*\))/\1(data \2, errors\3)/" graph/schema.resolvers.go

_add_omitempty:
	# Don't add omitempty for boolean field has it get remove if set to false!
	sed -i  '/bool /I!s/`\w* *json:"\([^`]*\)"`/`json:"\1,omitempty"`/' graph/model/models_gen.go

