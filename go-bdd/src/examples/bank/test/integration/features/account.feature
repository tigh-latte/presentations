Feature: Account creation
  Rules:
  - An account is successfully created with a valid details
  - If the password doesn't match the criteria, creation is rejected
  - If an email address is already in use, creation is rejected
  - If an invalid email address is given, creation is rejected

  @account @create
  Scenario: Account can be created successfully
    Given I run the SQL "reset.sql"
    When I make a POST request to "/api/v1/accounts" using "POST-account.json"
    Then the response status should be CREATED
    And the response body should match "POST-account.json"

  @account @create @error
  Scenario: CONFLICT on duplicate email
    Given I run the SQL "reset.sql"
    And I make a POST request to "/api/v1/accounts" using "POST-account.json"
    And the response status is CREATED
    When I make a POST request to "/api/v1/accounts" using "POST-account.json"
    Then the response status should be CONFLICT
    And the response body should match "POST-account_conflict.json"

  @account @create @error
  Scenario: BAD REQUEST on invalid id
    Given I run the SQL "reset.sql"
    When I make a GET request to "/api/v1/accounts/sajhfd"
    Then the response status should be BAD_REQUEST
    And the response body should match "GET-account_by_id_bad_request.json"
