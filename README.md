# randovec

- Adds random vector data to a weaviate instance

# Prerequisite

- golang -> 1.23
- weaviate cluster

## Local 

- Copy the example env file

  `cp .env.example .env`

- Populate all the values in the env file and source it   
 
  `source .env`

- Run 
  ```bash
    go mod download
    go run cmd/main.go
  ```

