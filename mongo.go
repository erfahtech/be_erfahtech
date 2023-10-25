package beurse

import (
	"context"
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

func GetDevicesByUserId(conn *mongo.Database, collectionname string, email string) ([]Device, error) {
    var devices []Device
    collection := conn.Collection(collectionname)

    // Menggunakan filter untuk mencari data perangkat yang sesuai dengan ID pengguna
    filter := bson.M{"user": email}

    cursor, err := collection.Find(context.TODO(), filter)
    if err != nil {
        return devices, err
    }

    defer cursor.Close(context.TODO())

    for cursor.Next(context.TODO()) {
        var device Device
        if err := cursor.Decode(&device); err != nil {
            return devices, err
        }
        devices = append(devices, device)
    }

    return devices, nil
}