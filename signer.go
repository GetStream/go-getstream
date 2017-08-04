package getstream

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"net/http"
	"strings"

	httpsig "gopkg.in/LeisureLink/httpsig.v1"
	"gopkg.in/dgrijalva/jwt-go.v3"
)

// Credits to https://github.com/hyperworks/go-getstream for the urlSafe and generateToken methods

// Signer is responsible for generating Tokens
type Signer struct {
	Key    string
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

// generateToken will use the secret of the signer and the message passed as an argument to generate a Token
func (s Signer) GenerateToken(message string) string {
	hash := sha1.New()
	hash.Write([]byte(s.Secret))
	key := hash.Sum(nil)
	mac := hmac.New(sha1.New, key)
	mac.Write([]byte(message))
	digest := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	return s.UrlSafe(digest)
}

func (s Signer) SignHTTP(request *http.Request) error {
	signer, err := httpsig.NewRequestSigner(s.Key, s.Secret, "hmac-sha256")
	if err != nil {
		return err
	}
	return signer.SignRequest(request, []string{}, nil)
}

func (s Signer) SignJWT(request *http.Request, context ScopeContext, action ScopeAction, feedID string) error {
	token, err := s.GenerateJWT(context, action, feedID)
	if err != nil {
		return err
	}
	request.Header.Set("Stream-Auth-Type", "jwt")
	request.Header.Set("Authorization", token)
	return nil
}

// GenerateFeedScopeToken returns a jwt
func (s Signer) GenerateJWT(context ScopeContext, action ScopeAction, feedID string) (string, error) {
	if feedID == "" {
		feedID = "*"
	}

	claims := jwt.MapClaims{
		"resource": context.Value(),
		"action":   action.Value(),
		"feed_id":  feedID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
