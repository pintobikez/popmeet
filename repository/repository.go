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
	FindInterestById(id int64) (*models.Interest, error)
	GetAllInterests() ([]*models.Interest, error)
	// User methods
	InsertUser(u *models.User) error
	UpdateUser(u *models.User) error
	FindUserById(id int64) (*models.User, error)
	FindUserByEmail(email string) (*models.User, error)
	// User profile methods
	InsertUserProfile(u *models.UserProfile, id int64) error
	UpdateUserProfile(u *models.UserProfile) error
	FindUserProfileByUserId(id int64) (*models.UserProfile, error)
	// Languages
	FindLanguageById(id int64) (*models.Language, error)
	GetAllLanguage() ([]*models.Language, error)
	// UserSecurity
	InsertUserSecurity(u *models.UserSecurity, id int64) error
	UpdateUserSecurity(u *models.UserSecurity) error
	FindSecurityInfoByUserId(id int64) (*models.UserSecurity, error)
	// LoginProvider
	FindLoginProviderById(id int64) (*models.LoginProvider, error)
	GetAllLoginProvider() ([]*models.LoginProvider, error)
	//
	UpdateLoginData(u *models.UserSecurity) error
}
