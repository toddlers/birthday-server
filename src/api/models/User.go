package models

import (
	"errors"
	"fmt"
	"html"
	"regexp"
	"strings"
	"time"

	"github.com/badoux/checkmail"
	"github.com/jinzhu/gorm"
)

// User with birthday
type User struct {
	ID        uint32    `gorm:"primary_key;auto_increment" json:"id"`
	Name      string    `gorm:"size:255;not null;unique" json:"name"`
	Birthday  string    `gorm:"size:255;not null" json:"birthday"`
	Email     string    `gorm:"size:100;not null;uniqe" json:"email"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (u *User) Prepare() {
	u.ID = 0
	u.Name = html.EscapeString(strings.TrimSpace(u.Name))
	u.Birthday = html.EscapeString(strings.TrimSpace(u.Birthday))
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
}
func (u *User) Validate(action string) error {
	switch strings.ToLower(action) {
	case "update":
		if u.Name == "" {
			return errors.New("Required Name")
		}
		if !u.IsUsernameCorrect() {
			return errors.New("Only alphanumeric usernames accepted")
		}
		if u.Birthday == "" {
			return errors.New("Required Birthday")
		}
		if !u.IsBirthdayCorrect() {
			return errors.New("Birthday should be before today")
		}
		if u.Email == "" {
			return errors.New("Required Email")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Invalid Email")
		}

		return nil
	default:
		if u.Name == "" {
			return errors.New("Required Name")
		}
		if !u.IsUsernameCorrect() {
			return errors.New("only alphanumeric usernames accepted")
		}
		if u.Birthday == "" {
			return errors.New("Required Birthday")
		}
		if !u.IsBirthdayCorrect() {
			return errors.New("Birthday should be before today")
		}
		if u.Email == "" {
			return errors.New("Required Email")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Invalid Email")
		}
		return nil
	}
}

func (u *User) SaveUser(db *gorm.DB) (*User, error) {
	err := db.Debug().Create(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func (u *User) FindAllUsers(db *gorm.DB) (*[]User, error) {
	users := []User{}
	err := db.Debug().Model(&User{}).Limit(10).Find(&users).Error
	if err != nil {
		return &[]User{}, err
	}
	return &users, err
}
func (u *User) FindUserByID(db *gorm.DB, uid uint32) (*User, error) {
	err := db.Debug().Model(&User{}).Where("id = ?", uid).Take(&u).Error
	if err != nil {
		return &User{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &User{}, errors.New("User Not Found")
	}
	return u, err
}

func (u *User) UpdateAUser(db *gorm.DB, uid uint32) (*User, error) {
	fmt.Println(u)
	db = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&User{}).UpdateColumns(
		map[string]interface{}{
			"name":       u.Name,
			"birthday":   u.Birthday,
			"email":      u.Email,
			"updated_at": time.Now(),
		},
	)

	if db.Error != nil {
		return &User{}, db.Error
	}
	err := db.Debug().Model(&User{}).Where("id = ?", uid).Take(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func (u *User) DeleteAUser(db *gorm.DB, uid uint32) (int64, error) {
	db = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&User{}).Delete(&User{})
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

func (u *User) CheckBirthday() map[string]string {
	days := u.CalculateDays(u.Birthday)
	birthdayMessage := make(map[string]string)
	if days == 0 {
		birthdayMessage["message"] = fmt.Sprintf("Hello %s !Happy Birthday", u.Name)
		return birthdayMessage
	}
	birthdayMessage["message"] = fmt.Sprintf("Hello %s! Your birthday is in %d days!", u.Name, days)
	return birthdayMessage
}

func (u *User) IsUsernameCorrect() bool {
	isAlpha := regexp.MustCompile(`^[A-Za-z]+$`).MatchString
	return isAlpha(u.Name)
}

func (u *User) IsBirthdayCorrect() bool {
	born, _ := time.Parse("2006-01-02", u.Birthday)
	today := time.Now()
	return born.Before(today)
}

func (u *User) CalculateDays(birthday string) int {
	//yyy-mm-dd
	born, _ := time.Parse("2006-01-02", birthday)
	today := time.Now()
	if (today.Month() == born.Month()) && (today.Day() == born.Day()) {
		return 0
	}
	bd := time.Date(today.Year(), born.Month(), born.Day(), 0, 0, 0, 0, time.UTC)
	if bd.Before(today) {
		bd = time.Date(today.Year()+1, born.Month(), born.Day(), 0, 0, 0, 0, time.UTC)
	}
	days := int(bd.Sub(today).Hours() / 24)
	return days
}
