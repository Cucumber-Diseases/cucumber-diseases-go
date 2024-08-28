package main

import (
	"errors"
	"strings"
	"time"
)

type Customer struct {
	FirstName string
	LastName  string
	Birthday  time.Time
}

func NewCustomer(firstName, lastName string, birthday time.Time) *Customer {
	return &Customer{
		FirstName: firstName,
		LastName:  lastName,
		Birthday:  birthday,
	}
}

func (c *Customer) FullName() string {
	return strings.ToLower(c.FirstName) + " " + c.LastName
}

func (c *Customer) Email() string {
	return strings.ToLower(c.FirstName) + "." + strings.ToLower(c.LastName) + "@mybank.com"
}

type CustomerService struct {
	customers []Customer
}

func NewCustomerService() *CustomerService {
	return &CustomerService{}
}

func (cs *CustomerService) AddCustomer(firstName, lastName string, birthday time.Time) error {
	if strings.TrimSpace(firstName) == "" || strings.TrimSpace(lastName) == "" {
		return errors.New("mandatory name parameter is missing")
	}

	if cs.CustomerExists(firstName, lastName) {
		return errors.New("customer already exists")
	}

	cs.customers = append(cs.customers, Customer{
		FirstName: firstName,
		LastName:  lastName,
		Birthday:  birthday,
	})
	return nil
}

func (cs *CustomerService) CustomerExists(firstName, lastName string) bool {
	for _, customer := range cs.customers {
		if hasSameName(customer, firstName, lastName) {
			return true
		}
	}
	return false
}

func (cs *CustomerService) RemoveCustomer(firstName, lastName string, birthday time.Time) {
	for i := 0; i < len(cs.customers); i++ {
		customer := cs.customers[i]
		if hasSameName(customer, firstName, lastName) && customer.Birthday.Equal(birthday) {
			cs.customers = append(cs.customers[:i], cs.customers[i+1:]...)
			i-- // Adjust index after removal
		}
	}
}

func (cs *CustomerService) SearchCustomer(firstName, lastName string) *Customer {
	for _, customer := range cs.searchCustomers(func(c Customer) bool {
		return hasSameName(c, firstName, lastName)
	}) {
		return &customer
	}
	return nil
}

func (cs *CustomerService) SearchCustomers() []Customer {
	return cs.searchCustomers(func(c Customer) bool {
		return true
	})
}

func (cs *CustomerService) SearchCustomersByName(firstName, lastName string) []Customer {
	return cs.searchCustomers(func(c Customer) bool {
		return hasSameName(c, firstName, lastName)
	})
}

func (cs *CustomerService) searchCustomers(match func(Customer) bool) []Customer {
	var result []Customer
	for _, customer := range cs.customers {
		if match(customer) {
			result = append(result, customer)
		}
	}
	return result
}

func hasSameName(customer Customer, firstName, lastName string) bool {
	return customer.FirstName == firstName && customer.LastName == lastName
}
