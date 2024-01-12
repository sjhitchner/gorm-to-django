package gorm

import (
	"fmt"
	"sync"
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type GORMSuite struct{}

var _ = Suite(&GORMSuite{})

var modelLookup = map[string]bool{
	"Category":        true,
	"Money":           true,
	"Event":           true,
	"Venue":           true,
	"Genre":           true,
	"Country":         true,
	"State":           true,
	"City":            true,
	"EventCategories": true,
}

func (s *GORMSuite) TestRecursiveName(c *C) {
	out, err := Parse("test")
	c.Assert(err, IsNil)

	var count int
	for o := range out {
		count++
		c.Assert(modelLookup[o.Name], Equals, true)
	}

	c.Assert(count, Equals, 8)
}

func (s *GORMSuite) TestPreprocess(c *C) {
	var wg sync.WaitGroup
	wg.Add(1)

	ch, err := Parse("test")
	c.Assert(err, IsNil)

	out, errCh := Preprocess(ch)

	go func() {
		defer wg.Done()
		for err := range errCh {
			c.Assert(err, IsNil)
		}
	}()

	var count int
	for o := range out {
		switch o.Name {
		case "Category":
			checkCategory(c, o)
		case "Event":
			checkEvent(c, o)
		case "Money":
			checkMoney(c, o)
		case "EventCategories":
			checkEventCategories(c, o)
		}
		c.Assert(modelLookup[o.Name], Equals, true)
		count++
	}
	c.Assert(count, Equals, 9)

	wg.Wait()
}

func checkEvent(c *C, s Struct) {
	fmt.Println(s)
	c.Assert(s.Name, Equals, "Event")
	c.Assert(len(s.Fields), Equals, 19)
}

func checkCategory(c *C, s Struct) {
	c.Assert(s.Name, Equals, "Category")
	c.Assert(len(s.Fields), Equals, 3)
}

func checkMoney(c *C, s Struct) {
	c.Assert(s.Name, Equals, "Money")
	c.Assert(len(s.Fields), Equals, 3)
}

func checkEventCategories(c *C, s Struct) {
	c.Assert(s.Name, Equals, "EventCategories")
	c.Assert(len(s.Fields), Equals, 3)
}
