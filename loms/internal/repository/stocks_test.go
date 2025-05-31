package repository

//
//import (
//	"testing"
//
//	"route256/loms/internal/models"
//
//	"github.com/stretchr/testify/require"
//)
//
//const (
//	filePath = "stock-test-data.json"
//)
//
//func TestNewStocks(t *testing.T) {
//	stocks := []models.Stock{
//		{SKU: 1, Count: 65534, Reserved: 10},
//	}
//
//	repo := NewStocks(filePath)
//
//	require.Equal(t, stocks[0].Count, repo.items[1].Count)
//	require.Equal(t, stocks[0].SKU, repo.items[1].SKU)
//	require.Equal(t, stocks[0].Reserved, repo.items[1].Reserved)
//}
//
//func TestNewStocks_FileNotFound(t *testing.T) {
//	_, err := NewStocks("nonexistent.json")
//	require.Error(t, err)
//}
//
//func TestStocks_Reserve_Ok(t *testing.T) {
//	repo := &Stocks{
//		items: models.StocksMap{
//			1: {SKU: 1, Count: 65534, Reserved: 10},
//		},
//	}
//
//	err := repo.Reserve(1, 5)
//	require.NoError(t, err)
//	require.Equal(t, uint32(15), repo.items[1].Reserved)
//}
//
//func TestStocks_Reserve_NotEnoughStock(t *testing.T) {
//	repo := &Stocks{
//		items: models.StocksMap{
//			1: {SKU: 1, Count: 5, Reserved: 3},
//		},
//	}
//
//	err := repo.Reserve(1, 5)
//	require.ErrorIs(t, err, models.ErrNotEnoughStockLeft)
//}
//
//func TestStocks_Reserve_ProductNotFound(t *testing.T) {
//	repo := &Stocks{
//		items: models.StocksMap{},
//	}
//
//	err := repo.Reserve(1, 5)
//	require.ErrorIs(t, err, models.ErrProductNotFound)
//}
//
//func TestStocks_Release_Ok(t *testing.T) {
//	repo := &Stocks{
//		items: models.StocksMap{
//			1: {SKU: 1, Count: 65534, Reserved: 5},
//		},
//	}
//
//	err := repo.Release(1, 3)
//	require.NoError(t, err)
//	require.Equal(t, uint32(2), repo.items[1].Reserved)
//}
//
//func TestStocks_Release_ProductNotFound(t *testing.T) {
//	repo := &Stocks{
//		items: models.StocksMap{},
//	}
//
//	err := repo.Release(1, 3)
//	require.ErrorIs(t, err, models.ErrProductNotFound)
//}
//
//// Sell
//
//func TestStocks_Sell_Ok(t *testing.T) {
//	repo := &Stocks{
//		items: models.StocksMap{
//			1: {SKU: 1, Count: 10, Reserved: 5},
//		},
//	}
//
//	err := repo.Sell(1, 3)
//	require.NoError(t, err)
//	require.Equal(t, uint32(7), repo.items[1].Count)
//	require.Equal(t, uint32(2), repo.items[1].Reserved)
//}
//
//func TestStocks_Sell_NotEnoughReserved(t *testing.T) {
//	repo := &Stocks{
//		items: models.StocksMap{
//			1: {SKU: 1, Count: 10, Reserved: 2},
//		},
//	}
//
//	err := repo.Sell(1, 3)
//	require.ErrorIs(t, err, models.ErrNotEnoughStockReserved)
//}
//
//func TestStocks_Sell_NotEnoughStock(t *testing.T) {
//	repo := &Stocks{
//		items: models.StocksMap{
//			1: {SKU: 1, Count: 2, Reserved: 4},
//		},
//	}
//
//	err := repo.Sell(1, 3)
//	require.ErrorIs(t, err, models.ErrNotEnoughStockLeft)
//}
//
//func TestStocks_Sell_ProductNotFound(t *testing.T) {
//	repo := &Stocks{
//		items: models.StocksMap{},
//	}
//
//	err := repo.Sell(1, 3)
//	require.ErrorIs(t, err, models.ErrProductNotFound)
//}
//
//// AvailableStocks
//
//func TestStocks_AvailableStocks_Ok(t *testing.T) {
//	repo := &Stocks{
//		items: models.StocksMap{
//			1: {SKU: 1, Count: 10, Reserved: 3},
//		},
//	}
//
//	available, err := repo.AvailableStocks(1)
//	require.NoError(t, err)
//	require.Equal(t, uint32(7), available)
//}
//
//func TestStocks_AvailableStocks_ProductNotFound(t *testing.T) {
//	repo := &Stocks{
//		items: models.StocksMap{},
//	}
//
//	_, err := repo.AvailableStocks(1)
//	require.ErrorIs(t, err, models.ErrProductNotFound)
//}
