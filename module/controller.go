package beurse

import (
	"context"
	"log"
	"os"
	"time"

	// "crypto/rand"
	// "encoding/hex"
	"errors"
	"fmt"

	"strings"

	"github.com/aiteung/atdb"
	"github.com/badoux/checkmail"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	// "go.mongodb.org/mongo-driver/mongo/options"
	// "golang.org/x/crypto/argon2"

	model "github.com/erfahtech/be_erfahtech/model"
)

func SetConnection(MONGOCONNSTRINGENV, dbname string) *mongo.Database {
	var DBmongoinfo = atdb.DBInfo{
		DBString: os.Getenv(MONGOCONNSTRINGENV),
		DBName:   dbname,
	}
	return atdb.MongoConnect(DBmongoinfo)
}

func Waktu(s string) time.Time {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	t, err := time.ParseInLocation(time.RFC3339, s, loc)
	if err != nil {
	  log.Fatal(err)
	}
	return t
  }

func IsPasswordValid(mongoconn *mongo.Database, collection string, userdata model.User) bool {
	filter := bson.M{"email": userdata.Email}
	res := atdb.GetOneDoc[model.User](mongoconn, collection, filter)
	return CheckPasswordHash(userdata.Password, res.Password)
}

func MongoConnect(MongoString, dbname string) *mongo.Database {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(os.Getenv(MongoString)))
	if err != nil {
		fmt.Printf("MongoConnect: %v\n", err)
	}
	return client.Database(dbname)
}

func GetAllDocs(db *mongo.Database, col string, docs interface{}) interface{} {
	collection := db.Collection(col)
	filter := bson.M{}
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		fmt.Println("Error GetAllDocs in colection", col, ":", err)
	}
	err = cursor.All(context.TODO(), &docs)
	if err != nil {
		fmt.Println(err)
	}
	return docs
}

func InsertOneDoc(db *mongo.Database, col string, doc interface{}) (insertedID primitive.ObjectID, err error) {
	result, err := db.Collection(col).InsertOne(context.Background(), doc)
	if err != nil {
		return insertedID, fmt.Errorf("kesalahan server : insert")
	}
	insertedID = result.InsertedID.(primitive.ObjectID)
	return insertedID, nil
}

func UpdateOneDoc(db *mongo.Database, col string, id primitive.ObjectID, doc interface{}) (err error) {
	filter := bson.M{"_id": id}
	result, err := db.Collection(col).UpdateOne(context.Background(), filter, bson.M{"$set": doc})
	if err != nil {
		fmt.Printf("UpdatePresensi: %v\n", err)
		return
	}
	if result.ModifiedCount == 0 {
		err = errors.New("no data has been changed with the specified id")
		return
	}
	return nil
}

func DeleteOneDoc(_id primitive.ObjectID, db *mongo.Database, col string) error {
	collection := db.Collection(col)
	filter := bson.M{"_id": _id}
	result, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return fmt.Errorf("error deleting data for ID %s: %s", _id, err.Error())
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("data with ID %s not found", _id)
	}

	return nil
}

func GetDocsByFilter(db *mongo.Database, collectionName string, filter bson.M) ([]bson.M, error) {
    var documents []bson.M

    ctx := context.TODO()
	collection := db.Collection(collectionName)

    cursor, err := collection.Find(ctx, filter)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    for cursor.Next(ctx) {
        var document bson.M
        if err := cursor.Decode(&document); err != nil {
            return nil, err
        }
        documents = append(documents, document)
    }

    if err := cursor.Err(); err != nil {
        return nil, err
    }

    return documents, nil
}

//User
func SignUp(db *mongo.Database, col string, insertedDoc model.User) error {
	objectId := primitive.NewObjectID()

	if insertedDoc.Username == "" || insertedDoc.Email == "" || insertedDoc.Password == "" {
		return fmt.Errorf("mohon untuk melengkapi data")
	}
	if err := checkmail.ValidateFormat(insertedDoc.Email); err != nil {
		return fmt.Errorf("email tidak valid")
	}
	userExists, _ := GetUserFromEmail(insertedDoc.Email, db)
	if insertedDoc.Email == userExists.Email {
		return fmt.Errorf("email sudah terdaftar")
	}
	if strings.Contains(insertedDoc.Password, " ") {
		return fmt.Errorf("password tidak boleh mengandung spasi")
	}
	if len(insertedDoc.Password) < 8 {
		return fmt.Errorf("password terlalu pendek")
	}

	hash, _ := HashPassword(insertedDoc.Password)
	// insertedDoc.Password = hash
	user := bson.M{
		"_id":      objectId,
		"username": insertedDoc.Username,
		"email":    insertedDoc.Email,
		"password": hash,
		// "role":     "user",
	}
	_, err := InsertOneDoc(db, col, user)
	if err != nil {
		return err
	}
	return nil
}

