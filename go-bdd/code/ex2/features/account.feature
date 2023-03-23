Feature: Account
  - User can withdraw from their account
  - User withdrawing surplus funds receivees an error

  @account @withdraw
  Scenario: Money can be withdrawn if account has sufficient funds
    Given I have an account with £100
    When I withdraw £50
    Then my remaining balance should be £50

  @account @withdraw @error
  Scenario: Withdrawing funds exceeding my balance errors
    Given I have an account with £50
    When I withdraw £100
    Then an error should state "insufficient funds"
    And my remaining balance should be £50

  @account @deposit
  Scenario: Money can be deposited
    Given I have an account with £100
    When I deposit £50
    Then my remaining balance should be £150

  @account @deposit @error
  Scenario: Depositing an empty sum errors
    Given I have an account with £50
    When I deposit £0
    Then an error should state "empty transaction"
    And my remaining balance should be £50
