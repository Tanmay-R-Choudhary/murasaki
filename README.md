# Murasaki

A Redis-inspired in-memory data store in Go. Named after the Japanese word for purple.

Built by following CodeCrafters' ["Build your own Redis"](https://github.com/codecrafters-io/build-your-own-redis/tree/main/stage_descriptions) course. All code is written by me with minimal AI assistance as a way to learn.

**Status:** in progress, on Base stage 6. Single node only.

## Goal

A single-node server first, then cluster support, taking inspiration from how Redis works rather than copying it.

## How it works

```
client ──MWP over TCP──► goroutine per connection (max 50) ──► deserialize
                                                                    │
                                              global mutex ──► execute ──► in-memory map
```

- **Commands:** `PING`, `ECHO`, `GET`, `SET`
- **Execution:** each connection is read in its own goroutine, but every command runs under one global mutex. Commands are serialized, as in Redis's single-threaded loop.
- **Store:** one global `map` of keys to typed values.

### Murasaki Wire Protocol (MWP)

MWP reuses RESP's framing. The only difference is that the client resolves types. RESP clients send everything as bulk strings.

`SET score 42`:

| Protocol | Bytes |
|---|---|
| RESP | `*3\r\n$3\r\nSET\r\n$5\r\nscore\r\n$2\r\n42\r\n` |
| MWP | `*3\r\n$3\r\nSET\r\n$5\r\nscore\r\n:42\r\n` |

There is no error type yet. Errors come back as strings starting with `ERR`.

## Run it

Requires Go (see `go.mod`).

```bash
go run ./server    # listens on localhost:8080
go run ./client    # reads commands from stdin, e.g.  SET name "murasaki"
```

## Roadmap

Follows the CodeCrafters order: finish Base, then lists, transactions, streams, persistence, pub/sub, sorted sets, geospatial, replication, AOF, bitmaps, auth. Cluster support comes after the single-node work; the coordination design is undecided.
