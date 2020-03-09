package services

import (
	"fmt"
	"net"
	"time"

	"github.com/go-pg/pg/v9"
	"github.com/matiss/go-graphql-server/models"
	"github.com/matiss/go-graphql-server/utils"
)

const (
	defaultListFetchSize = 10
)

type UserService struct {
	pgdb *pg.DB
}

func NewUserService(pgdb *pg.DB) *UserService {
	return &UserService{
		pgdb: pgdb,
	}
}

func (u *UserService) FindByEmail(email string) (*models.User, error) {
	user := &models.User{}

	_, err := u.pgdb.QueryOne(user, "SELECT * FROM users WHERE email = ?", email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserService) List(first *int32) ([]*models.User, error) {
	var fetchSize int32
	if first == nil || *first == 0 {
		fetchSize = defaultListFetchSize
	} else {
		fetchSize = *first
	}

	users := make([]*models.User, 0)

	_, err := u.pgdb.Query(&users, `SELECT * FROM users ORDER BY created_at DESC LIMIT ?;`, fetchSize)
	if err != nil {
		return nil, err
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

func (u *UserService) Count() (int, error) {
	var count int

	_, err := u.pgdb.Query(pg.Scan(&count), "SELECT count(*) FROM users")
	if err != nil {
		return 0, err
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