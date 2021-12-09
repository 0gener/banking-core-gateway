package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	jwtgo "github.com/auth0/go-jwt-middleware/validate/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

var ErrAuth0DomainMandatory = errors.New("auth0 domain is mandatory")
var ErrAuth0AudienceMandatory = errors.New("auth0 audience is mandatory")
var ErrUnexpectedIssuer = errors.New("token claims validation failed: unexpected issuer")
var ErrUnexpectedAudience = errors.New("token claims validation failed: unexpected audience")

const signatureAlgorithm = "RS256"

var _ jwtgo.CustomClaims = &CustomClaims{}

type CustomClaims struct {
	Scope string `json:"scope"`
	jwt.StandardClaims
}

func (c CustomClaims) Validate(_ context.Context) error {
	expectedIssuer := "https://" + "dev-oclvmm6.eu.auth0.com" + "/"
	if c.Issuer != expectedIssuer {
		return ErrUnexpectedIssuer
	}

	expectedAudience := "https://banking-core.com/"
	if c.Audience != expectedAudience {
		return ErrUnexpectedAudience
	}

	return nil
}

func (c CustomClaims) HasScope(expectedScope string) bool {
	result := strings.Split(c.Scope, " ")
	for i := range result {
		if result[i] == expectedScope {
			return true
		}
	}

	return false
}

type JwtMiddleware interface {
	EnsureValidToken() gin.HandlerFunc
}

type JwtMiddlewareOptions struct {
	Domain   string
	Audience string
}

type auth0JwtMiddleware struct {
	options *JwtMiddlewareOptions
}

func NewJwtMiddleware(options *JwtMiddlewareOptions) (*auth0JwtMiddleware, error) {
	if options.Domain == "" {
		return nil, ErrAuth0DomainMandatory
	}

	if options.Audience == "" {
		return nil, ErrAuth0AudienceMandatory
	}

	return &auth0JwtMiddleware{
		options,
	}, nil
}

func (j *auth0JwtMiddleware) EnsureValidToken() gin.HandlerFunc {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		certificate, err := getPEMCertificate(token)
		if err != nil {
			return token, err
		}

		return jwt.ParseRSAPublicKeyFromPEM([]byte(certificate))
	}

	customClaims := func() jwtgo.CustomClaims {
		return &CustomClaims{}
	}

	validator, err := jwtgo.New(
		keyFunc,
		signatureAlgorithm,
		jwtgo.WithCustomClaims(customClaims),
	)
	if err != nil {
		log.Fatalf("Failed to set up the jwt validator")
	}

	m := jwtmiddleware.New(validator.ValidateToken)

	return func(ctx *gin.Context) {
		var encounteredError = true
		var handler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
			encounteredError = false
			ctx.Request = r
			ctx.Next()
		}

		m.CheckJWT(handler).ServeHTTP(ctx.Writer, ctx.Request)

		if encounteredError {
			ctx.AbortWithStatusJSON(
				http.StatusUnauthorized,
				map[string]string{"message": "Failed to validate JWT."},
			)
		}
	}
}

type (
	jwks struct {
		Keys []jsonWebKeys `json:"keys"`
	}

	jsonWebKeys struct {
		Kty string   `json:"kty"`
		Kid string   `json:"kid"`
		Use string   `json:"use"`
		N   string   `json:"n"`
		E   string   `json:"e"`
		X5c []string `json:"x5c"`
	}
)

func getPEMCertificate(token *jwt.Token) (string, error) {
	response, err := http.Get("https://" + "dev-oclvmm6.eu.auth0.com" + "/.well-known/jwks.json")
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	var jwks jwks
	if err = json.NewDecoder(response.Body).Decode(&jwks); err != nil {
		return "", err
	}

	var cert string
	for _, key := range jwks.Keys {
		if token.Header["kid"] == key.Kid {
			cert = "-----BEGIN CERTIFICATE-----\n" + key.X5c[0] + "\n-----END CERTIFICATE-----"
			break
		}
	}

	if cert == "" {
		return cert, errors.New("unable to find appropriate key")
	}

	return cert, nil
}
