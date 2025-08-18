package helper

import (
	"fmt"
	"os"
	"time"

	jwtToken "github.com/golang-jwt/jwt/v5"
)

var Secret_key = os.Getenv("Secret_key")

type UserField struct {
	Name          string
	Email         string
	Refresh_token string
	jwtToken.RegisteredClaims
}

func GetClaimfromToken(tokenval string) (claimval *UserField, msg string) {
	token, err := jwtToken.ParseWithClaims(
		tokenval,
		&UserField{},
		func(token *jwtToken.Token) (interface{}, error) {
			return []byte(Secret_key), nil
		},
	)
	if err != nil {
		fmt.Println("Error is ", err.Error())
		msg = err.Error()
		return
	}
	claim, ok := token.Claims.(*UserField)
	if !ok {
		fmt.Println("Invalid token")
		msg = "Invalid token"
		return
	}
	if claim.ExpiresAt.Time.Before(time.Now()) {
		fmt.Println("Token is expired")
		msg = "Token is expired"
		return
	}
	return claim, msg
}
func TokenGeneration(name string, email string) (access_token string, refresh_token string, errval error) {

	newclaims := &UserField{
		Name:  name,
		Email: email,
		RegisteredClaims: jwtToken.RegisteredClaims{
			ExpiresAt: jwtToken.NewNumericDate(time.Now().Local().Add(24 * time.Hour)),
			Issuer:    "ParitoshRick",
			IssuedAt:  jwtToken.NewNumericDate(time.Now()),
		},
	}
	refreshclaims := &UserField{
		RegisteredClaims: jwtToken.RegisteredClaims{
			ExpiresAt: jwtToken.NewNumericDate(time.Now().Local().Add(48 * time.Hour)),
			Issuer:    "RcikRefreshToken",
		},
	}
	tokenGenerate_access := jwtToken.NewWithClaims(jwtToken.SigningMethodES256, newclaims)
	tokenGenerate_referesh := jwtToken.NewWithClaims(jwtToken.SigningMethodES256, refreshclaims)
	access_token, errval = tokenGenerate_access.SignedString(Secret_key)
	if errval != nil {
		fmt.Println("Error while generating the access token")
	}
	refresh_token, errval = tokenGenerate_referesh.SignedString(Secret_key)
	if errval != nil {
		fmt.Println("Error while generating the refresh token")
	}
	return access_token, refresh_token, errval

}
