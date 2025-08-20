package Model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserDb struct {
	Id            primitive.ObjectID `bson:"_id"`
	Name          *string            `bson:"mongoName" validate:"required,min=2,max=100"`
	Email         *string            `bson:"mongoEmail" validate:"required,email"`
	Phone         *string            `bson:"mongoPhone" validate:"required,len=10,numeric"`
	Password      *string            `bson:"mongoPwd" validate:"required,min=5"`
	CreatedAt     time.Time          `bson:"mongoCreatedAt"`
	UpdatedAt     time.Time          `bson:"mongoUpdatedAt"`
	Token_gen     *string            `bosn:"mongoToken"`
	Refresh_token *string            `bson:"mongoRefresh"`
	UserRole      *string            `bson:"mongoUserRole" validate:"required"`
	User_id       string             `bson:"mongouserId"`
}
