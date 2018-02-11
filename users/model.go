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
	Name      string        `json:"name" bson:"name"`
	Birthdate string        `json:"birthdate" bson:"birthdate"`
}

// GetOrCreateUser ... returns User by name or
// creates new one if user with this name doesn't exist
func GetOrCreateUser(username string) *User {
	user := FindUserByName(username)

	if user == nil {
		user = CreateUser(username)
		log.WithField("id", user.Id.Hex()).Debug("User created")
	} else {
		log.WithField("id", user.Id.Hex()).Debug("User exist")
	}

	return user
}

// FindUserByName ... queries user by name
func FindUserByName(username string) *User {
	u := User{}
	q := map[string]string{"name": username}
	db := service.GetDatabase()

	if err := db.C("users").Find(q).One(&u); err != nil {
		log.Debug("No user with that Username.")
		return nil
	}

	return &u
}

// CreateUser ... creates new user in user table
func CreateUser(username string) *User {
	db := service.GetDatabase()
	u := User{
		Id:        bson.NewObjectId(),
		Name:      username,
		Birthdate: "",
	}

	db.C("users").Insert(u)

	return &u
}

// SetBirthdate ... updates user's birthdate field
func (u *User) SetBirthdate(birthdate string) {
	db := service.GetDatabase()
	upd := bson.M{"$set": bson.M{"birthdate": birthdate}}

	err := db.C("users").UpdateId(u.Id, upd)
	utils.CheckErr(err, "Database SetBirthday error")

	u.Birthdate = birthdate
}
