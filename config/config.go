package config

import (
    "encoding/json"
    "io/ioutil"
    "github.com/dgrijalva/jwt-go"
)

type (
    Configuration struct {
        BaseUrl string `json:"baseurl"`
        Site string `json:"site"`
        WlanId string `json:"wlanid"`
        LoginUrl string
        WlanconfUrl string
    }

    JwtCustomClaims struct {
        Username string `json:"username"`
        Password string `json:"password"`
        jwt.StandardClaims
    }

)


func genUrls(c *Configuration) {
    c.LoginUrl = c.BaseUrl + "/api/login"
    c.WlanconfUrl = c.BaseUrl + "/api/s/" + c.Site + "/rest/wlanconf/" + c.WlanId
}

func LoadConfig () (*Configuration, error) {
    bytes, err := ioutil.ReadFile("./config.json")
    if err != nil {
        return &Configuration{}, err
    }

    var c Configuration
    err = json.Unmarshal(bytes, &c)
    if err != nil {
        return &Configuration{}, err
    }

    genUrls(&c)

    return &c, err
}
