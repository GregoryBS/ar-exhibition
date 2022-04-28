package user

import (
	"ar_exhibition/pkg/utils"
	"strings"

	"github.com/robbert229/jwt"
)

var JWTKey = utils.RandString(32)

const Field = "Identificator"

func CreateJWT(id int) (string, error) {
	algorithm := jwt.HmacSha256(JWTKey)
	claims := jwt.NewClaim()
	claims.Set(Field, id)
	return algorithm.Encode(claims)
}

func CheckJWT(header string) int {
	var id int
	if header != "" {
		token := strings.Split(header, " ")[1]
		algorithm := jwt.HmacSha256(JWTKey)

		if claims, err := algorithm.DecodeAndValidate(token); err == nil {
			if identificator, err := claims.Get(Field); err == nil {
				id = int(identificator.(float64))
			}
		} else {
			id = -1
		}
	}
	return id
}
