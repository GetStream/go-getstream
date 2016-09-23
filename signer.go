package getstream

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"strings"

	"gopkg.in/dgrijalva/jwt-go.v3"
)

// Credits to https://github.com/hyperworks/go-getstream for the urlSafe and generateToken methods

// Signer is responsible for generating Tokens
type Signer struct {
	Secret string
}

// SignFeed sets the token on a Feed
func (s Signer) SignFeed(feedID string) string {
	return s.GenerateToken(feedID)
}

func (s Signer) UrlSafe(src string) string {
	src = strings.Replace(src, "+", "-", -1)
	src = strings.Replace(src, "/", "_", -1)
	src = strings.Trim(src, "=")
	return src
}

// generateToken will use the Secret of the signer and the message passed as an argument to generate a Token
func (s Signer) GenerateToken(message string) string {
	hash := sha1.New()
	hash.Write([]byte(s.Secret))
	key := hash.Sum(nil)
	mac := hmac.New(sha1.New, key)
	mac.Write([]byte(message))
	digest := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	return s.UrlSafe(digest)
}

// GenerateFeedScopeToken returns a jwt
func (s Signer) GenerateFeedScopeToken(context ScopeContext, action ScopeAction, feedIDWithoutColon string) (string, error) {

	claims := jwt.MapClaims{
		"resource": context.Value(),
		"action":   action.Value(),
		// "aud":
		// "exp": time.Now().UTC().Add(time.Hour * 1),
		// "jti": uuid.New(),
		// "iat": time.Now(),
		// "iss":
		// "nbf": time.Now().Unix(),
		// "sub":
	}

	if feedIDWithoutColon != "" {
		claims["feed_id"] = feedIDWithoutColon
	} else {
		claims["feed_id"] = "*"
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(s.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// GenerateUserScopeToken returns a jwt
func (s Signer) GenerateUserScopeToken(context ScopeContext, action ScopeAction, userID string) (string, error) {

	claims := jwt.MapClaims{
		"resource": context.Value(),
		"action":   action.Value(),
		// "aud":
		// "exp": time.Now().UTC().Add(time.Hour * 1),
		// "jti": uuid.New(),
		// "iat": time.Now(),
		// "iss":
		// "nbf": time.Now().Unix(),
		// "sub":
	}

	if userID != "" {
		claims["user_id"] = userID
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(s.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
