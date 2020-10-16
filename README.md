# go-keycloak-demo

Demo app using keycloak to authenticate users in Go web apps.

## How to run

Run `docker-compose up --build`.

 1. Create a Client with the name `go-demo` [here](http://localhost:8080/auth/admin/master/console/#/create/client/master)
 1. Create a user [here](http://localhost:8080/auth/admin/master/console/#/create/user/master)
 1. Run via `curl --data "{\"username\": \"user\",\"password\": \"pass\"}" http://localhost:8000` or via Postman
 