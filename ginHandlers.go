package main

import (
	gin "github.com/gin-gonic/gin"
	"fmt"
	"net/http"
	"encoding/json"
	//"github.com/pquerna/ffjson/ffjson"
)

type ginHandlers struct {
	token SimpleToken

}

func (g *ginHandlers) searchHandler(c *gin.Context) {

	recruiter := g.token.Name
	authorizationlevel := g.token.Authorizationlevel
	fmt.Println("getAllHandler: recruiter", recruiter, "accessed database")

	var recruits []RecruitFront
	if authorizationlevel == "admin" {
		recruits, _ = dbi.searchRecruit("name", "", recruiter, authorizationlevel)
	} else {
		recruits, _ = dbi.searchRecruit("recruiter", recruiter, recruiter, authorizationlevel)
	}

	c.JSON(http.StatusOK, recruits)
}

func (g *ginHandlers) getAllHandler(c *gin.Context) {

	recruiter := g.token.Name
	authlevel := g.token.Authorizationlevel
	fmt.Println("getAllHandler: recruiter", recruiter, "accessed database")

	var recruits []RecruitFront
	if authlevel == "admin" {
		recruits, _ = dbi.searchRecruit("name", "", recruiter, authlevel)
	} else {
		recruits, _ = dbi.searchRecruit("recruiter", recruiter, recruiter, authlevel)
	}

	c.JSON(http.StatusOK, recruits)

}

func (g *ginHandlers) postHandler(c *gin.Context) {
	recruiter := "Jonas Sedin"

	var recruit RecruitFront
	err := c.BindJSON(&recruit)
	if err != nil {
		fmt.Println("postHandler: error decoding", err)
	}

	recrDB := recruit.ToRecruitDB(recruiter)
	dbi.addRecruit(recrDB)
}

func (g *ginHandlers) getHandler(c *gin.Context) {
	recruiter := g.token.Name

	id := c.Param("id")

	var recruits []RecruitFront

	recruits, err := dbi.getRecruit(id, recruiter)
	if err != nil {
		fmt.Println("getHandler: Error getting recruit",err)
	}

	if err != nil {
		fmt.Println("err:", err)
	}
	c.JSON(http.StatusOK, recruits)
}

func (g *ginHandlers) putHandler(c *gin.Context) {

	recruiter := g.token.Name

	var recruit RecruitFront
	decoder := json.NewDecoder(c.Request.Body)
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
	c.String(httpStatus, "")
}

func (g *ginHandlers) deleteHandler(c *gin.Context) {

	id := c.Param("id")

	err := dbi.removeRecruit(id)
	if err != nil {
		fmt.Println("deleteHandler: Error deleting", err)
		c.String(http.StatusNotFound,"")
		return
	}

	c.String(http.StatusOK,"")
}