func SignIn(db *mongo.Database, col string, insertedDoc model.User) (user model.User, Status bool, err error) {
	if insertedDoc.Email == "" || insertedDoc.Password == "" {
		return user, false, fmt.Errorf("mohon untuk melengkapi data")
	}
	if err = checkmail.ValidateFormat(insertedDoc.Email); err != nil {
		return user, false, fmt.Errorf("email tidak valid")
	}
	existsDoc, err := GetUserFromEmail(insertedDoc.Email, db)
	if err != nil {
		return
	}
	if !CheckPasswordHash(insertedDoc.Password, existsDoc.Password) {
		return user, false, fmt.Errorf("password salah")
	}

	return existsDoc, true, nil
}

func GetAllUser(mongoconn *mongo.Database, collection string) []model.User {
	user := atdb.GetAllDoc[[]model.User](mongoconn, collection)
	return user
}

func GetUserFromID(_id primitive.ObjectID, db *mongo.Database) (doc model.User, err error) {
	collection := db.Collection("user")
	filter := bson.M{"_id": _id}
	err = collection.FindOne(context.TODO(), filter).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return doc, fmt.Errorf("no data found for ID %s", _id)
		}
		return doc, fmt.Errorf("error retrieving data for ID %s: %s", _id, err.Error())
	}
	return doc, nil
}

func GetUserFromEmail(email string, db *mongo.Database) (doc model.User, err error) {
	collection := db.Collection("user")
	filter := bson.M{"email": email}
	err = collection.FindOne(context.TODO(), filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return doc, fmt.Errorf("email tidak ditemukan")
		}
		return doc, fmt.Errorf("kesalahan server")
	}
	return doc, nil
}

//Device
func GetAllDevice(mongoconn *mongo.Database, collection string) []model.Device {
	device := atdb.GetAllDoc[[]model.Device](mongoconn, collection)
	return device
}

func GetDevicesByUser(conn *mongo.Database, collectionname string, email string) ([]model.Device, error) {
	var devices []model.Device
	collection := conn.Collection(collectionname)

	// Menggunakan filter untuk mencari data perangkat yang sesuai dengan ID pengguna
	filter := bson.M{"user": email}

	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return devices, err
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var device model.Device
		if err := cursor.Decode(&device); err != nil {
			return devices, err
		}
		devices = append(devices, device)
	}

	return devices, nil
}

func UpdateDeviceByID(id primitive.ObjectID, db *mongo.Database, doc model.Device) error {
	filter := bson.M{"_id": id}
	result, err := db.Collection("devices").UpdateOne(context.Background(), filter, bson.M{"$set": doc})
	if err != nil {
		return fmt.Errorf("error updating data for ID %s: %s", id, err.Error())
	}

	if result.ModifiedCount == 0 {
        return errors.New("no data has been changed with the specified id")
    }

	return nil
}

func UpdateDeviceStatusByID(id primitive.ObjectID, db *mongo.Database, fieldName string, fieldValue interface{}) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // Adjust the timeout as needed
    defer cancel()

    filter := bson.M{"_id": id}
    update := bson.M{"$set": bson.M{fieldName: fieldValue}}

    result, err := db.Collection("devices").UpdateOne(ctx, filter, update)
    if err != nil {
        return fmt.Errorf("failed to update field %s for ID %s: %w", fieldName, id, err)
    }

    if result.ModifiedCount == 0 {
        return fmt.Errorf("no data has been changed with the specified ID %s", id)
    }

    return nil
}


func DeleteDeviceByID(id primitive.ObjectID, db *mongo.Database) error {
	collection := db.Collection("devices")
	filter := bson.M{"_id": id}
	result, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return fmt.Errorf("error deleting data for ID %s: %s", id, err.Error())
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("data with ID %s not found", id)
	}

	return nil
}

//History
func GetAllHistory(mongoconn *mongo.Database, collection string) []model.History {
	history := atdb.GetAllDoc[[]model.History](mongoconn, collection)
	return history
}

func GetHistoryByUser(conn *mongo.Database, collectionname string, email string) ([]model.History, error) {
	var history []model.History
	collection := conn.Collection(collectionname)

	// Menggunakan filter untuk mencari data histort yang sesuai dengan ID pengguna
	filter := bson.M{"user": email}

	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return history, err
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var datahistory model.History
		if err := cursor.Decode(&datahistory); err != nil {
			return history, err
		}
		history = append(history, datahistory)
	}

	return history, nil
}

func DeleteAllHistoryByUser(conn *mongo.Database, collectionname string, userId string) error {
    collection := conn.Collection(collectionname)
    filter := bson.M{"user": userId}
    _, err := collection.DeleteMany(context.Background(), filter)
    return err
}