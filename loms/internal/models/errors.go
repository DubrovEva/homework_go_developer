package models

import (
	"errors"
)

var (
	ErrProductNotFound = errors.New("product not found")
	ErrOrderNotFound   = errors.New("order not found")

	ErrOrderIsNotAwaitingPayment = errors.New("order isn't awaiting payment")
	ErrOrderIsNotNew             = errors.New("order isn't new")
	ErrOrderIsAlreadyPayed       = errors.New("order is already payed")

	ErrNotEnoughStockLeft     = errors.New("not enough stock left")
	ErrNotEnoughStockReserved = errors.New("not enough stock reserved")

	ErrOrderFailed             = errors.New("order failed")
	ErrNegativeCount           = errors.New("item count can't be negative")
	ErrNegativeStocksCountLeft = errors.New("stocks count can't be negative")
)
