# gRPC-Web + HAProxy example

This project demonstrates how to run a gRPC-Web app behind HAProxy, with both gRPC and REST endpoints.

## Prerequisites

- Node.js v22.x
- Go 1.25.x
- [Homebrew](https://brew.sh/) (for macOS users)

## Quick-start

### 1. Install protobuf & gRPC tools

```sh
brew install protobuf protoc-gen-grpc-web
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

### 2. Install Go dependencies and generate PB files

```sh
make tidy
```

### 4. Build the JS app

```sh
cd web
npm install
npm run build
```

### 5. Serve the frontend

```sh
cd dist
python3 -m http.server 8082
```

### 6. Build the API server and start it

```sh
# From the project root
make build

./bin/grpcwebexample
```

It serves on port 8080.

### 7. Generate a self-signed certificate

```sh
DOMAIN_NAME=acme.example ./create-self-signed-cert.sh
```

This generates a self-signed certificate and key in the `private/` directory. HAProxy will use these files for TLS termination.
**macOS users:**
To trust the certificate, drag `private/acme.example.crt` into the "System" keychain in Keychain Access. Double-click the certificate and set "Always Trust" to avoid browser warnings.

Additionally, add the following entry to your `/etc/hosts` file to map `acme.example` to localhost:

```
127.0.0.1 acme.example
```

### 8. Run HAProxy

```sh
haproxy -f haproxy.cfg
```

## API

You can test the endpoints as follows:

### REST

```sh
curl 'https://acme.example/hello?name=Tester'
```

### gRPC

```sh
grpcurl -d '{"name":"Tester"}' acme.example:443 hellopb.Greeter/SayHello
```

Both commands should return:

```json
{"message":"Hello, Tester!"}
```

### gRPC-Web

- **Web UI:**
  Open [https://acme.example/static/index.html](https://acme.example/static/index.html) in your browser.
