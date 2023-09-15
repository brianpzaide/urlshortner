# urlshortner
This is a simple url shortner REST api, built with Golang. It is backed by postgres database, follows reository pattern and uses unit test to ensure its quality. Built this REST api as a learning project.
## Run

### Docker
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

### Kubernetes
```
cd deployment

kubectl apply -f .
```
