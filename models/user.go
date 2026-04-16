package models


// User is the raw HN API user model.
type User struct {
	ID        string `json:"id"`
	Created   int64  `json:"created"`
	Karma     int    `json:"karma"`
	About     string `json:"about"`
	Submitted []int  `json:"submitted"`
}

// UserView is the display model for users.
type UserView struct {
	ID      string `toon:"id" json:"id"`
	Karma   int    `toon:"karma" json:"karma"`
	Created string `toon:"created" json:"created"`
	About   string `toon:"about,omitempty" json:"about,omitempty"`
}

// ToUserView converts an API User to a display model.
func (u *User) ToUserView(truncate func(string) string) UserView {
	about := u.About
	if truncate != nil {
		about = truncate(about)
	}
	return UserView{
		ID:      u.ID,
		Karma:   u.Karma,
		Created: formatTime(u.Created),
		About:   about,
	}
}

// UserOutput wraps a user result.
type UserOutput struct {
	User       UserView `toon:"user" json:"user"`
	Submitted  int      `toon:"submitted" json:"submitted"`
}
