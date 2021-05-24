.ONESHELL:
GOFLAGS ?= $(GOFLAGS:) -v
GOFLAGS_PROD ?= $(GOFLAGS:) -mod=vendor
GOBIN := $(PWD)/bin
RELEASE := "fractal6"
MOD := "zerogov/fractal6.go"
LANGS := en fr

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
	go build -trimpath $(GOFLAGS_PROD) -ldflags "-X $(MOD)/cmd.buildMode=PROD -X $(MOD)/web/auth.buildMode=PROD"  -o $(GOBIN)/$(RELEASE) main.go

vendor:
	go mod vendor

#
# Generate Graphql code and schema
#

gen: schema generate
genall: schema_all generate

dgraph:
	cd ../schema
	make dgraph # Do alter Dgraph
	cd -

schema:
	cd ../schema
	make schema # Do not alter Dgraph
	cd -

schema_all:
	# dgraph + schema
	cd ../schema
	make schema_all
	cd -

generate:
	go run ./scripts/gqlgen.go && \
		sed -i "s/\(func.*\)(\([^,]*\),\([^,]*\))/\1(data \2, errors\3)/" graph/schema.resolvers.go && \
		# Don't add omitempty for boolean field has it get remove if set to false!
		sed -i  '/bool /I!s/`\w* *json:"\([^`]*\)"`/`json:"\1,omitempty"`/' graph/model/models_gen.go
	# Or @DEBUG: why it doesnt work anymore ?
	#go generate ./...
	# go run github.com/99designs/gqlgen generate

#
# Generate Data
#

docs: $(LANGS)

$(LANGS):
	wildq -M -i toml -o json '.[] | {name:.name, tasks:.tasks[]|flatten }' ../doc/doc.$@.toml > data/quickdoc.$@.json_
	jq -s "." data/quickdoc.$@.json_ > data/quickdoc.$@.json
	rm -f data/quickdoc.$@.json_
