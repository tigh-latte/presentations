Feature: Account creation
  Rules:
  - An account is successfully created with a valid details
  - If the password doesn't match the criteria, creation is rejected
  - If an email address is already in use, creation is rejected
  - If an invalid email address is given, creation is rejected

  Scenario: Account can be created successfully
    Given I am unauthenticated
    When I make a POST request to "/accounts" using "post/create_account.json"
    Then the response code should be CREATED
    And the response body should match "post/create_account.json"

  Scenario: Invalid password is rejected
    Given I am unauthenticated
    When I make a POST request to "/accounts" using "post/create_account_bad_password.json"
    Then the response code should be BAD_REQUEST
    And the response body should match "post/create_account_bad_password.json"

  Scenario: Email address already in use is rejected
    Given I am unauthenticated
    When I make a POST request to "/accounts" using "post/create_account_dupe_email.json"
    Then the response code should be CONFLICT
    And the response body should match "post/create_account_dupe_email.json"

  Scenario: Invalid email address is rejected
    Given I am unauthenticated
    When I make a POST request to "/accounts" using "post/create_account_bad_email.json"
    Then the response code should be BAD_REQUEST
    And the response body should match "post/create_account_bad_email.json"
