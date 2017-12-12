package main

import (
	"gopkg.in/mgo.v2/bson"
	"fmt"
)

type RecruitFront struct {
	ID 			string		`json:"id"`
	Name 		string		`json:"name"`
	Tipper 		string		`json:"tipper"`
	Notes 		string		`json:"notes"`
	Reminder 	string		`json:"reminder"`
	NotIntTech 	string		`json:"notinttech"`
	WhyNotInt 	string		`json:"whynotint"`
}

func (rec RecruitFront) PrintRec(method string) {
	fmt.Println(method,"Name:",rec.Name,"ID:",rec.ID)
}

type RecruitDB struct {
	ID 			bson.ObjectId		`json:"_id"         bson:"_id,omitempty"`
	Recruiter 	string				`json:"recruiter"`
	Name 		string				`json:"name"`
	Tipper 		string				`json:"tipper"`
	Notes 		string				`json:"notes"`
	Reminder 	string				`json:"reminder"`
	NotIntTech 	string				`json:"notinttech"`
	WhyNotInt 	string				`json:"whynotint"`
}

func (rec RecruitDB) PrintRec(method string) {
	fmt.Println(method,"Name:",rec.Name,"ID:",rec.ID)
}

func (rec *RecruitFront) ToRecruitDB(recruiter string) *RecruitDB{
	var recDB = RecruitDB{}

	recDB.ID = bson.ObjectIdHex(rec.ID)
	recDB.Recruiter = recruiter
	recDB.Name = rec.Name
	recDB.Tipper = rec.Tipper
	recDB.Notes = rec.Notes
	recDB.Reminder = rec.Reminder
	recDB.NotIntTech = rec.NotIntTech
	recDB.WhyNotInt = rec.WhyNotInt

	return &recDB
}

func (recDB *RecruitDB) ToRecruitFront() *RecruitFront{
	var rec = RecruitFront{}

	rec.ID 			= recDB.ID.Hex()
	rec.Name 		= recDB.Name
	rec.Tipper 		= recDB.Tipper
	rec.Notes 		= recDB.Notes
	rec.Reminder 	= recDB.Reminder
	rec.NotIntTech 	= recDB.NotIntTech
	rec.WhyNotInt 	= recDB.WhyNotInt

	return &rec
}

type Credentials struct {
	Cid 	string 		`json:"cid"`
	Csecret string 		`json:"csecret"`
}

type DatabaseSettings struct {
	DatabaseURL 		string
	DatabaseName 		string
	CollectionName 		string
	DatabaseURLSec 		string
	DatabaseNameSec 	string
	CollectionNameSec 	string
}

type SimpleToken struct {
	Name 				string
	SubID 				string
	Authorizationlevel	string
}