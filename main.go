package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Nerzal/gocloak/v7"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/", getToken)

	e.Logger.Fatal(e.Start(":8000"))
}

func getToken(c echo.Context) error {
	var user User
	if err := json.NewDecoder(c.Request().Body).Decode(&user); err != nil{
		http.Error(c.Response(), "Forbidden", http.StatusForbidden)
		return err
	}

	client := gocloak.NewClient("http://localhost:8080")
	token, err := client.Login(context.TODO(), "go-demo", "", "master", user.Username, user.Password)
	if err != nil {
		http.Error(c.Response(), "Forbidden", http.StatusForbidden)
		return err
	}

	info, err := client.GetUserInfo(context.TODO(), token.AccessToken, "master")
	if err != nil {
		http.Error(c.Response(), "Forbidden", http.StatusForbidden)
		return err
	}

	fmt.Println("user_id: ", *info.Sub)
	if err := json.NewEncoder(c.Response()).Encode(token); err != nil{
		http.Error(c.Response(), "Forbidden", http.StatusForbidden)
		return err
	}

	return nil
}
