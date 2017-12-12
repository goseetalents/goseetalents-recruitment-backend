package main

import (
	"gopkg.in/mgo.v2"
	"fmt"
	"gopkg.in/mgo.v2/bson"
)

type DBInterface struct {
	DatabaseURL string
	DatabaseName string
	CollectionName string
	Session *mgo.Session
}

func (dbInterface *DBInterface) initializeDataBase(settings DatabaseSettings) {

	dbInterface.DatabaseURL 	= settings.DatabaseURL
	dbInterface.DatabaseName 	= settings.DatabaseName
	dbInterface.CollectionName 	= settings.CollectionName
	dbInterface.Session = &mgo.Session{}

	var err error
	dbInterface.Session, err = mgo.Dial(dbInterface.DatabaseURL )
	if err != nil {
		fmt.Println("initializedatabase: Error connecting to database,",err)
	}
	dbInterface.Session.SetMode(mgo.Monotonic,true)

	fmt.Println("initializedatabase: Database initialized...")
}

func (dbInterface *DBInterface) addRecruit(recruit *RecruitDB) error{

	// The struct cannot have an ID
	toDB := struct {
		Recruiter string		`json:"recruiter"`
		Name string				`json:"name"`
		Tipper string			`json:"tipper"`
		Notes string			`json:"notes"`
		Reminder string			`json:"reminder"`
		NotIntTech string		`json:"notinttech"`
		WhyNotInt string		`json:"whynotint"`
	}{
		recruit.Recruiter,
		recruit.Name,
		recruit.Tipper,
		recruit.Notes,
		recruit.Reminder,
		recruit.NotIntTech,
		recruit.WhyNotInt,
	}

	recruit.PrintRec("addRecruit: adding recruit")
	collect := dbInterface.Session.DB(dbInterface.DatabaseName).C(dbInterface.CollectionName)
	err := collect.Insert(toDB)
	if err != nil {
		fmt.Println("addRecruit: error adding", err)
		return err
	}

	return nil
}

func (dbInterface *DBInterface) updateRecruit(recruit *RecruitDB) error {

	recruit.PrintRec("updateRecruit:")

	collect := dbInterface.Session.DB(dbInterface.DatabaseName).C(dbInterface.CollectionName)
	err := collect.Update(bson.M{"_id" : recruit.ID}, recruit)
	if err != nil {
		return err
	}
	return nil
}

func (dbInterface *DBInterface) removeRecruit(recruitID string) error {

	fmt.Println("removeRecruit:",recruitID)
	collect := dbInterface.Session.DB(dbInterface.DatabaseName).C(dbInterface.CollectionName)
	err := collect.Remove(bson.M{"_id" : bson.ObjectIdHex(recruitID)})
	if err != nil {
		return err
	}

	return nil
}

func (dbInterface *DBInterface) searchRecruit(field string, value string, recruiter string, authLevel string) ([]RecruitFront,error) {

	var recruitDBs []RecruitDB
	var recruitFronts []RecruitFront

	collect := dbInterface.Session.DB(dbInterface.DatabaseName).C(dbInterface.CollectionName)
	if authLevel == "admin" {
		if field == "name" {
			//THIS IS A REGEX TEXT-SEARCH. THE INDEX HAS ALREADY BEEN ESTABLISHED IN DATABASE
			err := collect.Find(bson.M{"name": &bson.RegEx{Pattern: value, Options: "i"}}).Sort("name").All(&recruitDBs)
			if err != nil {
				fmt.Println("searchRecruit: error searching",err)
			}
		} else {
			collect.Find(bson.M{field: value}).Sort("name").All(&recruitDBs)
		}
	}else {
		if field == "name" {
			//THIS IS A REGEX TEXT-SEARCH. THE INDEX HAS ALREADY BEEN ESTABLISHED IN DATABASE
			err := collect.Find(bson.M{"name": &bson.RegEx{Pattern: value, Options: "i"}, "recruiter" : recruiter}).Sort("name").All(&recruitDBs)
			if err != nil {
				fmt.Println("searchRecruit: error searching" ,err)
			}
		} else {
			collect.Find(bson.M{field: value, "recruiter" : recruiter}).Sort("name").All(&recruitDBs)
		}
	}

	for _,recruit := range recruitDBs {
		rec := recruit.ToRecruitFront()
		recruitFronts = append(recruitFronts,*rec)
		rec.PrintRec("searchRecruit:")
	}

	return recruitFronts,nil
}

func (dbInterface *DBInterface) getRecruit(ID string, recruiter string) ([]RecruitFront,error) {

	var recruitDB RecruitDB
	var recruitFront []RecruitFront
	collect := dbInterface.Session.DB(dbInterface.DatabaseName).C(dbInterface.CollectionName)

	collect.FindId(bson.ObjectIdHex(ID)).One(&recruitDB)

	recruitFront = append(recruitFront, *recruitDB.ToRecruitFront())
	recruitFront[0].PrintRec("getRecruit:")

	return recruitFront,nil
}

func (dbInterface *DBInterface) close() {
	dbInterface.Session.Close()
}
