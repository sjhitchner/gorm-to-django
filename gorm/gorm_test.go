package gorm

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type GORMSuite struct{}

var _ = Suite(&GORMSuite{})

func (s *GORMSuite) TestPreprocess(c *C) {
	var wg sync.WaitGroup
	wg.Add(1)

	ch := structGenerator(c)

	out, errCh := Preprocess(ch)

	go func() {
		defer wg.Done()
		for err := range errCh {
			c.Assert(err, IsNil)
		}
	}()

	for o := range out {
		switch o.Name {
		case "Category":
			checkCategory(c, o)
		case "Event":
			checkEvent(c, o)
		case "Money":
			checkMoney(c, o)
		}
	}

	wg.Wait()
}

func checkEvent(c *C, s Struct) {
	c.Assert(s.Name, Equals, "Event")
	c.Assert(len(s.Fields), Equals, 17)
}

func checkCategory(c *C, s Struct) {
	fmt.Println(s)
	c.Assert(s.Name, Equals, "Category")
	c.Assert(len(s.Fields), Equals, 4)
}

func checkMoney(c *C, s Struct) {
	c.Assert(s.Name, Equals, "Money")
	c.Assert(len(s.Fields), Equals, 3)
}

func structGenerator(c *C) <-chan Struct {
	out := make(chan Struct)

	go func() {
		defer close(out)

		var strts []Struct
		c.Assert(json.NewDecoder(strings.NewReader(Structs)).Decode(&strts), IsNil)

		for _, s := range strts {
			out <- s
		}
	}()
	return out
}

var Structs = `[
{
  "IsModel": true,
  "Name": "Event",
  "Metadata": {
    "tablename": "events"
  },
  "Fields": [
    {
      "Name": "ID",
      "Type": "int64",
      "Tags": {
        "primaryKey": {
          "Name": "primaryKey",
          "Value": ""
        }
      }
    },
    {
      "Name": "EventAPI",
      "Type": "string",
      "Tags": {
        "index": {
          "Name": "index",
          "Value": ""
        }
      }
    },
    {
      "Name": "Name",
      "Type": "string",
      "Tags": null
    },
    {
      "Name": "Link",
      "Type": "string",
      "Tags": null
    },
    {
      "Name": "StartDate",
      "Type": "time.Time",
      "Tags": {
        "index": {
          "Name": "index",
          "Value": ""
        }
      }
    },
    {
      "Name": "EndDate",
      "Type": "*time.Time",
      "Tags": null
    },
    {
      "Name": "OnSaleDate",
      "Type": "*time.Time",
      "Tags": {
        "index": {
          "Name": "index",
          "Value": ""
        }
      }
    },
    {
      "Name": "DateConfirmed",
      "Type": "bool",
      "Tags": null
    },
    {
      "Name": "TimeConfirmed",
      "Type": "bool",
      "Tags": null
    },
    {
      "Name": "Type",
      "Type": "string",
      "Tags": {
        "index": {
          "Name": "index",
          "Value": ""
        }
      }
    },
    {
      "Name": "MinTicketPrice",
      "Type": "*Money",
      "Tags": {
        "embedded": {
          "Name": "embedded",
          "Value": ""
        },
        "embeddedPrefix": {
          "Name": "embeddedPrefix",
          "Value": "min_ticket_price_"
        }
      }
    },
    {
      "Name": "Status",
      "Type": "string",
      "Tags": null
    },
    {
      "Name": "GenreID",
      "Type": "int64",
      "Tags": null
    },
    {
      "Name": "Genre",
      "Type": "*Genre",
      "Tags": null
    },
    {
      "Name": "VenueID",
      "Type": "int64",
      "Tags": null
    },
    {
      "Name": "Venue",
      "Type": "*Venue",
      "Tags": null
    },
    {
      "Name": "Categories",
      "Type": "[]*Category",
      "Tags": {
		"constraint": {
          "Name": "constraint",
          "Value": "OnUpdate:CASCADE,OnDelete:CASCADE"
        },
        "foreignKey": {
          "Name": "foreignKey",
          "Value": "EventID"
        }
      }
    },
    {
      "Name": "CreatedAt",
      "Type": "*time.Time",
      "Tags": {
        "autoCreateTime": {
          "Name": "autoCreateTime",
          "Value": ""
        }
      }
    },
    {
      "Name": "UpdatedAt",
      "Type": "*time.Time",
      "Tags": {
        "autoUpdateTime": {
          "Name": "autoUpdateTime",
          "Value": ""
        }
      }
    },
    {
      "Name": "DeletedAt",
      "Type": "*time.Time",
      "Tags": null
    }
  ]
},
{
  "IsModel": false,
  "Name": "Money",
  "Metadata": {},
  "Fields": [
    {
      "Name": "Amount",
      "Type": "float64",
      "Tags": null
    },
    {
      "Name": "CurrencyCode",
      "Type": "string",
      "Tags": null
    },
    {
      "Name": "Display",
      "Type": "string",
      "Tags": null
    }
  ]
},
{
  "IsModel": true,
  "Name": "Category",
  "Metadata": {
    "tablename": "categories"
  },
  "Fields": [
    {
      "Name": "ID",
      "Type": "int64",
      "Tags": {
        "primaryKey": {
          "Name": "primaryKey",
          "Value": ""
        }
      }
    },
    {
      "Name": "Name",
      "Type": "string",
      "Tags": null
    },
    {
      "Name": "Role",
      "Type": "string",
      "Tags": null
    }
  ]
}
]`
