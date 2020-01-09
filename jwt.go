package ginx

import (
	"encoding/json"
	"net/http"
	"time"

	lib "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/gin-gonic/gin"
)

type jwtManager struct {
	secret    string
	exp       time.Duration
	alg       string
	header    string
	filterFun gin.HandlerFunc
}

type JWTAuthConfig struct {
	Secret    string
	Exp       time.Duration
	Alg       string
	Header    string
	FilterFun gin.HandlerFunc
}

func newJwtAuthDriver(cfg JWTAuthConfig) *jwtManager {
	if cfg.FilterFun == nil {
		cfg.FilterFun = defaultFilterRes
	}
	return &jwtManager{
		secret:    cfg.Secret,
		exp:       cfg.Exp,
		alg:       cfg.Alg,
		header:    cfg.Header,
		filterFun: cfg.FilterFun,
	}
}

const (
	contextJWTUserTokenKey = "auth_user_jwt_token"
	contextJWTUser         = "auth_user_jwt"
)

// Check the token of request header is valid or not.
func (app *jwtManager) Check(c *gin.Context) bool {
	token := c.Request.Header.Get(app.header)
	if token == "" {
		return false
	}
	authJwtToken, err := app.user(c)

	if err != nil {
		return false
	}

	c.Set(contextJWTUserTokenKey, authJwtToken)

	return authJwtToken.Valid
}

func (app *jwtManager) user(c *gin.Context) (*lib.Token, error) {
	var keyFun = func(token *lib.Token) (interface{}, error) {
		b := []byte(app.secret)
		return b, nil
	}

	return request.ParseFromRequest(c.Request, &request.MultiExtractor{
		&request.PostExtractionFilter{
			Extractor: request.HeaderExtractor{app.header},
			Filter: func(s string) (string, error) {
				return s, nil
			},
		},
		request.ArgumentExtractor{app.header},
	}, keyFun)
}

// User is get the auth user from token string of the request header which
// contains the user ID. The token string must start with "Bearer "
func (app *jwtManager) User(c *gin.Context, userPointer interface{}) {

	var (
		jwtToken     *lib.Token
		exist        bool
		valInterface interface{}
	)

	if valInterface, exist = c.Get(contextJWTUser); exist {
		userPointer = valInterface
		return
	}

	if valInterface, exist = c.Get(contextJWTUserTokenKey); !exist {
		var err error
		jwtToken, err = app.user(c)
		if err != nil {
			panic("invalid token")
		}
	} else {
		jwtToken = valInterface.(*lib.Token)
	}

	if claims, ok := jwtToken.Claims.(lib.MapClaims); ok && jwtToken.Valid {
		if err := json.Unmarshal([]byte(claims["user"].(string)), userPointer); err != nil {
			panic("invalid token")
		}
		c.Set(contextJWTUser, userPointer)
	} else {
		panic("invalid token")
	}
}

var timeZone, _ = time.LoadLocation("Asia/Shanghai")

func (app *jwtManager) Login(http *http.Request, w http.ResponseWriter, user interface{}) interface{} {

	token := lib.New(lib.GetSigningMethod(app.alg))
	// Set some claims
	userStr, err := json.Marshal(user)
	if err != nil {
		return nil
	}
	token.Claims = lib.MapClaims{
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

func (app *jwtManager) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !app.Check(c) {
			app.filterFun(c)
			c.Abort()
		}
		c.Next()
	}
}
