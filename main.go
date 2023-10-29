package beurse

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/aiteung/atdb"
	"github.com/whatsauth/watoken"
)

func GCFReturnStruct(DataStuct any) string {
	jsondata, _ := json.Marshal(DataStuct)
	return string(jsondata)
}

// func Decode(publicKey string, tokenstring string) (payload Payload, err error) {
// 	var token *paseto.Token
// 	var pubKey paseto.V4AsymmetricPublicKey
// 	pubKey, err = paseto.NewV4AsymmetricPublicKeyFromHex(publicKey) // this wil fail if given key in an invalid format
// 	if err != nil {
// 		fmt.Println("Decode NewV4AsymmetricPublicKeyFromHex : ", err)
// 	}
// 	parser := paseto.NewParser()                                // only used because this example token has expired, use NewParser() (which checks expiry by default)
// 	token, err = parser.ParseV4Public(pubKey, tokenstring, nil) // this will fail if parsing failes, cryptographic checks fail, or validation rules fail
// 	if err != nil {
// 		fmt.Println("Decode ParseV4Public : ", err)
// 	} else {
// 		json.Unmarshal(token.ClaimsJSON(), &payload)
// 	}
// 	return payload, err
// }

// func Encode(id string, privateKey string) (string, error) {
// 	token := paseto.NewToken()
// 	token.SetIssuedAt(time.Now())
// 	token.SetNotBefore(time.Now())
// 	token.SetExpiration(time.Now().Add(2 * time.Hour))
// 	token.SetString("id", id)
// 	secretKey, err := paseto.NewV4AsymmetricSecretKeyFromHex(privateKey)
// 	return token.V4Sign(secretKey, nil), err

// }

func GCFPostHandler(PASETOPRIVATEKEYENV, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	var Response Credential
	Response.Status = false
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	var datauser User
	err := json.NewDecoder(r.Body).Decode(&datauser)
	if err != nil {
		Response.Message = "error parsing application/json: " + err.Error()
	} else {
		if IsPasswordValid(mconn, collectionname, datauser) {
			Response.Status = true
			tokenstring, err := watoken.Encode(datauser.Email, os.Getenv(PASETOPRIVATEKEYENV))
			if err != nil {
				Response.Message = "Gagal Encode Token : " + err.Error()
			} else {
				Response.Message = "Selamat Datang " + datauser.Username
				Response.Token = tokenstring
			}
		} else {
			Response.Message = "Email atau Password Salah"
		}
	}

	return GCFReturnStruct(Response)
}

func GCFHandlerGetAll(MONGOCONNSTRINGENV, dbname, col string) string {
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	data := GetAllUser(mconn, col)
	return GCFReturnStruct(data)
}


func InsertUser(r *http.Request) string {
	var Response Credential
	var userdata User
	err := json.NewDecoder(r.Body).Decode(&userdata)
	if err != nil {
		Response.Message = "error parsing application/json: " + err.Error()
		return GCFReturnStruct(Response)
	}
	hash, _ := HashPassword(userdata.Password)
	userdata.Password = hash
	atdb.InsertOneDoc(SetConnection("MONGOSTRING", "db_urse"), "user", userdata)
	Response.Status = true
	Response.Message = "Akun berhasil dibuat untuk username: " + userdata.Username
	return GCFReturnStruct(Response)
}

func InsertDevice(PASETOPUBLICKEYENV string, r *http.Request) string {
	var Response Credential
	var devicedata Device
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
	mconn := SetConnection("MONGOSTRING", "db_urse")
	atdb.InsertOneDoc(mconn, "devices", devicedata)
	Response.Status = true
	Response.Message = "Device berhasil ditambahkan dengan nama: " + devicedata.Name
	return GCFReturnStruct(Response)
}


func GetDevices(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
// func GetDevices(PASETOPRIVATEKEYENV, MONGOCONNSTRINGENV, dbname, collectionname string, r string) string {
    var Response DeviceResponse
    Response.Status = false
    mconn := SetConnection(MONGOCONNSTRINGENV, dbname)

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
        devices, err := GetDevicesByUserId(mconn, collectionname, user.Id)
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