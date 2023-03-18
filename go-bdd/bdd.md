---
theme: ./theme.json
---
# Welcome!

BDD Development in Golang, using Cucumber

---
# What is BDD?

An extension of TDD (I know...) that emphasizes feature development around user stories.

I'm assuming we, for the most part, use user stories here.

---
# What is cucumber?

A test process the deals with application behaviour, and is typically a natural extension of BDD.

Tests are written in a natural (instead of a programming) language, potentially allowing non-technical members of the workforce to understand/contribute to the test suite.

So in Golange where you would likely:
```go
// ex1
func TestAccount_Withdrawal(t *testing.T) {
	tests := map[string]struct {
		balance  int64
		withdraw int64

		expBalance int64
		expErr     error
	}{
		"successful withdraw": {
			balance:    100,
			withdraw:   50,
			expBalance: 50,
		},
		"not enough funds": {
			balance:  50,
			withdraw: 100,

			expBalance: 50,
			expErr:     errs.ErrInsufficientFunds,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			account := ex1.Account{
				Balance: test.balance,
			}
			_, err := account.Withdraw(test.withdraw)

			assert.Equal(t, test.expBalance, account.Balance)
			assert.Equal(t, test.expErr, err)
		})
	}
}
```

---
In cucumber, you would write the tests in natural language:

```gherkin
Scenario: Money can be withdrawn if account has sufficient funds
  Given I have an account with £100
  When I withdraw £50
  Then my remaining balance should be £50

Scenario: Withdrawing funds exceeding my balance errors
  Given I have an account with £50
  When I withdraw £100
  Then an error should state "insufficient funds"
  And my remaining balance should be £50
```

---
# Example

Take the following golang unit test:

```go
// models.go
type Plate struct {
    Contents PlateItems
}

type PlateItems map[PlateItem]int

type PlateItem int

const (
    PlateItemSausage      PlateItem = iota
    PlateItemBacon
    PlateItemEggsPoached
    PlateItemToastSlice
)

type Person struct {}

func (p Person) EatFrom(plate *Plate, items PlateItems) {
    // definition
}

// person_test.go
func Test_Eat(t *testing.T) {
    person := models.Person{}

    plate := models.Plate{
        Contents: models.PlateItems{
            models.PlateItemSausage:     2,
            models.PlateItemBacon:       3,
            models.PlateItemEggsPoached: 2,
            models.PlateItemToastSlice:  2,
        }
    }

    person.EatFrom(&plate, models.PlateItems{
        models.PlateItemSausage:    1,
        models.PlateItemToastSlice: 2,
    })

    assert.Equal(t, plate, &models.PlateItems{
        models.PlateItemSausage:     1,
        models.PlateItemBacon:       3,
        models.PlateItemEggsPoached: 2,
    })
}
```

---
# A world in cucumber
The same test couple be written as follows:
```gherkin
Scenario: Eaten food is removed from plate
  Given I have a plate with:
    | item           | quantity |
    | sausage        | 2        |
    | bacon          | 3        |
    | poached egg    | 2        |
    | slice of toast | 2        |
  When I eat:
    | item           | quantity |
    | sausage        | 1        |
    | slice of toast | 2        |
  Then the plate should be left with:
    | item           | quantity |
    | sausage        | 1        |
    | bacon          | 3        |
    | poached egg    | 2        |

```

---
# What's wrong with Golang?

The problem isn't with golang in particular, but instead writing usage tests in a programming language.

To understand why this is a problem we should consider what makes a test suite a good.

---
# A good test suite

Follows the "even juniors" rule.

1. Is extendable.
      - Even juniors should be able to make meaning contributions.
2. Is approachable
      - Even juniors should be able to understand it.
3. Covers core usage cases.

4. Acts as documentation for the codebase.

---
# 1. Is understandable.

Code is easy to write, but hard to read.

In the same vein, navigation code that adheres to good principles involves lots of jumping about.

---
# Why?

Writing tests is boring.

Writing integration tests are especially boring.

Reading integration tests is worse.

---
# Problems with Golang based integration tests

As mentioned, the tests are hard to read.

Code is easy to write, but hard to read.

No matter how extensive the documentation, you will have codebases that
do the same thing in vastly different ways.

So naturally a test framework written in native code will have
the same problems.

---
# Example

Take this example found from medium.com after searching "golang integration tests":

```go
func (s *e2eTestSuite) Test_EndToEnd_GetAllArticle() {
    article := model.Article{
        Title:   "my-title",
        Author:  "my-author",
        Content: "my-content",
    }

    s.NoError(s.dbConn.Create(&article).Error)

    req, err := http.NewRequest(echo.GET, fmt.Sprintf("http://localhost:%d/articles", s.port), nil)
    s.NoError(err)

    req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

    client := http.Client{}
    response, err := client.Do(req)
    s.NoError(err)
    s.Equal(http.StatusOK, response.StatusCode)

    byteBody, err := ioutil.ReadAll(response.Body)
    s.NoError(err)

    s.Equal(`{"status":200,"message":"Success","data":[{"id":1,"title":"my-title","content":"my-content","author":"my-author"}]}`, strings.Trim(string(byteBody), "\n"))
    response.Body.Close()
}
```

Almost all of this is boilerplate that needs to be duplicated in future tests.

---
# Example

Let's make it less boilerplate via well defined functions:

```go
func (s *e2eTestSuite) Test_EndToEnd_GetAllArticle() {
    s.NoError(s.dbInsert(model.Article{
            Title:   "my-title",
            Author:  "my-author",
            Content: "my-content",
        }
    ))

    respBody, err := s.httpGetJSONMustOK("/articles")
    s.NoError(err)

    s.Equal(`{"status":200,"message":"Success","data":[{"id":1,"title":"my-title","content":"my-content","author":"my-author"}]}`, strings.Trim(string(respBody), "\n"))
    response.Body.Close()
}
```

This is definitely better, but now we have a highly specific `httpGetJSONMustOK` function which assignes headers
and requires a 200 status code. If we want granularity we have to implement more specfic use case functions.

---
# Example

Well what about...
```go
func (s *e2eTestSuite) Test_EndToEnd_GetAllArticle() {
    s.NoError(s.dbInsert(model.Article{
            Title:   "my-title",
            Author:  "my-author",
            Content: "my-content",
        }
    ))

    headers := http.Header{}
    headers.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

    respBody, statusCode, err := s.httpGet("/articles", headers)
    s.NoError(err)
    s.Equal(http.StatusOK, response.StatusCode)

    s.Equal(`{"status":200,"message":"Success","data":[{"id":1,"title":"my-title","content":"my-content","author":"my-author"}]}`, strings.Trim(string(respBody), "\n"))
    response.Body.Close()
}
```

---
# What about this?

```gherkin
Scenario: I can retrieve all articles
  Given I seed the database with "insert_article.sql"
  When I send a GET request to "/articles"
  Then the http response code should be OK
  And the response should match "articles.json"
```

---
# Or this
Scenario: I can retrieve all articles
  Given an aritcle in the database
  When I retrieve all articles
  Then the http response code should be OK
  And the response should match all articles
