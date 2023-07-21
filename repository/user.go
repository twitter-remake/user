package repository

import (
	"errors"
	"net/mail"
	"time"

	userpb "github.com/twitter-remake/user/proto/gen/go/user"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// User represents a generic user
type User struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	ScreenName       string    `json:"screen_name"`
	Email            string    `json:"email"`
	Bio              string    `json:"bio"`
	Location         string    `json:"location"`
	Website          string    `json:"website"`
	ProfileImageURL  string    `json:"profile_image_url"`
	ProfileBannerURL string    `json:"profile_banner_url"`
	BirthDate        time.Time `json:"birth_date"`
	FollowersCount   int       `json:"followers_count"`
	FollowingsCount  int       `json:"followings_count"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func (u User) PB() *userpb.User {
	return &userpb.User{
		Id:               u.ID,
		Name:             u.Name,
		ScreenName:       u.ScreenName,
		Email:            u.Email,
		Bio:              u.Bio,
		Location:         u.Location,
		Website:          u.Website,
		BirthDate:        timestamppb.New(u.BirthDate),
		ProfileImageUrl:  u.ProfileImageURL,
		ProfileBannerUrl: u.ProfileBannerURL,
		FollowersCount:   int32(u.FollowersCount),
		FollowingsCount:  int32(u.FollowingsCount),
		CreatedAt:        timestamppb.New(u.CreatedAt),
		UpdatedAt:        timestamppb.New(u.UpdatedAt),
	}
}

// Validate validates the user fields
func (u User) Validate() error {
	if u.ScreenName == "" {
		return errors.New("missing User.ScreenName")
	}

	if u.Email == "" {
		return errors.New("missing User.Email")
	}

	if _, err := mail.ParseAddress(u.Email); err != nil {
		return errors.New("invalid email address")
	}

	return nil
}
