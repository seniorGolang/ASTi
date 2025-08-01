// @asti version=2.0 author="Team" advanced=true
package examples

import (
	"context"
	"encoding/json"
	"math/big"
	"time"
)

// @asti name="AdvancedService" timeout=60 retry=5
type AdvancedService interface {
	// @asti method=CreateOrder timeout=30
	CreateOrder(ctx context.Context, userID string, amount *big.Float) (orderID string, err error)

	// @asti method=GetOrder timeout=10
	GetOrder(ctx context.Context, orderID string) (order Order, err error)

	// @asti method=UpdateOrder timeout=20
	UpdateOrder(ctx context.Context, orderID string, updates OrderUpdates) (order Order, err error)
}

// @asti type=OrderModel validation=strict
type Order struct {
	// @asti field=ID required=true
	ID string `json:"id"`

	// @asti field=UserID required=true
	UserID string `json:"userId"`

	// @asti field=Amount required=true
	Amount *big.Float `json:"amount"`

	// @asti field=CreatedAt required=true
	CreatedAt time.Time `json:"createdAt"`

	// @asti field=Status required=true
	Status string `json:"status"`

	// @asti field=Metadata optional=true
	Metadata json.RawMessage `json:"metadata,omitempty"`
}

// @asti type=OrderUpdatesModel
type OrderUpdates struct {
	// @asti field=Amount min=0
	Amount *big.Float `json:"amount,omitempty"`

	// @asti field=Status enum=pending,processing,completed,cancelled
	Status *string `json:"status,omitempty"`

	// @asti field=Metadata optional=true
	Metadata *json.RawMessage `json:"metadata,omitempty"`
}
