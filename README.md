48hr.dev
===

48hr.dev written in Go and built with Tailwind.

more coming soon...

# Development

## Run locally

1. Get your `key.json` service account details from firebase console and store it in this directory.
2. Get your firebase config from the console and cp `.env.example` to a `.env` file. Then, fill out that info correctly
3. `npx tailwindcss -o public/css/style.css --watch`
4. `go run main.go`

## Notes
* `src/` contains the CSS files. Do not edit `public/css/*`

### Generate SSL info for testing locally (https is required)

```sh
openssl genrsa -out server-key.pem 1024
openssl req -new -key server-key.pem -out server-csr.pem
openssl x509 -req -in server-csr.pem -signkey server-key.pem -out server-cert.pem
```
