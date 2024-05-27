package item

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/icrowley/fake"
)

func Test_ItemRepository_GetItems_ShouldReturnItems(t *testing.T) {
	dbConn, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer dbConn.Close()

	const pageSize = 5
	const page = 0

	columns := []string{"id", "name", "price", "manufacturer"}
	expectedItem := Item{ID: uuid.New(), Name: fake.ProductName(), Price: 23, Manufacturer: fake.Brand()}
	item2 := Item{ID: uuid.New(), Name: fake.ProductName(), Price: 4, Manufacturer: fake.Brand()}
	item3 := Item{ID: uuid.New(), Name: fake.ProductName(), Price: 5, Manufacturer: fake.Brand()}
	item4 := Item{ID: uuid.New(), Name: fake.ProductName(), Price: 11, Manufacturer: fake.Brand()}
	item5 := Item{ID: uuid.New(), Name: fake.ProductName(), Price: 100, Manufacturer: fake.Brand()}

	mock.ExpectQuery("SELECT id, name, price, manufacturer FROM cart ORDER BY id LIMIT \\$1 OFFSET \\$2").
		WithArgs(pageSize, page*pageSize).
		WillReturnRows(
			sqlmock.NewRows(columns).
				FromCSVString(convertObjectToCSV(expectedItem)).
				FromCSVString(convertObjectToCSV(item2)).
				FromCSVString(convertObjectToCSV(item3)).
				FromCSVString(convertObjectToCSV(item4)).
				FromCSVString(convertObjectToCSV(item5)),
		).
		RowsWillBeClosed()

	sut := NewRepository(dbConn)
	ctx := context.Background()

	result, err := sut.GetItems(ctx, page, pageSize)
	if err != nil {
		t.Fatalf("Error '%s' was not expected when fetching cart items", err)
	}

	if len(result) != pageSize {
		t.Fatalf("Unexpected number of items were given, '%d'. Expected '%d'.", len(result), pageSize)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func Test_ItemRepository_GetItems_WhenErrorOccurs_ShouldReturnError(t *testing.T) {
	dbConn, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer dbConn.Close()

	const pageSize = 5
	const page = 0
	expectedError := createError()

	mock.ExpectQuery("SELECT id, name, price, manufacturer FROM cart ORDER BY id LIMIT \\$1 OFFSET \\$2").
		WithArgs(pageSize, page*pageSize).
		WillReturnError(expectedError)

	sut := NewRepository(dbConn)
	ctx := context.Background()

	result, err := sut.GetItems(ctx, page, pageSize)
	if result != nil {
		t.Fatalf("Result '%s' was not expected when simulating a failed fetching cart item", err)
	}

	if !errors.Is(expectedError, err) {
		t.Fatalf("Expected failure '%s', but received '%s' when simulating a failed fetching cart item", expectedError, err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func Test_ItemRepository_GetItemByID_WhenItemExists_ShouldReturnItem(t *testing.T) {
	dbConn, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer dbConn.Close()

	columns := []string{"id", "name", "price", "manufacturer"}
	expectedItem := Item{ID: uuid.New(), Name: fake.ProductName(), Price: 23, Manufacturer: fake.Brand()}

	mock.ExpectQuery("SELECT id, name, price, manufacturer FROM cart WHERE id = \\$1").
		WithArgs(expectedItem.ID).
		WillReturnRows(sqlmock.NewRows(columns).FromCSVString(convertObjectToCSV(expectedItem))).
		RowsWillBeClosed()

	sut := NewRepository(dbConn)
	ctx := context.Background()

	result, err := sut.GetItemByID(ctx, expectedItem.ID)
	if err != nil {
		t.Fatalf("Error '%s' was not expected when fetching cart item", err)
	}

	if result != expectedItem {
		t.Fatalf("Unexpected item was given, '%+v'. Expected '%+v'.", result, expectedItem)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func Test_ItemRepository_GetItemByID_WhenItemDoesNotExist_ShouldReturnError(t *testing.T) {
	dbConn, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer dbConn.Close()

	expectedItemID := uuid.New()
	expectedError := createError()

	mock.ExpectQuery("SELECT id, name, price, manufacturer FROM cart WHERE id = \\$1").
		WithArgs(expectedItemID).
		WillReturnError(expectedError)

	sut := NewRepository(dbConn)
	ctx := context.Background()

	_, err = sut.GetItemByID(ctx, expectedItemID)
	if !errors.Is(expectedError, err) {
		t.Fatalf("Expected failure '%s', but received '%s' when simulating a failed fetching cart item", expectedError, err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func Test_ItemRepository_GetItemByID_WhenErrorOccurs_ShouldReturnError(t *testing.T) {
	dbConn, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer dbConn.Close()

	expectedItemID := uuid.New()
	expectedError := createError()

	mock.ExpectQuery("SELECT id, name, price, manufacturer FROM cart WHERE id = \\$1").
		WithArgs(expectedItemID).
		WillReturnError(expectedError)

	sut := NewRepository(dbConn)
	ctx := context.Background()

	_, err = sut.GetItemByID(ctx, expectedItemID)

	if !errors.Is(expectedError, err) {
		t.Fatalf("Expected failure '%s', but received '%s' when simulating a failed fetching cart item", expectedError, err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func Test_ItemRepository_AddItem_ShouldReturnInsertedItem(t *testing.T) {
	dbConn, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer dbConn.Close()

	columns := []string{"id"}
	expectedId := uuid.New()
	expectedItem := Item{ID: expectedId, Name: fake.ProductName(), Price: 23, Manufacturer: fake.Brand()}

	mock.ExpectQuery("INSERT INTO cart \\(name, price, manufacturer\\) VALUES \\(\\$1, \\$2, \\$3\\)").
		WithArgs(expectedItem.Name, expectedItem.Price, expectedItem.Manufacturer).
		WillReturnRows(sqlmock.NewRows(columns).FromCSVString(expectedId.String()))

	sut := NewRepository(dbConn)
	ctx := context.Background()

	result, err := sut.AddItem(ctx, expectedItem.Name, expectedItem.Price, expectedItem.Manufacturer)
	if err != nil {
		t.Fatalf("Error '%s' was not expected when adding an item to cart", err)
	}

	if result != expectedItem {
		t.Fatalf("Unexpected item was given, '%+v'. Expected '%+v'.", result, expectedItem)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func Test_ItemRepository_AddItem_WhenErrorOccurs_ShouldReturnError(t *testing.T) {
	dbConn, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer dbConn.Close()

	expectedItem := Item{ID: uuid.New(), Name: fake.ProductName(), Price: 23, Manufacturer: fake.Brand()}
	expectedError := createError()

	mock.ExpectQuery("INSERT INTO cart \\(name, price, manufacturer\\) VALUES \\(\\$1, \\$2, \\$3\\)").
		WithArgs(expectedItem.Name, expectedItem.Price, expectedItem.Manufacturer).
		WillReturnError(expectedError)

	sut := NewRepository(dbConn)
	ctx := context.Background()

	_, err = sut.AddItem(ctx, expectedItem.Name, expectedItem.Price, expectedItem.Manufacturer)
	if !errors.Is(expectedError, err) {
		t.Fatalf("Expected failure '%s', but received '%s' when simulating failure while adding cart item", expectedError, err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func Test_ItemRepository_UpdateItem_ShouldUpdateSpecificItem(t *testing.T) {
	dbConn, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer dbConn.Close()
	expectedItem := Item{ID: uuid.New(), Name: fake.ProductName(), Price: 23, Manufacturer: fake.Brand()}

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE cart SET name = \\$1, price = \\$2, manufacturer = \\$3 WHERE id = \\$4").
		WithArgs(expectedItem.Name, expectedItem.Price, expectedItem.Manufacturer, expectedItem.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	sut := NewRepository(dbConn)
	ctx := context.Background()

	result, err := sut.UpdateItem(ctx, &expectedItem)
	if err != nil {
		t.Fatalf("Result '%s' was not expected when simulating failure while updating cart item", err)
	}

	if result != expectedItem {
		t.Fatalf("Unexpected item was given, '%+v'. Expected '%+v'.", result, expectedItem)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func Test_ItemRepository_UpdateItem_WhenErrorOccurs_ShouldReturnError(t *testing.T) {
	dbConn, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer dbConn.Close()

	expectedItem := Item{ID: uuid.New(), Name: fake.ProductName(), Price: 23, Manufacturer: fake.Brand()}
	expectedError := createError()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE cart SET name = \\$1, price = \\$2, manufacturer = \\$3 WHERE id = \\$4").
		WithArgs(expectedItem.Name, expectedItem.Price, expectedItem.Manufacturer, expectedItem.ID).
		WillReturnError(expectedError)
	mock.ExpectRollback()

	sut := NewRepository(dbConn)
	ctx := context.Background()

	_, err = sut.UpdateItem(ctx, &expectedItem)
	if !errors.Is(expectedError, err) {
		t.Fatalf("Expected failure '%s', but received '%s' when simulating failure while updating cart item", expectedError, err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func Test_ItemRepository_RemoveItem_ShouldReturnID(t *testing.T) {
	dbConn, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer dbConn.Close()

	expectedItemID := uuid.New()

	mock.ExpectPrepare("DELETE FROM cart WHERE id = \\$1").
		ExpectExec().
		WithArgs(expectedItemID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	sut := NewRepository(dbConn)
	ctx := context.Background()

	result, err := sut.RemoveItem(ctx, expectedItemID)
	if err != nil {
		t.Fatalf("Result '%s' was not expected when simulating failure while removing cart item", err)
	}

	if result != expectedItemID {
		t.Fatalf("Unexpected id was given, '%d'. Expected '%d'.", result, expectedItemID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func Test_ItemRepository_RemoveItem_WhenErrorOccurs_ShouldReturnError(t *testing.T) {
	dbConn, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer dbConn.Close()

	expectedItemID := uuid.New()
	expectedError := createError()

	mock.ExpectPrepare("DELETE FROM cart WHERE id = \\$1").
		ExpectExec().
		WithArgs(expectedItemID).
		WillReturnError(expectedError)

	sut := NewRepository(dbConn)
	ctx := context.Background()

	result, err := sut.RemoveItem(ctx, expectedItemID)
	if result != expectedItemID {
		t.Fatalf("Unexpected id was given, '%d'. Expected '%d'.", result, expectedItemID)
	}

	if !errors.Is(expectedError, err) {
		t.Fatalf("Expected failure '%s', but received '%s' when simulating failure while removing cart item", expectedError, err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func convertObjectToCSV(item Item) string {
	return fmt.Sprintf("%s,%s,%d,%s", item.ID.String(), item.Name, item.Price, item.Manufacturer)
}

func createError() error {
	return fmt.Errorf("some error")
}
