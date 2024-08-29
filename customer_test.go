package main

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/cucumber/godog"
)

type CustomerTestSteps struct {
	customerService                                      *CustomerService
	firstName, lastName, secondFirstName, secondLastName string
	err                                                  error
	count                                                int
}

var DEFAULT_BIRTHDAY = time.Date(1995, time.January, 1, 0, 0, 0, 0, time.UTC)

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features"},
			TestingT: t, // Testing instance that will run subtests.
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

func (t *CustomerTestSteps) theCustomerNameIs(ctx context.Context, fn, ln string) error {
	t.firstName = fn
	t.lastName = ln
	return nil
}

func (t *CustomerTestSteps) theCustomerIsCreated(ctx context.Context) error {
	t.err = t.customerService.AddCustomer(t.firstName, t.lastName, DEFAULT_BIRTHDAY)
	return nil
}

func (t *CustomerTestSteps) theCustomerCreationShouldBeSuccessful(ctx context.Context) error {
	if t.err != nil {
		return fmt.Errorf("expected no error but got %v", t.err)
	}
	return nil
}

func (t *CustomerTestSteps) theSecondCustomerIsCreated(ctx context.Context) error {
	return nil
}

func (t *CustomerTestSteps) theCustomerCreationShouldFail(ctx context.Context) error {
	if t.err == nil {
		return fmt.Errorf("expected error but got nil")
	}

	if t.err.Error() != "mandatory name parameter is missing" {
		return fmt.Errorf("expected 'mandatory name parameter is missing' error but got '%s'", t.err.Error())
	}

	return nil
}

func (t *CustomerTestSteps) theSecondCustomerCreationShouldFail(ctx context.Context) error {
	err := t.customerService.AddCustomer(t.secondFirstName, t.secondLastName, DEFAULT_BIRTHDAY)
	if err == nil {
		return fmt.Errorf("expected error but got nil")
	}

	if err.Error() != "customer already exists" {
		return fmt.Errorf("expected 'customer already exists' error but got '%s'", err.Error())
	}

	return nil
}

func (t *CustomerTestSteps) thereAreNoCustomers(ctx context.Context) error {
	return nil
}

func (t *CustomerTestSteps) noCustomersExist(ctx context.Context) error {
	return nil
}

func (t *CustomerTestSteps) thereIsACustomer(ctx context.Context, table *godog.Table) error {
	row := table.Rows[0]
	t.customerService.AddCustomer(row.Cells[0].Value, row.Cells[1].Value, DEFAULT_BIRTHDAY)
	return nil
}

func (t *CustomerTestSteps) thereAreSomeCustomers(ctx context.Context, table *godog.Table) error {
	for i, row := range table.Rows {
		if i == 0 {
			continue // skip header...
		}

		t.customerService.AddCustomer(row.Cells[0].Value, row.Cells[1].Value, DEFAULT_BIRTHDAY)
	}
	return nil
}

func (t *CustomerTestSteps) allCustomersAreSearched(ctx context.Context) error {
	t.count = len(t.customerService.SearchCustomers())
	return nil
}

func (t *CustomerTestSteps) theCustomerSabineMustermannIsSearched(ctx context.Context) error {
	t.count = len(t.customerService.SearchCustomersByName("Sabine", "Mustermann"))
	return nil
}

func (t *CustomerTestSteps) theCustomerRoseSmithIsSearched(ctx context.Context) error {
	return nil
}

func (t *CustomerTestSteps) theCustomerCanBeFound(ctx context.Context) error {
	customer := t.customerService.SearchCustomer(t.firstName, t.lastName)
	if customer == nil {
		return fmt.Errorf("expected customer to be found but got nil")
	}
	return nil
}

func (t *CustomerTestSteps) theCustomerCanNotBeFound(ctx context.Context) error {
	customer := t.customerService.SearchCustomer(t.firstName, t.lastName)
	if customer != nil {
		return fmt.Errorf("expected customer not to be found but got %v", customer)
	}
	return nil
}

func (t *CustomerTestSteps) theCustomerSabineMustermannCanBeFound(ctx context.Context) error {
	customer := t.customerService.SearchCustomer("Sabine", "Mustermann")

	if customer.FirstName != "Sabine" {
		return fmt.Errorf("expected first name to be Sabine but got %v", customer.FirstName)
	}

	if customer.LastName != "Mustermann" {
		return fmt.Errorf("expected last name to be Mustermann but got %v", customer.LastName)
	}

	return nil
}

func (t *CustomerTestSteps) theNumberOfCustomersFoundIs(ctx context.Context, expectedCount int) error {
	if t.count != expectedCount {
		return fmt.Errorf("expected %d customers to be found but got %d", expectedCount, t.count)
	}
	return nil
}

func (t *CustomerTestSteps) theSecondCustomerCanBeFound(ctx context.Context) error {
	t.customerService.AddCustomer(t.secondFirstName, t.secondLastName, DEFAULT_BIRTHDAY)
	customer := t.customerService.SearchCustomer(t.secondFirstName, t.secondLastName)
	if customer == nil {
		return fmt.Errorf("expected customer to be found but got nil")
	}
	return nil
}

func (t *CustomerTestSteps) theSecondCustomerIs(ctx context.Context, fn, ln string) error {
	t.secondFirstName = fn
	t.secondLastName = ln
	return nil
}

func InitializeScenario(sc *godog.ScenarioContext) {

	t := CustomerTestSteps{
		customerService: NewCustomerService(),
		firstName:       "",
		lastName:        "",
		secondFirstName: "",
		secondLastName:  "",
		count:           0,
	}

	sc.Given(`the customer name is (\w*) (\w*)`, t.theCustomerNameIs)
	sc.Given(`^the customer is created$`, t.theCustomerIsCreated)
	sc.When(`^the customer is created$`, t.theCustomerIsCreated)
	sc.When(`an invalid customer is created`, t.theCustomerIsCreated)
	sc.When(`the second customer is created`, t.theCustomerIsCreated)
	sc.Then(`the customer creation should be successful`, t.theCustomerCreationShouldBeSuccessful)
	sc.Then(`the customer creation should fail`, t.theCustomerCreationShouldFail)
	sc.Then(`the second customer creation should fail`, t.theSecondCustomerCreationShouldFail)
	sc.Given(`there are no customers`, t.thereAreNoCustomers)
	sc.Given(`no customers exist`, t.noCustomersExist)
	sc.Given(`there is a customer`, t.thereIsACustomer)
	sc.Given(`there are some customers`, t.thereAreSomeCustomers)
	sc.When(`all customers are searched`, t.allCustomersAreSearched)
	sc.When(`the customer Sabine Mustermann is searched`, t.theCustomerSabineMustermannIsSearched)
	sc.When(`the customer Rose Smith is searched`, t.theCustomerRoseSmithIsSearched)
	sc.Then(`the customer can be found`, t.theCustomerCanBeFound)
	sc.Then(`the customer can not be found`, t.theCustomerCanNotBeFound)
	sc.Then(`the customer Sabine Mustermann can be found`, t.theCustomerSabineMustermannCanBeFound)
	sc.Then(`^the number of customers found is (\d+)$`, t.theNumberOfCustomersFoundIs)
	sc.Then(`^the second customer can be found$`, t.theSecondCustomerCanBeFound)
	sc.Given(`^the second customer is (\w+) (\w+)$`, t.theSecondCustomerIs)

}
