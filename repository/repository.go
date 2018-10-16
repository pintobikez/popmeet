package repository

import models "github.com/pintobikez/popmeet/api/structures"

type Repository interface {
	Connect() error
	Disconnect()
	Health() error
	// Interest methods
	FindInterestById(id uint32) (*models.Interest, error)
	GetAllInterestByUserId(id uint32) ([]*models.Interest, error)
	GetAllInterests() ([]*models.Interest, error)
	// User methods
	FindUserById(id uint32) (*models.User, error)
	FindUserProfileByUserId(id uint32) (*models.UserProfile, error)
	// Languages
	FindLanguageById(uint32 id) (*models.Language, error)
	GetAllLanguage() ([]*models.Language, error)
	// UserSecurity
	FindSecurityInfoByUserId(uint32 id) (*models.UserSecurity, error)
	// LoginProvider
	FindLoginProviderById(uint32 id) (*models.LoginProvider, error)
	GetAllLoginProvider() ([]*models.LoginProvider, error)
}
