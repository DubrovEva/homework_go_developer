// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: api.proto

package v1

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"google.golang.org/protobuf/types/known/anypb"
)

// ensure the imports are used
var (
	_ = bytes.MinRead
	_ = errors.New("")
	_ = fmt.Print
	_ = utf8.UTFMax
	_ = (*regexp.Regexp)(nil)
	_ = (*strings.Reader)(nil)
	_ = net.IPv4len
	_ = time.Duration(0)
	_ = (*url.URL)(nil)
	_ = (*mail.Address)(nil)
	_ = anypb.Any{}
	_ = sort.Sort
)

// Validate checks the field values on CreateOrderRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *CreateOrderRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on CreateOrderRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// CreateOrderRequestMultiError, or nil if none found.
func (m *CreateOrderRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *CreateOrderRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for UserId

	for idx, item := range m.GetItems() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, CreateOrderRequestValidationError{
						field:  fmt.Sprintf("Items[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, CreateOrderRequestValidationError{
						field:  fmt.Sprintf("Items[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return CreateOrderRequestValidationError{
					field:  fmt.Sprintf("Items[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return CreateOrderRequestMultiError(errors)
	}

	return nil
}

// CreateOrderRequestMultiError is an error wrapping multiple validation errors
// returned by CreateOrderRequest.ValidateAll() if the designated constraints
// aren't met.
type CreateOrderRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m CreateOrderRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m CreateOrderRequestMultiError) AllErrors() []error { return m }

// CreateOrderRequestValidationError is the validation error returned by
// CreateOrderRequest.Validate if the designated constraints aren't met.
type CreateOrderRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CreateOrderRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CreateOrderRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CreateOrderRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CreateOrderRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CreateOrderRequestValidationError) ErrorName() string {
	return "CreateOrderRequestValidationError"
}

// Error satisfies the builtin error interface
func (e CreateOrderRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCreateOrderRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CreateOrderRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CreateOrderRequestValidationError{}

// Validate checks the field values on CreateOrderResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *CreateOrderResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on CreateOrderResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// CreateOrderResponseMultiError, or nil if none found.
func (m *CreateOrderResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *CreateOrderResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for OrderId

	if len(errors) > 0 {
		return CreateOrderResponseMultiError(errors)
	}

	return nil
}

// CreateOrderResponseMultiError is an error wrapping multiple validation
// errors returned by CreateOrderResponse.ValidateAll() if the designated
// constraints aren't met.
type CreateOrderResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m CreateOrderResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m CreateOrderResponseMultiError) AllErrors() []error { return m }

// CreateOrderResponseValidationError is the validation error returned by
// CreateOrderResponse.Validate if the designated constraints aren't met.
type CreateOrderResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CreateOrderResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CreateOrderResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CreateOrderResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CreateOrderResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CreateOrderResponseValidationError) ErrorName() string {
	return "CreateOrderResponseValidationError"
}

// Error satisfies the builtin error interface
func (e CreateOrderResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCreateOrderResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CreateOrderResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CreateOrderResponseValidationError{}

// Validate checks the field values on OrderInfoRequest with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *OrderInfoRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on OrderInfoRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// OrderInfoRequestMultiError, or nil if none found.
func (m *OrderInfoRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *OrderInfoRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for OrderId

	if len(errors) > 0 {
		return OrderInfoRequestMultiError(errors)
	}

	return nil
}

// OrderInfoRequestMultiError is an error wrapping multiple validation errors
// returned by OrderInfoRequest.ValidateAll() if the designated constraints
// aren't met.
type OrderInfoRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m OrderInfoRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m OrderInfoRequestMultiError) AllErrors() []error { return m }

// OrderInfoRequestValidationError is the validation error returned by
// OrderInfoRequest.Validate if the designated constraints aren't met.
type OrderInfoRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e OrderInfoRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e OrderInfoRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e OrderInfoRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e OrderInfoRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e OrderInfoRequestValidationError) ErrorName() string { return "OrderInfoRequestValidationError" }

