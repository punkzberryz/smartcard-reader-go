# smartcard-reader-go

Go application read personal and nhso data from thai id card, inspired by [go-thai-smartcard](https://github.com/somprasongd/go-thai-smartcard)

## How to build

- Required version [Go](https://go.dev/dl/) version 1.18+
- Clone this repository
- Download all depencies with `go mod download`

> Linux install `sudo apt install build-essential libpcsclite-dev pcscd`

- Build with `go build -o bin/smartcard-reader-go ./main.go`

  > Windows `go build -o bin/smartcard-reader-go.exe ./main.go`

## How to run

Run from binary file that builded from the previous step.

### Configurations

|        ENV         | Default |                    Description                    |
| :----------------: | :-----: | :-----------------------------------------------: |
| **SMC_AGENT_PORT** | "9898"  |                    Server port                    |
| **SMC_SHOW_IMAGE** | "true"  | Enable or disable read face image from smartcard. |
| **SMC_SHOW_NHSO**  | "flase" | Enable or disable read nsho data from smartcard.  |
| **SMC_SHOW_LASER** | "flase" |  Enable or disable read laser id from smartcard.  |
|    **API_KEY**     | "1234"  |  APIKEY, used during query for smart-card result  |

### Client connect via REST API

To be updated..
