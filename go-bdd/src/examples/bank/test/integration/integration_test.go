package integration_test

import (
	"bytes"
	"context"
	"database/sql"
	"embed"
	"encoding/json"
	"flag"
	"io"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/tigh-latte/presentations/go-bdd/src/examples/bank/test/integration/data"
	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/pkg/errors"
	"github.com/yazgazan/jaydiff/diff"

	_ "github.com/lib/pq"
)

//go:embed sql/*
var sqlData embed.FS

//go:embed http/*
var httpData embed.FS

type (
	httpCodeKey struct{}
	httpRespKey struct{}
)

var opts = godog.Options{
	Output: colors.Colored(os.Stdout),
	Format: "pretty",
	Paths:  []string{"features"},
}

func init() {
	godog.BindFlags("godog.", flag.CommandLine, &opts)
}

type Suite struct {
	url   url.URL
	suite godog.TestSuite
	db    *sql.DB
	sql   data.DataStore
	reqs  data.DataStore
	resps data.DataStore

	http *http.Client
}

func Value[K, V any](ctx context.Context) V {
	return ctx.Value(*new(K)).(V)
}

func TestMain(m *testing.M) {
	flag.Parse()

	db, err := sql.Open("postgres", "postgresql://dev:dev@0.0.0.0/dev?sslmode=disable")
	if err != nil {
		panic(err)
	}

	s := &Suite{
		url: url.URL{
			Scheme: "http",
			Host:   "localhost:8081",
		},
		suite: godog.TestSuite{
			Name:    "bank",
			Options: &opts,
		},
		db: db,
		sql: &data.DataDir{
			Prefix: "sql",
			FS:     sqlData,
		},
		reqs: &data.DataDir{
			Prefix: "http/requests",
			FS:     httpData,
		},
		resps: &data.DataDir{
			Prefix: "http/responses",
			FS:     httpData,
		},
		http: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
	s.suite.ScenarioInitializer = s.initScenario

	os.Exit(s.suite.Run())
}

func (s *Suite) initScenario(sc *godog.ScenarioContext) {
	sc.Step(`^I run the SQL "([^"]+)"$`, s.iRunTheSQL)
	sc.Step(`^I make a (GET|POST) request to "([^"]+)"$`, s.iMakeARequestTo)
	sc.Step(`^I make a (GET|POST) request to "([^"]+)" using "([^"]+)"$`, s.iMakeARequestToUsing)
	sc.Step(`^the response status is (\w+)`, s.theResponseCodeShouldBe)
	sc.Step(`^the response status should be (\w+)`, s.theResponseCodeShouldBe)
	sc.Step(`^the response body should match "([^"]+)"`, s.theResponseBodyShouldMatch)
}

func (s *Suite) iRunTheSQL(ctx context.Context, file string) error {
	sql, err := s.sql.Load(file)
	if err != nil {
		return errors.Wrapf(err, "failed to load sql file '%s'", file)
	}

	_, err = s.db.ExecContext(ctx, string(sql))
	return errors.Wrapf(err, "failed to execute sql file '%s'", file)
}

func (s *Suite) iMakeARequestTo(ctx context.Context, verb, endpoint string) (context.Context, error) {
	req, err := http.NewRequest(verb, s.url.JoinPath(endpoint).String(), nil)
	if err != nil {
		return ctx, errors.Wrapf(err, "failed to make request '%s %s'", verb, endpoint)
	}

	resp, err := s.http.Do(req)
	if err != nil {
		return ctx, errors.Wrapf(err, "failed to perform request '%s %s'", verb, endpoint)
	}

	bb, err := io.ReadAll(resp.Body)
	if err != nil {
		return ctx, errors.Wrapf(err, "failed to read body for request '%s %s'", verb, endpoint)
	}

	ctx = context.WithValue(ctx, httpCodeKey{}, resp.StatusCode)
	return context.WithValue(ctx, httpRespKey{}, bb), nil
}

func (s *Suite) iMakeARequestToUsing(ctx context.Context, verb, endpoint, file string) (context.Context, error) {
	reqBody, err := s.reqs.Load(file)
	if err != nil {
		return ctx, errors.Wrapf(err, "failed to load req '%s'", file)
	}

	req, err := http.NewRequest(verb, s.url.JoinPath(endpoint).String(), bytes.NewReader(reqBody))
	if err != nil {
		return ctx, errors.Wrapf(err, "failed to make request '%s %s'", verb, endpoint)
	}

	resp, err := s.http.Do(req)
	if err != nil {
		return ctx, errors.Wrapf(err, "failed to perform request '%s %s'", verb, endpoint)
	}

	bb, err := io.ReadAll(resp.Body)
	if err != nil {
		return ctx, errors.Wrapf(err, "failed to read body for request '%s %s'", verb, endpoint)
	}

	ctx = context.WithValue(ctx, httpCodeKey{}, resp.StatusCode)
	return context.WithValue(ctx, httpRespKey{}, bb), nil
}

func (s *Suite) theResponseCodeShouldBe(ctx context.Context, status string) error {
	want, ok := map[string]int{
		"OK":          http.StatusOK,
		"CREATED":     http.StatusCreated,
		"BAD_REQUEST": http.StatusBadRequest,
		"CONFLICT":    http.StatusConflict,
	}[status]
	if !ok {
		return errors.Errorf("unrecognised response code '%s'", status)
	}
	statusCode := Value[httpCodeKey, int](ctx)
	if statusCode != want {
		return errors.Errorf("unexpected status code. want %d got %d", want, statusCode)
	}

	return nil
}

func (s *Suite) theResponseBodyShouldMatch(ctx context.Context, file string) error {
	respBody := Value[httpRespKey, []byte](ctx)

	expBody, err := s.resps.Load(file)
	if err != nil {
		return errors.Wrapf(err, "failed to load file '%s'", file)
	}

	var want, got interface{}
	if err = json.Unmarshal(respBody, &want); err != nil {
		return errors.Wrap(err, "failed to unmarshal response body")
	}
	if err = json.Unmarshal(expBody, &got); err != nil {
		return errors.Wrapf(err, "failed to unmarshal body in file '%s'", file)
	}

	d, err := diff.Diff(want, got)
	if err != nil {
		return errors.Wrap(err, "error calculating diff")
	}
	if d.Diff() == diff.Identical {
		return nil
	}
	report := d.StringIndent("", "", diff.Output{
		Indent:     "  ",
		ShowTypes:  true,
		JSON:       true,
		JSONValues: true,
	})

	return errors.Errorf(report)
}
