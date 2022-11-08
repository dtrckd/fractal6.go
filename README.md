# fractal6.go

Backend, API, Business-logic layer for [Fractale](https://fractale.co).

**Fractale** is a platform for self-organization. It is designed around the concept that an organization can be represented as a tree, and the principles of transparency, governance decentralization and authority distribution. A tree divides in branches and form leaves, likewise an organization divides in **Circles** that can have **Roles**. Both circles and roles have an associated descriptive document, called **Mandate**, intended to define its purpose and operating rules. Finally, the communication inside the organization is done trough **Tensions**, and make the link between users and organizations. You can think of it as an email, but more structured and more powerful.

Using Fractale for your organization offers the following capabilities and feature:
* Tree and Graph-Packing organisation navigation.
* Organization visibility define at circles level.
* ACL based on member role and circle governance rules.
* Ticketing management through Tensions.
* Discussion thread and subscription per tension.
* Journal history of events (including mandate updates).
* Email notifications and response.
* Labels system.
* Role templates system.
* GraphQL API.


## Requirements

* Redis 4+


## Install

#### From source

**Setup**

    git clone -b prod https://github.com/fractal6/fractal6.go
    cd fractal6.go

    # Install the client UI (Optional)
    # NOTE: This will install the client build for fractale.co.
    #       To point to your own instance, you need to rebuild it (see https://github.com/fractal6/fractal6-ui.elm/)
    #       Otherwise it will point to api.fractale.co
    make install_client

    # Start Redis (KV cache store)
    sudo systemctl restart redis-server  # or "systemctl restart redis" depending on your version

    # Setup Dgraph (database)
    make bootstrap
    ./bin/dgraph zero --config contrib/dgraph/config-zero.yml
    # Open a new terminal and run
    ./bin/dgraph alpha --config contrib/dgraph/config-alpha.yml

**Configure**

The server need a `config.toml` config file to run (in the project's root folder, i.e `fractal6.go/`).
You can use the following template:

```config.toml
[server]
instance_name = "Fractale"
domain = "fractale.co"
hostname = "localhost"
port = "8888"
jwt_secret = "my_jwt_secret"
prometheus_instrumentation = true
prometheus_credentials = "my_prom_secret"
client_version = "git hash used to build the client"

[mailer]
admin_email = "admin@mydomain.com"
# URL API
email_api_url = "https://..."
email_api_key = "..."
# SMTP api
# ...TODO...
# Postal validation creds
# postal default-dkim-record: Just the p=... part of the TXT record (without the semicolon at the end)
dkim_key = "..."
# webhook redirection for Postal alert.
matrix_postal_room = "!...:matrix.org"
matrix_token = "..."

[db]
hostname = "localhost"
port_graphql = "8080"
port_grpc = "9080"
api = "graphql"
admin = "admin"
dgraph_public_key = "public.pem"
dgraph_private_key = "private.pem"

[graphql]
complexity_limit = 200 # 50
introspection = false

[admin]
max_public_orgas = -1    # Maximum public organnization per user, -1 for unlimited
max_private_orgas = -1   # Maximum private organnization per user, -1 for unlimited
max_orga_reg = -1        # Maximum organnization per regular user, -1 for unlimited
max_orga_pro = -1        # Maximum organnization per pro user, -1 for unlimited
```

Finally, generate the certificate for dgraph authorization, and populate the schema:

    # Generate certs
    make certs

	# Copy public key for the Dgraph authorization at the end of the schema
    sed -i '$ d' schema/dgraph_schema.graphql
	cat public.pem | sed 's/$/\\\n/' | tr -d "\n" | head -c -2 | { read my; echo "# Dgraph.Authorization {\"Header\":\"X-Frac6-Auth\",\"Namespace\":\"https://YOUR_DOMAIN/jwt/claims\",\"Algo\":\"RS256\",\"VerificationKey\":\"$PUBKEY\"}"; }  >> schema/dgraph_schema.graphql

    # Update Dgraph schema
    curl -X POST http://localhost:8080/admin/schema --data-binary "@schema/dgraph_schema.graphql" | jq


**Launch for production**

    # Build
    go mod vendor
    make prod

    # Open a terminal and run (main server)
    ./bin/fractal6 api
    # Open a second terminal and run (message passing that manage event notifications)
    ./bin/fractal6 notifier


**Launch for development**

	go run main.go api
	go run main.go notifier


You can add users in Fractale with the following sub-command :

    ./bin/fractal6 adduser


Note that this command would be required to add users if the mailer is not enabled as the sign-up process has an email validation step. Once the mailer is setup, new users can be invited to organizations and roles from their email, or from their username if they have already sign-up.


## Contributing

You can open issues for bugs you've found or features you think are missing. You can also submit pull requests to this repository. To get started, take a look at [CONTRIBUTING.md](CONTRIBUTING.md).

You can follow Fractale organisation and roadmap at [o/f6](https://fractale.co/o/f6) and espacially [t/f6/tech](https://fractale.co/t/f6/tech).

IRC channel: #fractal6 on matrix.org

## License

Fractale is free, open-source software licensed under AGPLv3.
