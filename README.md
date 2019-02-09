[![Build Status](https://travis-ci.com/sheirys/wsreply.svg?branch=master)](https://travis-ci.com/sheirys/wsreply)
[![Go Report Card](https://goreportcard.com/badge/github.com/sheirys/wsreply)](https://goreportcard.com/report/github.com/sheirys/wsreply)
[![GoDoc](https://godoc.org/github.com/sheirys/wsreply?status.svg)](https://godoc.org/github.com/sheirys/wsreply)
[![codecov](https://codecov.io/gh/sheirys/wsreply/branch/master/graph/badge.svg)](https://codecov.io/gh/sheirys/wsreply)

# wsreply

wsreply is single exchange pub/sub websocket server.

## Recommended tools

Examples use [wscat](https://www.npmjs.com/package/wscat) websocket client.

## Server

Server entrypoint can be found on `cmd/server/main.go`. Points of interest about server;
* `ws://localhost/pub` endpoint for publishers.
* `ws://localhost/sub` endpoint for subscribers.
* Endpoins expect or produces messages in form:
```
{
    "op": <int>,         // op code
    "payload": <string>,    // message
}
```

Start server by using `go run cmd/server/main.go`. Explore more options with `-h` flag. All messages produced by publishers will be transformed by rules defined in `translate.go` file before broadcasted to subscribers. Whenever new subscriber joins the server, all publishers will be notified by op command `OpHasSubscribers`. If no subscribers left on server, all publishers will be notified by op command `OpNoSubscribers` When publisher joins server it should ask if there is any subscribers in server by using `OpSyncSubscribers` command.

*Possible op codes and op commands*

| Op code | Op command | Comment |
| --- | --- | --- |
| `0` | `OpNoSubscribers` | When no subscibers left on server all publishers will be notified by this command. |
| `1` | `OpHasSubscribers` | When new subscriber joins the server all publishers will be notified by this command. |
| `2` | `OpSyncSubscribers` | Publisher can send this command to ask server if there is any subscribers on server. Server will respond with `OpNoSubscribers` or `OpHasSubscribers` |
| `3` | `OpMessage` | Publisher should use this command if wants to broadcast message to subscribers |

## Publisher

When publisher joins server it should ask if there is any connected subscribers in the server by using `OpSyncSubscribers`. Example with `wscat`:

    $ wscat -c ws://localhost:8886/pub
    connected (press CTRL+C to quit)
    > {"op":2}
    < {"op":0,"payload":""}

Here we connected to server and published `OpSyncSubscribers` command (op code `2`). Server responded with `OpNoSubscribers` command (op code `0`). To broadcast message to subscribers you should use `OpMessage` command (op code `3`) and specify message text in `payload`. Example with `wscat`:

    $ wscat -c ws://localhost:8886/pub
    connected (press CTRL+C to quit)
    > {"op":3, "payload":"hello?!"}
    
Be noticed, that server will transform all `!` characters to `?` in payload before sending message to subscriber.

Automated publisher example can be found in `cmd/publisher/main.go`. This publisher will connect to server, asks if any subscribers are connected to server and if there is any subscribers will produce messages every 1s untis there is no subscribers left.

## Subscriber

Exampe subscribe with `wscat`:

    $ wscat -c ws://localhost:8886/sub
    connected (press CTRL+C to quit)
    < {"op":3,"payload":"hello??"}

Here we received `OpMessage` command (op code `3`) with message `hello??` that has been sent by publisher. Subscriber example can be found in `cmd/subscriber/main.go`.

## Known issues

* publisher and subsriber clients defined in `cmd/` does not disconnect when server is closed. However `wscat` does.