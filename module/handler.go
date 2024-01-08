package beurse

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aiteung/atapi"
	"github.com/aiteung/atmessage"
	model "github.com/erfahtech/be_erfahtech/model"
	"github.com/whatsauth/wa"
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
func GCFHandlerLoginWhatsauth(PASETOPRIVATEKEYENV, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
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
		dt := &wa.TextMessage{
			To:       dataUser.PhoneNumber,
			IsGroup:  false,
			Messages: "Halo kamu, " + dataUser.Username + "\n Berhasil Login di website ursmartecosystem.my.id \n",
		}
		atapi.PostStructWithToken[atmessage.Response]("Token", os.Getenv("TOKEN"), dt, "https://api.wa.my.id/api/send/message/text")
	}
	return GCFReturnStruct(Response)
}

func GCFHandlerGetAll(MONGOCONNSTRINGENV, dbname, col string, docs interface{}) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	data := GetAllDocs(conn, col, docs)
	return GCFReturnStruct(data)
}

func GCFGetUserByEmail(MONGOCONNSTRINGENV, PASETOPUBLICKEYENV, dbname, collectionname string, r *http.Request) string {
	var Response model.Credential
	Response.Status = false
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")

	profile, err := watoken.Decode(os.Getenv(PASETOPUBLICKEYENV), token)
	    if err != nil {
        Response.Message = "Error decoding token: " + err.Error()
        return GCFReturnStruct(Response)
    }
	filter := bson.M{"email": profile.Id}
	user, err := GetDocsByFilter(conn, collectionname, filter)
	if err != nil {
		var Response model.Credential
		Response.Status = false
		Response.Message = "Error fetching user: " + err.Error()
		return GCFReturnStruct(Response)
	}
	return GCFReturnStruct(user)
}

// Device

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

// func GCFPostDevice(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
// 	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
// 	var response model.Response
// 	var dataDevice model.Device
// 	token := r.Header.Get("Authorization")
// 	token = strings.TrimPrefix(token, "Bearer ")
// 	response.Status = false
// 	//
// 	user, err := watoken.Decode(os.Getenv(PASETOPUBLICKEYENV), token)
// 	if err != nil {
// 		response.Message = err.Error()
// 		return GCFReturnStruct(response)
// 	}

// 	data, err := InsertDevice(user.Id, conn, dataDevice)
// 	if err != nil {
// 		response.Message = err.Error()
// 		return GCFReturnStruct(response)
// 	}
// 	//
// 	response.Status = true
// 	response.Message = "Berhasil Menambahkan Device"
// 	responData := bson.M{
// 		"status":  response.Status,
// 		"message": response.Message,
// 		"data":    data,
// 	}
// 	return GCFReturnStruct(responData)
// }

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

// func GCFEditDevice(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
// 	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
// 	var response model.Response
// 	var dataDevice model.Device
// 	token := r.Header.Get("Authorization")
// 	token = strings.TrimPrefix(token, "Bearer ")
// 	response.Status = false
// 	//
// 	user, err := watoken.Decode(os.Getenv(PASETOPUBLICKEYENV), token)
// 	if err != nil {
// 		response.Message = err.Error()
// 		return GCFReturnStruct(response)
// 	}
// 	id := GetID(r)
// 	if id == "" {
// 		response.Message = "Wrong parameter"
// 		return GCFReturnStruct(response)
// 	}
// 	idparam, err := primitive.ObjectIDFromHex(id)
// 	if err != nil {
// 		response.Message = "Invalid id parameter"
// 		return GCFReturnStruct(response)
// 	}
// 	data, err := EditDevice( idparam, user.Id, conn, dataDevice)
// 	if err != nil {
// 		response.Message = err.Error()
// 		return GCFReturnStruct(response)
// 	}
// 	//
// 	response.Status = true
// 	response.Message = "Berhasil mengubah device"
// 	responData := bson.M{
// 		"status":  response.Status,
// 		"message": response.Message,
// 		"data":    data,
// 	}
// 	return GCFReturnStruct(responData)
// }

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

func GCFHandlerUpdateStatusDevice(PASETOPUBLICKEY, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
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
	_, err1 := watoken.Decode(os.Getenv(PASETOPUBLICKEY), token)

	if err1 != nil {
		Response.Message = "error parsing application/json2: " + err1.Error() + ";" + token
		return GCFReturnStruct(Response)
	}

	err = json.NewDecoder(r.Body).Decode(&dataDevice)
	if err != nil {
		Response.Message = "error parsing application/json3: " + err.Error()
		return GCFReturnStruct(Response)
	}
	err = UpdateDeviceStatusByID(idparam, conn, "status", dataDevice.Status)
	if err != nil {
		Response.Message = "error parsing application/json4: " + err.Error()
		return GCFReturnStruct(Response)
	}
	Response.Status = true
	Response.Message = "Status Device berhasil diupdate"
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

//otp

func GCFHandlerSendOTP(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	var Response model.Response
	Response.Status = false
	var dataUser model.User
	err := json.NewDecoder(r.Body).Decode(&dataUser)
	if err != nil {
		Response.Message = "error parsing application/json: " + err.Error()
		return GCFReturnStruct(Response)
	}
	_, err = SendOTP(conn, dataUser.Email)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	Response.Status = true
	Response.Message = "msg " + dataUser.Email
	return GCFReturnStruct(Response)
}

func GCFHandlerVerifyOTP(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	var Response model.Response
	Response.Status = false
	var dataOTP model.Otp
	err := json.NewDecoder(r.Body).Decode(&dataOTP)
	if err != nil {
		Response.Message = "error parsing application/json: " + err.Error()
		return GCFReturnStruct(Response)
	}
	otp, err := VerifyOTP(conn, dataOTP.Email, dataOTP.OTP)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	Response.Status = true
	Response.Message = otp
	return GCFReturnStruct(Response)
}

func GCFHandlerResetPassword(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	var Response model.Response
	Response.Status = false
	var data model.ResetPassword

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		Response.Message = "error parsing application/json: " + err.Error()
		return GCFReturnStruct(Response)
	}

	_, err = ResetPassword(conn, data.Email, data.OTP, data.Password)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	Response.Status = true
	Response.Message = "Berhasil, Silahkan Kembali ke Login " + data.Email
	return GCFReturnStruct(Response)
}