# urlshortner

urlshortner is a simple REST API built with Golang that shortens URLs. It is utilizes a PostgreSQL database, follows reository pattern and uses unit test to ensure its quality. This project was created as a learning exercise for building REST APIs.

## Getting Started

### Docker
To run the application with Docker, use the following commands:

```bash
docker-compose up -d
```
Ensure you have ```migrate``` installed to run migrations. If not installed, you can download and install it as follows:

```
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.14.1/migrate.linux-amd64.tar.gz | tar xvz

mv migrate.linux-amd64 $GOPATH/bin/migrate
```
Set the database environment variable:
```
DSN=postgres://postgres:mysecretpassword@localhost/urlshortner?sslmode=disable
```

creating the database tables
```
migrate -path=./migrations -database=$DSN up
```

### Kubernetes
For Kubernetes deployment, navigate to the deployment directory and apply the YAML files:
```
kubectl apply -f .
```
## Api endpoints
The API provides the following endpoints:

* GET ```/v1/{url_key}```: Redirects to the actual website.
 
* POST ```/register```: Registers a new user.

* POST ```/login```: Logs in a user.

protected endpoints(authentication required)

* GET ```/v1/user/urls/{url_key}```: Get stats for a specific URL.

* DELETE ```/v1/user/urls/{url_key}```: Delete a created URL.

* GET ```/v1/user/urls/```: List all created URLs.

* POST ```/v1/user/urls/```: Create a new short URL.
