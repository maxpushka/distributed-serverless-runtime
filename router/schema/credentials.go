package schema

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func CredentialsFromUser(user User) Credentials {
	return Credentials{
		Username: user.UserName,
	}
}
