package item

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
)

// Decimal ..
type Decimal int64

// Item ..
type Item struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Price        Decimal   `json:"price"`
	Manufacturer string    `json:"manufacturer"`
}

// Validate ..
func (item Item) Validate() error {
	return validation.ValidateStruct(&item,
		// Name cannot be blank
		validation.Field(&item.Name, validation.Required),
		// Price should be greater than 0
		validation.Field(&item.Price, validation.Required, validation.Min(99)),
		// Manufacturer cannot be blank
		validation.Field(&item.Manufacturer, validation.Required),
	)
}

// ItemDTO ..
type ItemDTO struct {
	Name         string  `json:"name"`
	Price        Decimal `json:"price"`
	Manufacturer string  `json:"manufacturer"`
}

// Validate ..
func (item ItemDTO) Validate() error {
	return validation.ValidateStruct(&item,
		// Name cannot be blank
		validation.Field(&item.Name, validation.Required),
		// Price should be greater than 0
		validation.Field(&item.Price, validation.Required, validation.Min(99)),
		// Manufacturer cannot be blank
		validation.Field(&item.Manufacturer, validation.Required),
	)
}
