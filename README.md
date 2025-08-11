# Book Management App

There are 2 different branches:
1. `main` (master branch)
2. `dev` (development/staging branch)
### Backend:

1. Set up Environment Variables: `cd server && cp .env.example`
2. Edit `.env` with your postgresql credentials, for example:

```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_NAME=bookdb
DB_SSLMODE=disable
```

3. Install dependencies: `go mod download`
4. Generate Swagger docs: `swag init -g cmd/api/main.go --output docs`
5. Run the Server: `go run cmd/api/main.go`

### Client:

1. Install the modules: `cd client && npm install`
2. Run the app: `npm run dev`
