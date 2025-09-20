package shared

import (
	"database/sql"
	"time"
)

type Message struct {
	Message  string `json:"message"`
	To       string `json:"to"`
	Incoming bool   `json:"incoming"`
	Username string `json:"username"`
}

type ResponceError struct {
	Error string `json:"error"`
}

type LoginRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type User struct {
	ID          []byte         `json:"id"`
	Username    string         `json:"username"`
	Displayname sql.NullString `json:"displayname"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}
