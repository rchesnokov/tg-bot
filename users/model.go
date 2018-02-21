package users

import (
	"gopkg.in/mgo.v2/bson"

	"github.com/rchesnokov/tg-bot/service"
	"github.com/rchesnokov/tg-bot/utils"

	log "github.com/sirupsen/logrus"
)

// User ... model
type User struct {
	Id        bson.ObjectId `json:"id" bson:"_id"`
	Username  string        `json:"username" bson:"username"`
	Realname  string        `json:"realname" bson:"realname"`
	Birthdate string        `json:"birthdate" bson:"birthdate"`
	Swearing  uint32        `json:"swearing" bson:"swearing"`
}

// GetOrCreateUser ... returns User by name or
// creates new one if user with this name doesn't exist
// func GetOrCreateUser(username string) *User {
// 	user := FindUserByName(username)

// 	if user == nil {
// 		user = CreateUser(username)
// 		log.WithField("id", user.Id.Hex()).Debug("User created")
// 	} else {
// 		log.WithField("id", user.Id.Hex()).Debug("User exist")
// 	}

// 	return user
// }

// FindByUsername ... queries user by name
func FindByUsername(username string) *User {
	u := User{}
	q := bson.M{"username": username}
	db := service.GetDatabase()

	if err := db.C("users").Find(q).One(&u); err != nil {
		log.Debug("No user with that Username.")
		return nil
	}

	return &u
}

// Create ... creates new user in user table
func Create(username string, realname string) *User {
	db := service.GetDatabase()
	u := User{
		Id:        bson.NewObjectId(),
		Username:  username,
		Realname:  realname,
		Birthdate: "",
	}

	db.C("users").Insert(u)

	return &u
}

// GetName ... returns user's realname if present or his username
func (u *User) GetName() string {
	var name string
	if u.Realname != "" {
		name = u.Realname
	} else {
		name = u.Username
	}

	return name
}

// SetBirthdate ... updates user's birthdate field
func (u *User) SetBirthdate(birthdate string) {
	db := service.GetDatabase()
	upd := bson.M{"$set": bson.M{"birthdate": birthdate}}

	err := db.C("users").UpdateId(u.Id, upd)
	utils.CheckErr(err, "Database SetBirthday error")

	u.Birthdate = birthdate
}

// SetSwearing ... set user's swearing field
func (u *User) SetSwearing(count int) {
	c := u.Swearing + uint32(count)
	db := service.GetDatabase()
	upd := bson.M{"$set": bson.M{"swearing": c}}

	err := db.C("users").UpdateId(u.Id, upd)
	utils.CheckErr(err, "Database SetSwearing error")

	u.Swearing = c
}
