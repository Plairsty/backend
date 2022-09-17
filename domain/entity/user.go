package entity

import (
	"html"
	"plairsty/backend/infrastructure/security"
	"regexp"
	"strings"
	"time"
)

const (
	emailRegex = `^[a-zA-Z0-9]+(?:\.[a-zA-Z0-9]+)*@[a-zA-Z0-9]+(?:\.[a-zA-Z0-9]+)*$`
)

type User struct {
	ID               uint64           `gorm:"primary_key;auto_increment" json:"id"`
	UserName         string           `gorm:"size:255;not null;unique" json:"user_name"`
	FirstName        string           `gorm:"size:100;not null;" json:"first_name"`
	MiddleName       string           `gorm:"size:100;not null;" json:"middle_name"`
	LastName         string           `gorm:"size:100;not null;" json:"last_name"`
	DateOfBirth      time.Time        `json:"date_of_birth"`
	Email            string           `gorm:"size:100;not null;unique" json:"email"`
	Phone            string           `gorm:"size:100;not null;" json:"phone"`  // Parent Phone
	Mobile           string           `gorm:"size:100;not null;" json:"mobile"` // Personal Phone
	HashedPassword   string           `gorm:"size:250;not null;" json:"password"`
	Salt             string           `gorm:"size:100;not null;" json:"salt"`
	Role             string           `gorm:"size:100;not null;default:'user'" json:"role"`
	Active           bool             `gorm:"not null;default:false" json:"active"` // If user haven't logged in for first time, this will be false
	LastLogin        *time.Time       `json:"last_login"`
	LoginIPAddresses []LoginIPAddress `gorm:"foreignKey:UserID" json:"login_ip_addresses"`
	CreatedAt        time.Time        `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt        time.Time        `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt        *time.Time       `json:"deleted_at,omitempty"`
}

type LoginIPAddress struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	UserID    uint64    `gorm:"not null" json:"user_id"` // Foreign Key (belongs to - User)
	IPAddress string    `gorm:"size:100;not null;" json:"ip_address"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
}

type PublicUser struct {
	ID        uint64 `gorm:"primary_key;auto_increment" json:"id"`
	UserName  string `gorm:"size:255;not null;unique" json:"user_name"`
	FirstName string `gorm:"size:100;not null;" json:"first_name"`
	LastName  string `gorm:"size:100;not null;" json:"last_name"`
}

type Users []User

func (user *User) BeforeSave() error {
	salt := security.GenerateRandomString(32)
	saltyHashedPassword, err := security.SaltedHash(user.HashedPassword, []byte(salt))
	if err != nil {
		return err
	}
	user.Salt = salt
	user.HashedPassword = string(saltyHashedPassword)
	return nil
}

// publicUsers So that we don't expose the user's email address and password to the world
func (users Users) publicUsers() []interface{} {
	result := make([]interface{}, len(users))
	for index, user := range users {
		result[index] = user.PublicUser()
	}
	return result
}

// PublicUser So that we don't expose the user's email address and password to the world
func (user *User) PublicUser() interface{} {
	return PublicUser{
		ID:        user.ID,
		UserName:  user.UserName,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}
}

func (user *User) Prepare() {
	user.UserName = html.EscapeString(strings.TrimSpace(user.UserName))
	user.FirstName = html.EscapeString(strings.TrimSpace(user.FirstName))
	user.MiddleName = html.EscapeString(strings.TrimSpace(user.MiddleName))
	user.LastName = html.EscapeString(strings.TrimSpace(user.LastName))
	user.Email = html.EscapeString(strings.TrimSpace(user.Email))
	user.Phone = html.EscapeString(strings.TrimSpace(user.Phone))
	user.Mobile = html.EscapeString(strings.TrimSpace(user.Mobile))
	user.Role = html.EscapeString(strings.TrimSpace(user.Role))
	user.Active = false
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
}

func (user *User) Validate(action string) map[string]string {
	var errorMessages = make(map[string]string)

	switch strings.ToLower(action) {
	default:
		if user.UserName == "" {
			errorMessages["Required_UserName"] = "Required UserName"
		}
		if user.FirstName == "" {
			errorMessages["Required_FirstName"] = "Required FirstName"
		}
		if user.LastName == "" {
			errorMessages["Required_LastName"] = "Required LastName"
		}
		if user.Email == "" {
			errorMessages["Required_Email"] = "Required Email"
		}
		if user.Phone == "" {
			errorMessages["Required_Phone"] = "Required Phone"
		}
		if user.Mobile == "" {
			errorMessages["Required_Mobile"] = "Required Mobile"
		}
		if user.HashedPassword == "" {
			errorMessages["Required_Password"] = "Required Password"
		}
		if user.Role == "" {
			errorMessages["Required_Role"] = "Required Role"
		}
		// Check email is valid
		if user.Email != "" {
			if isCorrectEmailFormat := ValidateEmail(user.Email); !isCorrectEmailFormat {
				errorMessages["invalid_email"] = "invalid email"
			}
		}

	case "login":
		if user.UserName == "" {
			errorMessages["Required_UserName"] = "Required UserName"
		}
		if user.HashedPassword == "" {
			errorMessages["Required_Password"] = "Required Password"
		}
		if user.Email != "" {
			if isCorrectEmailFormat := ValidateEmail(user.Email); !isCorrectEmailFormat {
				errorMessages["invalid_email"] = "invalid email"
			}
		}

	case "update":
		if user.Email == "" {
			errorMessages["email_required"] = "email required"
		}
		if user.Email != "" {
			if isCorrectEmailFormat := ValidateEmail(user.Email); !isCorrectEmailFormat {
				errorMessages["invalid_email"] = "invalid email"
			}
		}
	}
	return errorMessages
}

func ValidateEmail(email string) bool {
	if email == "" {
		return false
	}
	emailRegex := regexp.MustCompile(emailRegex)
	result, err := regexp.MatchString(emailRegex.String(), email)
	if result == false || err != nil {
		return false
	}
	return true
}
