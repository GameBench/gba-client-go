[GBA getting started](https://docs.gamebench.net/automation-interface-usage/http-api/#getting-started)

Please note this will only work with GBA version v1.5.0 or greater.

### Create a client

```go
import gba "github.com/GameBench/gba-client-go"

func main() {
	config := &gba.Config{BaseUrl: "http://localhost:8000", Username: "ade@gamebench.net", Password: ""}
	client := gba.New(config)
}
```

Alternatively, use env vars for configuration

```
GBA_BASE_URL=
GBA_USERNAME=
GBA_PASSWORD=
```

### List devices

```go
client.listDevices()
```
