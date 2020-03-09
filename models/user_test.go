package models

import (
	"testing"
)

func TestHashedPassword(t *testing.T) {
	user := User{
		Password: "secret123456",
	}

	user.HashedPassword()

	if user.Password == "secret123456" {
		t.Errorf("Hasing failed, got: %s", user.Password)
	}
}

func TestComparePassword(t *testing.T) {
	user := User{
		Password: "$2a$10$MZiaf/ozEVthSFNR5vtwyu3cUKPG5Hmid7sg6GWFQHpVQufTZovKS",
	}

	result := user.ComparePassword("secret123456")

	if !result {
		t.Errorf("Password compare failed")
	}
}
