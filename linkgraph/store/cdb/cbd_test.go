package cdb

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/Waqas-Shah-42/Links-R-Us-2/linkgraph/graph/graphtest"
	gc "gopkg.in/check.v1"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var _ = gc.Suite(new(CockroachDbGraphTestSuite))

func Test(t *testing.T) { gc.TestingT(t) }

type CockroachDbGraphTestSuite struct {
	graphtest.SuiteBase
	db *sql.DB
}

func (s *CockroachDbGraphTestSuite) SetUpSuite(c *gc.C) {
	dsn := os.Getenv("CDB_DSN")
	if dsn == "" {
		c.Skip("Missing CDB_DSN envvar; skipping cockroachdb-backed graph test suite")
	}

	//migrating
	path, err := os.Getwd()
	if err != nil {
    	log.Println(err)
	}
	fmt.Println("Current working directory: %w", path)


	m, err := migrate.New("file://migrations", dsn)
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil {
		log.Fatal(err)
	}

	g, err := NewCockroachDbGraph(dsn)
	c.Assert(err, gc.IsNil)
	s.SetGraph(g)
	s.db = g.db
}

func (s *CockroachDbGraphTestSuite) SetUpTest(c *gc.C) {
	s.flushDB(c)
}

func (s *CockroachDbGraphTestSuite) TearDownSuite(c *gc.C) {
	if s.db != nil {
		s.flushDB(c)
		c.Assert(s.db.Close(), gc.IsNil)
	}
}

func (s *CockroachDbGraphTestSuite) flushDB(c *gc.C) {
	_, err := s.db.Exec("DELETE FROM links")
	c.Assert(err, gc.IsNil)
	_, err = s.db.Exec("DELETE FROM edges")
	c.Assert(err, gc.IsNil)
}
