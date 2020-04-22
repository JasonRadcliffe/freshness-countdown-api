package users

type user struct {
	UserID int    `json:"UserID"`
	Email  string `json:"Email"`
}

//OauthUser is what will be populated upon receiving confirmation from Oauth Provider.
type OauthUser struct {
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	UserID        int
}
