# Build Instructions

## Building with Custom Login URL

The login URL can be customized at build time using Go's `-ldflags`:

```bash
# Build with custom login URL
go build -ldflags "-X github.com/jawahars16/jebi/internal/core.LoginURL=https://my-custom-jebi.com/auth/login" -o jebi

# Example for production build
go build -ldflags "-X github.com/jawahars16/jebi/internal/core.LoginURL=https://app.jebi.dev/auth/login" -o jebi

# Example for development
go build -ldflags "-X github.com/jawahars16/jebi/internal/core.LoginURL=http://localhost:3000/auth/login" -o jebi
```

## Default Configuration

If no custom URL is provided, the default login URL is: `https://app.jebi.dev/auth/login`

## Authentication Flow

1. CLI starts a local HTTP server on a random port (e.g., `http://127.0.0.1:8080`)
2. Opens browser to: `{LoginURL}?source=cli&callback_port={port}`
3. User authenticates in the browser
4. Web application sends authentication result to: `http://127.0.0.1:{port}/callback`
5. CLI receives the tokens and saves them securely

## Web Application Integration

The web application should:

1. Detect `source=cli` parameter
2. Extract `callback_port` parameter
3. After successful authentication, POST the auth result to `http://127.0.0.1:{callback_port}/callback`

Example POST body:

```json
{
  "username": "user@example.com",
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "optional_refresh_token",
  "expires_in": 3600
}
```
