// @asti version=1.0 author="Team"
package examples

import (
	"context"

	"github.com/seniorGolang/asti/examples/sample/dto"
)

// @asti name="UserService" timeout=30 retry=3
type UserService interface {
	// @asti method=CreateUser timeout=10
	CreateUser(ctx context.Context, name string, email string) (userID dto.SomeUser, err error)

	// @asti method=GetUser timeout=5
	GetUser(ctx context.Context, userID string) (user User, err error)

	// @asti method=UpdateUser timeout=15
	UpdateUser(ctx context.Context, userID string, updates UserUpdates) (user User, err error)

	// @asti method=DeleteUser timeout=8
	DeleteUser(ctx context.Context, userID string) (err error)
}

// @asti type=UserModel validation=strict
type User struct {
	// @asti field=ID required=true
	ID string `json:"id"`

	// @asti field=Name required=true maxLength=100
	Name string `json:"name"`

	// @asti field=Email validation=email
	Email string `json:"email"`

	// @asti field=Age min=0 max=150
	Age int `json:"age"`
}

// @asti type=UserUpdatesModel
type UserUpdates struct {
	// @asti field=Name maxLength=100
	Name *string `json:"name,omitempty"`

	// @asti field=Email validation=email
	Email *string `json:"email,omitempty"`

	// @asti field=Age min=0 max=150
	Age *int `json:"age,omitempty"`
}
