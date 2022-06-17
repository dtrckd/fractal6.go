# Build 

Warning, this will depend of files located in th seperated repository `schema` who contains all the graphql schemas.

    make genall

Build only the gqlgen code

    make generate

# Config file

The server need a `config.toml` config file. You can use the folowing template:

```
[server]
host = "localhost"
port = "8888"
prometheus_instrumentation = true
prometheus_credentials = my_secret
client_version = "1c555fa"
maintainer_email = "admin@email.com"
jwt_secret = "my_secret"
email_api_url = "url_api_email"
email_api_key = "url_api_key"

[db]
host = "localhost"
port_graphql = "8080"
port_grpc = "9080"
api = "graphql"
admin = "admin"
dgraph_public_key = ""
dgraph_private_key = ""

[graphql]
complexity_limit = 200 # 50
introspection = false
```

# Environement variable

export EMAIL_API_URL=https://postal/api/v1/send/message
export EMAIL_API_KEY=
export JWT_SECRET=
export DGRAPH_PUBLIC_KEY=$(cat public.pem)      # fish: set DGRAPH_PUBLIC_KEY (cat public.pem | string split0)
export DGRAPH_PRIVATE_KEY=$(cat private.pem)    # fish: set DGRAPH_PRIVATE_KEY (cat private.pem | string split0)
environement variable 


