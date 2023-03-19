---
theme: ./theme.json
title: Go BDD
author: Tighearnán Carroll
extensions:
  - image_ueberzug
  - terminal
  - file_loader
---
# Welcome!

BDD Development in Golang, using Cucumber

---
# What is BDD?

An extension of TDD (I know...) that emphasizes feature development around user stories.

I'm assuming we, for the most part, use user stories here.

---
# What is cucumber?

A test process written in the gherkin language, which expresses tests in natual language.

This brings the advantage of:

<!-- stop -->

- Non-technical members of an organisation can understand / contribute to test suites
<!-- stop -->

- Developers who are less familiar / experienced can understand / contribute to test suites
<!-- stop -->

- The tests themselves can become firm, unambiguous acceptance requirements


---
# What is cucumber?

Gherkin has the following keywords:

```gherkin
Given a precondition
 And another precondition
When I do an action
Then I should have some result
 And I should have some other result
```

---
# What is cucumber?

Gherkin has the following keywords:

```gherkin
Given I seed the database with "some_test_data.sql"
 And I login as user "admin"
When I send a GET request to "/contacts"
Then the response code should be OK
 And the response body should match "some_truth_source.json"
```

---
# What is cucumber?
Cucumber tests have three entities:
Feature -> Scenario -> Steps

A feature has scenarios. And a scenario has steps.

<!-- stop -->

| cucumber | golang                            |
|----------|-----------------------------------|
| feature  | some\_test.go                     |
| scenario | `func Test_Something(t *testing.T)` |
| step     | `s.httpGet(req)`                  |

---
# What is cucumber?

Features a store in txt files, typically having in files with a `.feature` extension.

```terminal-ex
command: bash -il
rows: 30
init_text: Simple package structure
init_codeblock_lang: bash
```

---
# What is cucumber?

Cucumber works by regex matching steps with functions, and passes captured expressions into the associated function:

```gherkin
Given I login as user "admin"
```

```go
// implementation
sc.Step(`^I login as user "([^"]+)"$`, func(user string) error {
    // perform login
})
```

---
# What is cucumber?

Imagine we an account entity, and wanted to test its `Withdraw` functionality.

When funds are withdrawn, the account balance is updated to reflect this withdrawal.

If the funds requested exceed the balance, the balance isn't touched, instead we relay an error.

<!-- stop -->

The unit tests for this service may look like this:

```file
path: code/ex1/account_test.go
lang: go
```

---
# What is cucumber?
The same service, using a cucumber framework, would look like so:

```file
path: code/ex2/features/account.feature
lang: gherkin
```

<!-- stop -->

Much easier read.

---

So how do we do this?

<!-- stop -->

Surely there must be code _somewhere_???

<!-- stop -->

The answer:

https://github.com/cucumber/godog

![7](src/godog_logo.png)

---

# Godog

Godog is the golang implementation of cucumber.

As you'd guess, behind the scenes there _is_ code.

When a scenario is initialised, its available steps are defined via regex matching, and associate that step with a function:

```go
// ex2

func InitScenario(sc *godog.ScenarioContext) {
    ctx := struct {
        records []string
    }{}

    sc.Step(`^I say "([^"]+)" (\d+) times$`, func(phrase string, times int64) {
        for i := 0; i < times; i++ {
            ctx.records = append(ctx.records, phrase)
        }
    })

    sc.Step(`^I should have (\d+) records of "([^"]+)"$`, func(expected int64, phrase string) error {
        if expected != len(ctx.records) {
            return errors.Errorf("expected %d records, got %d", expected, len(ctx.records))
        }

        for i, record := range ctx.records {
            if record != phrase {
                return errors.Errorf("expected record %d to be %s, got %s", i, phrase, record)
            }
        }
    })
}
```

```gherkin
Scenario: I can speak
  Given I say "hello world" 7 times
  Then I should have 7 records of "hello world"
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
