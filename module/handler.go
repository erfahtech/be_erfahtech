package beurse

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"

	model "github.com/erfahtech/be_erfahtech/model"
	"github.com/whatsauth/watoken"
)

func GCFReturnStruct(DataStuct any) string {
	jsondata, _ := json.Marshal(DataStuct)
	return string(jsondata)
}

func GCFHandlerSignup(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	var Response model.Credential
	Response.Status = false
	var dataUser model.User
	err := json.NewDecoder(r.Body).Decode(&dataUser)
	if err != nil {
		Response.Message = "error parsing application/json: " + err.Error()
		return GCFReturnStruct(Response)
	}
	err = SignUp(conn, collectionname, dataUser)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	Response.Status = true
	Response.Message = "Halo " + dataUser.Username
	return GCFReturnStruct(Response)
}

func GCFHandlerLogin(PASETOPRIVATEKEYENV, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	var Response model.Credential
	Response.Status = false
	var dataUser model.User
	err := json.NewDecoder(r.Body).Decode(&dataUser)
	if err != nil {
		Response.Message = "error parsing application/json: " + err.Error()
		return GCFReturnStruct(Response)
	}
	user, status1, err := SignIn(conn, collectionname, dataUser)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	Response.Status = true
	tokenstring, err := watoken.Encode(dataUser.Email, os.Getenv(PASETOPRIVATEKEYENV))
	if err != nil {
		Response.Message = "Gagal Encode Token : " + err.Error()
	} else {
		Response.Message = "Selamat Datang " + user.Email + " di USE" + strconv.FormatBool(status1)
		Response.Token = tokenstring
	}
	return GCFReturnStruct(Response)
}

func GCFGetDevice(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
		var Response model.DeviceResponse
		Response.Status = false
		conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	
		// Menyimpan token dari request
		token := r.Header.Get("Authorization")
		token = strings.TrimPrefix(token, "Bearer ")
	
		// token := r
		// token = strings.TrimPrefix(token, "Bearer ")
	
		// Decode token untuk mendapatkan ID pengguna
		user, err := watoken.Decode(os.Getenv(PASETOPUBLICKEYENV), token)
		if err != nil {
			Response.Message = "Error decoding token: " + err.Error()
		} else {
			// Mengambil data perangkat berdasarkan ID pengguna
			devices, err := GetDevicesByUserId(conn, collectionname, user.Id)
			if err != nil {
				Response.Message = "Error fetching devices: " + err.Error()
			} else {
				Response.Status = true
				Response.Message = "Device data successfully retrieved"
				Response.Data = devices
			}
		}
	
		return GCFReturnStruct(Response)
	}