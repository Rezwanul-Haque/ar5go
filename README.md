# Clean Architecture with JWT & Cache based authentication in golang

## Run Seeder
```terminal
go mod vendor

go run main.go seed

or

go run main.go seed --truncate=true or -t=true
```

## Local
```terminal
go mod vendor

go run main.go serve
```

## Docker
```terminal
make development
```