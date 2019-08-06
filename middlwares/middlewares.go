package middlewares

import (
    "../config"
    "net/http"
    "bytes"
    "encoding/json"
    "github.com/dgrijalva/jwt-go"
    "github.com/labstack/echo"
    "net/http/cookiejar"
)

type (
    LoginPayload struct {
       Username string `json:"username"`
       Password string `json:"password"`

    }

    Middlewares struct {
       Config *config.Configuration
    }
)

func New(config *config.Configuration) Middlewares {
    return Middlewares{config}
}


func (m *Middlewares) TokenAuth(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        token := c.Get("user").(*jwt.Token)
        claims := token.Claims.(*config.JwtCustomClaims)

        cookieJar, _ := cookiejar.New(nil)
        client := &http.Client{
            Jar: cookieJar,
        }

        loginPayload := &LoginPayload{
            claims.Username,
            claims.Password,
        }

        buf := new(bytes.Buffer)
        json.NewEncoder(buf).Encode(loginPayload)

        resp, err := client.Post(m.Config.LoginUrl, "application/json", buf)

        if err != nil {
            return err
        }

        defer resp.Body.Close()

        if resp.StatusCode == 200 {
            c.Set("httpclient", client)
            return next(c)
        }

        return echo.ErrUnauthorized
    }
}
