package models

type Interest struct {
	ID   uint32 `json:"id" validate:"required,numeric"`
	Name string `json:"name" validate:"required,alpha,min=1,max=255"`
}

type Language struct {
	ID       uint32 `json:"id" validate:"required,numeric"`
	Name     string `json:"name" validate:"required,alpha,min=1,max=40"`
	NameIso2 string `json:"name_iso2" validate:"required,alpha,len=2"`
	NameIso3 string `json:"name_iso3" validate:"required,alpha,len=3"`
}

type Event struct {
	ID        uint32  `json:"id" validate:"required,numeric"`
	CreatedAt string  `json:"created_at"`
	StartDate string  `json:"start_date"`
	EndDate   string  `json:"end_date"`
	Location  string  `json:"name" validate:"required,alpha,min=1,max=255"`
	Active    bool    `json:"active" validate:"required"`
	CreatedBy *User   `json:"created_by"`
	Users     []*User `json:"users"`
}

type EventUsers struct {
	EventID   uint32 `json:"id" validate:"required,numeric"`
	CreatedBy *User  `json:"user"`
}

type User struct {
	ID        int8          `json:"id" validate:"required,numeric"`
	Email     string        `json:"email" validate:"required,email"`
	Name      string        `json:"name" validate:"required,alpha,min=1,max=255"`
	CreatedAt string        `json:"created_at"`
	UpdatedAt string        `json:"updated_at"`
	Active    bool          `json:"active" validate:"required"`
	Profile   *UserProfile  `json:"profile"`
	Security  *UserSecurity `json:"security"`
}

type UserProfile struct {
	ID        uint32      `json:"id" validate:"required,numeric"`
	Language  *Language   `json:"language"`
	Sex       string      `json:"sex" validate:"required,male|female"`
	AgeRange  string      `json:"age_range" validate:"required,18-25|26-32|33-39|40-46|47-53|54-60|61-70|+70"`
	UpdatedAt string      `json:"updated_at"`
	Interests []*Interest `json:"interests"`
}

type LoginProvider struct {
	ID              uint32 `json:"id"`
	Name            string `json:"name"`
	WebClientid     string `json:"web_clientid"`
	WebSecret       string `json:"web_secret"`
	AndroidClientid string `json:"android_clientid"`
	AndroidSecret   string `json:"android_secret"`
	IphoneClientid  string `json:"iphone_clientid"`
	IphoneSecret    string `json:"iphone_secret"`
	UpdatedAt       string `json:"updated_at"`
}

type UserSecurity struct {
	ID             uint32         `json:"id" validate:"required,numeric"`
	Provider       *LoginProvider `json:"login_provider,omitempty"`
	HashedPassword string         `json:"password" validate:"required"`
	Salt           string         `json:"salt"`
	LastMachine    string         `json:"last_machine"`
	Token          string         `json:"token"`
	LastLogin      string         `json:"last_login"`
	UpdatedAt      string         `json:"updated_at"`
}
