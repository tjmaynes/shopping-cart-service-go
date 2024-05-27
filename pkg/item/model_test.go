package item

import (
	"errors"
	"github.com/google/uuid"
	"testing"

	validation "github.com/go-ozzo/ozzo-validation"
)

func Test_Item_Validate_WhenGivenValidItem_ShouldReturnNoErrors(t *testing.T) {
	item := Item{ID: uuid.New(), Name: "Some Product Name", Price: 23, Manufacturer: "Some Manufacturer"}

	if err := item.Validate(); err != nil {
		var e validation.InternalError
		if errors.As(err, &e) {
			t.Errorf("Received internal errors on validate Item, %s", e.InternalError())
		}
	}
}

func Test_Item_Validate_WhenGivenBadItems_ShouldReturnErrors(t *testing.T) {
	invalidItem := Item{ID: uuid.New(), Name: "", Price: -1, Manufacturer: ""}

	err := invalidItem.Validate()
	expectedErrors := "manufacturer: cannot be blank; name: cannot be blank; price: must be no less than 99."
	if err.Error() != expectedErrors {
		t.Errorf("Expected %s, Received %s", expectedErrors, err)
	}
}
