package models

import (
	"time"

"go.mongodb.org/mongo-driver/bson/primitive"
)


// this will be the data fields for the user which are editors
type User struct{
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`


	// identity of the user

	UserName string `bson:"username" json:"username" validate:"required,min=5,max=20"`
	Email string `bson:"email" json:"email" validate:"required,email"`
	Phone string `bson:"phone" json:"phone" validate:"required"`


	PasswordHash string `bson:"password_hash" json:"password_hash"`
	
	// role

	// by default it will be editor only will set that in the backend auth
	Role string `bson:"role" json:"role"`

	/* visibility required data for the hiring page 

	// is hiring listed is managed by admin 

	// few rules described below for the profile which will be visible on the hiring page 

	

	1. The admin must turn is hiring listed to true 
	2. They should allow their profile to show on hiting page 
	3. The employment status should not be working it should be open to work
	*/

	ShowOnHiringPage bool `bson:"show_on_hiring_page"  json:"show_on_hiring_page"`
	IsHiringListed bool `bson:"is_hiring_listed" json:"is_hiring_listed"`
	EmploymentStatus string  `bson:"employment_status" json:"employment_status"`



	// verification

	IsEmailVerified bool `bson:"is_email_verified" json:"is_email_verified"`
	IsPhoneVerified bool `bson:"is_phone_verified" json:"is_phone_verified"`

	//otp check and method

	OtpHash string `bson:"otp_hash,omitempty" json:"otp_hash"`
	OtpExpiry time.Time `bson:"otp_expiry,omitempty" json:"-"`
 OtpPurpose string `bson:"otp_purpose" json:"otp_purpose"`

	//meta data for the user profile display

	ProfileImage string `bson:"profile_image,omitempty" json:"profile_image,omitempty"`
	Skills []string `bson:"skills,omitempty" json:"skills,omitempty"`
	Bio string `bson:"bio,omitempty" json:"bio,omitempty"`
	PortFolio string `bson:"portfolio,omitempty" json:"portfolio,omitempty"`



	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`

	// will use this in the websockets chat for real time status 


	LastSeen *time.Time `bson:"last_seen,omitempty" json:"last_seen,omitempty"`


}