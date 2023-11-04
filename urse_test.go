package beurse

import (
	"fmt"
	"testing"

	"github.com/aiteung/atdb"
	"github.com/whatsauth/watoken"
	"go.mongodb.org/mongo-driver/bson"

	model "github.com/erfahtech/be_erfahtech/model"
	module "github.com/erfahtech/be_erfahtech/module"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var db = module.MongoConnect("MONGOSTRING", "db_urse")

func TestGeneratePasswordHash(t *testing.T) {
	password := "secret"
	hash, _ := module.HashPassword(password) // ignore error for the sake of simplicity

	fmt.Println("Password:", password)
	fmt.Println("Hash:    ", hash)

	match := module.CheckPasswordHash(password, hash)
	fmt.Println("Match:   ", match)
}
func TestGeneratePrivateKeyPaseto(t *testing.T) {
	privateKey, publicKey := watoken.GenerateKey()
	fmt.Println("Ini Private", privateKey)
	fmt.Println("Ini Public", publicKey)
	hasil, err := watoken.Encode("urse", privateKey)
	fmt.Println("Ini Hasil", hasil, err)
}

func TestHashFunction(t *testing.T) {
	mconn := module.SetConnection("MONGOSTRING", "db_urse")
	var userdata model.User
	userdata.Email = "dito@gmail.com"
	userdata.Password = "secret"

	filter := bson.M{"email": userdata.Email}
	res := atdb.GetOneDoc[model.User](mconn, "user", filter)
	fmt.Println("Mongo User Result: ", res)
	hash, _ := module.HashPassword(userdata.Password)
	fmt.Println("Hash Password : ", hash)
	match := module.CheckPasswordHash(userdata.Password, res.Password)
	fmt.Println("Match:   ", match)

}

func TestIsPasswordValid(t *testing.T) {
	mconn := module.SetConnection("MONGOSTRING", "db_urse")
	var userdata model.User
	userdata.Email = "dito@gmail.com"
	userdata.Password = "secret"

	anu := module.IsPasswordValid(mconn, "user", userdata)
	fmt.Println(anu)
}

func TestSignUp(t *testing.T) {
	var doc model.User
	doc.Username = "Erdito Nausha Adam"
	doc.Email = "dito@gmail.com"
	doc.Password = "secret"

	err := module.SignUp(db, "user", doc)
	if err != nil {
		t.Errorf("Error inserting document: %v", err)
	} else {
		fmt.Println("Data berhasil disimpan dengan nama :", doc.Username)
	}
}

func TestLogIn(t *testing.T) {
	var doc model.User
	doc.Email = "dito@gmail.com"
	doc.Password = "secret"
	user, Status, err := module.SignIn(db, "user", doc)
	fmt.Println("Status :", Status)
	if err != nil {
		t.Errorf("Error getting document: %v", err)
	} else {
		fmt.Println("Selamat Datang :", user)
	}
}


func TestInsertUser(*testing.T){
	var userdata model.User 
	mconn := module.SetConnection("MONGOSTRING", "db_urse")
	userdata.Username = "fatwa"
	userdata.Password = "secretcuy"

	hash, _ := module.HashPassword(userdata.Password)
	userdata.Password = hash
	nama:=atdb.InsertOneDoc(mconn, "user", userdata)
	fmt.Println(nama)
}

func TestGetAllUser(*testing.T){	
	mconn := module.SetConnection("MONGOSTRING", "db_urse")	
	user := module.GetAllUser(mconn, "user")
	fmt.Println(user)
}

func TestGetAllDevice(*testing.T){	
	mconn := module.SetConnection("MONGOSTRING", "db_urse")	
	device := module.GetAllDevice(mconn, "devices")
	fmt.Println(device)
}

func TestInsertDevice(*testing.T){
	var devicedata model.Device
	mconn := module.SetConnection("MONGOSTRING", "db_urse")
	token,_:=watoken.Decode("c49482e6de1fa07a349f354c2277e11bc7115297a40a1c09c52ef77b905d07c4","v4.public.eyJleHAiOiIyMDIzLTEwLTI0VDEwOjI3OjI2WiIsImlhdCI6IjIwMjMtMTAtMjRUMDg6Mjc6MjZaIiwiaWQiOiJlcmZhaEBnbWFpbC5jb20iLCJuYmYiOiIyMDIzLTEwLTI0VDA4OjI3OjI2WiJ98pBh-mjEoJlp-4vOVFrfzBcFZzzVsavflcv-wQWfGAVNDGL3A4ebwfNwzG91OnRWHDLbM17VghkQa578tLMhAg")
	devicedata.Name = "Lampu"
	devicedata.Topic = "test/lampu"
	devicedata.User = token.Id
	nama:=atdb.InsertOneDoc(mconn, "devices", devicedata)
	fmt.Println(nama)
}

func TestGetDevicesByUser(*testing.T){
	token,_:=watoken.Decode("c49482e6de1fa07a349f354c2277e11bc7115297a40a1c09c52ef77b905d07c4","v4.public.eyJleHAiOiIyMDIzLTEwLTMwVDAyOjM5OjMwWiIsImlhdCI6IjIwMjMtMTAtMzBUMDA6Mzk6MzBaIiwiaWQiOiJlcmZhaEBnbWFpbC5jb20iLCJuYmYiOiIyMDIzLTEwLTMwVDAwOjM5OjMwWiJ9TRYrR-Ffd_4e1yMaSgkWrcffu7ebEcPmq8VG3_8-MnfNt8cqIStVVbr-0qk5IQom5B3btqK42DhDurCweQu3Ag")
	mconn := module.SetConnection("MONGOSTRING", "db_urse")
	devices,_:=module.GetDevicesByUser(mconn, "devices", token.Id)
	fmt.Println(devices)
}

func TestGetDevicesByEmail(*testing.T){
	var userdata model.User
	userdata.Email = "erfah@gmail.com"
	mconn := module.SetConnection("MONGOSTRING", "db_urse")
	filter := bson.M{"user": userdata.Email}
	devices,_:=module.GetDocsByFilter(mconn, "devices", filter)
	fmt.Println(devices)
}

func TestUpdateDevice(t *testing.T) {
	user,_:=watoken.Decode("c49482e6de1fa07a349f354c2277e11bc7115297a40a1c09c52ef77b905d07c4","v4.public.eyJleHAiOiIyMDIzLTExLTA0VDA3OjQyOjU5WiIsImlhdCI6IjIwMjMtMTEtMDRUMDU6NDI6NTlaIiwiaWQiOiJkaXRvQGdtYWlsLmNvbSIsIm5iZiI6IjIwMjMtMTEtMDRUMDU6NDI6NTlaIn1nAaEIarEleWlY2-Q40BYeHRXJzvyk9qnKFi1xjLtFqRoTrveB-MaS5UCwyMlppZM3hwCiVmyE9cYc128lBEQD")
	var doc model.Device
	doc.Name = "Ac @2"
	doc.Topic = "kamar/ac@2"
	doc.User = user.Id
	id, err := primitive.ObjectIDFromHex("6543b1f219d472b85816dad8")
	doc.ID = id
	if err != nil {
		fmt.Printf("Data tidak berhasil diubah dengan id")
	} else {

		err = module.UpdateDevice(db, doc)
		if err != nil {
			t.Errorf("Error updateting document: %v", err)
		} else {
			fmt.Println("Data berhasil diubah dengan id :", doc.ID)
		}
	}

}

func TestDelete(t *testing.T) {
	var doc model.Device
	id, err := primitive.ObjectIDFromHex("6545b7461dde927263a00ccd")
	doc.ID = id
	if err != nil {
		fmt.Printf("Data tidak berhasil dihapus dengan id : ")
	} else {
		err = module.DeleteDevice(db, doc)
		if err != nil {
			t.Errorf("Error updating document: %v", err)
		} else {
			fmt.Println("Data berhasil dihapus dengan id : ", doc.ID)
		}
	}
}