package item

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"testing"

	"github.com/icrowley/fake"
)

func Test_ItemService_GetItems_WhenItemsExist_ShouldReturnAllItems(t *testing.T) {
	items := []Item{
		{ID: uuid.New(), Name: fake.ProductName(), Price: 23, Manufacturer: fake.Brand()},
		{ID: uuid.New(), Name: fake.ProductName(), Price: 4, Manufacturer: fake.Brand()},
		{ID: uuid.New(), Name: fake.ProductName(), Price: 5, Manufacturer: fake.Brand()},
		{ID: uuid.New(), Name: fake.ProductName(), Price: 11, Manufacturer: fake.Brand()},
		{ID: uuid.New(), Name: fake.ProductName(), Price: 100, Manufacturer: fake.Brand()},
	}

	const pageSize = 10
	const page = 0
	var pageSizeCalled int64
	var pageCalled int64

	mockRepository := &RepositoryMock{
		GetItemsFunc: func(ctx context.Context, pageSize int64, page int64) ([]Item, error) {
			pageSizeCalled = pageSize
			pageCalled = page
			return items, nil
		},
	}

	ctx := context.Background()
	sut := NewService(mockRepository)

	results, err := sut.GetItems(ctx, pageSize, page)
	if err != nil {
		t.Fatalf("Should not have failed!")
	}

	if len(results) != len(items) {
		t.Errorf("Expected an array of cart items of size %d. Got %d", len(items), len(results))
	}

	callsToSend := len(mockRepository.GetItemsCalls())
	if callsToSend != 1 {
		t.Errorf("Send was called %d times", callsToSend)
	}

	if pageSizeCalled != pageSize {
		t.Errorf("Unexpected recipient: %d", pageSizeCalled)
	}

	if pageCalled != page {
		t.Errorf("Unexpected recipient: %d", pageCalled)
	}
}

func Test_ItemService_GetItemByID_WhenItemExists_ShouldReturnItem(t *testing.T) {
	id := uuid.New()
	item := Item{ID: id, Name: fake.ProductName(), Price: 23, Manufacturer: fake.Brand()}
	var idCalled uuid.UUID

	mockRepository := &RepositoryMock{
		GetItemByIDFunc: func(ctx context.Context, id uuid.UUID) (Item, error) {
			idCalled = id
			return item, nil
		},
	}

	ctx := context.Background()
	sut := NewService(mockRepository)

	result, err := sut.GetItemByID(ctx, id)
	if err != nil {
		t.Fatalf("Should not have failed!")
	}

	if result != item {
		t.Errorf("Expected cart items %+v. Got %+v", item, result)
	}

	callsToSend := len(mockRepository.GetItemByIDCalls())
	if callsToSend != 1 {
		t.Errorf("Send was called %d times", callsToSend)
	}

	if idCalled != id {
		t.Errorf("Unexpected recipient: %d", id)
	}
}

func Test_ItemService_GetItemByID_WhenItemDoesNotExist_ShouldReturnError(t *testing.T) {
	testError := createError()

	mockRepository := &RepositoryMock{
		GetItemByIDFunc: func(ctx context.Context, id uuid.UUID) (Item, error) {
			return Item{}, testError
		},
	}

	ctx := context.Background()
	sut := NewService(mockRepository)

	_, err := sut.GetItemByID(ctx, uuid.New())
	if !errors.Is(err, testError) {
		t.Errorf("Expected error message %s. Got %s", testError, err)
	}

	callsToSend := len(mockRepository.GetItemByIDCalls())
	if callsToSend != 1 {
		t.Errorf("Send was called %d times", callsToSend)
	}
}

func Test_ItemService_AddItem_WhenGivenValidItem_ShouldReturnItem(t *testing.T) {
	var itemCalled *Item
	expectedItem := Item{
		ID:           uuid.New(),
		Name:         fake.ProductName(),
		Price:        Decimal(99),
		Manufacturer: fake.Brand(),
	}

	mockRepository := &RepositoryMock{
		AddItemFunc: func(ctx context.Context, name string, price Decimal, manufacturer string) (Item, error) {
			itemCalled = &expectedItem
			return expectedItem, nil
		},
	}

	ctx := context.Background()
	sut := NewService(mockRepository)

	result, err := sut.AddItem(ctx, &ItemDTO{
		Name:         expectedItem.Name,
		Price:        expectedItem.Price,
		Manufacturer: expectedItem.Manufacturer,
	})
	if err != nil {
		t.Fatalf("Should not have failed!")
	}

	if result != *itemCalled {
		t.Errorf("Expected cart item: %+v. Got %+v", itemCalled, result)
	}

	callsToSend := len(mockRepository.AddItemCalls())
	if callsToSend != 1 {
		t.Errorf("Send was called %d times", callsToSend)
	}
}

func Test_ItemService_AddItem_WhenGivenInvalidItem_ShouldReturnError(t *testing.T) {
	mockRepository := &RepositoryMock{
		AddItemFunc: func(ctx context.Context, name string, price Decimal, manufacturer string) (Item, error) {
			return Item{}, nil
		},
	}

	ctx := context.Background()
	sut := NewService(mockRepository)

	expectedErrorMessage := "price: must be no less than 99."

	_, err := sut.AddItem(ctx, &ItemDTO{Name: fake.ProductName(), Price: 23, Manufacturer: fake.Brand()})
	if err.Error() != expectedErrorMessage {
		t.Errorf("Error unexpected error message %s was given", err)
	}

	callsToSend := len(mockRepository.AddItemCalls())
	if callsToSend != 0 {
		t.Errorf("Send was called %d times", callsToSend)
	}
}