// Error satisfies the builtin error interface
func (e OrderInfoRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sOrderInfoRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = OrderInfoRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = OrderInfoRequestValidationError{}

// Validate checks the field values on OrderInfoResponse with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *OrderInfoResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on OrderInfoResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// OrderInfoResponseMultiError, or nil if none found.
func (m *OrderInfoResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *OrderInfoResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if all {
		switch v := interface{}(m.GetOrder()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, OrderInfoResponseValidationError{
					field:  "Order",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, OrderInfoResponseValidationError{
					field:  "Order",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetOrder()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return OrderInfoResponseValidationError{
				field:  "Order",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return OrderInfoResponseMultiError(errors)
	}

	return nil
}

// OrderInfoResponseMultiError is an error wrapping multiple validation errors
// returned by OrderInfoResponse.ValidateAll() if the designated constraints
// aren't met.
type OrderInfoResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m OrderInfoResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m OrderInfoResponseMultiError) AllErrors() []error { return m }

// OrderInfoResponseValidationError is the validation error returned by
// OrderInfoResponse.Validate if the designated constraints aren't met.
type OrderInfoResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e OrderInfoResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e OrderInfoResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e OrderInfoResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e OrderInfoResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e OrderInfoResponseValidationError) ErrorName() string {
	return "OrderInfoResponseValidationError"
}

// Error satisfies the builtin error interface
func (e OrderInfoResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sOrderInfoResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = OrderInfoResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = OrderInfoResponseValidationError{}

// Validate checks the field values on PayOrderRequest with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *PayOrderRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on PayOrderRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// PayOrderRequestMultiError, or nil if none found.
func (m *PayOrderRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *PayOrderRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for OrderId

	if len(errors) > 0 {
		return PayOrderRequestMultiError(errors)
	}

	return nil
}

// PayOrderRequestMultiError is an error wrapping multiple validation errors
// returned by PayOrderRequest.ValidateAll() if the designated constraints
// aren't met.
type PayOrderRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m PayOrderRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m PayOrderRequestMultiError) AllErrors() []error { return m }

// PayOrderRequestValidationError is the validation error returned by
// PayOrderRequest.Validate if the designated constraints aren't met.
type PayOrderRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e PayOrderRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e PayOrderRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e PayOrderRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e PayOrderRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e PayOrderRequestValidationError) ErrorName() string { return "PayOrderRequestValidationError" }

// Error satisfies the builtin error interface
func (e PayOrderRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sPayOrderRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = PayOrderRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = PayOrderRequestValidationError{}

// Validate checks the field values on PayOrderResponse with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *PayOrderResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on PayOrderResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// PayOrderResponseMultiError, or nil if none found.
func (m *PayOrderResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *PayOrderResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if len(errors) > 0 {
		return PayOrderResponseMultiError(errors)
	}

	return nil
}

// PayOrderResponseMultiError is an error wrapping multiple validation errors
// returned by PayOrderResponse.ValidateAll() if the designated constraints
// aren't met.
type PayOrderResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m PayOrderResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m PayOrderResponseMultiError) AllErrors() []error { return m }

// PayOrderResponseValidationError is the validation error returned by
// PayOrderResponse.Validate if the designated constraints aren't met.
type PayOrderResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e PayOrderResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e PayOrderResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e PayOrderResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e PayOrderResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e PayOrderResponseValidationError) ErrorName() string { return "PayOrderResponseValidationError" }

// Error satisfies the builtin error interface
func (e PayOrderResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sPayOrderResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = PayOrderResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = PayOrderResponseValidationError{}

// Validate checks the field values on CancelOrderRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *CancelOrderRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on CancelOrderRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// CancelOrderRequestMultiError, or nil if none found.
func (m *CancelOrderRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *CancelOrderRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for OrderId

	if len(errors) > 0 {
		return CancelOrderRequestMultiError(errors)
	}

	return nil
}

// CancelOrderRequestMultiError is an error wrapping multiple validation errors
// returned by CancelOrderRequest.ValidateAll() if the designated constraints
// aren't met.
type CancelOrderRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m CancelOrderRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m CancelOrderRequestMultiError) AllErrors() []error { return m }

// CancelOrderRequestValidationError is the validation error returned by
// CancelOrderRequest.Validate if the designated constraints aren't met.
type CancelOrderRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CancelOrderRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CancelOrderRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CancelOrderRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CancelOrderRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CancelOrderRequestValidationError) ErrorName() string {
	return "CancelOrderRequestValidationError"
}

