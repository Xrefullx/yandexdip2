package sqllte

import (
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

const (
	dbFile = "test.db"
)

type TestSuite struct {
	suite.Suite
	storage *Storage
}

func (s *TestSuite) SetupSuite() {
	st, err := NewStorage(dbFile)
	s.Require().NoError(err)

	s.storage = st
}

func (s *TestSuite) TearDownSuite() {
	s.storage.Close()
	s.Require().NoError(os.Remove(dbFile))
}

func TestSuite_SQLliteStorage(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
