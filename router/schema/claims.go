package schema

import "github.com/golang-jwt/jwt"

type Claims struct {
	UserId   int    `json:"userId"`
	UserName string `json:"userName"`
	jwt.StandardClaims
}

func (claims Claims) ToUser() User {
	return User{
		UserId:   claims.UserId,
		UserName: claims.UserName,
	}
}
