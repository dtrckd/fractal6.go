.ONESHELL:
GOFLAGS ?= $(GOFLAGS:) -v
GOFLAGS_PROD ?= $(GOFLAGS:) -mod=vendor
GOBIN := $(PWD)/bin
RELEASE := "fractal6"
MOD := "fractale/fractal6.go"

# TODO: versioning
# LDFLAGS see versioning, hash etc...

.PHONY: build prod vendor schema
default: build

#
# Build commands
#

run_api:
	# DO NOT FORGET TO SET THE FOLLOWING ENV VARIABLE
	# * EMAIL_API_URL
	# * EMAIL_API_KEY
	# * JWT_SECRET
	# * DGRAPH_PUBLIC_KEY
	# * DGRAPH_PRIVATE_KEY
	go run main.go api

run_notifier:
	# DO NOT FORGET TO SET THE FOLLOWING ENV VARIABLE
	# * EMAIL_API_URL
	# * EMAIL_API_KEY
	# * DGRAPH_PUBLIC_KEY
	# * DGRAPH_PRIVATE_KEY
	go run main.go notifier

build:
	go build $(GOFLAGS) -o $(GOBIN)/$(RELEASE) main.go

prod:
	go build -trimpath $(GOFLAGS_PROD) \
		-ldflags "-X $(MOD)/cmd.buildMode=PROD -X $(MOD)/web/auth.buildMode=PROD -X $(MOD)/db.buildMode=PROD" \
		-o $(GOBIN)/$(RELEASE) main.go

vendor:
	go mod vendor

#
# Generate Graphql code and schema
#

genall: dgraph schema generate
gen: schema generate

dgraph:
	cd ../schema
	make dgraph # Do alter Dgraph
	cd -

schema:
	cd ../schema
	make schema # Do not alter Dgraph
	cd -
	mkdir -p schema/
	cp ../schema/gen/*.graphql schema/

generate:
	# go generate ./... | go run github.com/99designs/gqlgen generate
	# We add "omitempty" for each generate type's literal except for Bool and Int to prevent
	# loosing data (when literal are set to false/0 values) when marshalling.
	go run ./scripts/gqlgen.go && \
		sed -i "s/\(func.*\)(\([^,]*\),\([^,]*\))/\1(data \2, errors\3)/"   graph/schema.resolvers.go && \
		sed -i '/\W\(bool\|int\)\W/I!s/`\w* *json:"\([^`]*\)"`/`json:"\1,omitempty"`/' graph/model/models_gen.go


#
# Generate Data
#

docs:
	cp ../doc/_data/* data

# Utils

show_query:
	rg "Gqlgen" graph/schema.resolvers.go -B 2 |grep func |sed "s/^func[^)]*)\W*\([^(]*\).*/\1/" | sort

install:
	# Redis
	#curl https://packages.redis.io/gpg | sudo apt-key add -
	#echo "deb https://packages.redis.io/deb $(lsb_release -cs) main" | sudo tee /etc/apt/sources.list.d/redis.list
	# -- official way
	#curl -fsSL https://packages.redis.io/gpg | sudo gpg --dearmor -o /usr/share/keyrings/redis-archive-keyring.gpg
	#echo "deb [signed-by=/usr/share/keyrings/redis-archive-keyring.gpg] https://packages.redis.io/deb $(lsb_release -cs) main" | sudo tee /etc/apt/sources.list.d/redis.list
	#sudo apt-get update
	sudo apt-get install redis


rsa:
	# Dgraph Authorization
	#ssh-keygen -t rsa -P "" -b 2048 -m PEM -f jwtRS256.key
	#ssh-keygen -e -m PEM -f jwtRS256.key > jwtRS256.key.pub
	openssl genrsa -out private.pem 2048
	openssl rsa -in private.pem -pubout -out public.pem
	# Copy public key to the Dgraph authorization object
	# cat public.pem | sed 's/$/\\\n/' | tr -d "\n" |  xclip -selection clipboard;
