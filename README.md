# AntDB

AntDB is a lightweight Redis-like in-memory key-value database written in Go.

## Features

- RESP (Redis Serialization Protocol) support
- Redis 1.X compatible commands (WIP)
- TTL handling with background eviction
- Transaction support (`MULTI`/`EXEC`)

## Quickstart

Run the server locally:

```bash
git clone https://github.com/ghosind/antdb.git
cd antdb
go run cmd/main.go
```

Connect with the standard `redis-cli`:

```bash
redis-cli PING
# PONG
redis-cli SET mykey value
# OK
redis-cli GET mykey
# "value"
```

## Contributing

Contributions are welcome. Please open issues for bugs or feature requests and submit PRs for changes. Keep changes small and focused; add tests for new behaviors when possible.

## License

This project is licensed under the MIT License. See the `LICENSE` file for details.
