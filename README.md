# ChronoTube

## Instructions to run

### Prerequisites
1. Golang
2. Goose
  Install goose using
```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

Step 1: Create .env
```bash
  cp .example.env .env
```
Add your YOUTUBE_API_KEY into the .env only the initial API KEY the extra ones can be specified as a function parameter inside

Step 2: Run the database
for Linux/OS X users
```bash
make db-up
```
for windows users
```bash
docker compose up -d
```

Step 3: Run the migrations
for Linux/OS X users
```bash
make migrate-up
```
for Windows users
```powershell
cd db/schema && goose postgres "host=localhost port=<port> user=<user> password=<password> dbname=<dbname> sslmode=disable" up && cd ../..
```
Fill in the respective fields from the .env

Step 4: Go mod tidy
```bash
go mod tidy
```

Step 5:
For linux users just run 
```bash
make run
```
for windows users
```bash
go run cmd/main.go
```

You should get a running instance logging about Fetch Videos service

the videos can be seen from the /videos GET route with empty next page parameter to show first page
```curl
curl --location --request GET 'http://localhost:8080/videos' \
--header 'Content-Type: application/json' \
--data '{
    "next_page": ""
}'
```
