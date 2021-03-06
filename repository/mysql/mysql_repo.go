package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pintobikez/popmeet/api/models"
	cnfs "github.com/pintobikez/popmeet/config/structures"
	serror "github.com/pintobikez/popmeet/errors"
	"strconv"
)

const (
	IsEmpty = "%s is empty"
)

type Client struct {
	config *cnfs.DatabaseConfig
	db     *sql.DB
	tx     *sql.Tx
}

func New(cnfg *cnfs.DatabaseConfig) (*Client, error) {
	if cnfg == nil {
		return nil, fmt.Errorf("Client configuration not loaded")
	}

	return &Client{config: cnfg}, nil
}

// Connects to the mysql database
func (r *Client) Connect() error {

	urlString, err := r.buildStringConnection()
	if err != nil {
		return err
	}

	r.db, err = sql.Open("mysql", urlString)
	if err != nil {
		return err
	}
	return nil
}

// Disconnects from the mysql database
func (r *Client) Disconnect() {
	r.db.Close()
}

// InsertUser Creates a new record in the user table
func (r *Client) InsertUser(u *models.User) error {
	var err error
	// start a transaction
	r.tx, err = r.db.Begin()
	if err != nil {
		r.tx = nil
		return err
	}
	defer r.deferRollback()

	stmt, err := r.tx.Prepare("INSERT INTO `user` VALUES (null,?,?,now(),now(),1)")
	if err != nil {
		return fmt.Errorf("Error in insert user prepared statement: %s", err.Error())
	}

	res, err := stmt.Exec(u.Email, u.Name)
	defer stmt.Close()

	if err != nil {
		return fmt.Errorf("Error in insert user %s, email: %s %s", u.Name, u.Email, err.Error())
	}

	u.ID, _ = res.LastInsertId()
	if err = r.InsertUserSecurity(u.Security, u.ID); err != nil {
		return err
	}

	r.commit()

	return nil
}

// UpdateUser Updates the given user in the user table
func (r *Client) UpdateUser(u *models.User) error {
	var err error
	// start a transaction
	r.tx, err = r.db.Begin()
	if err != nil {
		r.tx = nil
		return err
	}
	defer r.deferRollback()

	stmt, err := r.tx.Prepare("UPDATE `user` SET email=?,name=?,updated_at=now() WHERE id=?")
	if err != nil {
		return fmt.Errorf("Error in update user prepared statement: %s", err.Error())
	}

	_, err = stmt.Exec(u.Email, u.Name, u.ID)
	defer stmt.Close()

	fmt.Printf("%v", u.Profile)

	if err != nil {
		return fmt.Errorf("Could not update userID %d : %s", u.ID, err.Error())
	}

	// UPDATE SECURITY
	if u.Security != nil {
		if err = r.UpdateUserSecurity(u.Security); err != nil {
			return err
		}
	}
	// UPDATE PROFILE
	if u.Profile != nil && u.Profile.ID > 0 {
		if u.Profile.Interests == nil {
			u.Profile.Interests = []*models.Interest{}
		}
		if err = r.UpdateUserProfile(u.Profile); err != nil {
			return err
		}
	}
	if u.Profile != nil && u.Profile.ID <= 0 {
		if u.Profile.Interests == nil {
			u.Profile.Interests = []*models.Interest{}
		}
		if err = r.InsertUserProfile(u.Profile, u.ID); err != nil {
			return err
		}
	}

	r.commit()

	return nil
}

