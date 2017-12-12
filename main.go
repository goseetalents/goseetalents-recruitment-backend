package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/gin-gonic/gin"
	"time"
	//"cloud.google.com/go/storage"
	graceful "gopkg.in/tylerb/graceful.v1"

)


// -------- Settings -----------
var authStruct AuthSimple
const httpProtocol = "http"
var hostName = ":3000"
var hostNameURL = httpProtocol + "://" + hostName + "/"
var dbi = &DBInterface{}
var dbSettings = DatabaseSettings{
	"localhost",
	"LabDB",
	"Recruits",
	"LabDB",
	"LabDB",
	"accounts",
}
// -----------------------------


func startHTTPServer() {
	router := mux.NewRouter()
	router.HandleFunc("/login", 			authStruct.loginHandler)
	router.HandleFunc("/callback", 		authStruct.callbackHandler)

	router.HandleFunc("/applicants", 		accessHandler(searchHandler)).Methods(http.MethodGet)
	router.HandleFunc("/applicants/",		accessHandler(getAllHandler)).Methods(http.MethodGet)
	router.HandleFunc("/applicants/",		accessHandler(postHandler)).Methods(http.MethodPost)
	router.HandleFunc("/applicants/{id}/",	accessHandler(getHandler)).Methods(http.MethodGet)
	router.HandleFunc("/applicants/{id}/",	accessHandler(putHandler)).Methods(http.MethodPut)
	router.HandleFunc("/applicants/{id}/",	accessHandler(deleteHandler)).Methods(http.MethodDelete)

	handler := cors.AllowAll().Handler(router)
	fmt.Println("Main: Listening on 3000...")
	if httpProtocol == "http" {
		err := http.ListenAndServe(hostName, handler)
		fmt.Println(err)
	} else {
		certDir := "C:/Users/Jonas/GoglandProjects/src/mongolab/cred"
		http.ListenAndServeTLS(hostName, certDir + "/cert.pem", certDir + "/key.pem",router)
	}
}

func startGinServer() {
	gin.SetMode(gin.TestMode)
	g := ginHandlers{}
	router := gin.New()

	router.GET("/applicants", 			g.searchHandler)
	router.GET("/applicants/", 			g.getAllHandler)
	router.POST("/applicants",			g.postHandler)
	router.GET("/applicants/:id/", 		g.getHandler)
	router.PUT("/applicants/:id/", 		g.putHandler)
	router.DELETE("/applicants/:id/",	g.deleteHandler)


	srv := &graceful.Server{
		Timeout: 0,
		Server: &http.Server{
			Addr:           ":8080",
			Handler:        router,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
		},
		BeforeShutdown: func() bool {
			fmt.Println("shutdown completed")
			return true
		},
	}
	fmt.Println("listening on 8080...")
	err := srv.ListenAndServe()

	if err != nil {
		fmt.Println("error starting server:",err)
	}
}

func main() {
	dbi.initializeDataBase(dbSettings)
	defer dbi.close()

	authStruct.initialize()
	startHTTPServer()

}

type handler func(http.ResponseWriter, *http.Request, SimpleToken)

func accessHandler(funct handler) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request){

		sToken := authStruct.userFromSession(r)
		fmt.Println("accessHandler: Checking access permisson...", sToken.Name)

		if sToken.Name == "" {
			http.Redirect(w, r, hostNameURL+ "login", http.StatusFound)
			return
		}
		if sToken.Authorizationlevel != "unauthorized" {
			funct(w, r, sToken)
		} else {
			w.WriteHeader(http.StatusForbidden)
		}
	}
}

func getAllHandler(w http.ResponseWriter, r *http.Request, token SimpleToken) {

	recruiter := token.Name
	fmt.Println("getAllHandler: recruiter", recruiter, "accessed database")

	var recruits []RecruitFront
	if token.Authorizationlevel == "admin" {
		recruits, _ = dbi.searchRecruit("name", "", recruiter, token.Authorizationlevel)
	} else {
		recruits, _ = dbi.searchRecruit("recruiter", recruiter, recruiter, token.Authorizationlevel)
	}

	respBody, _ := json.MarshalIndent(recruits, "", " ")
	JSONResponse(w, respBody, http.StatusOK)
}

func getHandler(w http.ResponseWriter, r *http.Request, token SimpleToken) {

	recruiter := token.Name

	vars := mux.Vars(r)
	id, _ := vars["id"]

	var recruits []RecruitFront

	recruits, err := dbi.getRecruit(id, recruiter)
	if err != nil {
		fmt.Println("getHandler: Error getting recruit",err)
	}

	respBody, _ := json.MarshalIndent(recruits, "", " ")
	JSONResponse(w, respBody, http.StatusOK)
}

func putHandler(w http.ResponseWriter, r *http.Request, token SimpleToken) {

	recruiter := token.Name

	var recruit RecruitFront
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&recruit)
	if err != nil {
		fmt.Println("putHandler: Error decoding:", err)
	}

	recrDB := recruit.ToRecruitDB(recruiter)
	err = dbi.updateRecruit(recrDB)

	var httpStatus int
	if err != nil {
		fmt.Println("putHandler: Error updating", err)
		httpStatus = http.StatusNotFound
	} else {
		httpStatus = http.StatusCreated
	}
	w.WriteHeader(httpStatus)
}

func searchHandler(w http.ResponseWriter, r *http.Request, token SimpleToken) {

	fmt.Println("INCOMING", r.Method, "REQUEST")

	recruiter := token.Name

	queryParam := r.URL.Query()
	var respBody []byte

	// A for range loop is used to retrieve key and value
	for key,val := range queryParam {
		recruit, _ := dbi.searchRecruit(key, val[0], recruiter, token.Authorizationlevel)
		respBody, _ = json.Marshal(recruit)
		fmt.Println(len(recruit))
	}

	JSONResponse(w, respBody, http.StatusOK)
}

func postHandler(w http.ResponseWriter, r *http.Request, token SimpleToken) {

	recruiter := token.Name

	var recruit RecruitFront
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&recruit)
	if err != nil {
		fmt.Println("postHandler: error decoding", err)
	}

	recrDB := recruit.ToRecruitDB(recruiter)
	dbi.addRecruit(recrDB)

	w.WriteHeader(http.StatusCreated)
}

func deleteHandler(w http.ResponseWriter, r *http.Request, _ SimpleToken) {

	vars := mux.Vars(r)
	id, _ := vars["id"]

	err := dbi.removeRecruit(id)
	if err != nil {
		fmt.Println("deleteHandler: Error deleting", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func JSONResponse(w http.ResponseWriter, json []byte, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	w.Write(json)
}
