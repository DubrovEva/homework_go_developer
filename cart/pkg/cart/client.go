package cart

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"route256/cart/internal/models"
)

type Client struct {
	client *http.Client

	addr string
}

func NewClient(addr string) *Client {
	client := &http.Client{Transport: http.DefaultTransport}

	return &Client{
		client: client,
		addr:   addr,
	}
}

func (c *Client) DeleteProduct(userID, sku int64) error {
	request, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodDelete,
		fmt.Sprintf("%s/user/%d/cart/%d", c.addr, userID, sku),
		http.NoBody,
	)
	if err != nil {
		return fmt.Errorf("can't create request to delete product: %w", err)
	}

	response, err := c.client.Do(request)
	if err != nil {
		return fmt.Errorf("can't delete product: %w", err)
	}

	if response.StatusCode != http.StatusNoContent {
		return fmt.Errorf("http query failed")
	}

	return nil
}

func (c *Client) AddProduct(userID, sku int64) error {
	request, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodPost,
		fmt.Sprintf("%s/user/%d/cart/%d", c.addr, userID, sku),
		bytes.NewReader([]byte(`{"count": 1}`)),
	)
	if err != nil {
		return fmt.Errorf("can't create request to add product: %w", err)
	}

	response, err := c.client.Do(request)
	if err != nil {
		return fmt.Errorf("can't add product: %w", err)
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("http query failed")
	}

	return nil
}

func (c *Client) Cart(userID int64) (*models.Cart, error) {
	request, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodGet,
		fmt.Sprintf("%s/user/%d/cart", c.addr, userID),
		http.NoBody,
	)
	if err != nil {
		return nil, fmt.Errorf("can't create request to get cart: %w", err)
	}

	response, err := c.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("can't get cart: %w", err)
	}

	if response.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http query failed")
	}

	cart := &models.Cart{}
	err = json.NewDecoder(response.Body).Decode(cart)
	if err != nil {
		return nil, fmt.Errorf("can't decode cart: %w", err)
	}

	return cart, nil
}

func (c *Client) ClearCart(userID int64) error {
	request, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodDelete,
		fmt.Sprintf("%s/user/%d/cart", c.addr, userID),
		http.NoBody,
	)
	if err != nil {
		return fmt.Errorf("can't create request to clear cart: %w", err)
	}

	response, err := c.client.Do(request)
	if err != nil {
		return fmt.Errorf("can't clear cart: %w", err)
	}

	if response.StatusCode != http.StatusNoContent {
		return fmt.Errorf("http query failed")
	}

	return nil
}
