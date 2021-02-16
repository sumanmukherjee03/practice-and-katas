package users

import (
	"strings"

	"github.com/sumanmukherjee03/practice-and-katas/go/rest-api-project/bookstore_users-api/utils/crypto_utils"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (l *LoginRequest) PrepBeforeSubmit() {
	l.Email = strings.TrimSpace(strings.ToLower(l.Email))
	l.Password = crypto_utils.GetMd5(l.Password)
}
