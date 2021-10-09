package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/KK-1729/Appointy-Instagram_API-task/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var collection = ConnectDB()
const postsUrl = "localhost:8080/posts/users/{id}"

//MongoDB Connection
func ConnectDB() *mongo.Collection {
	clientOptions := options.Client().ApplyURI("mongodb+srv://Karthik:Karthik123@userandpost.novhg.mongodb.net/myFirstDatabase?retryWrites=true&w=majority")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Database connection successful")
	collection := client.Database("userandpost").Collection("users")
	return collection

}

//HELPER FUNCTIONS
func GetPasswordHash(password []byte) string {
	hash, err := bcrypt.GenerateFromPassword(password, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}

func GetResponse(url string) (b []byte, e error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln(err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

type ErrorResponse struct {
	StatusCode   int    `json:"status"`
	ErrorMessage string `json:"message"`
}

func GetError(err error, w http.ResponseWriter) {
	log.Fatal(err.Error())
	var response = ErrorResponse {
		ErrorMessage: err.Error(),
		StatusCode: http.StatusInternalServerError,
	}
	message, _ := json.Marshal(response)
	w.WriteHeader(response.StatusCode)
	w.Write(message)
}

//Route Functions
func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user models.User
	_ = json.NewDecoder(r.Body).Decode(&user)
	user.Password = GetPasswordHash([]byte(user.Password))
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	result, err := collection.InsertOne(ctx, user)
	if err != nil {
		GetError(err, w)
		return
	}
	err2 := json.NewEncoder(w).Encode(result)
	if err2 != nil {
		return
	}
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user models.User
	var params = mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	filter := bson.M{"_id": id}
	err := collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		GetError(err, w)
		return
	}
	err2 := json.NewEncoder(w).Encode(user)
	if err2 != nil {
		return
	}
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var post models.Post
	_ = json.NewDecoder(r.Body).Decode(&post)
	result, err := collection.InsertOne(context.TODO(), post)
	if err != nil {
		GetError(err, w)
		return
	}
	err2 := json.NewEncoder(w).Encode(result)
	if err2 != nil {
		return
	}
}

func GetPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var post models.Post
	var params = mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	filter := bson.M{"_id": id}
	err := collection.FindOne(context.TODO(), filter).Decode(&post)
	if err != nil {
		GetError(err, w)
		return
	}
	err2 := json.NewEncoder(w).Encode(post)
	if err2 != nil {
		return
	}
}

func AllPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user models.User
	var params = mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	filter := bson.M{"_id": id}
	err := collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		GetError(err, w)
		return
	}
	stringObjectId := user.UID.Hex()
	t := models.UserResp{}
	t.UserInfo.Posts.PageFormat.HasNextPage = true
	for t.UserInfo.Posts.PageFormat.HasNextPage {
		postUrl := strings.Replace(postsUrl, "{id}", stringObjectId , 1)
		postUrl += t.UserInfo.Posts.PageFormat.EndCursor
		v, err2 := GetResponse(postUrl)
		if err2 != nil {
			GetError(err2, w)
			return
		}
		err3 := json.NewEncoder(w).Encode(v)
		if err3 != nil {
			GetError(err3, w)
			return
		}
	}
	return
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/users", CreateUser).Methods("POST")
	r.HandleFunc("/users/{id}", GetUser).Methods("GET")
	r.HandleFunc("/posts", CreatePost).Methods("POST")
	r.HandleFunc("/posts/{id}", GetPost).Methods("GET")
	r.HandleFunc("/posts/users/{id}", AllPosts).Methods("GET")

	log.Fatal(http.ListenAndServe("localhost:8080", r))
}

