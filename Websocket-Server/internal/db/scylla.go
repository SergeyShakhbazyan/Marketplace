package db

import (
	"github.com/gocql/gocql"
	"log"
)

func Connection() *gocql.Session {
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "messenger_keyspace"
	userSession, err := cluster.CreateSession()
	if err != nil {
		log.Fatalln(err)
	}

	return userSession
}
