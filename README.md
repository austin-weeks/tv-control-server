[![Tests](https://github.com/austin-weeks/tv-control-server/actions/workflows/Tests.yml/badge.svg)](https://github.com/austin-weeks/tv-control-server/actions/workflows/Tests.yml)
[![CD](https://github.com/austin-weeks/tv-control-server/actions/workflows/CD.yml/badge.svg)](https://github.com/austin-weeks/tv-control-server/actions/workflows/CD.yml)

# 📺 Samsung TV Control Server
A server that provides endpoints for controlling your Samsung TV without needing to reach for your remote. Currently, endpoints for increasing and decreasing brightness are provided.

## Installation
Run the install script to download the latest version for your OS.
```bash
curl -fsSL https://github.com/austin-weeks/tv-control-server/install.sh | sh
```

To install an older pre-built binary, see [Releases](/releases).

Otherwise, build from source with the [Go toolchain](https://go.dev/).

## Usage

### Configuration

App settings must be provided in a `config.json` file.
- `tv_ip` - ***required*** - IP address of your TV
- `app_name` - Name of the app as displayed on your TV - *default: Gopher Remote*
- `app_port` - Port to run the server - *default: 1234*
- `token_file` - Storage location of generated tv authentication token - *default: .tv_token*
- `tv_port` - Port of your TV (8002 for Samsung TV's) - *default: 8002*
- `client_password` - Optional password required to authenticate with app - *default: none*
- `brightness_location` - Location of brightness settings in your TV's quick menu - *default: 3*
- `initial_delay_ms` - Delay in ms to wait after connecting to TV to send commands. Useful for TVs with a laggy UI - *default: 2000*
```json
{
    "tv_ip": "127.0.0.1",
    "app_name": "My Remote",
    "app_port": "5173",
    "token_file": "token.txt",
    "tv_port": "8002",
    "client_password": "password123",
    "brightness_location": 1,
    "initial_delay_ms": 4000
}
```

### Endpoints

`GET /increase-brightness` - Increase the brightness of your TV

`GET /decrease-brightness` - Decrease the brightness of your TV

#### Headers
- `"Authorization: <password>"` - client password if set in `config.json`
- `"Adjustment: <value>"` - number by which to increase/decrease TV brightness
