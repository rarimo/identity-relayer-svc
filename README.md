# identity-relayer-svc

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

The Relayer service as a part of Rarimo cross-chain system designed to finalize transferring flow by submitting final
transaction to the target chain.

The goal of identity relayer is to observe and fetch information about new signatures for
`IDENTITY_STATE_TRANSFER` and `IDENTITY_GIST_TRANSFER` operations and after submit the state transit transactions to
configured EVM chain by request.

For more information about how the PolygonID identity transfer works
visit: [rarimo-core docs](https://rarimo.github.io/rarimo-core/docs/common/bridging/002-identity.html).

----

## Build

You can use the image from GitHub registry or build the executable by yourself.

Build command:

```
go build .
```

Also, you can use the Dockerfile inside the repository.

----

## Configuration

Service configuration consists of two parts:

### Environment

```shell
export KV_VIPER_FILE=/config.yaml
```

### Config file

```yaml
log:
  disable_sentry: true
  level: debug

# The port to run on
listener:
  addr: :8000

# PostgreSQL DB connect
db:
  url: "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"

# Rarimo core RPCs
core:
  addr: tcp://validator:26657

cosmos:
  addr: validator:9090

evm:
  chains:
    - name: "Ethereum"
      ## Address of the modified state contract on target chain
      contract_address: ""
      ## Private key HEX (without leading 0x) that will pay the tx fee
      submitter_private_key: ""
      ##  RPC address. Example https://mainnet.infura.io/v3/11111
      rpc:
      ## Target chain id
      chain_id: 1

relay:
  # Flag the indicates should service iterate over all existing transfer operation and fill the database
  catchup_disabled: true
```

----

## Run

Use the `relayer-svc migrate up && relayer-svc run all` command to perform database migrations and run the service.

Explore the simple docker-compose file to run described services:

```yaml
version: "3.7"

services:
  relayer-db:
    image: postgres:13
    restart: unless-stopped
    environment:
      - POSTGRES_USER=relayer
      - POSTGRES_PASSWORD=relayer
      - POSTGRES_DB=relayer
      - PGDATA=/pgdata
    volumes:
      - relayer-data:/pgdata

  relayer:
    image: path/to/image:hash
    restart: on-failure
    ports:
      - "8000:8000"
    depends_on:
      - relayer-db
    volumes:
      - ./config/relayer.yaml:/config.yaml
    environment:
      - KV_VIPER_FILE=/config.yaml
    entrypoint: sh -c "relayer-svc migrate up && relayer-svc run all"


volumes:
  relayer-data:
```

----

## Using service

1. Execute the POST `/integrations/relayer/state/relay` request with the following body to perform state publishing:

```json
{
  "chain": "The name of chain submit to according to the service configuration",
  "hash": "The state hash (the same as on state contract) in 0x... hex format",
  "waitConfirm": true
}
```

2. Execute the POST `/integrations/relayer/gist/relay` request with the following body to perform gist publishing:

```json
{
  "chain": "The name of chain submit to according to the service configuration",
  "hash": "The GIST hash (the same as on state contract) in 0x... hex format",
  "waitConfirm": true
}
```

`"waitConfirm": true` - indicates if request should wait until transaction will be included into the block.
Default - `false`.

The response will be:

* Code 200, successful relay, tx hash in body.
* Code 404, state is not transferred yet, wait a little and repeat request.
* Code 400, state has be relayed before.


