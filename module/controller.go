package beurse

import (
	"bytes"
	"context"
	"crypto/rand"
	"log"
	"math/big"
	"net/http"
	"os"
	"regexp"
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

func ValidatePhoneNumber(phoneNumber string) (bool, error) {
	// Define the regular expression pattern for numeric characters
	numericPattern := `^[0-9]+$`

	// Compile the numeric pattern
	numericRegexp, err := regexp.Compile(numericPattern)
	if err != nil {
		return false, err
	}
	// Check if the phone number consists only of numeric characters
	if !numericRegexp.MatchString(phoneNumber) {
		return false, nil
	}

	// Define the regular expression pattern for "62" followed by 6 to 12 digits
	pattern := `^62\d{6,13}$`

	// Compile the regular expression
	regexpPattern, err := regexp.Compile(pattern)
	if err != nil {
		return false, err
	}

	// Test if the phone number matches the pattern
	isValid := regexpPattern.MatchString(phoneNumber)

	return isValid, nil
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

func GetOTPbyEmail(email string, db *mongo.Database) (doc model.Otp, err error) {
	collection := db.Collection("otp")
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

	if insertedDoc.Username == "" || insertedDoc.Email == "" || insertedDoc.Password == "" || insertedDoc.PhoneNumber == ""{
		return fmt.Errorf("mohon untuk melengkapi data")
	}
	valid, _ := ValidatePhoneNumber(insertedDoc.PhoneNumber)
	if !valid {
		return fmt.Errorf("nomor telepon tidak valid")
	}
	numberphoneExists, _ := GetUserFromPhoneNumber(insertedDoc.PhoneNumber, db)
    if insertedDoc.PhoneNumber == numberphoneExists.PhoneNumber {
        return fmt.Errorf("nomor telepon sudah terdaftar")
    }
	if err := checkmail.ValidateFormat(insertedDoc.Email); err != nil {
		return fmt.Errorf("email tidak valid")
	}
	emailExists, _ := GetUserFromEmail(insertedDoc.Email, db)
	if insertedDoc.Email == emailExists.Email {
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
		"phonenumber": insertedDoc.PhoneNumber,
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

func GetUserFromPhoneNumber(phonenumber string, db *mongo.Database) (doc model.User, err error) {
    collection := db.Collection("user")
    filter := bson.M{"phonenumber": phonenumber}
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

func GetDeviceByID(_id primitive.ObjectID, db *mongo.Database) (doc model.Device, err error) {
	collection := db.Collection("devices")
	filter := bson.M{"_id": _id}
	err = collection.FindOne(context.Background(), filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return doc, fmt.Errorf("no data found for ID %s", _id)
		}
		return doc, fmt.Errorf("error retrieving data for ID %s: %s", _id, err.Error())
	}
	return doc, nil
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

func InsertDevice(iduser string, db *mongo.Database, doc model.Device) (bson.M, error) {

	if doc.Name == "" || doc.Topic == "" {
		return bson.M{}, fmt.Errorf("mohon untuk melengkapi data")
	}
	user, err := GetUserFromEmail(iduser, db)
	if err != nil {
		return bson.M{}, fmt.Errorf("user tidak ditemukan")
	}

	device := bson.M{
		"_id": primitive.NewObjectID(),
		"name": doc.Name,
		"topic": doc.Topic,
		"user": user.Email,
		"status": false,
	}
	_, err = InsertOneDoc(db, "devices", device)
	if err != nil {
		return bson.M{}, err
	}
	return device, nil
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

func EditDevice(idparam primitive.ObjectID, iduser string, db *mongo.Database, doc model.Device) (bson.M, error) {

	if doc.Name == "" || doc.Topic == "" {
		return bson.M{}, fmt.Errorf("mohon untuk melengkapi data")
	}
	user, err := GetUserFromEmail(iduser, db)
	if err != nil {
		return bson.M{}, fmt.Errorf("user tidak ditemukan")
	}
	device, err := GetDeviceByID(idparam, db)
	if err != nil {
		return bson.M{}, fmt.Errorf("device tidak ditemukan")
	}

	data := bson.M{
		"name": doc.Name,
		"topic": doc.Topic,
		"user": user.Email,
	}
	err = UpdateOneDoc(db, "devices", device.ID, data)
	if err != nil {
		return bson.M{}, err
	}
	return data, nil
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

// OTP
func OtpGenerate() (string, error) {
	randomNumber, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return "", err
	}
	// Format the random number as a 6-digit string
	otp := fmt.Sprintf("%06d", randomNumber)

	return otp, nil
}

func GenerateExpiredAt() int64 {
	currentTime := time.Now()

	// Add 5 minutes
	newTime := currentTime.Add(5 * time.Minute)
	return newTime.Unix()
}

func SendOTP(db *mongo.Database, email string) (string, error) {
	// GET OTP
	otp, _ := OtpGenerate()

	// GET EXPIRED AT
	expiredAt := GenerateExpiredAt()

	// get user by email
	existsDoc, err := GetUserFromEmail(email, db)
	if err != nil {
		return "", fmt.Errorf("email tidak ditemukan1")
	}
	if existsDoc.Email == "" {
		return "", fmt.Errorf("email tidak ditemukan2")
	}

	// save otp to db
	// objectId := primitive.NewObjectID()
	otpDoc := bson.M{
		// "_id":       objectId,
		"email":     email,
		"otp":       otp,
		"expiredat": expiredAt,
		"status":    false,
	}

	// get otp by email
	_, err = GetOTPbyEmail(email, db)

	if err != nil {
		if err.Error() == "email tidak ditemukan" {
			// return "", fmt.Errorf("error getting OTP from email: %s", err.Error())
			// insert new OTP
			_, err = db.Collection("otp").InsertOne(context.Background(), otpDoc)
			if err != nil {
				return "", fmt.Errorf("error inserting OTP: %s", err.Error())
			}
			return otp, nil
		} else {
			return "", fmt.Errorf("error Get OTP: %s", err.Error())
		}
	} else {
		// update existing OTP
		filter := bson.M{"email": email}
		update := bson.M{"$set": otpDoc}
		_, err = db.Collection("otp").UpdateOne(context.Background(), filter, update)
		if err != nil {
			return "", fmt.Errorf("error updating OTP: %s", err.Error())
		}
	}

	// postapi
	url := "https://api.wa.my.id/api/send/message/text"

	// Data yang akan dikirimkan dalam format JSON
	jsonStr := []byte(`{
        "to": "` + existsDoc.PhoneNumber + `",
        "isgroup": false,
        "messages": "Berikut kode Otp reset password akun ursmartecosystem.my.id atas nama *` + email + `* adalah *` + otp + `*"
    }`)

	// Membuat permintaan HTTP POST
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return "", err
	}

	// Menambahkan header ke permintaan
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Token", "v4.public.eyJleHAiOiIyMDIzLTEyLTIyVDE0OjQzOjUzKzA3OjAwIiwiaWF0IjoiMjAyMy0xMS0yMlQxNDo0Mzo1MyswNzowMCIsImlkIjoiNjI4NTE2MTk5MjA1MyIsIm5iZiI6IjIwMjMtMTEtMjJUMTQ6NDM6NTMrMDc6MDAifYnf6Vf8-q6QpR20Pso6RWvhq50jonNOV_ucf9-ppSLCem6BDSjdjymp1bQjw81eZlrV0VzKd0arb-0YSl0EiAE")
	req.Header.Set("Content-Type", "application/json")

	// Melakukan permintaan HTTP POST
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return "", err
	}
	defer resp.Body.Close()

	// Menampilkan respons dari server
	fmt.Println("Response Status:", resp.Status)
	return "success", nil
}

func VerifyOTP(db *mongo.Database, email, otp string) (string, error) {
	// get otp by email
	otpDoc, err := GetOTPbyEmail(email, db)
	if err != nil {
		return "", fmt.Errorf("error Get OTP: %s", err.Error())
	}

	// check otp
	if otpDoc.OTP != otp {
		return "", fmt.Errorf("kode otp tidak valid")
	}

	// check expired at
	if otpDoc.ExpiredAt < time.Now().Unix() {
		return "", fmt.Errorf("kode otp telah kadaluarsa")
	}

	//update otp
	filter := bson.M{"email": email}
	update := bson.M{"$set": bson.M{"status": true}}
	_, err = db.Collection("otp").UpdateOne(context.Background(), filter, update)

	if err != nil {
		return "", fmt.Errorf("error updating OTP: %s", err.Error())
	}

	return otp, nil
}

func ResetPassword(db *mongo.Database, email, otp, password string) (string, error) {
	// get user by email
	existsDoc, err := GetUserFromEmail(email, db)
	if err != nil {
		return "", fmt.Errorf("email tidak ditemukan1")
	}
	if existsDoc.Email == "" {
		return "", fmt.Errorf("email tidak ditemukan2")
	}

	// check otp
	docOtp, err := GetOTPbyEmail(email, db)
	if err != nil {
		return "", fmt.Errorf("error Get OTP: %s", err.Error())
	}
	if docOtp.OTP != otp || !docOtp.Status {
		return "", fmt.Errorf("kode otp tidak valid")
	}

	// hash password
	hash, _ := HashPassword(password)

	// update password
	filter := bson.M{"email": email}
	update := bson.M{"$set": bson.M{"password": hash}}
	_, err = db.Collection("user").UpdateOne(context.Background(), filter, update)
	if err != nil {
		return "", fmt.Errorf("error updating password: %s", err.Error())
	}

	// update otp
	filter = bson.M{"email": email}
	update = bson.M{"$set": bson.M{"status": false}}
	_, err = db.Collection("otp").UpdateOne(context.Background(), filter, update)
	if err != nil {
		return "", fmt.Errorf("error updating password: %s", err.Error())
	}

	return "success", nil
}