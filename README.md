microblog
===

# Features
    * TODO

# Development

## Run locally

1. Get your `key.json` service account details from firebase console and store it in this directory.
2. Get your firebase config from the console and cp `.env.example` to a `.env` file. Then, fill out that info correctly
3. `npx tailwindcss -o public/css/style.css --watch`
4. `go run main.go`

## Notes
* `src/` contains the CSS files. Do not edit `public/css/*`
* `public/js/*` can/should be edited. We can clean this up later.
