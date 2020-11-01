package iamiam

const (
	// EmailProfile is a profile for returning the email.
	EmailProfile string = "email"
	// SimpleProfile is a profile for returning email, firstname and lastname.
	SimpleProfile string = "simple"
)

// UserInfo contains info for profile creation.
type UserInfo struct {
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

// CreateEmailProfile creates email user profile.
func (u *UserInfo) CreateEmailProfile() *UserInfo {
	return &UserInfo{Email: u.Email}
}

// CreateSimpleProfile creates simple user profile.
func (u *UserInfo) CreateSimpleProfile() *UserInfo {
	return &UserInfo{Email: u.Email, FirstName: u.FirstName, LastName: u.LastName}
}
