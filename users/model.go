package users

import (
	"database/sql"

	"github.com/rchesnokov/tg-bot/service"
	"github.com/rchesnokov/tg-bot/utils"

	log "github.com/sirupsen/logrus"
)

// User ... model
type User struct {
	ID        uint
	Name      string
	Birthdate sql.NullString
}

// GetUser ... returns User by name or
// creates new one if user with this name doesn't exist
func GetUser(username string) *User {
	user := FindUserByName(username)
	if user == nil {
		user = CreateUser(username)
		log.Debugf("User created, id is %d", user.ID)
	} else {
		log.Debugf("User exist, id is %d", user.ID)
	}

	return user
}

// FindUserByName ... queries user by name
func FindUserByName(username string) *User {
	var (
		id        uint
		name      string
		birthdate sql.NullString
	)

	db := service.GetDatabase()
	err := db.QueryRow("select id, username, birthdate::date from users where username = $1", username).Scan(&id, &name, &birthdate)

	switch {
	case err == sql.ErrNoRows:
		log.Debug("No user with that Username.")
		return nil
	case err != nil:
		log.Fatal(err)
	}

	return &User{
		id,
		name,
		birthdate,
	}
}

// CreateUser ... creates new user in user table
func CreateUser(username string) *User {
	var newUserID uint

	db := service.GetDatabase()
	err := db.QueryRow("INSERT INTO users(username) VALUES($1) returning id;", username).Scan(&newUserID)
	utils.CheckErr(err, "Database CreateUser error")

	return &User{
		ID:   newUserID,
		Name: username,
		Birthdate: sql.NullString{
			String: "",
			Valid:  false,
		},
	}
}

// SetBirthdate ... updates user's birthdate field
func (u *User) SetBirthdate(birthdate string) {
	db := service.GetDatabase()
	_, err := db.Exec("UPDATE users SET birthdate = $1 WHERE id = $2;", birthdate, u.ID)
	utils.CheckErr(err, "Database SetBirthday error")
	u.Birthdate = sql.NullString{
		String: birthdate + "T00:00:00Z",
		Valid:  true,
	}
}