// GetUserById Get an User by a given id
func (r *Client) GetUserById(id int64) (*models.User, error) {
	var found bool
	resp := &models.User{}

	err := r.db.QueryRow("SELECT IF(COUNT(*),'true','false') FROM user WHERE id=?", id).Scan(&found)
	if err != nil {
		return resp, err
	}

	if !found {
		return resp, fmt.Errorf("User with id %d not found", id)
	}

	err = r.db.QueryRow("SELECT id, email, name, created_at, updated_at, active FROM user WHERE id=?", id).Scan(&resp.ID, &resp.Email, &resp.Name, &resp.CreatedAt, &resp.UpdatedAt, &resp.Active)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// GetUserByEmail Get an User by a given email
func (r *Client) GetUserByEmail(email string) (*models.User, error) {
	var found bool
	resp := &models.User{}

	err := r.db.QueryRow("SELECT IF(COUNT(*),'true','false') FROM user WHERE email=?", email).Scan(&found)
	if err != nil {
		return resp, err
	}

	if !found {
		return resp, fmt.Errorf("User with email %s not found", email)
	}

	err = r.db.QueryRow("SELECT id, email, name, created_at, updated_at, active FROM user WHERE email=?", email).Scan(&resp.ID, &resp.Email, &resp.Name, &resp.CreatedAt, &resp.UpdatedAt, &resp.Active)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// FindUserByEmail Checks if a given email exist
func (r *Client) FindUserByEmail(email string) (bool, error) {
	var found bool

	err := r.db.QueryRow("SELECT IF(COUNT(*),'true','false') FROM user WHERE email=? and active=1", email).Scan(&found)
	if err != nil {
		return false, err
	}

	return found, nil
}

// FindUserById Check if the user exists and its active
func (r *Client) FindUserById(id int64) (bool, error) {

	var found bool
	err := r.db.QueryRow("SELECT IF(COUNT(*),'true','false') FROM user WHERE id=? and active=true", id).Scan(&found)
	if err != nil {
		return false, err
	}

	return found, nil
}

// InsertUserProfile Creates a new record in the user_profile table
func (r *Client) InsertUserProfile(u *models.UserProfile, id int64) error {
	var err error
	var stmt *sql.Stmt

	if r.tx != nil {
		stmt, err = r.tx.Prepare("INSERT INTO `user_profile` VALUES (null,?,?,?,?,now())")
	} else {
		stmt, err = r.db.Prepare("INSERT INTO `user_profile` VALUES (null,?,?,?,?,now())")
	}

	if err != nil {
		return fmt.Errorf("Error in insert user_profile prepared statement: %s", err.Error())
	}

	res, err := stmt.Exec(id, u.Language.ID, u.AgeRange, u.Sex)
	defer stmt.Close()

	if err != nil {
		return fmt.Errorf("Error in insert user_profile for user id: %d", id, err.Error())
	}
	u.ID, _ = res.LastInsertId()

	//Update User Interests
	if u.Interests != nil {
		if err = r.UpdateUserInterests(u.Interests, u.ID); err != nil {
			return err
		}
	}

	return nil
}

// UpdateUserProfile Updates the given user in the user_profile table
func (r *Client) UpdateUserProfile(u *models.UserProfile) error {
	var err error
	var stmt *sql.Stmt

	if r.tx != nil {
		stmt, err = r.tx.Prepare("UPDATE `user_profile` SET fk_language=?,sex=?,age_range=?,updated_at=now() WHERE id=?")
	} else {
		stmt, err = r.db.Prepare("UPDATE `user_profile` SET fk_language=?,sex=?,age_range=?,updated_at=now() WHERE id=?")
	}

	_, err = stmt.Exec(u.Language.ID, u.Sex, u.AgeRange, u.ID)
	defer stmt.Close()

	if err != nil {
		return fmt.Errorf("Error in update user_profile userID %d : %s", u.ID, err.Error())
	}

	//Update User Interests
	if err = r.UpdateUserInterests(u.Interests, u.ID); err != nil {
		return err
	}

	return nil
}

// GetUserProfileByUserId Get the UserProfile by a given User id
func (r *Client) GetUserProfileByUserId(id int64) (*models.UserProfile, error) {
	var found bool
	var fkLanguage int64
	resp := &models.UserProfile{}

	err := r.db.QueryRow("SELECT IF(COUNT(*),'true','false') FROM user_profile WHERE fk_user=?", id).Scan(&found)
	if err != nil {
		return resp, err
	}

	if !found {
		return resp, fmt.Errorf("UserProfile for user with id %d not found", id)
	}

	err = r.db.QueryRow("SELECT id,fk_language,age_range,sex,updated_at FROM user_profile WHERE fk_user=?", id).Scan(&resp.ID, &fkLanguage, &resp.AgeRange, &resp.Sex, &resp.UpdatedAt)
	if err != nil {
		return resp, err
	}

	resp.Language, err = r.GetLanguageById(fkLanguage)
	if err != nil {
		return resp, err
	}

	resp.Interests, err = r.GetAllInterestByUserProfileId(id)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// InsertUserSecurity Creates a new record in the user_security table
func (r *Client) InsertUserSecurity(u *models.UserSecurity, id int64) error {

	var err error
	var stmt *sql.Stmt

	if r.tx != nil {
		stmt, err = r.tx.Prepare("INSERT INTO `user_security` VALUES (null,?,?,?,?,now(),now())")
	} else {
		stmt, err = r.db.Prepare("INSERT INTO `user_security` VALUES (null,?,?,?,?,now(),now())")
	}

	if err != nil {
		return fmt.Errorf("Error in insert user_security prepared statement: %s", err.Error())
	}

	var hash string
	var lp int64

	if u.Provider != nil {
		lp = u.Provider.ID
	}
	if u.Hash != "" {
		hash = u.Hash
	}

	res, err := stmt.Exec(id, lp, hash, u.LastMachine)
	defer stmt.Close()

	if err != nil {
		return fmt.Errorf("Error in insert user_security for user id: %d", id, err.Error())
	}
	u.ID, _ = res.LastInsertId()

	return nil
}

// UpdateUserSecurity Updates the given user in the user_security table
func (r *Client) UpdateUserSecurity(u *models.UserSecurity) error {

	var err error
	var stmt *sql.Stmt

	if r.tx != nil {
		stmt, err = r.tx.Prepare("UPDATE `user_security` SET fk_login_provider=?,hash=?,updated_at=now() WHERE id=?")
	} else {
		stmt, err = r.db.Prepare("UPDATE `user_security` SET fk_login_provider=?,hash=?,updated_at=now() WHERE id=?")
	}

	if err != nil {
		return fmt.Errorf("Error in update user security prepared statement: %s", err.Error())
	}

	_, err = stmt.Exec(u.Provider.ID, u.Hash, u.ID)
	defer stmt.Close()

	if err != nil {
		return fmt.Errorf("Error in update security for userID %d : %s", u.ID, err.Error())
	}

	return nil
}

// UpdateLoginData Updates the Data for the login stats
func (r *Client) UpdateLoginData(u *models.UserSecurity) error {

	stmt, err := r.db.Prepare("UPDATE `user_security` SET last_login_date=now(),last_machine=? WHERE id=?")
	if err != nil {
		return fmt.Errorf("Error in update user security prepared statement: %s", err.Error())
	}

	_, err = stmt.Exec(u.LastMachine, u.ID)
	defer stmt.Close()

	if err != nil {
		return fmt.Errorf("Error in update security for userID %d : %s", u.ID, err.Error())
	}

	return nil
}

// GetSecurityInfoByUserId Get the UserSecurity by a given User id
func (r *Client) GetSecurityInfoByUserId(id int64) (*models.UserSecurity, error) {
	var found bool
	var fkProvider int64
	resp := &models.UserSecurity{Provider: &models.LoginProvider{}}

	err := r.db.QueryRow("SELECT IF(COUNT(*),'true','false') FROM user_security WHERE fk_user=?", id).Scan(&found)
	if err != nil {
		return resp, err
	}

	if !found {
		return resp, fmt.Errorf("UserSecurity for user with id %d not found", id)
	}

	err = r.db.QueryRow("SELECT id,fk_login_provider,hash,last_machine,last_login_date,updated_at FROM user_security WHERE fk_user=?", id).
		Scan(&resp.ID, &fkProvider, &resp.Hash, &resp.LastMachine, &resp.LastLogin, &resp.UpdatedAt)
	if err != nil {
		return resp, err
	}

	resp.Provider, err = r.GetLoginProviderById(fkProvider)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// GetInterestById Get an Interest by a given id
func (r *Client) GetInterestById(id int64) (*models.Interest, error) {
	var found bool
	resp := &models.Interest{}

	err := r.db.QueryRow("SELECT IF(COUNT(*),'true','false') FROM interest WHERE id=?", id).Scan(&found)
	if err != nil {
		return resp, err
	}

	if !found {
		return resp, fmt.Errorf("Interest with id %d not found", id)
	}

	err = r.db.QueryRow("SELECT id, name FROM interest WHERE id=?", id).Scan(&resp.ID, &resp.Name)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// GetAllInterests Gets all interests
func (r *Client) GetAllInterests() ([]*models.Interest, error) {

	var resp []*models.Interest

	rows, err := r.db.Query("SELECT id, name from interest")
	if err != nil {
		return resp, err
	}

	for rows.Next() {
		var n = new(models.Interest)

		err = rows.Scan(&n.ID, &n.Name)
		if err != nil {
			defer rows.Close()
			return resp, fmt.Errorf("Error reading rows: %s", err.Error())
		}

		resp = append(resp, n)
	}

	rows.Close()
	if len(resp) == 0 {
		return resp, fmt.Errorf("No Interests found")
	}

	return resp, nil
}

// UpdateUserInterests Udaptes All User interests
func (r *Client) UpdateUserInterests(interests []*models.Interest, id int64) error {
	var err error
	var stmtd *sql.Stmt
	var stmti *sql.Stmt

	if r.tx != nil {
		stmtd, err = r.tx.Prepare("DELETE FROM `users_profile_interests` WHERE fk_user_profile=?")
		stmti, err = r.tx.Prepare("INSERT INTO `users_profile_interests` VALUES (?,?)")
	} else {
		stmtd, err = r.db.Prepare("DELETE FROM `users_profile_interests` WHERE fk_user_profile=?")
		stmti, err = r.db.Prepare("INSERT INTO `users_profile_interests` VALUES (?,?)")
	}

	if err != nil {
		return fmt.Errorf("Error in deleting user interests prepared statement: %s", err.Error())
	}

	_, err = stmtd.Exec(id)
	defer stmtd.Close()
	if err != nil {
		return fmt.Errorf("Error in deleting user interests for userID %d : %s", id, err.Error())
	}

	if interests != nil && len(interests) > 0 {
		for _, i := range interests {
			_, err = stmti.Exec(i.ID, id)
			defer stmti.Close()
			if err != nil {
				return fmt.Errorf("Error in inserting user interests for userID %d : %s", id, err.Error())
			}
		}
		stmti.Close()
	}

	return nil
}

// GetAllInterestByUserProfileId Gets all interests of a given user
func (r *Client) GetAllInterestByUserProfileId(id int64) ([]*models.Interest, error) {

	var resp []*models.Interest

	rows, err := r.db.Query("SELECT it.id, it.name FROM users_profile_interests upi INNER JOIN interest it on upi.fk_interest=it.id WHERE upi.fk_user_profile=?", id)
	if err != nil {
		return resp, err
	}

	for rows.Next() {
		var n = new(models.Interest)

		err = rows.Scan(&n.ID, &n.Name)
		if err != nil {
			defer rows.Close()
			return resp, fmt.Errorf("Error reading rows: %s", err.Error())
		}

		resp = append(resp, n)
	}

	rows.Close()

	return resp, nil
}

// GetAllLanguage Gets all languages
func (r *Client) GetAllLanguage() ([]*models.Language, error) {

	var resp []*models.Language

	rows, err := r.db.Query("SELECT id, name, name_iso2, name_iso3 from language")
	if err != nil {
		return resp, err
	}

	for rows.Next() {
		var n = new(models.Language)

		err = rows.Scan(&n.ID, &n.Name, &n.NameIso2, &n.NameIso3)
		if err != nil {
			defer rows.Close()
			return resp, fmt.Errorf("Error reading rows: %s", err.Error())
		}

		resp = append(resp, n)
	}

	rows.Close()
	if len(resp) == 0 {
		return resp, fmt.Errorf("No Languages found")
	}

	return resp, nil
}

// GetLanguageById Gets a Language by its Id
func (r *Client) GetLanguageById(id int64) (*models.Language, error) {

	var found bool
	resp := &models.Language{}

	err := r.db.QueryRow("SELECT IF(COUNT(*),'true','false') FROM language WHERE id=?", id).Scan(&found)
	if err != nil {
		return resp, err
	}

	if !found {
		return resp, fmt.Errorf("Language with id %d not found", id)
	}

	err = r.db.QueryRow("SELECT id, name, name_iso2, name_iso3 FROM language WHERE id=?", id).Scan(&resp.ID, &resp.Name, &resp.NameIso2, &resp.NameIso3)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// GetAllLoginProvider Gets all login providers
func (r *Client) GetAllLoginProvider() ([]*models.LoginProvider, error) {

	var resp []*models.LoginProvider

	rows, err := r.db.Query("SELECT id,name,web_clientid,web_secret,android_clientid,android_secret,iphone_clientid,iphone_secret,updated_at from language")
	if err != nil {
		return resp, err
	}

	for rows.Next() {
		var n = new(models.LoginProvider)

		err = rows.Scan(&n.ID, &n.Name, &n.WebClientid, &n.WebSecret, &n.AndroidClientid, &n.AndroidSecret, &n.IphoneClientid, &n.IphoneSecret, &n.UpdatedAt)
		if err != nil {
			return resp, fmt.Errorf("Error reading rows: %s", err.Error())
		}

		resp = append(resp, n)
	}

	rows.Close()
	if len(resp) == 0 {
		return resp, fmt.Errorf("No Login providers found")
	}

	return resp, nil
}

// GetLoginProviderById Gets a LoginProvider by its Id
func (r *Client) GetLoginProviderById(id int64) (*models.LoginProvider, error) {

	var found bool
	resp := &models.LoginProvider{}

	err := r.db.QueryRow("SELECT IF(COUNT(*),'true','false') FROM login_provider WHERE id=?", id).Scan(&found)
	if err != nil {
		return resp, err
	}

	if !found {
		return resp, fmt.Errorf("Login Provider with id %d not found", id)
	}

	err = r.db.QueryRow("SELECT id,name,web_clientid,web_secret,android_clientid,android_secret,iphone_clientid,iphone_secret,updated_at FROM login_provider WHERE id=?", id).
		Scan(&resp.ID, &resp.Name, &resp.WebClientid, &resp.WebSecret, &resp.AndroidClientid, &resp.AndroidSecret, &resp.IphoneClientid, &resp.IphoneSecret, &resp.UpdatedAt)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// InsertEvent Inserts and event into event table
func (r *Client) InsertEvent(ev *models.Event) error {

	stmt, err := r.db.Prepare("INSERT INTO `event` VALUES (null,now(),?,?,?,?,?,?,?)")
	if err != nil {
		return fmt.Errorf("Error in insert event prepared statement: %s", err.Error())
	}

	res, err := stmt.Exec(ev.StartDate, ev.EndDate, ev.Location, ev.Latitude, ev.Longitude, ev.Active, ev.CreatedBy.ID)
	defer stmt.Close()

	if err != nil {
		return fmt.Errorf("Error in insert event for user id %d", ev.CreatedBy.ID, err.Error())
	}

	ev.ID, _ = res.LastInsertId()

	return nil
}

// UpdateEvent Update the given event in event table
func (r *Client) UpdateEvent(ev *models.Event) error {

	stmt, err := r.tx.Prepare("UPDATE `event` SET location=?,latitude=?,longitude=?,start_datetime=?,end_datetime=?,active=? WHERE id=?")
	if err != nil {
		return fmt.Errorf("Error in update user prepared statement: %s", err.Error())
	}

	_, err = stmt.Exec(ev.Location, ev.Latitude, ev.Longitude, ev.StartDate, ev.EndDate, ev.StartDate, ev.Active, ev.ID)
	defer stmt.Close()

	if err != nil {
		return fmt.Errorf("Could not update eventID %d : %s", ev.ID, err.Error())
	}

	return nil
}

// GetEventById Gets an event by a given id
func (r *Client) GetEventById(id int64) (*models.Event, error) {

	var found bool
	var fkCreatedBy int64
	ev := &models.Event{}

	err := r.db.QueryRow("SELECT IF(COUNT(*),'true','false') FROM event WHERE id=?", id).Scan(&found)
	if err != nil {
		return ev, err
	}

	if !found {
		return ev, fmt.Errorf("Event with id %d not found", id)
	}

	err = r.db.QueryRow("SELECT id,created_at,start_datetime,end_datetime,location,latitude,longitude,active,fk_created_by FROM event WHERE id=?", id).
		Scan(&ev.ID, &ev.CreatedAt, &ev.StartDate, &ev.EndDate, &ev.Location, &ev.Latitude, &ev.Longitude, &ev.Active, &fkCreatedBy)
	if err != nil {
		return ev, err
	}

	ev.CreatedBy, err = r.GetUserById(fkCreatedBy)
	if err != nil {
		return &models.Event{}, err
	}

	// Get the users in the event
	ev.Users = []*models.User{}

	rows, err := r.db.Query("SELECT u.id,u.email,u.name,u.created_at,u.updated_at,u.active,up.ID as pid,up.age_range,up.sex,up.updated_at as udate,la.ID as lid,la.name as lname,la.name_iso2 as lname2,la.name_iso3 as lname3 FROM event_users as eu INNER JOIN user as u on eu.fk_user=u.id LEFT JOIN user_profile up on u.ID=up.fk_user LEFT JOIN language la on up.fk_language=la.ID WHERE eu.fk_event=?", id)
	if err != nil {
		// there are no users at the events
		return ev, nil
	}

	// fill in users and interests
	for rows.Next() {
		var u = new(models.User)
		var p = new(models.UserProfile)
		var l = new(models.Language)

		err = rows.Scan(&u.ID, &u.Email, &u.Name, &u.CreatedAt, &u.UpdatedAt, &u.Active, &p.ID, &p.AgeRange, &p.Sex, &p.UpdatedAt, &l.ID, &l.Name, &l.NameIso2, &l.NameIso3)
		if err != nil {
			defer rows.Close()
			return ev, fmt.Errorf("Error reading rows: %s", err.Error())
		}
		p.Language = l
		u.Profile = p

		if p.Interests, err = r.GetAllInterestByUserProfileId(u.ID); err != nil {
			defer rows.Close()
			return ev, err
		}

		ev.Users = append(ev.Users, u)
	}

	rows.Close()

	return ev, nil

}

// GetUserEventsByUserId Gets the events of a given user id
func (r *Client) GetUserEventsByUserId(id int64) ([]*models.Event, error) {

	var evs []*models.Event

	rows, err := r.db.Query("SELECT id,created_at,start_datetime,end_datetime,location,latitude,longitude,active FROM event WHERE fk_created_by=?", id)
	if err != nil {
		return evs, err
	}

	// fill in users and interests
	for rows.Next() {
		var ev = new(models.Event)

		err = rows.Scan(&ev.ID, &ev.CreatedAt, &ev.StartDate, &ev.EndDate, &ev.Location, &ev.Latitude, &ev.Longitude, &ev.Active)
		if err != nil {
			defer rows.Close()
			return evs, fmt.Errorf("Error reading rows: %s", err.Error())
		}

		evs = append(evs, ev)
	}

	rows.Close()

	return evs, nil
}

// FindEventById Check if the event exists and its active
func (r *Client) FindEventById(id int64) (bool, error) {

	var found bool
	err := r.db.QueryRow("SELECT IF(COUNT(*),'true','false') FROM event WHERE id=? and active=true", id).Scan(&found)
	if err != nil {
		return false, err
	}

	return found, nil
}

// AddUserToEvent Adds a user to an event
func (r *Client) AddUserToEvent(idEvent int64, idUser int64) error {

	var found bool
	err := r.db.QueryRow("SELECT IF(COUNT(*),'true','false') FROM event WHERE id=? and fk_created_by=?", idEvent, idUser).Scan(&found)
	if err != nil {
		return err
	}
	if found {
		return fmt.Errorf("%d", serror.ErrorCantAddUSerToEvent)
	}

	stmt, err := r.db.Prepare("INSERT INTO `event_users` VALUES (?,?)")
	if err != nil {
		return fmt.Errorf("Error in adding user to event prepared statement: %s", err.Error())
	}

	_, err = stmt.Exec(idEvent, idUser)
	defer stmt.Close()

	if err != nil {
		return fmt.Errorf("Error adding user %d to event %d - %s", idUser, idEvent, err.Error())
	}

	return nil
}

// RemoveUserFromEvent Removes a user from an event
func (r *Client) RemoveUserFromEvent(idEvent int64, idUser int64) error {

	stmt, err := r.db.Prepare("DELETE FROM `event_users` WHERE fk_event=? AND fk_user=?")
	if err != nil {
		return fmt.Errorf("Error in removing user from event prepared statement: %s", err.Error())
	}

	_, err = stmt.Exec(idEvent, idUser)
	defer stmt.Close()

	if err != nil {
		return fmt.Errorf("Error removing user %d from event %d - %s", idUser, idEvent, err.Error())
	}

	return nil
}

// Health Endpoint of the Client
func (r *Client) Health() error {

	str, err := r.buildStringConnection()
	if err != nil {
		return err
	}

	db, err := sql.Open("mysql", str)
	if err != nil {
		return err
	}

	db.Close()
	return nil
}

// buildStringConnection builds the string connection to connect to the mysql server
func (r *Client) buildStringConnection() (string, error) {
	// [username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
	if r.config == nil {
		return "", fmt.Errorf("Client configuration not loaded")
	}
	if r.config.User == "" {
		return "", fmt.Errorf(IsEmpty, "User")
	}
	if r.config.Pw == "" {
		return "", fmt.Errorf(IsEmpty, "Password")
	}
	if r.config.Host == "" {
		return "", fmt.Errorf(IsEmpty, "Host")
	}
	if r.config.Port <= 0 {
		return "", fmt.Errorf(IsEmpty, "Port")
	}
	if r.config.Schema == "" {
		return "", fmt.Errorf(IsEmpty, "Schema")
	}

	stringConn := r.config.User + ":" + r.config.Pw
	stringConn += "@tcp(" + r.config.Host + ":" + strconv.Itoa(r.config.Port) + ")"
	stringConn += "/" + r.config.Schema + "?charset=utf8&parseTime=True"

	return stringConn, nil
}

//deferRollback default defer to rollback transactions on error
func (r *Client) deferRollback() {
	if r.tx != nil {
		r.tx.Rollback()
		r.tx = nil
	}
}

//commit commits and set the transaction to nil
func (r *Client) commit() {
	if r.tx != nil {
		r.tx.Commit()
		r.tx = nil
	}
}
