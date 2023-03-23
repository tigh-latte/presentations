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

![22](src/imgs/scientist.png)

---
# What is BDD?

An extension of TDD (I know...) that emphasizes feature development around user stories.

I'm assuming we, for the most part, use user stories here.

---
# What is cucumber?

A test process written in the gherkin language, which expresses tests in natual language.

This brings the advantage of:

<!-- stop -->

- Non-technical members of an organisation can understand & contribute to test suites
<!-- stop -->

- Developers who are less familiar / experienced can understand & contribute to test suites
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
Cucumber tests have three main entities:

## Steps
A single instruction. These are written in natural language, and are akin to well define functions in an integration test suite, such as `s.createUser(dbConn)`.

<!-- stop -->

```gherkin
Given I login as "admin"
```
```gherkin
Then the response body should match "account_created.json"
```

<!-- stop -->

## Scenario
A collection of steps structured to test some functionality, or test case. Scenarios are given a description declaring what exactly they are testing.

Think of a scenario as a `Test` function in your integration test suite, `func TestAccount_Create(t *testing.T)`. Like a test function, they are independent and can be executed concurrently.

<!-- stop -->

```file
path: src/examples/bank/test/integration/features/account.feature
lang: gherkin
lines:
    start: 8
    end: 16
```

<!-- stop -->

## Features
The feature being tested.  Features give a description of what should be tested, and can declare rules / criteria for documentation purposes.

Features are akin to `_test.go` files in your integration test suite.

<!-- stop -->

```file
path: src/examples/bank/test/integration/features/account.feature
lang: gherkin
lines:
    start: 0
    end: 7
```

---
# What is cucumber?

Putting this all this together we get:
```file
path: src/examples/bank/test/integration/features/account.feature
lang: gherkin
```

---
# Ok, where does go come into this?

Obviously we can't just run these feature files, we need to somehow link our natural langauge with code. Ideally, we'd use a framework provided from the cucumber group.

Luckily, we have just that.

<!-- stop -->

![15](src/imgs/godog_logo.png) 

https://github.com/cucumber/godog

---

# Using godog

There are two ways to bake godog into your test suite.

1. Using \*testing.T
<!-- stop -->

1. Using \*testing.M
<!-- stop-->

Both of these are fine, but for this talk we're doing to focus on \*testing.M as it is, in my opinion, easier for integration testing.

---
# What is cucumber?

Imagine we an account entity, and wanted to test its `Withdraw` functionality.

When funds are withdrawn, the account balance is updated to reflect this withdrawal.

If the funds requested exceed the balance, the balance isn't touched, instead we relay an error.

<!-- stop -->

The unit tests for this service may look like this:

```file
path: code/ex1/account_test.go
transform: sed 's/\t/    /g'
lang: go
lines:
  start: 9
  end: null
```

---
# What is cucumber?
The same service, using a cucumber framework, would look like so:

```file
path: code/ex2/features/account.feature
lang: gherkin
transform: 'sed /@account/d'
lines:
   start: 0
   end: 14
```
<!-- stop -->

Much easier read.

---

# Show me the golang!

## Getting started

To begin, we create our `TestMain` function and build our test suite struct:

```file
path: code/ex2/account_test.go
lang: go
transform: sed 's/\t/   /g;/flag/d'
lines:
  start: 21
  end: 35
```

---
# Show me the golang!

## Wiring up our steps

Inside the scenario initialiser we map all naturual langauge steps to a corresponding golang function.

<!-- stop -->

```file
path: code/ex2/account_test.go
lang: go
transform: sed 's/\t/   /g;/BeforeStep\|time.Sleep\|})\|deposit/d;'
lines:
  start: 39
  end: 45
```

---
# Show me the golang!

## Wiring up our steps

We can also define `BeforeStep`, `AfterStep`, `BeforeScenario` and `AfterScenario` hooks.

```file
path: code/ex2/account_test.go
lang: go
transform: sed 's/\t/   /g;/deposit/d;'
lines:
  start: 39
  end: 48
```

---

# Show me the golang!
## Wiring up our steps

Step functions can:
- take `context.Context` as their first argument. If this is present, godog will build and supply a context.
- return `context.Context` as their first return. If this is present, godog will pass the returned context to future steps.
- take typed parameters.

```file
path: code/ex2/account_test.go
lang: go
transform: sed 's/\t/   /g'
lines:
  start: 49
  end: 87
```

---
# Show me the golang!

## And that's it!

We now have a test suite we can run.

<!-- stop -->

```terminal-ex
command: zsh -il
rows: 30
init_text: cd code/ex2/; go test -v ./account_test.go
init_wait: '> '
init_codeblock_lang: zsh
```

---
# Running a subset
With the `testing.M` approach, all test execution is done via the godog api.

Luckily godog exposes a nice way to run a subset of tests.

<!-- stop -->

Let's say we added new `Deposit` functionality to our service, and wanted to run only these deposit tests:

