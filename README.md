# urlshortner
Simple url shortner service for learning purposes.
## Run

Run with Docker
```bash
docker-compose up
```
Run migrations

needs to have `migrate` installed
```
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.14.1/migrate.linux-amd64.tar.gz | tar xvz

mv migrate.linux-amd64 $GOPATH/bin/migrate
```
set database env variable
```
DSN=postgres://postgres:mysecretpassword@localhost/urlshortner?sslmode=disable
```
creating the tables
```
migrate -path=./migrations -database=$DSN up
```
### TODO
- [ ] Add tests
