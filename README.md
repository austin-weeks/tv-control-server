# 📺 Samsung TV Control Server
A server that provides endpoints for controlling your Samsung TV without needing to reach for your remote. Currently, endpoints for increasing and decreasing brightness are provided.

## Installation
For pre-built binaries, see [Releases](/releases).

Otherwise, build from source with the [Go toolchain](https://go.dev/).

## Usage
You must create a `.env` file in the executable's directory with the following variables:
```shell
TV_IP="127.0.0.1" # Your Samsung TV's IP Address
CLIENT_PW="password" # an optional password to authenticate with your server
```

### Endpoints

`GET /increase-brightness` - Increase the brightness of your TV
`GET /decrease-brightness` - Decrease the brightness of your TV

Both endpoints require the following **headers**:
- `"Authorization: <password>"` - client password as set in `.env`
- `"Adjustment: <value>"` - number by which to increase/decrease TV brightness

The server will listen on port `:1234`. HTTPS is **not** supported.

### Apple Shortcuts
One way to use the app is by sending a request to your server via a shortcut. See below for an example.

<picture>
    <img src="/images/shortcut.png" alt="Using the app via Apple shortcuts.">
</picture>