```file
path: code/ex2/features/account.feature
lang: gherkin
transform: sed '/ *\(When\|Then\|And\)/d;s/^.*Given.*$/   \# Scenario definition/g'
```

---
# Running a subset

We can then selectively run these scenarios with the `--godog.tags=` flag:

<!-- stop -->

```terminal-ex
command: zsh -il
rows: 30
init_text: cd code/ex2/; go test -v ./account_test.go --godog.tags=deposit
init_wait: '> '
init_codeblock_lang: zsh
```

---
# Running randomly

Scenarios can be executed in a random order using `--godog.random`, ensuring tests are not dependent on one another

<!-- stop -->

```terminal-ex
command: zsh -il
rows: 30
init_text: cd code/ex2/; go test -v ./account_test.go --godog.random
init_wait: '> '
init_codeblock_lang: zsh
```

---
# Running concurrently

Scenarios can be executed concurrently using `--godog.concurrency=num_procs`:

<!-- stop -->

```terminal-ex
command: zsh -il
rows: 30
init_text: cd code/ex2/; go test -v ./account_test.go --godog.concurrency=4
init_wait: '> '
init_codeblock_lang: zsh
```

---
# So this relatively contrived example is cool and all, but like, we're still a long way away from integration tests!

Some basic functions of integration test suites are:
1. Direct database access to prime data
<!-- stop -->

2. Some http or grpc client to perform requests
<!-- stop -->

3. Test data, in the form of structs in the code, or files on the filesystem.

<!-- stop -->

So let's add these!

---

# Build a global suite

We're going to want a test suite struct that holds global clients to be used across scenarios:

```file
path: src/examples/bank/test/integration/integration_test.go
lang: go
transform: sed 's/\t/    /g'
lines:
  start: 47
  end: 58
```

<!-- stop -->

```file
path: src/examples/bank/test/integration/data/data.go
lang: go
transform: sed 's/\t/    /g'
lines:
  start: 7
  end: null
```

<!-- stop -->

```file
path: src/examples/bank/test/integration/integration_test.go
lang: go
transform: sed 's/\t/    /g'
lines:
  start: 25
  end: 31
```

---

# Let's build this suite

```file
path: src/examples/bank/test/integration/integration_test.go
lang: go
transform: sed 's/\t/    /g;s/err/_/;/lookatme:ignore/d'
lines:
  start: 63
  end: 97
```

---

# Define our steps

```file
path: src/examples/bank/test/integration/integration_test.go
lang: go
transform: sed 's/\t/    /g;;/lookatme:ignore/d'
lines:
  start: 97
  end: 107
```

<!-- stop -->

Notice that the step funcs now have `Suite` as a receiver. This is so they can access the db connection and http client.

---

# Running SQL from a database connection

```gherkin
Given I run the SQL "init_dev_data.sql"
```

<!-- stop -->

```file
path: src/examples/bank/test/integration/integration_test.go
lang: go
transform: sed 's/\t/    /g;;/lookatme:ignore/d'
lines:
  start: 106
  end: 117
```

---

# Setting headers for a http request

```gherkin
Given the headers:
  | key           | value      |
  | Authorization | Bearer dev |
  | X-Request-ID  | 12345      |
```

<!-- stop -->

```file
path: src/examples/bank/test/integration/integration_test.go
lang: go
transform: sed 's/\t/    /g;/lookatme:ignore/d'
lines:
  start: 206
  end: 222
```

---

# Making a http request

```gherkin
When I make a POST request to "/api/v1/myendpoint" using "POST-my-request.json"
```

<!-- stop -->

```file
path: src/examples/bank/test/integration/integration_test.go
lang: go
transform: sed 's/\t/    /g;;/lookatme:ignore/d'
lines:
  start: 121
  end: 155
```

---

# Checking the response code (incomplete)
```gherkin
Then the response status should be BAD_REQUEST
```

<!-- stop -->

```file
path: src/examples/bank/test/integration/integration_test.go
lang: go
transform: sed 's/\t/    /g;;/lookatme:ignore/d'
lines:
  start: 154
  end: 173
```

---

# Finally, checking the response body
```gherkin
Then the response body should match "errs/bad-request.json"
```

<!-- stop -->

```file
path: src/examples/bank/test/integration/integration_test.go
lang: go
transform: sed 's/\t/    /g;/lookatme:ignore/d'
lines:
  start: 173
  end: 206
```

---

# Filesystem walkthrough
```terminal-ex
command: zsh -il
rows: 34
init_text: cd src/examples/bank
init_wait: '> '
init_codeblock_lang: zsh
```

---

# What's possible

HTTP and DB connections are cool and all, but we can take this much further.

## S3 actions

```gherkin
Given I put the files into buckets:
  | bucket       | file         |
  | my-bucket    | file.txt     |
  | other-bucket | holyhell.jpg |
```

<!-- stop -->

