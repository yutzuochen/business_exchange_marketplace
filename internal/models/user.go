package models

import (
	"database/sql"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID              int64      `json:"id" db:"id"`
	UUID            string     `json:"uuid" db:"uuid"`
	Email           string     `json:"email" db:"email"`
	PasswordHash    string     `json:"-" db:"password_hash"`
	FirstName       string     `json:"first_name" db:"first_name"`
	LastName        string     `json:"last_name" db:"last_name"`
	Phone           *string    `json:"phone" db:"phone"`
	CompanyName     *string    `json:"company_name" db:"company_name"`
	Bio             *string    `json:"bio" db:"bio"`
	ProfileImageURL *string    `json:"profile_image_url" db:"profile_image_url"`
	WebsiteURL      *string    `json:"website_url" db:"website_url"`
	LinkedinURL     *string    `json:"linkedin_url" db:"linkedin_url"`
	City            *string    `json:"city" db:"city"`
	State           *string    `json:"state" db:"state"`
	Country         string     `json:"country" db:"country"`
	ZipCode         *string    `json:"zip_code" db:"zip_code"`
	EmailVerified   bool       `json:"email_verified" db:"email_verified"`
	IsActive        bool       `json:"is_active" db:"is_active"`
	IsPremium       bool       `json:"is_premium" db:"is_premium"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at" db:"updated_at"`
	LastLoginAt     *time.Time `json:"last_login_at" db:"last_login_at"`
}

type CreateUserRequest struct {
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
	FirstName       string `json:"first_name" validate:"required,max=100"`
	LastName        string `json:"last_name" validate:"required,max=100"`
	Phone           string `json:"phone" validate:"omitempty,max=20"`
	CompanyName     string `json:"company_name" validate:"omitempty,max=255"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UpdateProfileRequest struct {
	FirstName   string  `json:"first_name" validate:"required,max=100"`
	LastName    string  `json:"last_name" validate:"required,max=100"`
	Phone       *string `json:"phone" validate:"omitempty,max=20"`
	CompanyName *string `json:"company_name" validate:"omitempty,max=255"`
	Bio         *string `json:"bio" validate:"omitempty,max=1000"`
	WebsiteURL  *string `json:"website_url" validate:"omitempty,url"`
	LinkedinURL *string `json:"linkedin_url" validate:"omitempty,url"`
	City        *string `json:"city" validate:"omitempty,max=100"`
	State       *string `json:"state" validate:"omitempty,max=100"`
	ZipCode     *string `json:"zip_code" validate:"omitempty,max=20"`
}

// HashPassword generates a bcrypt hash of the password
func (u *User) HashPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hash)
	return nil
}

// CheckPassword verifies if the provided password matches the hash
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

// FullName returns the user's full name
func (u *User) FullName() string {
	return u.FirstName + " " + u.LastName
}

// LocationString returns a formatted location string
func (u *User) LocationString() string {
	if u.City != nil && u.State != nil {
		return *u.City + ", " + *u.State
	}
	if u.City != nil {
		return *u.City
	}
	if u.State != nil {
		return *u.State
	}
	return ""
}

// GetUserByID retrieves a user by ID
func GetUserByID(db *sql.DB, id int64) (*User, error) {
	user := &User{}
	query := `
		SELECT id, uuid, email, password_hash, first_name, last_name, phone, company_name,
			   bio, profile_image_url, website_url, linkedin_url, city, state, country,
			   zip_code, email_verified, is_active, is_premium, created_at, updated_at, last_login_at
		FROM users WHERE id = ? AND is_active = TRUE
	`
	err := db.QueryRow(query, id).Scan(
		&user.ID, &user.UUID, &user.Email, &user.PasswordHash, &user.FirstName, &user.LastName,
		&user.Phone, &user.CompanyName, &user.Bio, &user.ProfileImageURL, &user.WebsiteURL,
		&user.LinkedinURL, &user.City, &user.State, &user.Country, &user.ZipCode,
		&user.EmailVerified, &user.IsActive, &user.IsPremium, &user.CreatedAt,
		&user.UpdatedAt, &user.LastLoginAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserByEmail retrieves a user by email
func GetUserByEmail(db *sql.DB, email string) (*User, error) {
	user := &User{}
	query := `
		SELECT id, uuid, email, password_hash, first_name, last_name, phone, company_name,
			   bio, profile_image_url, website_url, linkedin_url, city, state, country,
			   zip_code, email_verified, is_active, is_premium, created_at, updated_at, last_login_at
		FROM users WHERE email = ? AND is_active = TRUE
	`
	err := db.QueryRow(query, email).Scan(
		&user.ID, &user.UUID, &user.Email, &user.PasswordHash, &user.FirstName, &user.LastName,
		&user.Phone, &user.CompanyName, &user.Bio, &user.ProfileImageURL, &user.WebsiteURL,
		&user.LinkedinURL, &user.City, &user.State, &user.Country, &user.ZipCode,
		&user.EmailVerified, &user.IsActive, &user.IsPremium, &user.CreatedAt,
		&user.UpdatedAt, &user.LastLoginAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// CreateUser creates a new user in the database
func CreateUser(db *sql.DB, req *CreateUserRequest) (*User, error) {
	user := &User{
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Country:   "US",
		IsActive:  true,
	}

	if req.Phone != "" {
		user.Phone = &req.Phone
	}
	if req.CompanyName != "" {
		user.CompanyName = &req.CompanyName
	}

	// Hash password
	if err := user.HashPassword(req.Password); err != nil {
		return nil, err
	}

	query := `
		INSERT INTO users (email, password_hash, first_name, last_name, phone, company_name, country, is_active)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	result, err := db.Exec(query, user.Email, user.PasswordHash, user.FirstName, user.LastName,
		user.Phone, user.CompanyName, user.Country, user.IsActive)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return GetUserByID(db, id)
}

// UpdateUser updates user information
func (u *User) Update(db *sql.DB, req *UpdateProfileRequest) error {
	query := `
		UPDATE users SET 
			first_name = ?, last_name = ?, phone = ?, company_name = ?, bio = ?,
			website_url = ?, linkedin_url = ?, city = ?, state = ?, zip_code = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`
	_, err := db.Exec(query, req.FirstName, req.LastName, req.Phone, req.CompanyName,
		req.Bio, req.WebsiteURL, req.LinkedinURL, req.City, req.State, req.ZipCode, u.ID)
	return err
}

// UpdateLastLogin updates the user's last login timestamp
func (u *User) UpdateLastLogin(db *sql.DB) error {
	query := `UPDATE users SET last_login_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := db.Exec(query, u.ID)
	return err
}
