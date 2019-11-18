[GBA getting started](https://docs.gamebench.net/automation-interface-usage/http-api/#getting-started)

Please note this will only work with GBA version v1.5.0 or greater.

### Create a client

```go
import "github.com/GameBench/gba-client-go"

func main() {
	config := &gba.Config{BaseUrl: "http://localhost:8000"}
	client := gba.New(config)
}
```

Alternatively, use env vars for configuration

```
GBA_BASE_URL=
```

### List devices

```go
client.listDevices()
```
