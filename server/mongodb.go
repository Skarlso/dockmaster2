package main

import (
	"log"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//MongoDBConnection Encapsulates a connection to a database
type MongoDBConnection struct {
	session *mgo.Session
}

//Save will save a contaier using mongodb as a storage medium
func (mdb MongoDBConnection) Save(a Agent) error {
	mdb.session = mdb.GetSession()
	defer mdb.session.Close()
	db := mdb.session.DB("dockmaster").C("containers")
	db.Remove(bson.M{"agentid": a.AgentID})

	now := time.Now()
	//This is necessary to convert time.Now() to UTC which is CET by default
	date := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), now.Nanosecond(), time.UTC)
	a.CreatedAt = date
	err := db.Insert(a)
	if err != nil {
		return err
	}
	return nil
}

//removeOldData removes old data from the database on save
func (mdb MongoDBConnection) startCleansing() {
	log.Println("Looking for old Documents to Cleanse...")
	mdb.session = mdb.GetSession()
	defer mdb.session.Close()
	db := mdb.session.DB("dockmaster").C("containers")
	agents := []Agent{}
	iter := db.Find(nil).Iter()
	iter.All(&agents)
	now := time.Now()
	//time.Now stays in CET even after time.Now().UTC(). Which means this is needed to force it to UTC.
	date := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), now.Nanosecond(), time.UTC)
	for _, a := range agents {
		compareNow := a.CreatedAt.Add(time.Second * time.Duration(a.ExpireAfterSeconds)).UTC()
		if compareNow.Unix() < date.Unix() {
			log.Println("Cleansing Old Document:", a)
			db.Remove(bson.M{"agentid": a.AgentID})
		}
	}
}

//Load will load the contaier using mongodb as a storage medium
func (mdb MongoDBConnection) Load() (a []Agent, err error) {
	mdb.session = mdb.GetSession()
	defer mdb.session.Close()
	c := mdb.session.DB("dockmaster").C("containers")

	iter := c.Find(nil).Iter()
	err = iter.All(&a)

	return a, err
}

//Delete bulk deletes containers
func (mdb MongoDBConnection) Delete(a Agent) error {
	mdb.session = mdb.GetSession()
	defer mdb.session.Close()
	db := mdb.session.DB("dockmaster").C("containers")
	err := db.Remove(bson.M{"agentid": a.AgentID})
	if err != nil {
		return err
	}
	return nil
}

//GetSession return a new session if there is no previous one
func (mdb *MongoDBConnection) GetSession() *mgo.Session {
	if mdb.session != nil {
		return mdb.session.Copy()
	}
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	return session
}
