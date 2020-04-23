package user

//User type is the struct in the Domain that contains all the fields for what a User is.
type User struct {
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

//Contains methods and validators that a user would know about themselves
//
