FORMAT: 1A

# UP CoCo Admin

A REST API for some handy CoCo data.

## Group API

### /etcd-all

#### Returns all etcd keys on the server. [GET]

+ Request

    + Headers

            Accept: application/json

    + Body

+ Response 200 (application/json)

    Returns all etcd keys on the server as a json map.

    + Body

            {}

    + Schema

            {
              "type": "object"
            }

+ Request

    + Headers

            Accept: application/json

    + Body

+ Response 500 (application/json)

    The provided etcd url is invalid or not responding.

    + Body

            client: response is invalid json. The endpoint is probably not valid etcd cluster endpoint.

