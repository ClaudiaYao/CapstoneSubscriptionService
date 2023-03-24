package auth

import (
	"fmt"
	"log"
	"testing"
	"time"
)

func TestCreateJWT(t *testing.T) {

	jwtMaker, err := NewJWTMaker("tokenData")
	if err != nil {
		log.Panic(err)
	}

	userID := "user6"
	token, err := jwtMaker.CreateJWTToken(userID, time.Duration(time.Second)*1800)

	if err != nil {
		log.Panic(err)
	}
	// Auth token is printed here to be used for local testing of APIs.
	fmt.Println("=====================================")
	fmt.Printf("Authorization Token for testing: %s\n", token)
	fmt.Println("=====================================")
}

func TestVerifyJWT(t *testing.T) {

	jwtVerifier, err := NewJWTVerifier("tokenData")
	if err != nil {
		log.Panic(err)
	}
	isValid, err := jwtVerifier.IsValidToken("eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InVzZXI2IiwiaXNzdWVkX2F0IjoiMjAyMy0wMy0yNFQxMToyMjoyOC44NjA0NzUrMDg6MDAiLCJleHBpcmVkX2F0IjoiMjAyMy0wMy0yNFQxMToyNzoyOC44NjA0NzUrMDg6MDAifQ.RM7Ho4BsGPEzuCt-XFkEH2cJqrJWpfNkRQml6MXt5l4f0O2tHJc9Ef5rMrPKyDz4-7p7OXId4XewwUrF-gzdlQ")

	if err != nil {
		log.Panic(err)
	}
	// Auth token is printed here to be used for local testing of APIs.
	if isValid {
		fmt.Println("it is valid")
	} else {
		fmt.Println("invalid token")
	}
}
