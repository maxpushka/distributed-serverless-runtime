package schema

import "github.com/golang-jwt/jwt"

type User struct {
	UserId   int    `json:"userId"`
	UserName string `json:"userName"`
}

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

func (user User) ToClaims(claims jwt.StandardClaims) Claims {
	return Claims{
		UserId:         user.UserId,
		UserName:       user.UserName,
		StandardClaims: claims,
	}
}
