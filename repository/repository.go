package repository

import "github.com/pintobikez/popmeet/api/models"

type Repository interface {
	Connect() error
	Disconnect()
	Health() error
	// User - Interests
	UpdateUserInterests(interests []*models.Interest, id int64) error
	GetAllInterestByUserProfileId(id int64) ([]*models.Interest, error)
	// Interest methods
	GetInterestById(id int64) (*models.Interest, error)
	GetAllInterests() ([]*models.Interest, error)
	// User methods
	InsertUser(u *models.User) error
	UpdateUser(u *models.User) error
	GetUserById(id int64) (*models.User, error)
	FindUserById(id int64) (bool, error)
	FindUserByEmail(email string) (bool, error)
	GetUserByEmail(email string) (*models.User, error)
	// User profile methods
	InsertUserProfile(u *models.UserProfile, id int64) error
	UpdateUserProfile(u *models.UserProfile) error
	GetUserProfileByUserId(id int64) (*models.UserProfile, error)
	// Languages
	GetLanguageById(id int64) (*models.Language, error)
	GetAllLanguage() ([]*models.Language, error)
	// UserSecurity
	InsertUserSecurity(u *models.UserSecurity, id int64) error
	UpdateUserSecurity(u *models.UserSecurity) error
	GetSecurityInfoByUserId(id int64) (*models.UserSecurity, error)
	// LoginProvider
	GetLoginProviderById(id int64) (*models.LoginProvider, error)
	GetAllLoginProvider() ([]*models.LoginProvider, error)
	//User login updates
	UpdateLoginData(u *models.UserSecurity) error
	// Event methods
	AddUserToEvent(idEvent int64, idUser int64) error
	RemoveUserFromEvent(idEvent int64, idUser int64) error
	InsertEvent(u *models.Event) error
	UpdateEvent(u *models.Event) error
	FindEventById(id int64) (bool, error)
	GetEventById(id int64) (*models.Event, error)
	GetUserEventsByUserId(id int64) ([]*models.Event, error)
}
