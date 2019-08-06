package handlers

import (
    "net/http"
    "../config"
    "bytes"
    "encoding/json"
    "time"
    "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)


type (
    WlanPayload struct {
        Name string `json:"name" form:"name" query:"name"`
        Password string `json:"x_passphrase" form:"x_passphrase" query:"x_passphrase"`
    }

    LoginPayload struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }

    Handlers struct {
        Config *config.Configuration
    }
)

func New(config *config.Configuration) Handlers {
    return Handlers{config}
}

func (h *Handlers) Login(c echo.Context) error {
    u := new(LoginPayload)
    if err := c.Bind(u); err != nil {
        return echo.ErrUnauthorized
    }

    buf := new(bytes.Buffer)
    json.NewEncoder(buf).Encode(u)

    uresp, err := http.Post(h.Config.LoginUrl, "application/json", buf)
    if err != nil {
        return err
    }

    defer uresp.Body.Close()

    if uresp.StatusCode != 200 {
        return echo.ErrUnauthorized
    }

    claims := &config.JwtCustomClaims{
        u.Username,
        u.Password,
        jwt.StandardClaims{
            ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

    t, err := token.SignedString([]byte("43q4cqerasdfcd"))
    if err != nil {
        return err
    }

    return c.JSON(http.StatusOK, echo.Map{
        "token": t,
    })
}

func (h *Handlers) Wlanconf(c echo.Context) error {
    wlan := new(WlanPayload)
    if err := c.Bind(wlan); err != nil {
        return echo.ErrUnauthorized
    }

    client := c.Get("httpclient").(*http.Client)

    buf := new(bytes.Buffer)
    json.NewEncoder(buf).Encode(wlan)
    req, err := http.NewRequest("PUT", h.Config.WlanconfUrl, buf)
    req.Header.Set("Content-Type", "application/json")

    if err != nil {
        return err
    }

    resp, err := client.Do(req)
    if err != nil {
        return err
    }

    if resp.StatusCode != 200 {
        return echo.ErrUnauthorized
    }

    return c.JSON(http.StatusOK, echo.Map{
        "status": "ok",
    })
}
