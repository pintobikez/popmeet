package repository

import "github.com/pintobikez/popmeet/api/models"

type Repository interface {
	Connect() error
	Disconnect()
	Health() error
	// Interest methods
	FindInterestById(id int64) (*models.Interest, error)
	GetAllInterestByUserId(id int64) ([]*models.Interest, error)
	GetAllInterests() ([]*models.Interest, error)
	// User methods
	InsertUser(u *models.User) error
	FindUserById(id int64) (*models.User, error)
	FindUserProfileByUserId(id int64) (*models.UserProfile, error)
	// Languages
	FindLanguageById(id int64) (*models.Language, error)
	GetAllLanguage() ([]*models.Language, error)
	// UserSecurity
	InsertUserSecurity(u *models.UserSecurity, id int64) error
	FindSecurityInfoByUserId(id int64) (*models.UserSecurity, error)
	// LoginProvider
	FindLoginProviderById(id int64) (*models.LoginProvider, error)
	GetAllLoginProvider() ([]*models.LoginProvider, error)
}