func Test_ItemService_UpdateCartItem_WhenGivenValidItem_ShouldReturnItem(t *testing.T) {
	id := uuid.New()
	expectedItem := Item{
		ID:           id,
		Name:         fake.ProductName(),
		Price:        Decimal(99),
		Manufacturer: fake.Brand(),
	}

	mockRepository := &RepositoryMock{
		UpdateItemFunc: func(ctx context.Context, item *Item) (Item, error) {
			return expectedItem, nil
		},
	}

	ctx := context.Background()
	sut := NewService(mockRepository)

	result, err := sut.UpdateItem(ctx, &expectedItem)
	if err != nil {
		t.Fatalf("Should not have failed!")
	}

	if result != expectedItem {
		t.Errorf("Expected cart item: %+v. Got %+v", expectedItem, result)
	}

	callsToSend := len(mockRepository.UpdateItemCalls())
	if callsToSend != 1 {
		t.Errorf("Send was called %d times", callsToSend)
	}
}

func Test_ItemService_UpdateCartItem_WhenGivenInvalidItem_ShouldReturnServiceError(t *testing.T) {
	invalidItem := Item{
		ID:           uuid.New(),
		Name:         fake.ProductName(),
		Price:        Decimal(25),
		Manufacturer: fake.Brand(),
	}

	mockRepository := &RepositoryMock{
		UpdateItemFunc: func(ctx context.Context, item *Item) (Item, error) {
			return Item{}, nil
		},
	}

	ctx := context.Background()
	sut := NewService(mockRepository)

	_, serviceError := sut.UpdateItem(ctx, &invalidItem)

	if serviceError.StatusCode() != InvalidItem {
		t.Errorf("Error unexpected error message %s was given", serviceError.Message())
	}

	callsToSend := len(mockRepository.UpdateItemCalls())
	if callsToSend != 0 {
		t.Errorf("Send was called %d times", callsToSend)
	}
}

func Test_ItemService_UpdateCartItem_WhenUnknownErrorOccurs_ShouldReturnServiceError(t *testing.T) {
	invalidItem := Item{
		ID:           uuid.New(),
		Name:         fake.ProductName(),
		Price:        Decimal(99),
		Manufacturer: fake.Brand(),
	}

	err := errors.New("unknown error occurred")

	mockRepository := &RepositoryMock{
		UpdateItemFunc: func(ctx context.Context, item *Item) (Item, error) {
			return Item{}, err
		},
	}

	ctx := context.Background()
	sut := NewService(mockRepository)

	_, serviceError := sut.UpdateItem(ctx, &invalidItem)

	if serviceError.StatusCode() != UnknownException {
		t.Errorf("Error unexpected error message %s was given", serviceError.Message())
	}

	callsToSend := len(mockRepository.UpdateItemCalls())
	if callsToSend != 1 {
		t.Errorf("Send was called %d times", callsToSend)
	}
}

func Test_ItemService_RemoveCartItem_WhenItemExists_ShouldReturnItemID(t *testing.T) {
	var idCalled uuid.UUID
	deletedItem := Item{
		ID:           uuid.New(),
		Name:         fake.ProductName(),
		Price:        Decimal(99),
		Manufacturer: fake.Brand(),
	}

	mockRepository := &RepositoryMock{
		RemoveItemFunc: func(ctx context.Context, id uuid.UUID) (uuid.UUID, error) {
			idCalled = deletedItem.ID
			return deletedItem.ID, nil
		},
	}

	ctx := context.Background()
	sut := NewService(mockRepository)

	result, err := sut.RemoveItem(ctx, deletedItem.ID)
	if err != nil {
		t.Fatalf("Should not have failed!")
	}

	if result != idCalled {
		t.Errorf("Expected cart item: %d. Got %d", idCalled, result)
	}

	callsToSend := len(mockRepository.RemoveItemCalls())
	if callsToSend != 1 {
		t.Errorf("Send was called %d times", callsToSend)
	}
}

func Test_ItemService_RemoveCartItem_WhenUnknownErrorOccurs_ShouldReturnServiceError(t *testing.T) {
	deletedItem := Item{
		ID:           uuid.New(),
		Name:         fake.ProductName(),
		Price:        Decimal(99),
		Manufacturer: fake.Brand(),
	}

	unknownError := errors.New("unknown error")

	mockRepository := &RepositoryMock{
		RemoveItemFunc: func(ctx context.Context, id uuid.UUID) (uuid.UUID, error) {
			return deletedItem.ID, unknownError
		},
	}

	ctx := context.Background()
	sut := NewService(mockRepository)

	_, serviceError := sut.RemoveItem(ctx, deletedItem.ID)
	if serviceError.StatusCode() != UnknownException {
		t.Errorf("Error unexpected error message %s was given", serviceError.Message())
	}

	callsToSend := len(mockRepository.UpdateItemCalls())
	if callsToSend != 0 {
		t.Errorf("Send was called %d times", callsToSend)
	}
}
