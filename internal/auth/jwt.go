package auth

import (

	"github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
)

func InitializeJWTMiddleware(appSecret string) *jwtmiddleware.JWTMiddleware{
	return jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(appSecret), nil
		},
		// When set, the middleware verifies that tokens are signed with the specific signing algorithm
		// If the signing method is not constant the ValidationKeyGetter callback can be used to implement additional checks
		// Important to avoid security issues described here: https://auth0.com/blog/2015/03/31/critical-vulnerabilities-in-json-web-token-libraries/
		SigningMethod: jwt.SigningMethodHS256,
	})
}
