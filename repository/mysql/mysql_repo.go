package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	models "github.com/pintobikez/popmeet/api/structures"
	cnfs "github.com/pintobikez/popmeet/config/structures"
	"strconv"
)

const (
	IsEmpty = "%s is empty"
)

type Client struct {
	config *cnfs.DatabaseConfig
	db     *sql.DB
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

// FindUserById Find an User by a given id
func (r *Client) FindUserById(id uint32) (*models.User, error) {
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

// FindUserProfileByUserId Find the UserProfile by a given User id
func (r *Client) FindUserProfileByUserId(id uint32) (*models.UserProfile, error) {
	var found bool
	var fkLanguage uint32
	resp := &models.UserProfile{}

	err := r.db.QueryRow("SELECT IF(COUNT(*),'true','false') FROM user_profile WHERE fk_user=?", id).Scan(&found)
	if err != nil {
		return resp, err
	}

	if !found {
		return resp, fmt.Errorf("UserProfile for user with id %d not found", id)
	}

	err = r.db.QueryRow("SELECT id,fk_language,age_range,sex,updated_at FROM user WHERE fk_user=?", id).Scan(&resp.ID, &fkLanguage, &resp.AgeRange, &resp.Sex, &resp.UpdatedAt)
	if err != nil {
		return resp, err
	}

	resp.Language, err = r.FindLanguageById(resp.ID)
	if err != nil {
		return resp, err
	}

	resp.Interests, err = r.GetAllInterestByUserId(id)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// FindSecurityInfoByUserId Find the UserSecurity by a given User id
func (r *Client) FindSecurityInfoByUserId(id uint32) (*models.UserSecurity, error) {
	var found bool
	var fkProvider uint32
	resp := &models.UserSecurity{}

	err := r.db.QueryRow("SELECT IF(COUNT(*),'true','false') FROM user_security WHERE fk_user=?", id).Scan(&found)
	if err != nil {
		return resp, err
	}

	if !found {
		return resp, fmt.Errorf("UserSecurity for user with id %d not found", id)
	}

	err = r.db.QueryRow("SELECT id,fk_login_provider,hashed_password,salt,last_machine,token,last_login_date,updated_at FROM user_security WHERE fk_user=?", id).
		Scan(&resp.ID, &fkProvider, &resp.HashedPassword, &resp.Salt, &resp.LastMachine, &resp.Token, &resp.LastLogin, &resp.UpdatedAt)
	if err != nil {
		return resp, err
	}

	resp.Provider, err = r.FindLoginProviderById(fkProvider)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// FindInterestById Find an Interest by a given id
func (r *Client) FindInterestById(id uint32) (*models.Interest, error) {
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

// GetAllInterests Gets all interests
func (r *Client) GetAllInterestByUserId(id uint32) ([]*models.Interest, error) {

	var resp []*models.Interest

	rows, err := r.db.Query("SELECT int.id, int.name from users_profile_interests as upi inner join interest as int on upi.fk_interest=int.id WHERE upi.fk_user=?", id)
	if err != nil {
		return resp, err
	}

	for rows.Next() {
		var n = new(models.Interest)

		err = rows.Scan(&n.ID, &n.Name)
		if err != nil {
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

// FindLanguageById Gets a Language by its Id
func (r *Client) FindLanguageById(id uint32) (*models.Language, error) {

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

// FindLoginProviderById Gets a LoginProvider by its Id
func (r *Client) FindLoginProviderById(id uint32) (*models.LoginProvider, error) {

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