// Error satisfies the builtin error interface
func (e CancelOrderRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCancelOrderRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CancelOrderRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CancelOrderRequestValidationError{}

// Validate checks the field values on CancelOrderResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *CancelOrderResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on CancelOrderResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// CancelOrderResponseMultiError, or nil if none found.
func (m *CancelOrderResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *CancelOrderResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if len(errors) > 0 {
		return CancelOrderResponseMultiError(errors)
	}

	return nil
}

// CancelOrderResponseMultiError is an error wrapping multiple validation
// errors returned by CancelOrderResponse.ValidateAll() if the designated
// constraints aren't met.
type CancelOrderResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m CancelOrderResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m CancelOrderResponseMultiError) AllErrors() []error { return m }

// CancelOrderResponseValidationError is the validation error returned by
// CancelOrderResponse.Validate if the designated constraints aren't met.
type CancelOrderResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CancelOrderResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CancelOrderResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CancelOrderResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CancelOrderResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CancelOrderResponseValidationError) ErrorName() string {
	return "CancelOrderResponseValidationError"
}

// Error satisfies the builtin error interface
func (e CancelOrderResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCancelOrderResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CancelOrderResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CancelOrderResponseValidationError{}

// Validate checks the field values on StocksInfoRequest with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *StocksInfoRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on StocksInfoRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// StocksInfoRequestMultiError, or nil if none found.
func (m *StocksInfoRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *StocksInfoRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Sku

	if len(errors) > 0 {
		return StocksInfoRequestMultiError(errors)
	}

	return nil
}

// StocksInfoRequestMultiError is an error wrapping multiple validation errors
// returned by StocksInfoRequest.ValidateAll() if the designated constraints
// aren't met.
type StocksInfoRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m StocksInfoRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m StocksInfoRequestMultiError) AllErrors() []error { return m }

// StocksInfoRequestValidationError is the validation error returned by
// StocksInfoRequest.Validate if the designated constraints aren't met.
type StocksInfoRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e StocksInfoRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e StocksInfoRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e StocksInfoRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e StocksInfoRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e StocksInfoRequestValidationError) ErrorName() string {
	return "StocksInfoRequestValidationError"
}

// Error satisfies the builtin error interface
func (e StocksInfoRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sStocksInfoRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = StocksInfoRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = StocksInfoRequestValidationError{}

// Validate checks the field values on StocksInfoResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *StocksInfoResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on StocksInfoResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// StocksInfoResponseMultiError, or nil if none found.
func (m *StocksInfoResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *StocksInfoResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Count

	if len(errors) > 0 {
		return StocksInfoResponseMultiError(errors)
	}

	return nil
}

// StocksInfoResponseMultiError is an error wrapping multiple validation errors
// returned by StocksInfoResponse.ValidateAll() if the designated constraints
// aren't met.
type StocksInfoResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m StocksInfoResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m StocksInfoResponseMultiError) AllErrors() []error { return m }

// StocksInfoResponseValidationError is the validation error returned by
// StocksInfoResponse.Validate if the designated constraints aren't met.
type StocksInfoResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e StocksInfoResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e StocksInfoResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e StocksInfoResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e StocksInfoResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e StocksInfoResponseValidationError) ErrorName() string {
	return "StocksInfoResponseValidationError"
}

// Error satisfies the builtin error interface
func (e StocksInfoResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sStocksInfoResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = StocksInfoResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = StocksInfoResponseValidationError{}
