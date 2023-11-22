package beurse

import (
	"fmt"
	"testing"
	"time"

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
	data:=atdb.InsertOneDoc(mconn, "user", userdata)
	fmt.Println(data)
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
	data:=atdb.InsertOneDoc(mconn, "devices", devicedata)
	fmt.Println(data)
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
    user, err := watoken.Decode("c49482e6de1fa07a349f354c2277e11bc7115297a40a1c09c52ef77b905d07c4", "v4.public.eyJleHAiOiIyMDIzLTExLTA0VDEwOjQzOjAxWiIsImlhdCI6IjIwMjMtMTEtMDRUMDg6NDM6MDFaIiwiaWQiOiJkaXRvQGdtYWlsLmNvbSIsIm5iZiI6IjIwMjMtMTEtMDRUMDg6NDM6MDFaIn3ErasDBJ8ZPB0cronqu5S2WqSJ7fyy1YHXM2ovx_B8hrLfekxMtxCDCze8onBf8E02puRACVmq-P8wrxR1X9cC")
    if err != nil {
        t.Errorf("Error decoding token: %v", err)
        return
    }

    var doc model.Device
    doc.Name = "Ac @4"
    doc.Topic = "kamar/ac@4"
    doc.User = user.Id

    id, err := primitive.ObjectIDFromHex("6543b1f219d472b85816dad8")
    doc.ID = id
    if err != nil {
        t.Errorf("Error converting ID: %v", err)
        return
    }

    err = module.UpdateDeviceByID(id, db, doc)
    if err != nil {
        t.Errorf("Error updating document: %v", err)
    } else {
        fmt.Println("Data berhasil diubah dengan ID:", doc.ID)
    }
}
func TestUpdateDeviceStatus(t *testing.T) {
    var doc model.Device
    doc.Status = false

    id, err := primitive.ObjectIDFromHex("653e8ce4550665ec0bcfb9df")
    doc.ID = id
    if err != nil {
        t.Errorf("Error converting ID: %v", err)
        return
    }

    err = module.UpdateDeviceStatusByID(id, db, "status", doc.Status)
    if err != nil {
        t.Errorf("Error updating document: %v", err)
    } else {
        fmt.Println("Data berhasil diubah dengan ID:", doc.ID)
    }
}


func TestDelete(t *testing.T) {
	var doc model.Device
	conn := module.SetConnection("MONGOSTRING", "db_urse")
	id, err := primitive.ObjectIDFromHex("65439fa2aa2593eebaa09e88")
	doc.ID = id
	if err != nil {
		fmt.Printf("Data tidak berhasil dihapus dengan id : ")
	} else {
		err = module.DeleteDeviceByID(id, conn)
		if err != nil {
			t.Errorf("Error updating document: %v", err)
		} else {
			fmt.Println("Data berhasil dihapus dengan id : ", doc.ID)
		}
	}
}

func TestInsertHistory(t *testing.T) {
	currentTime := time.Now()
	waktu := module.Waktu(currentTime.Format(time.RFC3339))

	var historydata model.History
	mconn := module.SetConnection("MONGOSTRING", "db_urse")
	// token, _ := watoken.Decode("c49482e6de1fa07a349f354c2277e11bc7115297a40a1c09c52ef77b905d07c4", "v4.public.eyJleHAiOiIyMDIzLTExLTEzVDEwOjA2OjQ1WiIsImlhdCI6IjIwMjMtMTEtMTNUMDg6MDY6NDVaIiwiaWQiOiJkaXRvQGdtYWlsLmNvbSIsIm5iZiI6IjIwMjMtMTEtMTNUMDg6MDY6NDVaIn3W8ZmzU2TEiWjsvDqrnfMRtPXQRMWI2t2UUl5Y5oxUp-IwXQCMYHo6kt-A3yqjFamgWNOKq6aIkEovhbuqpGoC")
	const email = "dito@gmail.com"

	historydata.Name = "Lampu3"
	historydata.Topic = "test/lampu5"
	historydata.Payload = "1"
	historydata.User = email
	fmt.Println("Waktu sekarang:", waktu)
	historydata.CreatedAt = waktu
	fmt.Println("Waktu yang akan dimasukkan ke dalam MongoDB:", historydata.CreatedAt)


	data := atdb.InsertOneDoc(mconn, "history", historydata)
	fmt.Println(data)
}


