package ginx

import (
	"encoding/json"
	jwt_lib "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type jwtManager struct {
	secret string
	exp    time.Duration
	alg    string
	header string
}

func newJwtAuthDriver(secret, alg, header string, exp time.Duration) *jwtManager {
	return &jwtManager{
		secret: secret,
		exp:    exp,
		alg:    alg,
		header: header,
	}
}

// Check the token of request header is valid or not.
func (app *jwtManager) Check(c *gin.Context) bool {
	token := c.Request.Header.Get(app.header)
	if token == "" {
		return false
	}
	var keyFun jwt_lib.Keyfunc
	keyFun = func(token *jwt_lib.Token) (interface{}, error) {
		b := []byte(app.secret)
		return b, nil
	}

	authJwtToken, err := request.ParseFromRequest(c.Request, &request.MultiExtractor{
		&request.PostExtractionFilter{
			Extractor: request.HeaderExtractor{app.header},
			Filter: func(s string) (string, error) {
				return s, nil
			},
		},
		request.ArgumentExtractor{app.header},
	}, keyFun)

	if err != nil {
		//logger.Info(logger.E{
		//	Title:    "jwt auth check",
		//	Function: "jwt.Check",
		//	Error:    err,
		//})
		return false
	}

	c.Set("User", map[string]interface{}{
		"token": authJwtToken,
	})

	return authJwtToken.Valid
}

// User is get the auth user from token string of the request header which
// contains the user ID. The token string must start with "Bearer "
func (app *jwtManager) User(c *gin.Context) interface{} {

	var jwtToken *jwt_lib.Token
	if jwtUser, exist := c.Get("User"); !exist {
		tokenStr := c.Request.Header.Get(app.header)
		if tokenStr == "" {
			panic("非法token")
		}
		var err error
		jwtToken, err = jwt_lib.Parse(tokenStr, func(token *jwt_lib.Token) (interface{}, error) {
			return []byte(app.secret), nil
		})
		if err != nil {
			panic("非法token")
		}
	} else {
		jwtToken = jwtUser.(map[string]interface{})["token"].(*jwt_lib.Token)
	}

	if claims, ok := jwtToken.Claims.(jwt_lib.MapClaims); ok && jwtToken.Valid {
		var user map[string]interface{}
		if err := json.Unmarshal([]byte(claims["user"].(string)), &user); err != nil {
			panic("非法token")
		}
		c.Set("User", map[string]interface{}{
			"token": jwtToken,
			"user":  user,
		})
		return user
	} else {
		panic("非法token")
	}
}

var timeZone, _ = time.LoadLocation("Asia/Shanghai")

func (app *jwtManager) Login(http *http.Request, w http.ResponseWriter, user map[string]interface{}) interface{} {

	token := jwt_lib.New(jwt_lib.GetSigningMethod(app.alg))
	// Set some claims
	userStr, err := json.Marshal(user)
	if err != nil {
		return nil
	}
	token.Claims = jwt_lib.MapClaims{
		"user": string(userStr),
		"exp":  time.Now().In(timeZone).Add(app.exp).Unix(),
	}
	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString([]byte(app.secret))
	if err != nil {
		return nil
	}

	return tokenString
}

func (app *jwtManager) Logout(http *http.Request, w http.ResponseWriter) bool {
	return true
}
