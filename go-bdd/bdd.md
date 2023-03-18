---
theme: ./theme.json
---
# Welcome!

BDD Development in Golang, using Cucumber

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
# 1. Is understandable.

Consider:
```go
// interfaces.go
type Randomable [T any] interface {
    Randomise() T
}

func Randomise[T Randomable[T]]() T {
    return (*new(T)).Randomise()
}

func RandomiseMany[T Randomable[T]](total int) []T {
    tt := make([]T, len(total))
    for i := 0; i < total; i++ {
        tt[i] = Randomise[T]()
    }

    return tt
}

// models.go
type Person struct {
    Name    string  `json:"name"`
    Age     int16   `json:"age"`
    Address Address `json:"address"`
}

var (p Person) Randomise() Person {
    p.Name = gofakeit.Name()
    p.Age = rand.Intn(100) + 1

    addr := gofakeit.Address()
    p.Address = Address {
        Street:   addr.Street,
        PostCode: addr.Zip,
    }

    return p
}

// person_test.go
func (t *testSuite) Test_CreatePerson() {
    person := testutil.Random[models.Person]()
    // Make http request with person
}
```

---
# Why?

Writing tests is boring.

---
# Why?

Writing tests is boring.

Writing integration tests are especially boring.

---
# Why?

Writing tests is boring.

Writing integration tests are especially boring.

Reading integration tests is worse.

---
# Problems with Golang based integration tests

As mentioned, the tests are hard to read.

---
# Problems with Golang based integration tests

As mentioned, the tests are hard to read.

Code is easy to write, but hard to read.

---
# Problems with Golang based integration tests

As mentioned, the tests are hard to read.

Code is easy to write, but hard to read.

No matter how extensive the documentation, you will have codebases that
do the same thing in vastly different ways.

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
