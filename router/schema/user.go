package schema

import "github.com/golang-jwt/jwt"

type User struct {
	UserId   int    `json:"userId"`
	UserName string `json:"userName"`
}

func (user User) ToClaims(claims jwt.StandardClaims) Claims {
	return Claims{
		UserId:         user.UserId,
		UserName:       user.UserName,
		StandardClaims: claims,
	}
}
