package main

import (
	"context"
	"encoding/json"
	"github.com/Nerzal/gocloak/v7"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

type LoginUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type NewUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

var (
	port         = ":8000"
	keycloak     = "http://localhost:8080"
	clientId     = "go-demo"
	clientSecret = ""
	realm        = "master"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/login", login)
	e.POST("/register", register)

	e.Logger.Fatal(e.Start(port))
}

func login(c echo.Context) error {
	var user LoginUser
	if err := json.NewDecoder(c.Request().Body).Decode(&user); err != nil {
		http.Error(c.Response(), "Forbidden", http.StatusForbidden)
		return err
	}

	client := gocloak.NewClient(keycloak)
	token, err := client.Login(context.TODO(), clientId, clientSecret, realm, user.Username, user.Password)
	if err != nil {
		http.Error(c.Response(), "Forbidden", http.StatusForbidden)
		return err
	}

	if err := json.NewEncoder(c.Response()).Encode(token); err != nil {
		http.Error(c.Response(), "InternalServerError", http.StatusInternalServerError)
		return err
	}

	return c.NoContent(http.StatusOK)
}

func register(c echo.Context) error {
	client := gocloak.NewClient(keycloak)

	token, err := client.LoginAdmin(context.TODO(), "admin", "admin", realm)
	if err != nil {
		http.Error(c.Response(), "InternalServerError", http.StatusInternalServerError)
		return err
	}

	var user NewUser
	if err := json.NewDecoder(c.Request().Body).Decode(&user); err != nil {
		http.Error(c.Response(), "Forbidden", http.StatusForbidden)
		return err
	}

	enabled := true
	newUser := gocloak.User{
		Email:    &user.Email,
		Username: &user.Username,
		Enabled:  &enabled,
	}

	userId, err := client.CreateUser(context.TODO(), token.AccessToken, realm, newUser)
	if err != nil {
		http.Error(c.Response(), "Could not create user", http.StatusForbidden)
		return err
	}

	err = client.SetPassword(context.TODO(), token.AccessToken, userId, realm, user.Password, false)
	if err != nil {
		http.Error(c.Response(), "Could not set password", http.StatusForbidden)
		return err
	}

	return c.NoContent(http.StatusCreated)
}
