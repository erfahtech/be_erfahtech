package beurse

import (
	"os"

	"github.com/aiteung/atdb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetConnection(MONGOCONNSTRINGENV, dbname string) *mongo.Database {
	var DBmongoinfo = atdb.DBInfo{
		DBString: os.Getenv(MONGOCONNSTRINGENV),
		DBName:   dbname,
	}
	return atdb.MongoConnect(DBmongoinfo)
}

func IsPasswordValid(mongoconn *mongo.Database, collection string, userdata User) bool {
	filter := bson.M{"email": userdata.Email}
	res := atdb.GetOneDoc[User](mongoconn, collection, filter)
	return CheckPasswordHash(userdata.Password, res.Password)
}