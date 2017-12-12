package main

import (
	"gopkg.in/mgo.v2/bson"
	"strconv"
)

type DBSim struct {
	database 			[]RecruitDB
	databaseFront		[]RecruitFront
	nrOfRecruits		int
}

func (dbi *DBSim) initializeDataBase(settings DatabaseSettings) {

	recDB := RecruitDB{	"1212121212121212",
						"Jonas Sedin",
						"Hames Kisnani",
						"...",
						"...",
						"...",
						"...",
						"1212",
	}
	dbi.databaseFront = append(dbi.databaseFront, *recDB.ToRecruitFront())
	dbi.database = append(dbi.database, recDB)
}

func (dbi *DBSim) addRecruit(recruit *RecruitDB) error {

	recruit.ID = bson.ObjectIdHex(strconv.Itoa(dbi.nrOfRecruits))
	dbi.database = append(dbi.database,*recruit)
	dbi.nrOfRecruits++
	return nil
}

func (dbi *DBSim) updateRecruit(recruit *RecruitDB) error {

	recruitID := recruit.ID.String()

	for key, val := range dbi.database {
		if val.ID.String() == recruitID {
			dbi.database[key] = *recruit
		}
	}

	return nil
}

func (dbi *DBSim) removeRecruit(recruitID string) error {

	for key,val := range dbi.database {
		if val.ID.String() == recruitID {
			dbi.database = append(dbi.database[:key],dbi.database[key+1:]...)
		}
	}
	dbi.nrOfRecruits--
	return nil
}

func (dbi *DBSim) searchRecruit(field string, value string, recruiter string, authLevel string) ([]RecruitFront,error) {


	return dbi.databaseFront,nil
}

func (dbi *DBSim) getRecruit(ID string, recruiter string) ([]RecruitFront,error) {


	return dbi.databaseFront,nil
}

func (dbi *DBSim) close() {}