func TestGetAllHistory(*testing.T){
	mconn := module.SetConnection("MONGOSTRING", "db_urse")	
	history := module.GetAllHistory(mconn, "history")
	fmt.Println(history)
}

func TestGetHistoryByUser(*testing.T){
	// token,_:=watoken.Decode("c49482e6de1fa07a349f354c2277e11bc7115297a40a1c09c52ef77b905d07c4","v4.public.eyJleHAiOiIyMDIzLTExLTEzVDEwOjA2OjQ1WiIsImlhdCI6IjIwMjMtMTEtMTNUMDg6MDY6NDVaIiwiaWQiOiJkaXRvQGdtYWlsLmNvbSIsIm5iZiI6IjIwMjMtMTEtMTNUMDg6MDY6NDVaIn3W8ZmzU2TEiWjsvDqrnfMRtPXQRMWI2t2UUl5Y5oxUp-IwXQCMYHo6kt-A3yqjFamgWNOKq6aIkEovhbuqpGoC")
	email := "dito@gmail.com"
	mconn := module.SetConnection("MONGOSTRING", "db_urse")
	history,_:=module.GetHistoryByUser(mconn, "history", email)
	fmt.Println(history)
}

func TestDeleteAllHistory(t *testing.T){
	// token,_:=watoken.Decode("c49482e6de1fa07a349f354c2277e11bc7115297a40a1c09c52ef77b905d07c4","v4.public.eyJleHAiOiIyMDIzLTExLTEyVDA1OjEzOjM4WiIsImlhdCI6IjIwMjMtMTEtMTJUMDM6MTM6MzhaIiwiaWQiOiJkaXRvQGdtYWlsLmNvbSIsIm5iZiI6IjIwMjMtMTEtMTJUMDM6MTM6MzhaIn2tOpGc5ISmksdBUgsD_l7qpWVAqqCmQIbC3Cd9sW82sVxaNagaqyNQwRb5t_E_7tn1dvv78Ndw9Pe85fIueRUF")
	email := "dito@gmail.com"
	mconn := module.SetConnection("MONGOSTRING", "db_urse")
	err := module.DeleteAllHistoryByUser(mconn, "history", email)
	if err != nil {
		t.Errorf("Failed to delete history: %v", err)
	} else {
		fmt.Println("Successfully deleted all history for user:", email)
	}
}

func TestWaktu(t *testing.T){
	// s := "2022-03-23T07:00:00+01:00"
	time := module.Waktu(time.Now().Format(time.RFC3339))
	fmt.Println(time)
}

//test otp
func TestGenerateOTP(t *testing.T) {
	var email = ""
	otp, _ := module.OtpGenerate()
	var expiredAt = module.GenerateExpiredAt()
	var doc model.Otp
	doc.Email = email
	doc.OTP = otp
	doc.ExpiredAt = expiredAt
	fmt.Println(otp)
	fmt.Println(expiredAt)
}

func TestSendOTP(t *testing.T) {
	var email = ""
	otp, _ := module.OtpGenerate()
	var expiredAt = module.GenerateExpiredAt()
	var doc model.Otp
	doc.Email = email
	doc.OTP = otp
	doc.ExpiredAt = expiredAt
	fmt.Println(otp)
	fmt.Println(expiredAt)
	otp, err := module.SendOTP(db, "email@gmail.com")
	if err != nil {
		fmt.Println("Error sending otp: ", err)
	} else {
		fmt.Println("Data berhasil dikirim :", otp)
	}

}

func TestCekOTP(t *testing.T) {
	var email = "email@gmail.com"
	otp := "6453"
	otp, err := module.VerifyOTP(db, email, otp)
	if err != nil {
		fmt.Println("Error sending otp: ", err)
	} else {
		fmt.Println("Data berhasil dikirim :", otp)
	}
}

func TestUpdatePassword(t *testing.T) {
	var email = "email@gmail.com"
	otp := "6453"
	password := "daniaw"
	message, err := module.ResetPassword(db, email, otp, password)
	if err != nil {
		fmt.Println("Error sending otp: ", err)
	} else {
		fmt.Println("Data berhasil dikirim :", message)
	}
}

func TestValidatePhoneNumber(t *testing.T) {
	phoneNumber := "62812345690"
	isValid, err := module.ValidatePhoneNumber(phoneNumber)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if isValid {
		fmt.Println("Phone number is valid.")
	} else {
		fmt.Println("Phone number is not valid.")
	}
}