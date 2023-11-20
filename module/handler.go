package beurse

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	model "github.com/erfahtech/be_erfahtech/model"
	"github.com/whatsauth/watoken"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GCFReturnStruct(DataStuct any) string {
	jsondata, _ := json.Marshal(DataStuct)
	return string(jsondata)
}

func GetID(r *http.Request) string {
    return r.URL.Query().Get("id")
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
	user, _, err := SignIn(conn, collectionname, dataUser)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	Response.Status = true
	tokenstring, err := watoken.Encode(dataUser.Email, os.Getenv(PASETOPRIVATEKEYENV))
	if err != nil {
		Response.Message = "Gagal Encode Token : " + err.Error()
	} else {
		Response.Message = "Selamat Datang " + user.Username + " di USE "
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
	devicedata.Status = false
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
	var Response model.Response
	Response.Status = false
	var dataDevice model.Device

	// Get the "id" parameter from the URL
	id := GetID(r)
    if id == "" {
        Response.Message = "Missing 'id' parameter in the URL"
        GCFReturnStruct(Response)
    }

	// Convert the ID string to primitive.ObjectID
	idparam, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		Response.Message = "Invalid id parameter"
		GCFReturnStruct(Response)
	}

	// get token from header
	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	if token == "" {
		Response.Message = "error parsing application/json1:"
		return GCFReturnStruct(Response)
	}

	// decode token
	user, err1 := watoken.Decode(os.Getenv(PASETOPUBLICKEY), token)

	if err1 != nil {
		Response.Message = "error parsing application/json2: " + err1.Error() + ";" + token
		return GCFReturnStruct(Response)
	}

	// Set the user ID in dataDevice
	dataDevice.User = user.Id // Assuming "UserID" is the field where you want to store the user ID in dataDevice

	err = json.NewDecoder(r.Body).Decode(&dataDevice)
	if err != nil {
		Response.Message = "error parsing application/json3: " + err.Error()
		return GCFReturnStruct(Response)
	}
	err = UpdateDeviceByID(idparam, conn, dataDevice)
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
	var Response model.Response
	Response.Status = false

	// Get the "id" parameter from the URL
	id := GetID(r)
	if id == "" {
		Response.Message = "Missing 'id' parameter in the URL"
		return GCFReturnStruct(Response)
	}

	// Convert the ID string to primitive.ObjectID
	idparam, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		Response.Message = "Invalid id parameter"
		GCFReturnStruct(Response)
	}

	// Get token from header
	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	if token == "" {
		Response.Message = "error parsing application/json1:"
		return GCFReturnStruct(Response)
	}

	// Decode token
	_, err1 := watoken.Decode(os.Getenv(PASETOPUBLICKEY), token)

	if err1 != nil {
		Response.Message = "error parsing application/json2: " + err1.Error() + ";" + token
		return GCFReturnStruct(Response)
	}

	// Delete the device based on the ID
	err = DeleteDeviceByID(idparam, conn)
	if err != nil {
		Response.Message = "error deleting device: " + err.Error()
		return GCFReturnStruct(Response)
	}

	Response.Status = true
	Response.Message = "Device berhasil dihapus"
	return GCFReturnStruct(Response)
}

//History

func GCFInsertHistory(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	var Response model.Credential
	var historydata model.History
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	err := json.NewDecoder(r.Body).Decode(&historydata)
	if err != nil {
		Response.Message = "error parsing application/json: " + err.Error()
		return GCFReturnStruct(Response)
	}

	user, err := watoken.Decode(os.Getenv(PASETOPUBLICKEYENV), token)
	if err != nil {
		Response.Message = "Error decoding token: " + err.Error()
		return GCFReturnStruct(Response)
	}
	time := Waktu(time.Now().Format(time.RFC3339))
	historydata.User = user.Id
	historydata.CreatedAt = time
	InsertOneDoc(conn, "history", historydata)
	Response.Status = true
	Response.Message = "History berhasil ditambahkan dengan nama: " + historydata.Name + " dan payload: " + historydata.Payload
	return GCFReturnStruct(Response)
}

func GCFGetHistory(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	var Response model.HistoryResponse
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
		// Mengambil data history berdasarkan ID pengguna
		history, err := GetHistoryByUser(conn, collectionname, user.Id)
		if err != nil {
			Response.Message = "Error fetching history: " + err.Error()
		} else {
			Response.Status = true
			Response.Message = "History data successfully retrieved"
			Response.Data = history
		}
	}
	
	return GCFReturnStruct(Response)
}

func GCFDeleteAllHistory (PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	var Response model.Response
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
		// Mengambil data history berdasarkan ID pengguna
		err := DeleteAllHistoryByUser(conn, collectionname, user.Id)
		if err != nil {
			Response.Message = "Error deleting history: " + err.Error()
		} else {
			Response.Status = true
			Response.Message = "History data successfully deleted"
		}
	}
	
	return GCFReturnStruct(Response)
}