package repository

import (
	"sync"
	"testing"

	"go.uber.org/goleak"

	"route256/cart/internal/models"
)

const (
	countPerGoroutine  = 10
	numProductsPerUser = 10
	numOperations      = 5000
	numGoroutines      = 100
	numProducts        = 100
	numReaders         = 100
	numUsers           = 50
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func TestConcurrentAddSameProduct(t *testing.T) {
	cart := NewCart()
	product := &models.Product{SKU: 123}

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			cart.AddProduct(defaultUserID, product, countPerGoroutine)
		}()
	}

	wg.Wait()

	products := cart.Cart(defaultUserID)
	expectedCount := uint32(numGoroutines * countPerGoroutine)
	if products[product.SKU] != expectedCount {
		t.Errorf("Expected product count to be %d, got %d", expectedCount, products[product.SKU])
	}
}

func TestConcurrentAddDifferentProducts(t *testing.T) {
	cart := NewCart()

	var wg sync.WaitGroup
	wg.Add(numProducts)

	for i := 0; i < numProducts; i++ {
		sku := int64(i)
		go func(s int64) {
			defer wg.Done()
			product := &models.Product{SKU: s}
			cart.AddProduct(defaultUserID, product, 1)
		}(sku)
	}

	wg.Wait()

	products := cart.Cart(defaultUserID)
	if len(products) != numProducts {
		t.Errorf("Expected %d products in cart, got %d", numProducts, len(products))
	}

	for i := 0; i < numProducts; i++ {
		sku := int64(i)
		if products[sku] != 1 {
			t.Errorf("Expected product with SKU %d to have count 1, got %d", sku, products[sku])
		}
	}
}

func TestConcurrentReadCart(t *testing.T) {
	cart := NewCart()

	for i := 0; i < numProducts; i++ {
		product := &models.Product{SKU: int64(i)}
		cart.AddProduct(defaultUserID, product, 1)
	}

	var wg sync.WaitGroup
	wg.Add(numReaders)

	for i := 0; i < numReaders; i++ {
		go func() {
			defer wg.Done()
			products := cart.Cart(defaultUserID)
			if len(products) != numProducts {
				t.Errorf("Expected %d products in cart, got %d", numProducts, len(products))
			}
		}()
	}

	wg.Wait()
}

func TestConcurrentReadAndWrite(t *testing.T) {
	cart := NewCart()

	for i := 0; i < numProducts/2; i++ {
		product := &models.Product{SKU: int64(i)}
		cart.AddProduct(defaultUserID, product, 1)
	}

	var wg sync.WaitGroup
	wg.Add(numProducts)

	for i := 0; i < numProducts/2; i++ {
		go func() {
			defer wg.Done()
			_ = cart.Cart(defaultUserID)
		}()
	}

	for i := numProducts / 2; i < numProducts; i++ {
		sku := int64(i)
		go func(s int64) {
			defer wg.Done()
			product := &models.Product{SKU: s}
			cart.AddProduct(defaultUserID, product, 1)
		}(sku)
	}

	wg.Wait()

	products := cart.Cart(defaultUserID)
	if len(products) != numProducts {
		t.Errorf("Expected %d products in cart, got %d", numProducts, len(products))
	}
}

func TestConcurrentDeleteProducts(t *testing.T) {
	cart := NewCart()

	for i := 0; i < numProducts; i++ {
		product := &models.Product{SKU: int64(i)}
		cart.AddProduct(defaultUserID, product, 1)
	}

	var wg sync.WaitGroup
	wg.Add(numProducts)

	for i := 0; i < numProducts; i++ {
		sku := int64(i)
		go func(s int64) {
			defer wg.Done()
			product := &models.Product{SKU: s}
			cart.DeleteProduct(defaultUserID, product)
		}(sku)
	}

	wg.Wait()

	products := cart.Cart(defaultUserID)
	if len(products) > 0 {
		t.Logf("After concurrent deletion, cart still has %d products", len(products))
	}
}

func TestConcurrentMultipleUsers(t *testing.T) {
	cart := NewCart()

	var wg sync.WaitGroup
	wg.Add(defaultUserID * numProductsPerUser)

	for userID := int64(1); userID <= defaultUserID; userID++ {
		for j := 0; j < numProductsPerUser; j++ {
			uid := userID
			sku := int64(j)
			go func(userId int64, productSku int64) {
				defer wg.Done()
				product := &models.Product{SKU: productSku}
				cart.AddProduct(userId, product, 1)
			}(uid, sku)
		}
	}

	wg.Wait()

	for userID := int64(1); userID <= defaultUserID; userID++ {
		products := cart.Cart(userID)
		if len(products) != numProductsPerUser {
			t.Errorf("User %d expected to have %d products, got %d", userID, numProductsPerUser, len(products))
		}
	}
}

func TestConcurrentClear(_ *testing.T) {
	cart := NewCart()

	for i := 0; i < numProducts; i++ {
		product := &models.Product{SKU: int64(i)}
		cart.AddProduct(defaultUserID, product, 1)
	}

	var wg sync.WaitGroup
	wg.Add(numOperations * 3)
	for i := 0; i < numOperations; i++ {
		go func() {
			defer wg.Done()
			_ = cart.Cart(defaultUserID)
		}()
	}

	for i := 0; i < numOperations; i++ {
		sku := int64(i % numProducts)
		go func(s int64) {
			defer wg.Done()
			product := &models.Product{SKU: s}
			cart.AddProduct(defaultUserID, product, 1)
		}(sku)
	}

	for i := 0; i < numOperations; i++ {
		idx := i
		go func(index int) {
			defer wg.Done()
			if index%10 == 0 {
				cart.Clear(defaultUserID)
			} else {
				sku := int64(index % numProducts)
				product := &models.Product{SKU: sku}
				cart.DeleteProduct(defaultUserID, product)
			}
		}(idx)
	}

	wg.Wait()
}

func TestStressConcurrency(_ *testing.T) {
	cart := NewCart()

	var wg sync.WaitGroup
	wg.Add(numOperations)

	for i := 0; i < numOperations; i++ {
		idx := i
		go func(index int) {
			defer wg.Done()

			userID := int64(index%numUsers) + 1
			sku := int64(index % 1000)
			product := &models.Product{SKU: sku}

			switch index % 4 {
			case 0:
				cart.AddProduct(userID, product, 1)
			case 1:
				cart.DeleteProduct(userID, product)
			case 2:
				_ = cart.Cart(userID)
			case 3:
				if index%10 == 0 {
					cart.Clear(userID)
				}
			}
		}(idx)
	}

	wg.Wait()
}
