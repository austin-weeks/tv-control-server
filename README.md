![CI](https://github.com/austin-weeks/tv-control-server/actions/workflows/CI.yml/badge.svg)
![CD](https://github.com/austin-weeks/tv-control-server/actions/workflows/CD.yml/badge.svg)

# 📺 Samsung TV Control Server
A server that provides endpoints for controlling your Samsung TV without needing to reach for your remote. Currently, endpoints for increasing and decreasing brightness are provided.

## Installation
For pre-built binaries, see [Releases](/releases).

Otherwise, build from source with the [Go toolchain](https://go.dev/).

## Usage
You must set an enviroment variable for your TV's IP address and optionally, a client password. Variables will be read from a `.env` file if present.

```shell
TV_IP="127.0.0.1" # Your Samsung TV's IP Address
CLIENT_PW="password" # an optional password to authenticate with your server
```

### Endpoints

`POST /increase-brightness` - Increase the brightness of your TV
`POST /decrease-brightness` - Decrease the brightness of your TV

Both endpoints require the following **headers**:
- `"Authorization: <password>"` - client password as set in `.env`
- `"Adjustment: <value>"` - number by which to increase/decrease TV brightness

The server will listen on port `1234`. HTTPS is currently **not** supported.
