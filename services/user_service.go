package services

import (
	"fmt"
	"net"
	"time"

	"github.com/go-pg/pg/v9"
	"github.com/matiss/go-graphql-server/models"
	"github.com/matiss/go-graphql-server/utils"
)

type UserService struct {
	pgdb *pg.DB
}

func NewUserService(pgdb *pg.DB) *UserService {
	return &UserService{
		pgdb: pgdb,
	}
}

// Find user by ID
func (u *UserService) Find(ID int32) (*models.User, error) {
	user := &models.User{}

	_, err := u.pgdb.QueryOne(user, "SELECT * FROM users WHERE id = ?", ID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// FindByEmail finds user by email
func (u *UserService) FindByEmail(email string) (*models.User, error) {
	user := &models.User{}

	_, err := u.pgdb.QueryOne(user, "SELECT * FROM users WHERE email = ?", email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserService) List(limit int, offset int, userRole int) ([]*models.User, error) {
	fetchSize := limit

	// Limit fetch size
	if fetchSize == 0 {
		fetchSize = defaultListFetchSize
	} else if fetchSize > maxListFetchSize {
		fetchSize = maxListFetchSize
	}

	users := make([]*models.User, 0)

	if userRole > 0 {
		// Fetch users by role
		_, err := u.pgdb.Query(&users, `SELECT * FROM users WHERE role = ? ORDER BY created_at DESC LIMIT ? OFFSET ?;`, userRole, fetchSize, offset)
		if err != nil {
			return nil, err
		}
	} else {
		// Fetch all users
		_, err := u.pgdb.Query(&users, `SELECT * FROM users ORDER BY created_at DESC LIMIT ? OFFSET ?;`, fetchSize, offset)
		if err != nil {
			return nil, err
		}
	}

	return users, nil
}

func (u *UserService) Create(user *models.User) (*models.User, error) {
	// Validate email address
	if valid := user.ValidateEmail(); !valid {
		return nil, fmt.Errorf("Invalid email address")
	}

	// Validate password
	if valid := user.ValidatePassword(); !valid {
		return nil, fmt.Errorf("Invalid password")
	}

	// Generate password hash
	user.HashedPassword()

	// Set role
	user.Role = int(utils.AuthUser)

	// Make sure ID is clear
	user.ID = 0

	// Update timesstamps
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now
	user.LoginTime = now

	// Save user
	err := u.pgdb.Insert(user)
	if err != nil {
		return nil, err
	}

	// Remove password hash
	user.Password = ""

	return user, nil
}

func (u *UserService) ComparePassword(email string, password string) (*models.User, error) {
	user, err := u.FindByEmail(email)
	if err != nil || !user.Active() {
		return nil, fmt.Errorf("Could not find user")
	}

	if result := user.ComparePassword(password); !result {
		return nil, fmt.Errorf("Password did not match!")
	}

	return user, nil
}

func (u *UserService) UpdateLogin(user *models.User, ipText string) error {
	ip := net.ParseIP(ipText)

	now := time.Now()
	user.LoginIP = ip
	user.LoginTime = now
	user.UpdatedAt = now

	_, err := u.pgdb.QueryOne(user, `
		UPDATE users
		SET login_ip = ?login_ip, login_count = login_count + 1, login_at = ?login_at, updated_at = ?updated_at
		WHERE id = ?id RETURNING login_count`,
		user)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserService) Count(role int) (int, error) {
	var count int

	if role > 0 {
		// Count users by role
		_, err := u.pgdb.Query(pg.Scan(&count), "SELECT count(*) FROM users WHERE role = ?;", role)
		if err != nil {
			return 0, err
		}
	} else {
		// Count all users
		_, err := u.pgdb.Query(pg.Scan(&count), "SELECT count(*) FROM users;")
		if err != nil {
			return 0, err
		}
	}

	return count, nil
}

// Performs soft delete
func (u *UserService) Delete(email string) error {
	_, err := u.pgdb.Exec(`
		UPDATE users
		SET deleted_at = ? WHERE email = ?`, time.Now(), email)
	if err != nil {
		return err
	}

	return nil
}
