package main

import (
	"net/http"
	"fmt"
	"time"
)


type AuthSimple struct {
	Name 		string
	SubID 		string
	Accesslevel string
}

func (auth *AuthSimple) initialize() {
	auth.Name = "Jonas Sedin"
	auth.SubID = "0101010101010101"
	auth.Accesslevel = "authorized"
}

func (auth *AuthSimple) callbackHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("callbackHandler empty... ")
}

func (auth *AuthSimple) loginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("loginHandler empty... redirecting to main page")
	time.Sleep(time.Millisecond * 500)
	http.Redirect(w, r, hostNameURL + "applicants/", http.StatusFound)
}

func (auth *AuthSimple) userFromSession(r *http.Request) SimpleToken {
	simpleToken := SimpleToken{	auth.Name,
							   	auth.SubID,
							   	auth.Accesslevel}
	return simpleToken
}
