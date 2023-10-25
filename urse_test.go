package beurse

import (
	"fmt"
	"testing"

	"github.com/aiteung/atdb"
	"github.com/whatsauth/watoken"
	"go.mongodb.org/mongo-driver/bson"
)

func TestGeneratePasswordHash(t *testing.T) {
	password := "secret"
	hash, _ := HashPassword(password) // ignore error for the sake of simplicity

	fmt.Println("Password:", password)
	fmt.Println("Hash:    ", hash)

	match := CheckPasswordHash(password, hash)
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
	mconn := SetConnection("MONGOSTRING", "urse")
	var userdata User
	userdata.Username = "dito"
	userdata.Password = "secret"

	filter := bson.M{"username": userdata.Username}
	res := atdb.GetOneDoc[User](mconn, "user", filter)
	fmt.Println("Mongo User Result: ", res)
	hash, _ := HashPassword(userdata.Password)
	fmt.Println("Hash Password : ", hash)
	match := CheckPasswordHash(userdata.Password, res.Password)
	fmt.Println("Match:   ", match)

}

func TestIsPasswordValid(t *testing.T) {
	mconn := SetConnection("MONGOSTRING", "urse")
	var userdata User
	userdata.Username = "dani"
	userdata.Password = "secretoo"

	anu := IsPasswordValid(mconn, "user", userdata)
	fmt.Println(anu)
}

func TestInsertUser(*testing.T){
	var userdata User 
	mconn := SetConnection("MONGOSTRING", "db_urse")
	userdata.Username = "fatwa"
	userdata.Password = "secretcuy"

	hash, _ := HashPassword(userdata.Password)
	userdata.Password = hash
	nama:=atdb.InsertOneDoc(mconn, "user", userdata)
	fmt.Println(nama)
}

func TestInsertDevice(*testing.T){
	var devicedata Device
	mconn := SetConnection("MONGOSTRING", "db_urse")
	token,_:=watoken.Decode("c49482e6de1fa07a349f354c2277e11bc7115297a40a1c09c52ef77b905d07c4","v4.public.eyJleHAiOiIyMDIzLTEwLTI0VDEwOjI3OjI2WiIsImlhdCI6IjIwMjMtMTAtMjRUMDg6Mjc6MjZaIiwiaWQiOiJlcmZhaEBnbWFpbC5jb20iLCJuYmYiOiIyMDIzLTEwLTI0VDA4OjI3OjI2WiJ98pBh-mjEoJlp-4vOVFrfzBcFZzzVsavflcv-wQWfGAVNDGL3A4ebwfNwzG91OnRWHDLbM17VghkQa578tLMhAg")
	devicedata.Name = "Lampu"
	devicedata.Topic = "test/lampu"
	devicedata.User = token.Id
	nama:=atdb.InsertOneDoc(mconn, "devices", devicedata)
	fmt.Println(nama)
}





