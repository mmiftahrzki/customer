package customer

import "errors"

var errCustomerAlreadyExists = errors.New("customer already exists")
var errCustomerNotFound = errors.New("customer not found")
var errNotImplemented = errors.New("not implemented")
