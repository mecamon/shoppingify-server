# This is a project to fulfill the backend for the Shoppingify challenge.

## Description


## External libraries/packages
- [golang-migrate/migrate](https://github.com/golang-migrate/migrate/blob/master/GETTING_STARTED.md)
- [google-uuid](https://github.com/google/uuid)
- [nick-snyder i18n](https://github.com/nicksnyder/go-i18n)
- [chi router](https://go-chi.io/#/)
- [golang-jwt](https://github.com/golang-jwt/jwt)
- [golang-crypto-bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt)
- [pgx](https://github.com/jackc/pgx)
- [Cloudinary](https://cloudinary.com/documentation/go_integration)
- [Swag/swaggo](https://github.com/swaggo/swag)
- [http-swagger](https://github.com/swaggo/http-swagger)

## Other development and testing tools in use
- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)

## Instructions to set a fast local dev environment
1. Make sure you have docker and docker compose in your computer.
2. From the root of the project run: docker-compose -f docker-compose.yml up --build
3. The previous step will run a postgres DB on a container, the migrations and the rest API.
4. IMPORTANT---If you stop the project and you are receiving errors when trying to run it again use: "docker system prune" to delete the cache.

## Instructions to run tests
1. To run units tests just go to the root of the project and run: go test ./... (go must be installed in your computer)
2. To run integration tests, from the root of the project run: docker-compose -f docker-compose-test.yml up --build --always-recreate-deps
3. IMPORTANT---If you stop the project and you are receiving errors when trying to run it again use: "docker system prune" to delete the cache.
