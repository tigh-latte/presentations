Feature: Account creation
  Rules:
  - An account is successfully created with a valid details
  - If an email address is already in use, creation is rejected
  - If an api key is not given, request is rejected
  - If an invalid user id is given, request is rejected

  @account @create
  Scenario: Account can be created successfully
    Given I run the SQL "reset.sql"
     And the headers:
       | key           | value      |
       | Authorization | Bearer dev |
    When I make a POST request to "/api/v1/accounts" using "POST-account.json"
    Then the response status should be CREATED
    And the response body should match "POST-account.json"

  @account @create @error
  Scenario: CONFLICT on duplicate email
    Given I run the SQL "reset.sql"
     And the headers:
       | key           | value      |
       | Authorization | Bearer dev |
    And I make a POST request to "/api/v1/accounts" using "POST-account.json"
    And the response status is CREATED
    When I make a POST request to "/api/v1/accounts" using "POST-account.json"
    Then the response status should be CONFLICT
    And the response body should match "errs/conflict.json"

  @account @create @error
  Scenario: BAD REQUEST on invalid id
    Given I run the SQL "reset.sql"
     And the headers:
       | key           | value      |
       | Authorization | Bearer dev |
    When I make a GET request to "/api/v1/accounts/sajhfd"
    Then the response status should be BAD_REQUEST
    And the response body should match "errs/bad_request.json"

  @account @create @error
  Scenario: UNAUTHORIZED on invalid auth header
    Given I run the SQL "reset.sql"
     And the headers:
       | key           | value      |
       | Authorization | Bearer wow |
    When I make a GET request to "/api/v1/accounts/sajhfd"
    Then the response status should be UNAUTHORIZED
    And the response body should match "errs/unauthorized.json"
