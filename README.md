# MVP Boilerplate

## Deployment

Step 1: Go to xx.xx.xx.xx server using remote protocol (RDP) 

Step 2: Go to boilerplate_backend Folder  `G:\boilerplate\boilerplate_backend`

Step 3: To deploy this project run

```
./run.sh
```

>Note : To Seed the database for the first time. Modify the docker compose file to make the serve command to seed and run `./run.sh`


## Run Locally

Go to the project directory

```bash
  cd boilerplate_backend
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


