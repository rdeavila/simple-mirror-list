# Simple Mirror List

Environment variables

```bash
export DC01_MIRROR="http://10.10.1.1/almalinux"
export DC01_PREFIX="10.10"
export DC02_MIRROR="http://192.168.1.1/almalinux"
export DC02_PREFIX="192.168"
export DEFAULT_MIRROR="https://mirrors.rda.run/almalinux"
```

Run on dev environment

```bash
go run main.go
```

Compile new binary

```bash
CGO_ENABLED=0 go build -trimpath -buildvcs=false -ldflags "-s -w" -o bin/sml
```

Create podman image

```bash
podman build -t sml .
```

Run container

```bash
podman run -it --rm -p 8080:8080 localhost/sml:latest
```
