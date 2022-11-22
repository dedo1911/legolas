package users

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha1"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/registration"
)

const accountsFolder = "accounts"

type User struct {
	Email        string                 `json:"email"`
	Registration *registration.Resource `json:"registration"`
	PrivateKey   []byte                 `json:"key"`
	key          crypto.PrivateKey
}

func getEmailHash(email string) string {
	hasher := sha1.New()
	hasher.Write([]byte(email))
	return hex.EncodeToString(hasher.Sum(nil))
}

func NewUser(email string) (*User, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}

	return &User{
		Email: email,
		key:   privateKey,
	}, nil
}

func GetOrCreateUser(email string) (*User, error) {
	user, err := readUser(email)
	if err == nil {
		return user, nil
	}
	if !os.IsNotExist(err) {
		return nil, err
	}
	// Create new user
	newUser, err := NewUser(email)
	if err != nil {
		return nil, err
	}
	return newUser, nil
}

func (u *User) Register(client *lego.Client) error {
	if u.Registration != nil {
		return nil
	}
	reg, err := client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
	if err != nil {
		return err
	}
	u.Registration = reg
	// Store updated user data
	if err := u.writeUser(); err != nil {
		return err
	}
	return nil
}

func readUser(email string) (*User, error) {
	hash := getEmailHash(email)
	userPath := filepath.Join(accountsFolder, hash+".json")
	userData, err := os.ReadFile(userPath)
	if err != nil {
		return nil, err
	}
	var user User
	if err := json.Unmarshal(userData, &user); err != nil {
		return nil, err
	}
	privateKey, err := x509.ParseECPrivateKey(user.PrivateKey)
	if err != nil {
		return nil, err
	}
	user.key = privateKey
	return &user, nil
}

func (u *User) writeUser() error {
	os.MkdirAll(accountsFolder, 0600)
	privateKey, err := x509.MarshalECPrivateKey(u.key.(*ecdsa.PrivateKey))
	if err != nil {
		return err
	}
	u.PrivateKey = privateKey
	hash := getEmailHash(u.Email)
	userPath := filepath.Join(accountsFolder, hash+".json")
	userData, err := json.MarshalIndent(*u, "", "\t")
	if err != nil {
		return err
	}
	return os.WriteFile(userPath, userData, 0600)
}

func (u *User) GetEmail() string {
	return u.Email
}
func (u User) GetRegistration() *registration.Resource {
	return u.Registration
}
func (u *User) GetPrivateKey() crypto.PrivateKey {
	return u.key
}
