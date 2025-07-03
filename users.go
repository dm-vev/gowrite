package gowrite

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type UsersService struct {
	Client *AppwriteClient
}

// User represents an Appwrite user.
type User struct {
	ID     string                 `json:"$id"`
	Email  string                 `json:"email"`
	Phone  string                 `json:"phone"`
	Name   string                 `json:"name"`
	Status bool                   `json:"status"`
	Labels []string               `json:"labels"`
	Mfa    bool                   `json:"mfa"`
	Data   map[string]interface{} `json:"-"`
}

// Preferences represent user preferences.
type Preferences map[string]interface{}

func NewUsers(client *AppwriteClient) *UsersService {
	return &UsersService{client}
}

func (u *User) UnmarshalJSON(b []byte) error {
	type Alias User
	aux := &struct{ *Alias }{Alias: (*Alias)(u)}
	if err := json.Unmarshal(b, &aux); err != nil {
		return err
	}
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	delete(raw, "$id")
	delete(raw, "email")
	delete(raw, "phone")
	delete(raw, "name")
	delete(raw, "status")
	delete(raw, "labels")
	delete(raw, "mfa")
	u.Data = raw
	return nil
}

// CreateUser creates a new user with a plain text password.
func (s *UsersService) CreateUser(userID, email, phone, password, name string) (*User, error) {
	payload := map[string]interface{}{
		"userId":   userID,
		"email":    email,
		"phone":    phone,
		"password": password,
		"name":     name,
	}
	resp, err := s.Client.sendRequest("POST", "/users", payload)
	if err != nil {
		return nil, err
	}
	var user User
	if err = json.Unmarshal(resp, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UsersService) CreateAnonymousUser(userID string) (*User, error) {
	payload := map[string]interface{}{
		"userId": userID,
	}
	resp, err := s.Client.sendRequest("POST", "/users", payload)
	if err != nil {
		return nil, err
	}
	var user User
	if err = json.Unmarshal(resp, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

// CreateArgon2User creates a user with Argon2 hashed password.
func (s *UsersService) CreateArgon2User(userID, email, password, name string) (*User, error) {
	payload := map[string]interface{}{
		"userId":   userID,
		"email":    email,
		"password": password,
		"name":     name,
	}
	resp, err := s.Client.sendRequest("POST", "/users/argon2", payload)
	if err != nil {
		return nil, err
	}
	var user User
	if err = json.Unmarshal(resp, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

// CreateBcryptUser creates a user with Bcrypt hashed password.
func (s *UsersService) CreateBcryptUser(userID, email, password, name string) (*User, error) {
	payload := map[string]interface{}{
		"userId":   userID,
		"email":    email,
		"password": password,
		"name":     name,
	}
	resp, err := s.Client.sendRequest("POST", "/users/bcrypt", payload)
	if err != nil {
		return nil, err
	}
	var user User
	if err = json.Unmarshal(resp, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

// CreateMD5User creates a user with MD5 hashed password.
func (s *UsersService) CreateMD5User(userID, email, password, name string) (*User, error) {
	payload := map[string]interface{}{
		"userId":   userID,
		"email":    email,
		"password": password,
		"name":     name,
	}
	resp, err := s.Client.sendRequest("POST", "/users/md5", payload)
	if err != nil {
		return nil, err
	}
	var user User
	if err = json.Unmarshal(resp, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

// CreatePHPassUser creates a user with PHPass hashed password.
func (s *UsersService) CreatePHPassUser(userID, email, password, name string) (*User, error) {
	payload := map[string]interface{}{
		"userId":   userID,
		"email":    email,
		"password": password,
		"name":     name,
	}
	resp, err := s.Client.sendRequest("POST", "/users/phpass", payload)
	if err != nil {
		return nil, err
	}
	var user User
	if err = json.Unmarshal(resp, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

// CreateScryptModifiedUser creates a user with Scrypt Modified hashed password.
func (s *UsersService) CreateScryptModifiedUser(userID, email, password, passwordSalt, passwordSaltSeparator, passwordSignerKey, name string) (*User, error) {
	payload := map[string]interface{}{
		"userId":                userID,
		"email":                 email,
		"password":              password,
		"passwordSalt":          passwordSalt,
		"passwordSaltSeparator": passwordSaltSeparator,
		"passwordSignerKey":     passwordSignerKey,
		"name":                  name,
	}
	resp, err := s.Client.sendRequest("POST", "/users/scrypt-modified", payload)
	if err != nil {
		return nil, err
	}
	var user User
	if err = json.Unmarshal(resp, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

// CreateScryptUser creates a user with Scrypt hashed password.
func (s *UsersService) CreateScryptUser(userID, email, password, passwordSalt string, passwordCpu, passwordMemory, passwordParallel, passwordLength int, name string) (*User, error) {
	payload := map[string]interface{}{
		"userId":           userID,
		"email":            email,
		"password":         password,
		"passwordSalt":     passwordSalt,
		"passwordCpu":      passwordCpu,
		"passwordMemory":   passwordMemory,
		"passwordParallel": passwordParallel,
		"passwordLength":   passwordLength,
		"name":             name,
	}
	resp, err := s.Client.sendRequest("POST", "/users/scrypt", payload)
	if err != nil {
		return nil, err
	}
	var user User
	if err = json.Unmarshal(resp, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

// CreateSHAUser creates a user with SHA hashed password.
func (s *UsersService) CreateSHAUser(userID, email, password, passwordVersion, name string) (*User, error) {
	payload := map[string]interface{}{
		"userId":          userID,
		"email":           email,
		"password":        password,
		"passwordVersion": passwordVersion,
		"name":            name,
	}
	resp, err := s.Client.sendRequest("POST", "/users/sha", payload)
	if err != nil {
		return nil, err
	}
	var user User
	if err = json.Unmarshal(resp, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUser retrieves a user by ID.
func (s *UsersService) GetUser(userID string) (*User, error) {
	path := fmt.Sprintf("/users/%s", userID)
	resp, err := s.Client.sendRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	var user User
	if err = json.Unmarshal(resp, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserPreferences retrieves user preferences.
func (s *UsersService) GetUserPreferences(userID string) (Preferences, error) {
	path := fmt.Sprintf("/users/%s/prefs", userID)
	resp, err := s.Client.sendRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	var prefs Preferences
	if err = json.Unmarshal(resp, &prefs); err != nil {
		return nil, err
	}
	return prefs, nil
}

// ListUsers lists project users.
func (s *UsersService) ListUsers(queries []string, search string) ([]*User, error) {
	q := url.Values{}
	for _, qs := range queries {
		q.Add("queries[]", qs)
	}
	if search != "" {
		q.Add("search", search)
	}
	path := "/users"
	if encoded := q.Encode(); encoded != "" {
		path += "?" + encoded
	}
	resp, err := s.Client.sendRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	var result struct {
		Users []*User `json:"users"`
	}
	if err = json.Unmarshal(resp, &result); err != nil {
		return nil, err
	}
	return result.Users, nil
}

// UpdateEmail updates user email.
func (s *UsersService) UpdateEmail(userID, email string) (*User, error) {
	payload := map[string]interface{}{"email": email}
	path := fmt.Sprintf("/users/%s/email", userID)
	resp, err := s.Client.sendRequest("PATCH", path, payload)
	if err != nil {
		return nil, err
	}
	var user User
	if err = json.Unmarshal(resp, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateEmailVerification sets user email verification status.
func (s *UsersService) UpdateEmailVerification(userID string, verified bool) (*User, error) {
	payload := map[string]interface{}{"emailVerification": verified}
	path := fmt.Sprintf("/users/%s/verification", userID)
	resp, err := s.Client.sendRequest("PATCH", path, payload)
	if err != nil {
		return nil, err
	}
	var user User
	if err = json.Unmarshal(resp, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateMFA enables or disables MFA for a user.
func (s *UsersService) UpdateMFA(userID string, mfa bool) (*User, error) {
	payload := map[string]interface{}{"mfa": mfa}
	path := fmt.Sprintf("/users/%s/mfa", userID)
	resp, err := s.Client.sendRequest("PATCH", path, payload)
	if err != nil {
		return nil, err
	}
	var user User
	if err = json.Unmarshal(resp, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateName updates user name.
func (s *UsersService) UpdateName(userID, name string) (*User, error) {
	payload := map[string]interface{}{"name": name}
	path := fmt.Sprintf("/users/%s/name", userID)
	resp, err := s.Client.sendRequest("PATCH", path, payload)
	if err != nil {
		return nil, err
	}
	var user User
	if err = json.Unmarshal(resp, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdatePassword updates user password.
func (s *UsersService) UpdatePassword(userID, password string) (*User, error) {
	payload := map[string]interface{}{"password": password}
	path := fmt.Sprintf("/users/%s/password", userID)
	resp, err := s.Client.sendRequest("PATCH", path, payload)
	if err != nil {
		return nil, err
	}
	var user User
	if err = json.Unmarshal(resp, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdatePhone updates user phone number.
func (s *UsersService) UpdatePhone(userID, number string) (*User, error) {
	payload := map[string]interface{}{"number": number}
	path := fmt.Sprintf("/users/%s/phone", userID)
	resp, err := s.Client.sendRequest("PATCH", path, payload)
	if err != nil {
		return nil, err
	}
	var user User
	if err = json.Unmarshal(resp, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdatePhoneVerification sets phone verification status.
func (s *UsersService) UpdatePhoneVerification(userID string, verified bool) (*User, error) {
	payload := map[string]interface{}{"phoneVerification": verified}
	path := fmt.Sprintf("/users/%s/verification/phone", userID)
	resp, err := s.Client.sendRequest("PATCH", path, payload)
	if err != nil {
		return nil, err
	}
	var user User
	if err = json.Unmarshal(resp, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUserLabels replaces user labels.
func (s *UsersService) UpdateUserLabels(userID string, labels []string) (*User, error) {
	payload := map[string]interface{}{"labels": labels}
	path := fmt.Sprintf("/users/%s/labels", userID)
	resp, err := s.Client.sendRequest("PUT", path, payload)
	if err != nil {
		return nil, err
	}
	var user User
	if err = json.Unmarshal(resp, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUserPreferences sets user preferences.
func (s *UsersService) UpdateUserPreferences(userID string, prefs Preferences) (Preferences, error) {
	payload := map[string]interface{}{"prefs": prefs}
	path := fmt.Sprintf("/users/%s/prefs", userID)
	resp, err := s.Client.sendRequest("PATCH", path, payload)
	if err != nil {
		return nil, err
	}
	var out Preferences
	if err = json.Unmarshal(resp, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// UpdateUserStatus updates user status (enabled or disabled).
func (s *UsersService) UpdateUserStatus(userID string, status bool) (*User, error) {
	payload := map[string]interface{}{"status": status}
	path := fmt.Sprintf("/users/%s/status", userID)
	resp, err := s.Client.sendRequest("PATCH", path, payload)
	if err != nil {
		return nil, err
	}
	var user User
	if err = json.Unmarshal(resp, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

// DeleteUser deletes a user by ID.
func (s *UsersService) DeleteUser(userID string) error {
	path := fmt.Sprintf("/users/%s", userID)
	_, err := s.Client.sendRequest("DELETE", path, nil)
	return err
}
