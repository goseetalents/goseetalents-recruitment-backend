package main

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type AccessStruct struct {
	regAccount []Account
	DataBaseName string
	CollectionName string

	Session *mgo.Session
}

type Account struct {
	Name string			`json:"name"`
	SubID string		`json:"subid"`
	AccessLevel string	`json:"accesslevel"`
}

func (acc *AccessStruct) initialize(session *mgo.Session, settings DatabaseSettings) {
	acc.DataBaseName 	= settings.DatabaseNameSec
	acc.CollectionName 	= settings.CollectionNameSec

	var err error
	acc.Session = session
	if err != nil {
		fmt.Println("Error:",err)
	}
	acc.Session.SetMode(mgo.Monotonic,true)

	var accounts []Account
	collect := acc.Session.DB(acc.DataBaseName).C(acc.CollectionName)
	collect.Find(nil).All(&accounts)
	acc.regAccount = append(acc.regAccount,accounts...)
}

func (acc *AccessStruct) checkAddAccess(name string, subID string) {
	for _, account := range acc.regAccount {
		if account.Name == name && account.SubID == "" {
			account.SubID = subID
			collect := acc.Session.DB(acc.DataBaseName).C(acc.CollectionName)
			fmt.Println("Added account:", account.Name)
			fmt.Println(collect.Update(bson.M{"name" : name}, account))
		}
	}
}

func (acc *AccessStruct) checkAccess(name string, subID string) string {

	for _, account := range acc.regAccount {
		if account.Name == name && account.SubID == subID {
			return account.AccessLevel
		}
	}
	return "unauthorized"
}
