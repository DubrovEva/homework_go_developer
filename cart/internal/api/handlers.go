package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/asaskevich/govalidator"

	"route256/cart/internal/models"
)

func (r *Router) addProduct(writer http.ResponseWriter, req *http.Request) {
	sku, err := r.parseSKU(req)
	if err != nil {
		ErrorResponse(writer, http.StatusBadRequest, err.Error())

		return
	}

	userID, err := r.parseUserID(req)
	if err != nil {
		ErrorResponse(writer, http.StatusBadRequest, err.Error())

		return
	}

	var addProductRequest AddProductRequest

	defer req.Body.Close()

	if err = json.NewDecoder(req.Body).Decode(&addProductRequest); err != nil {
		ErrorResponse(writer, http.StatusBadRequest, err.Error())

		return
	}

	result, err := govalidator.ValidateStruct(addProductRequest)
	if err != nil {
		ErrorResponse(writer, http.StatusBadRequest, err.Error())

		return
	}

	if !result {
		ErrorResponse(writer, http.StatusBadRequest, "count must be more than zero")

		return
	}

	if err := r.cartService.AddProduct(req.Context(), userID, sku, addProductRequest.Count); err != nil {
		if errors.Is(err, models.ErrProductNotFound) {
			ErrorResponse(writer, http.StatusPreconditionFailed, models.ErrProductNotFound.Error())
		} else {
			ErrorResponse(writer, http.StatusInternalServerError, err.Error())
		}

		return
	}

	writer.WriteHeader(http.StatusOK)

}

func (r *Router) deleteProduct(writer http.ResponseWriter, req *http.Request) {
	sku, err := r.parseSKU(req)
	if err != nil {
		ErrorResponse(writer, http.StatusBadRequest, err.Error())

		return
	}

	userID, err := r.parseUserID(req)
	if err != nil {
		ErrorResponse(writer, http.StatusBadRequest, err.Error())

		return
	}

	if err := r.cartService.DeleteProduct(req.Context(), userID, sku); err != nil {
		if errors.Is(err, models.ErrProductNotFound) {
			ErrorResponse(writer, http.StatusPreconditionFailed, models.ErrProductNotFound.Error())
		} else {
			ErrorResponse(writer, http.StatusInternalServerError, err.Error())
		}

		return
	}

	writer.WriteHeader(http.StatusNoContent)
}

func (r *Router) cart(writer http.ResponseWriter, req *http.Request) {
	userID, err := r.parseUserID(req)
	if err != nil {
		ErrorResponse(writer, http.StatusBadRequest, err.Error())

		return
	}

	cart, err := r.cartService.Cart(req.Context(), userID)
	if err != nil {
		ErrorResponse(writer, http.StatusInternalServerError, err.Error())

		return
	}

	if len(cart.Items) == 0 {
		writer.WriteHeader(http.StatusNotFound)

		return
	}

	writer.Header().Add("Content-Type", "application/json")
	if err = json.NewEncoder(writer).Encode(&cart); err != nil {
		ErrorResponse(writer, http.StatusInternalServerError, err.Error())

		return
	}

	writer.WriteHeader(http.StatusOK)
}

func (r *Router) checkout(writer http.ResponseWriter, req *http.Request) {
	userID, err := r.parseUserID(req)
	if err != nil {
		ErrorResponse(writer, http.StatusBadRequest, err.Error())

		return
	}

	orderID, err := r.cartService.Checkout(req.Context(), userID)
	if err != nil {
		ErrorResponse(writer, http.StatusInternalServerError, err.Error())

		return
	}

	response := OrderIDResponse{
		OrderID: orderID,
	}

	writer.Header().Add("Content-Type", "application/json")
	if err = json.NewEncoder(writer).Encode(&response); err != nil {
		ErrorResponse(writer, http.StatusInternalServerError, err.Error())

		return
	}

	writer.WriteHeader(http.StatusOK)
}

func (r *Router) clearCart(writer http.ResponseWriter, req *http.Request) {
	userID, err := r.parseUserID(req)
	if err != nil {
		ErrorResponse(writer, http.StatusBadRequest, err.Error())

		return
	}

	r.cartService.ClearCart(userID)

	writer.WriteHeader(http.StatusNoContent)
}

func (r *Router) parseSKU(req *http.Request) (int64, error) {
	skuRaw := req.PathValue("sku")
	if skuRaw == "" {
		return 0, fmt.Errorf("sku is empty")
	}

	sku, err := strconv.Atoi(skuRaw)
	if err != nil {
		return 0, fmt.Errorf("can't parse sku: %w", err)
	}

	if sku <= 0 {
		return 0, fmt.Errorf("sku must be more than zero")
	}

	return int64(sku), nil
}

func (r *Router) parseUserID(req *http.Request) (int64, error) {
	userIDRaw := req.PathValue("user_id")
	if userIDRaw == "" {
		return 0, fmt.Errorf("user_id is empty")
	}

	userID, err := strconv.Atoi(userIDRaw)
	if err != nil {
		return 0, fmt.Errorf("can't parse user_id: %w", err)
	}

	if userID <= 0 {
		return 0, fmt.Errorf("user_id must be more than zero")
	}

	return int64(userID), nil
}
