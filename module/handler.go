package beurse

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"

	model "github.com/erfahtech/be_erfahtech/model"
	"github.com/whatsauth/watoken"
	"go.mongodb.org/mongo-driver/bson"
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
		Response.Message = "Selamat Datang " + user.Email + " di USE " + strconv.FormatBool(status1)
		Response.Token = tokenstring
	}
	return GCFReturnStruct(Response)
}

func GCFHandlerGetAll(MONGOCONNSTRINGENV, dbname, col string, docs interface{}) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	data := GetAllDocs(conn, col, docs)
	return GCFReturnStruct(data)
}

//Device

func GCFInsertDevice(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	var Response model.Credential
	var devicedata model.Device
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	err := json.NewDecoder(r.Body).Decode(&devicedata)
	if err != nil {
		Response.Message = "error parsing application/json: " + err.Error()
		return GCFReturnStruct(Response)
	}

	user, err := watoken.Decode(os.Getenv(PASETOPUBLICKEYENV), token)
	    if err != nil {
        Response.Message = "Error decoding token: " + err.Error()
        return GCFReturnStruct(Response)
    }

	devicedata.User = user.Id
	InsertOneDoc(conn, "devices", devicedata)
	Response.Status = true
	Response.Message = "Device berhasil ditambahkan dengan nama: " + devicedata.Name
	return GCFReturnStruct(Response)
}

func GCFGetDevice(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	var Response model.DeviceResponse
	Response.Status = false
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	
	// Menyimpan token dari request
	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	
	// Decode token untuk mendapatkan ID pengguna
	user, err := watoken.Decode(os.Getenv(PASETOPUBLICKEYENV), token)
	if err != nil {
		Response.Message = "Error decoding token: " + err.Error()
	} else {
		// Mengambil data perangkat berdasarkan ID pengguna
		devices, err := GetDevicesByUser(conn, collectionname, user.Id)
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

func GCFGetDeviceByEmail(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	var userdata model.User
	var Response model.DeviceResponse
	Response.Status = false
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	filter := bson.M{"user": userdata.Email}

	devices, err := GetDocsByFilter(conn, collectionname, filter)
	if err != nil {
		var Response model.Credential
		Response.Status = false
		Response.Message = "Error fetching devices: " + err.Error()
		return GCFReturnStruct(Response)
	}
	return GCFReturnStruct(devices)
}

func GCFHandlerUpdateDevice(PASETOPUBLICKEY, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	var Response model.DeviceResponse
	Response.Status = false
	var dataDevice model.Device

	// get token from header
	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	if token == "" {
		Response.Message = "error parsing application/json1:"
		return GCFReturnStruct(Response)
	}

	// decode token
	_, err1 := watoken.Decode(os.Getenv(PASETOPUBLICKEY), token)

	if err1 != nil {
		Response.Message = "error parsing application/json2: " + err1.Error() + ";" + token
		return GCFReturnStruct(Response)
	}

	err := json.NewDecoder(r.Body).Decode(&dataDevice)
	if err != nil {
		Response.Message = "error parsing application/json3: " + err.Error()
		return GCFReturnStruct(Response)
	}
	err = UpdateDevice(conn, dataDevice)
	if err != nil {
		Response.Message = "error parsing application/json4: " + err.Error()
		return GCFReturnStruct(Response)
	}
	Response.Status = true
	Response.Message = "Device berhasil diupdate"
	return GCFReturnStruct(Response)
}

func GCFHandlerDeleteDevice(PASETOPUBLICKEY, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	var Response model.DeviceResponse
	Response.Status = false
	var dataDevice model.Device

	// get token from header
	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	if token == "" {
		Response.Message = "error parsing application/json1:"
		return GCFReturnStruct(Response)
	}

	// decode token
	_, err1 := watoken.Decode(os.Getenv(PASETOPUBLICKEY), token)

	if err1 != nil {
		Response.Message = "error parsing application/json2: " + err1.Error() + ";" + token
		return GCFReturnStruct(Response)
	}

	err := json.NewDecoder(r.Body).Decode(&dataDevice)
	if err != nil {
		Response.Message = "error parsing application/json3: " + err.Error()
		return GCFReturnStruct(Response)
	}
	err = DeleteDevice(conn, dataDevice)
	if err != nil {
		Response.Message = "error parsing application/json4: " + err.Error()
		return GCFReturnStruct(Response)
	}
	Response.Status = true
	Response.Message = "Device berhasil dihapus"
	return GCFReturnStruct(Response)
}

		
	
		