package user

//User type is the struct in the Domain that contains all the fields for what a User is.
type User struct {
	UserID       int    `json:"UserID"`
	Email        string `json:"Email"`
	FirstName    string `json:"FirstName"`
	LastName     string `json:"LastName"`
	FullName     string `json:"FullName"`
	CreatedDate  string `json:"TimeCreated"`
	AccessToken  string `json:"AccessToken"`
	RefreshToken string `json:"RefreshToken"`
	AlexaUserID  string `json:"alexa_user_id"`
	TempMatch    string `json:"TempMatch"`
}

//OauthUser is what will be populated upon receiving confirmation from Oauth Provider.
type OauthUser struct {
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"email_verified"`
	FirstName     string `json:"given_name"`
	LastName      string `json:"family_name"`
	FullName      string `json:"name"`
	UserID        int
}

//Users is a slice of type User
type Users []User

//Contains methods and validators that a user would know about themselves
//