```go
func (s *Suite) IPutFilesIntoBuckets(ctx context.Context, table *godog.Table) error {
    for _, row := range table.Rows[1:] {
        bucket := row.Cells[0].Value
        file := row.Cells[1].Value

        bb, err := s.testData.Load(file)
        if err != nil {
            return errors.Wrapf(err, "failed to load file '%s'", row.Cells[1].Value)
        }
        if _, err = s.s3.PutObject(ctx, &s3.PutObjectInput{
                Bucket: aws.String(bucket),
                Body:   bytes.NewReader(bb),
                Key:    aws.String(file),
        }); err != nil {
            return errors.Wrap(err, "failed to put file")
        }
    }
    return nil
}
```

---

# What's possible

## S3 actions

```gherkin
Given I delete these files from buckets:
  | my-bucket    | file.txt     |
  | other-bucket | holyhell.jpg |
```

```go
func (s *Suite) IDeleteFilesFromS3(ctx context.Context, table *godog.Table) error {
    for _, row := range table.Rows[1:] {
        file := row.Cells[0].Value
        bucket := row.Cells[1].Value

        list, err := s3Client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
            Bucket: aws.String(bucket),
            Prefix: aws.String(file),
        })
        if err != nil {
            return errors.Wrap(err, "failed to list files to be deleted")
        }

        objects := make([]types.ObjectIdentifier, len(list.Contents))
        for i, v := range list.Contents {
            objects[i] = types.ObjectIdentifier{
                Key: v.Key,
            }
        }

        if len(objects) == 0 {
            return nil
        }

        if _, err := s.s3.DeleteObjects(ctx, &s3.DeleteObjectsInput{
            Bucket: aws.String(file),
            Delete: &types.Delete{
                Objects: objects,
            },
        }); err != nil {
            return errors.Wrap(err, "failed to put file")
        }
    }
    return nil
}
```

---

# What's possible

## RabbitMQ fun

We can get some really nice queue monitoring functionality with a little bit of setup work.

<!-- stop -->

Up front, we declare what queues on which exchanges we want the test suite to subscribe to

```go
var rmqSubs = []RabbitMQSubscription{{
    Exchange:   "accounts=exchange",
    RoutingKey: "bank.account.created",
}}
```

<!-- stop -->

Then we build a map of exchange to queuegroup subscriptions

```go
    m := map[string][]string{}
    for _, sub := range rmqSubs {
        if _, ok := m[sub.Exchange]; !ok {
            m[sub.Exchange] = []string{}
        }
        m[sub.Exchange] = append(m[sub.Exchange], sub.RoutingKey)
    }

    for exch, routingKeys := range m {
        qg := // create queuegroup for routingKeys

        s.rmqExchMap.Store(exch, qg)
    }
```

<!-- stop -->

THEN, the queue group has its own internal mapping of routing keys to channels.

```go
type queueGroup struct {
    routingKeys []string
    chans       syncmap.NewMap[string, chan amqp.Delivery](),
}

func (q *queueGroup) Connect() {
    for _, rk := range q.routingKeys {
        ch, _ := rmq.Consume(/* args */)
        q.chans.Store(rk, ch)
    }
}

func (q *queueGroup) C(name string) <-chan amqp.Delivery {
    ch, _ := q.chans.Get(name)
    return ch
}
```

---

# What's possible

## Reading a message from a queue

```gherkin
When I wait for a message on "accounts-exchange/bank.account.created"
```

<!-- stop -->

```go
func (s *Suite) IReadARabbitMQMessage(ctx, rk string) (context.Context, error) {
    qg, ok := s.rmqExchMap.Get(exch)
    if !ok {
        return errors.Errorf("exchange '%s' not registed: %#v", exch, exchangeMap)
    }

    tCtx, cancel := context.WithTimeout(ctx, 10 * time.Second)
    defer cancel()

    select {
    case msg = <-qg.C(rk):
        return context.WithValue(ctx, rmqMsgKey{}, msg)
    case <-tCtx.Done():
        return ctx, errors.New("timed out")
    }

    return ctx, nil
}
```

# What's possible

## Check the amount of messages in a queue

```gherkin
Then there should be 7 messages on "accounts-exchange/bank.account.updated"
```

```go
func (s *Suite) ThereShouldNBeRabbitMQMessages(ctx context.Context, exp int, exch, rk string) error {
    qg, ok := exchangeMap.Get(exch)
    if !ok {
        return errors.Errorf("exchange '%s' not registed: %#v", exch, exchangeMap)
    }

    total := len(qg.C(rk))
    if exp != total {
        return errors.Errorf("unexpected number of messages on queue '%s'. wanted '%d' got '%d'", rk, exp, total)
    }

    return nil
}
```

---

# What's possible

## Watch this space

SO, I am currently working on open sourcing a library that provides a lot of this functionality out of the box.

Hopefully V2 of this presentation will have a working example of most of these concepts.

---

# THANK YOU

That's a wrap.

<!-- stop -->

Holy hell.
