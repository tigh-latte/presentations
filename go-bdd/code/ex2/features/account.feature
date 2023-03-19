Feature: Account
  - User can withdraw from their account
  - User withdrawing surplus funds receivees an error

  Scenario: Money can be withdrawn if account has sufficient funds
    Given I have an account with £100
    When I withdraw £50
    Then my remaining balance should be £50

  Scenario: Withdrawing funds exceeding my balance errors
    Given I have an account with £50
    When I withdraw £100
    Then an error should state "insufficient funds"
    And my remaining balance should be £50
