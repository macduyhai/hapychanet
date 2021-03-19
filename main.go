package main

import (
	"RestfullApi_Mqtt/msgmqtt"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

// {
// 	"keycode": "",
// 	"date": "2021-02-25 20:58:29",
// 	"action_type": "update",
// 	"detected_image_url": "https://static.hanet.ai/face/upload/C21024B185/2021/02/25/b70263f0-a4a6-4a56-b916-d3c90806ca32.jpg",
// 	"placeID": "1576",
// 	"deviceID": "C21024B185",
// 	"deviceName": "hapyc",
// 	"personName": "",
// 	"aliasID": "",
// 	"data_type": "log",
// 	"personID": "",
// 	"id": "a92d9c5c-0c25-4d11-a683-9bf300cac4e1",
// 	"time": 1614261509000,
// 	"personType": "2",
// 	"placeName": "Welcome",
// 	"hash": "ae8e075381382baa0511b0183a695067"
//   }

type hanetMsg struct {
	Keycode            string `json:"keycode"`
	Date               string `json:"date"`
	Action_type        string `json:"action_type"`
	Detected_image_url string `json:"detected_image_url"`
	PlaceID            string `json:"placeID"`
	DeviceID           string `json:"deviceID"`
	DeviceName         string `json:"deviceName"`
	PersonName         string `json:"personName"`
	AliasID            string `json:"aliasID"`
	Data_type          string `json:"data_type"`
	PersonID           string `json:"personID"`
	Id                 string `json:"id"`
	Time               int    `json:"time"`
	PersonType         string `json:"personType"`
	PlaceName          string `json:"placeName"`
	Hash               string `json:"hash"`
}

type device struct {
	Key string `json:"key"`
	Mac string `json:"mac"`
	Id  int    `json:"id"`
}
type customClaims struct {
	Payload string `json:"payload"`
	jwt.StandardClaims
}

var Device device

// JWT
var jwtSecretKey = []byte("eNhomKou0CMJ694nK281vghbb6UtIQB2")
var msg = "{\"sensor\":\"gps\",\"time\":1351824120}"

// CreateJWT func will used to create the JWT while signing in and signing out
func CreateJWT(payload string) (response string, err error) {
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := customClaims{
		Payload: payload,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			Issuer:    "nameOfWebsiteHere",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecretKey)
	if err == nil {
		return tokenString, nil
	}
	return "", err
}

// VerifyToken func will used to Verify the JWT Token while using APIS
func VerifyToken(tokenString string) (tokenstr string, err error) {
	var claims customClaims

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecretKey, nil
	})

	if token != nil {
		return claims.Payload, nil
	}
	return "", err
}

func homeLink(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	//fmt.Println(string(reqBody))
	if err != nil {
		fmt.Fprintf(w, "Post form not true.! ")
	} else {
		// fmt.Fprintf(w, `{"error_code":10000}`)
	}
	var message hanetMsg
	err1 := json.Unmarshal(reqBody, &message)
	if err1 != nil {
		log.Println(err)
	}
	//fmt.Println(message)
	if message.PersonName != "" {
		fmt.Println("PersonName: " + message.PersonName)
		s := "{\"id\":" + message.Id + "," + "\"value\":\"1\"}"
		fmt.Println(s)
		msgmqtt.PublishData("hapyc", s)
	} else {
		fmt.Println("Unknow")
	}

	fmt.Fprintf(w, "Welcome Hapyc!")
}

// AddDevice asds
func AddDevice(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "AddDevice coming soon!")
}

// GetlistDevice asds
func GetlistDevice(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "GetlistDevice coming soon!")
}
func DeleteDevice(w http.ResponseWriter, r *http.Request) {

}
func GetSttDevice(w http.ResponseWriter, r *http.Request) {

}

// ControDevice asds
func ControDevice(w http.ResponseWriter, r *http.Request) {
	fmt.Println("API control device")
	stt := mux.Vars(r)["stt"]
	reqBody, err := ioutil.ReadAll(r.Body)
	fmt.Println(string(reqBody))
	if err != nil {
		fmt.Fprintf(w, "Post form not true.! Control device not support")
	} else {
		// fmt.Fprintf(w, `{"error_code":10000}`)
	}

	// json.Unmarshal(reqBody, &Device)
	err1 := json.Unmarshal(reqBody, &Device)
	if err1 != nil {
		log.Println(err)
	}
	fmt.Println(Device)
	mac := Device.Mac
	key := Device.Key
	id := Device.Id
	if key == "lvJvDWKiv0" {
		fmt.Printf("Mac device:%s\t-Trạng thái: %s \n", mac, stt)
		if stt == "1" {
			fmt.Println("Bật đèn")
			s := "{\"id\":" + strconv.Itoa(id) + "," + "\"value\":\"1\"}"
			fmt.Println(s)
			msgmqtt.PublishData(mac, s)
		} else if stt == "0" {
			fmt.Println("Tắt đèn")
			s := "{\"id\":" + strconv.Itoa(id) + "," + "\"value\":\"0\"}"
			fmt.Println(s)
			msgmqtt.PublishData(mac, s)
		}
		fmt.Println("------------------------------------------")
		fmt.Fprintf(w, `{"error_code":10000}`)
	} else {
		fmt.Println("Sai Key .! Vui long check lai API")
		fmt.Fprintf(w, `{"error_code":10002,"alert":"Key not true"}`)
	}

}

func main() {
	fmt.Println("====> START MAIN <=====")
	msgmqtt.MqttBegin()
	// time.Sleep(10)
	// s, _ := CreateJWT(msg)
	// fmt.Println(s)
	// fmt.Println("======")
	// payload, err := VerifyToken(s)
	// if err != nil {
	// 	fmt.Print("Paload:")
	// 	fmt.Println(payload)
	// } else {
	// 	fmt.Println("Error decode")
	// 	fmt.Println(err)
	// }
	fmt.Println("======")
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	router.HandleFunc("/add-device", AddDevice).Methods("POST")
	router.HandleFunc("/control-device/{stt}", ControDevice).Methods("POST")
	router.HandleFunc("/get-list-device", GetlistDevice).Methods("GET")
	log.Fatal(http.ListenAndServe(":9999", router))
}
