package benchmarks

import (
	"fmt"
	"testing"

	"route256/cart/pkg/cart"
	"route256/cart/tests"
)

func BenchmarkSample(b *testing.B) {
	if !tests.IsContainerRun() {
		b.Skip()

		return
	}

	client := cart.NewClient(tests.CartServiceAddress)

	for i := 0; i < b.N; i++ {
		if err := client.AddProduct(tests.DefaultUserID, tests.DefaultSKU); err != nil {
			b.Fatal(fmt.Errorf("can't add product: %w", err))
		}
	}
}
