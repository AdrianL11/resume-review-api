package mongodb

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Role string

const (
	OwnerRole     Role = "owner"
	Administrator      = "administrator"
	User               = "user"
)

type ChangePassword struct {
	Password string `bson:"password" json:"password"`
}

type InsertProfile struct {
	Password     string `bson:"password" json:"password"`
	FirstLogin   bool   `bson:"first_login" json:"first_login"`
	FirstName    string `bson:"first_name" json:"first_name"`
	LastName     string `bson:"last_name" json:"last_name"`
	Country      string `bson:"country" json:"country"`
	ProfileImage string `bson:"profile_image" json:"profile_image"`
}

type Profile struct {
	ID               primitive.ObjectID  `bson:"_id" json:"_id"`
	EmailAddress     string              `bson:"email_address" json:"email_address"`
	FirstLogin       bool                `bson:"first_login" json:"first_login"`
	FirstName        string              `bson:"first_name" json:"first_name"`
	LastName         string              `bson:"last_name" json:"last_name"`
	Country          string              `bson:"country" json:"country"`
	ProfileImage     string              `bson:"profile_image" json:"profile_image"`
	Role             Role                `bson:"role" json:"role"`
	ActiveUser       bool                `bson:"active_user" json:"active_user"`
	CreatedBy        primitive.ObjectID  `bson:"created_by" json:"created_by"`
	CreationDate     primitive.Timestamp `bson:"creation_date" json:"creation_date"`
	DeactivatedBy    primitive.ObjectID  `bson:"deactivated_by" json:"deactivated_by"`
	DeactivationDate primitive.Timestamp `bson:"deactivation_date" json:"deactivation_date"`
}

type Session struct {
	ID             string              `bson:"session_id" json:"session_id"`
	LoggedInIP     string              `bson:"loggedin_ip" json:"loggedin_ip"`
	LastSeenIP     string              `bson:"lastseen_ip" json:"lastseen_ip"`
	Expiration     primitive.Timestamp `bson:"expiration" json:"expiration"`
	LastSeen       primitive.Timestamp `bson:"last_seen" json:"last_seen"`
	Active         bool                `bson:"is_active" json:"is_active"`
	RevocationDate primitive.Timestamp `bson:"revocation_date" json:"revocation_date"`
	UserAgent      string              `bson:"user_agent" json:"user_agent"`
	UserId         primitive.ObjectID  `bson:"user_id" json:"user_id"`
}

type ForgotPassword struct {
	Token        string              `bson:"token" json:"token"`
	UserId       primitive.ObjectID  `bson:"user_id" json:"user_id"`
	CreationIP   string              `bson:"creation_ip" json:"creation_ip"`
	CreationDate primitive.Timestamp `bson:"creation_date" json:"creation_date"`
	Expiration   primitive.Timestamp `bson:"expiration" json:"expiration"`
	Active       bool                `bson:"is_active" json:"is_active"`
}
