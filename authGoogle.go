package main

import (
	"io/ioutil"
	"fmt"
	"os"
	"golang.org/x/oauth2"
	"encoding/json"
	"net/http"
	"github.com/gorilla/sessions"
	"context"
	"github.com/satori/go.uuid"
	"encoding/gob"
	"golang.org/x/oauth2/google"
	"gopkg.in/mgo.v2"
)

type AuthGoogle struct {
	Oauth2 *oauth2.Config
	OauthState string
	SessionStore sessions.Store
	Access *AccessStruct
}

func (auth *AuthGoogle) initRegister() {
	// Gob encoding for gorilla/sessions
	gob.Register(&oauth2.Token{})
	gob.Register(&UserGoogle{})
}

type UserGoogle struct {
	Sub string 				`json:"sub"`
	Name string 			`json:"name"`
	GivenName string 		`json:"given_name"`
	FamilyName string 		`json:"family_name"`
	Email string 			`json:"email"`
	EmailVerified bool 		`json:"email_verified"`
}

func (auth *AuthGoogle) initialize(mgoSession *mgo.Session, settings DatabaseSettings) {

	var c Credentials
	file, err := ioutil.ReadFile("./cred/googlecredentials.json")
	if err != nil {
		fmt.Printf("File error: %v\n", err)
		os.Exit(1)
	}
	json.Unmarshal(file, &c)

	auth.Oauth2 = &oauth2.Config{
		RedirectURL:    hostNameURL + "callback",
		ClientID:     	c.Cid,
		ClientSecret: 	c.Csecret,
		Scopes:       	[]string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     	google.Endpoint,
	}

	cookieStore := sessions.NewCookieStore([]byte("something-very-secret"))
	cookieStore.Options = &sessions.Options{
		HttpOnly: true,
	}
	auth.SessionStore = cookieStore

	auth.initRegister()

	var access AccessStruct
	auth.Access = &access
	auth.Access.initialize(mgoSession, settings)
}

func (auth *AuthGoogle) loginHandler(w http.ResponseWriter, r *http.Request) {
	sessionID := uuid.NewV4().String()
	fmt.Println("sessionID for login",sessionID)

	oathFlowSession, err := auth.SessionStore.New(r,sessionID)

	if err != nil {
		fmt.Println("Error creating flowSession:",err)
	}
	oathFlowSession.Options.MaxAge = 1*60

	err = oathFlowSession.Save(r,w)
	if err != nil {
		fmt.Println("Error saving session1:",err)
	}

	url := auth.Oauth2.AuthCodeURL(sessionID,oauth2.ApprovalForce,oauth2.AccessTypeOnline)
	http.Redirect(w, r, url, http.StatusFound)
}


func (auth *AuthGoogle) callbackHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("state for auth",r.FormValue("state"))


	code := r.FormValue("code")
	tok, err := auth.Oauth2.Exchange(context.Background(),code)
	if err != nil {
		fmt.Println("Error exchanging:",err)
	}

	session, err := auth.SessionStore.New(r,"default")
	if err != nil {
		fmt.Println("err:",err)
	}

	client := auth.Oauth2.Client(oauth2.NoContext, tok)
	email, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	defer email.Body.Close()

	var user UserGoogle
	decoder := json.NewDecoder(email.Body)
	err = decoder.Decode(&user)

	if err != nil {
		fmt.Println("error parsing json:",err)
	}

	session.Values["user-info"] = user
	session.Values["oauth_token"] = tok
	if err := session.Save(r,w); err != nil {
		fmt.Println("Error saving session2:",err)
	}
	auth.Access.checkAddAccess(user.Name, user.Sub)

	fmt.Println("User",user.Name,"(ID:", user.Sub,")logged in...")
	http.Redirect(w, r, hostNameURL + "/applicants/",http.StatusFound)
}

func (auth *AuthGoogle) userFromSession(r *http.Request) SimpleToken {
	session, err := auth.SessionStore.Get(r,"default")
	if err != nil {
		return SimpleToken{}
	}

	tok, ok := session.Values["oauth_token"].(*oauth2.Token)
	if !ok || !tok.Valid() {
		return SimpleToken{}
	}
	user, ok := session.Values["user-info"].(*UserGoogle)
	if !ok {
		return SimpleToken{}
	}

	simpleToken := SimpleToken{
		user.Name,
		user.Sub,
		auth.Access.checkAccess(user.Name, user.Sub),
	}

	return simpleToken
}
