# Fonto

To be implemented

## Build & run

1) Create `.env` file and put there your database credentials and other secrets
```dotenv
POSTGRES_USER=username
POSTGRES_PASSWORD=password
ACCESS_TOKEN_SECRET=access_secret
REFRESH_TOKEN_SECRET=refresh_secret
```

2) Run your PostgreSQL database and put database address to [app.yml](./config/app.yml)
3) Run `go build -o fonto cmd/app/main.go && ./fonto`