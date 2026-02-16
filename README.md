# gophermc

`gophermc` is a Go library for building Minecraft clients and bots.

It focuses on practical client operations:
- Server status ping and latency checks.
- Join/login flow across many protocol versions.
- Chat and player action packets.
- Event-driven loops for long-running bots.

## Features

- Multi-version protocol support (legacy to modern versions).
- High-level client API with configurable options.
- Event stream for chat, keep-alive, ready, and disconnect handling.
- Configuration-state handling for modern Minecraft protocol flows.

## Installation

```bash
go get github.com/obeliskdev/gophermc
```

## Quick Start: Server Status + Ping

```go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/obeliskdev/gophermc"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := gophermc.NewClient(
		gophermc.WithAddr("127.0.0.1:25565"),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	statusJSON, latency, err := client.GetStatus(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(statusJSON)
	fmt.Println("latency:", latency)
}
```

## Quick Start: Join and Send Chat

```go
package main

import (
	"context"
	"log"
	"time"

	"github.com/obeliskdev/gophermc"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	client, err := gophermc.NewClient(
		gophermc.WithAddr("127.0.0.1:25565"),
		gophermc.WithUsername("GopherBot"),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Destroy()

	if err := client.Join(ctx); err != nil {
		log.Fatal(err)
	}

	if err := client.Chat("hello from gophermc"); err != nil {
		log.Fatal(err)
	}
}
```

## Event-Driven Bot Pattern

```go
events, err := client.JoinAndListen(ctx, 128)
if err != nil {
	log.Fatal(err)
}

defer client.Destroy()

for {
	select {
	case ev, ok := <-events:
		if !ok {
			return
		}
		switch e := ev.(type) {
		case gophermc.ReadyEvent:
			log.Println("ready as", e.Username)
		case gophermc.ChatMessageEvent:
			log.Printf("<%s> %s", e.Sender, e.Message)
		case gophermc.KeepAliveEvent:
			// keepalive responses are handled automatically
			_ = e
		case gophermc.DisconnectEvent:
			log.Println("disconnect:", e.Reason)
			return
		}
	case <-ctx.Done():
		return
	}
}
```

## Client Options

Common options:
- `WithAddr("host:port")`
- `WithTCPAddr(*net.TCPAddr)`
- `WithUsername("name")`
- `WithUUID(uuid.UUID)`
- `WithVersion(protocol.Version)`
- `WithServerHostname("virtual-host")`
- `WithBrand("brand")`
- `WithPrivateKey(*rsa.PrivateKey)`
- `WithConn(conn, version)`

## Core Methods

- `GetStatus(ctx)`
- `Ping()`
- `Join(ctx)`
- `JoinAndListen(ctx, eventBuffer)`
- `Chat(message)`
- `SetPosition(...)`
- `Destroy()` for graceful shutdown

## Testing

```bash
go test ./...
```

Some tests are integration-focused and require a reachable Minecraft server.

## License

MIT. See `LICENSE`.

## Credits

Packet/version mappings are generated from `minecraft-data`.
