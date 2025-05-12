package notmain

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func notmain() {
	password := "admin"

	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	fmt.Println(string(hashed))
}
