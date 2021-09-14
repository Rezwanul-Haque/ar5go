# Clean Architecture without Cache based authentication in golang

## Run Locally

Go to the project directory

```bash
  cd clean_backend
```

Install dependencies

```bash
  go mod vendor
```
### Seedding database

```bash
  go run main.go seed

  ```
 Use the below command to truncate then seed database  
  ```
  go run main.go seed --truncate=true

  or

  go run main.go seed -t=true
```

Start the server Locally

```bash
  go run main.go serve
```

Start the server using Docker

```bash
  make development
```


