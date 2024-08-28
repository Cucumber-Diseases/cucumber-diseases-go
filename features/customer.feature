Feature: Customer

  Rule: Create a customer with user name

    Scenario: Should successfully create new customer
      Given the customer first name is "Max"
      And the customer last name is "Mustermann"
      And the customer's birthday is 2000/03/19
      When the customer is created
      Then the customer creation should be successful

    Scenario: Should find newly created customer
      Given there are no customers
      And the customer name is Rose Smith
      When the customer is created
      And all customers are searched
      Then the number of customers found is 1


  Rule: The first name and last name must be provided

    Scenario Outline: Cannot create customer without first or lastname
      Given the customer name is <firstname> <lastname>
      When an invalid customer is created
      Then the customer creation should fail
      Examples:
        | description    | firstname | lastname   |
        | only firstname | Max       |            |
        | just lastname  |           | Mustermann |
        | neither nor    |           |            |


  Rule: Customer name must be unique

    Scenario: Should be able to create two customer with different names
      Given the customer name is Max Mustermann
      When the customer is created
      Given the second customer is Sabine Mustermann
      When the second customer is created
      Then the second customer can be found

    Scenario: Cannot create two customer with the same name
      Given the customer name is Max Mustermann
      And the second customer is Max Mustermann
      When the customer is created
      Then the second customer creation should fail


  Rule: Existing customers can be searched

    Scenario: Should find an existing customer
      Given there is a customer
        | Sabine | Mustermann |
      Then the customer Sabine Mustermann can be found

    Scenario: Should find customers by name
      Given the customer name is Sabine Mustermann
      And the customer is created
      When the customer Sabine Mustermann is searched
      Then the number of customers found is 1

    Scenario: Cannot find customer that does not exist
      Given there are no customers
      When all customers are searched
      Then the number of customers found is 0

    Scenario: Should find multiple customers
      Given there are some customers
        | firstname | lastname   |
        | Max       | Mustermann |
        | Sabine    | Mustermann |
        | Horst     | Mustermann |
      When all customers are searched
      Then the number of customers found is 3