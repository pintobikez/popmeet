package models

import "time"

type EventSearch struct {
	Location    string      `json:"location" validate:"required,excludesall=!@#?,min=1,max=255"`
	Longitude   float64     `json:"longitude" validate:"required,numeric"`
	Latitude    float64     `json:"latitude" validate:"required,numeric"`
	SearchRange int64       `json:"range" validate:"required,numeric"`
	StartDate   time.Time   `json:"start_date" validate:"omitempty,required"`
	EndDate     time.Time   `json:"end_date" validate:"omitempty,required,gtfield=StartDate"`
	Interests   []*Interest `json:"interests,omitempty" validate:"omitempty,required,dive"`
	Sex         string      `json:"sex,omitempty" validate:"omitempty,required,oneof=male female"`
	AgeRange    string      `json:"age_range,omitempty" validate:"omitempty,required,oneof=18-25 26-32 33-39 40-46 47-53 54-60 61-70 +70"`
}

type LoginUser struct {
	Email    string `json:"email" validate:"required,email"`
	Provider int64  `json:"login_provider" validate:"required,numeric"`
	Password string `json:"password" validate:"omitempty,required"`
}

type NewUser struct {
	Email    string `json:"email" validate:"required,email"`
	Name     string `json:"name" validate:"required,excludesall=!@#?,min=1,max=255"`
	Provider int64  `json:"login_provider" validate:"required,numeric"`
	Password string `json:"password,omitempty" validate:"omitempty,required"`
}

type NewEvent struct {
	StartDate time.Time `json:"start_date" validate:"required"`
	EndDate   time.Time `json:"end_date" validate:"required,gtfield=StartDate"`
	Location  string    `json:"location" validate:"required,excludesall=!@#?,min=1,max=255"`
	Longitude float64   `json:"longitude" validate:"required,numeric"`
	Latitude  float64   `json:"latitude" validate:"required,numeric"`
	Active    bool      `json:"active" validate:"required"`
	CreatedBy int64     `json:"created_by validate:"required,numeric"`
}

type Interest struct {
	ID   int64  `json:"id" validate:"required,numeric"`
	Name string `json:"name,omitempty" validate:"omitempty,required,alpha,min=1,max=255"`
}

type Language struct {
	ID       int64  `json:"id" validate:"required,numeric"`
	Name     string `json:"name,omitempty" validate:"omitempty,required,alpha,min=1,max=40"`
	NameIso2 string `json:"name_iso2,omitempty" validate:"omitempty,required,alpha,len=2"`
	NameIso3 string `json:"name_iso3,omitempty" validate:"omitempty,required,alpha,len=3"`
}

type Event struct {
	ID        int64     `json:"id" validate:"required,numeric"`
	CreatedAt time.Time `json:"created_at"`
	StartDate time.Time `json:"start_date" validate:"required"`
	EndDate   time.Time `json:"end_date" validate:"required,gtfield=StartDate"`
	Location  string    `json:"location" validate:"required,excludesall=!@#?,min=1,max=255"`
	Longitude float64   `json:"longitude" validate:"required,numeric"`
	Latitude  float64   `json:"latitude" validate:"required,numeric"`
	Active    bool      `json:"active" validate:"required"`
	CreatedBy *User     `json:"created_by"`
	Users     []*User   `json:"users"`
}

type EventUsers struct {
	EventID   int64 `json:"id" validate:"required,numeric"`
	CreatedBy *User `json:"user"`
}

type User struct {
	ID        int64         `json:"id" validate:"required,numeric"`
	Email     string        `json:"email,omitempty" validate:"omitempty,required,email"`
	Name      string        `json:"name" validate:"required,excludesall=!@#?,min=1,max=255"`
	CreatedAt time.Time     `json:"created_at,omitempty"`
	UpdatedAt time.Time     `json:"updated_at,omitempty"`
	Active    bool          `json:"active,omitempty" validate:"omitempty,required"`
	Profile   *UserProfile  `json:"profile,omitempty" validate:"omitempty,required,dive"`
	Security  *UserSecurity `json:"security,omitempty" validate:"omitempty,required,dive"`
}

type UserProfile struct {
	ID        int64       `json:"id,omitempty" validate:"omitempty,required,numeric"`
	Language  *Language   `json:"language" validate:"required,dive"`
	Sex       string      `json:"sex" validate:"required,oneof=male female"`
	AgeRange  string      `json:"age_range" validate:"required,oneof=18-25 26-32 33-39 40-46 47-53 54-60 61-70 +70"`
	UpdatedAt time.Time   `json:"updated_at,omitempty"`
	Interests []*Interest `json:"interests,omitempty" validate:"omitempty,required,dive"`
}

type LoginProvider struct {
	ID              int64     `json:"id" validate:"required,numeric"`
	Name            string    `json:"name,omitempty"`
	WebClientid     string    `json:"web_clientid,omitempty"`
	WebSecret       string    `json:"web_secret,omitempty"`
	AndroidClientid string    `json:"android_clientid,omitempty"`
	AndroidSecret   string    `json:"android_secret,omitempty"`
	IphoneClientid  string    `json:"iphone_clientid,omitempty"`
	IphoneSecret    string    `json:"iphone_secret,omitempty"`
	UpdatedAt       time.Time `json:"updated_at,omitempty"`
}

type UserSecurity struct {
	ID          int64          `json:"id" validate:"required,numeric"`
	Provider    *LoginProvider `json:"login_provider" validate:"required,dive"`
	Hash        string         `json:"-"`
	LastMachine string         `json:"last_machine,omitempty"`
	LastLogin   time.Time      `json:"last_login,omitempty"`
	UpdatedAt   time.Time      `json:"updated_at,omitempty"`
}
