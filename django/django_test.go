package django

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/sjhitchner/gorm-to-django/gorm"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type DjangoSuite struct {
	gen *Generator
}

var _ = Suite(&DjangoSuite{})

func (s *DjangoSuite) SetUpSuite(c *C) {
	in := structGenerator(c)

	gen, err := New()
	c.Assert(err, IsNil)
	c.Assert(gen.Build(in), IsNil)

	s.gen = gen
}

func (s *DjangoSuite) TestModelGenerate(c *C) {
	c.Assert(s.gen.GenerateModels(os.Stdout), IsNil)
}

// TestModel is a test method for the DjangoSuite
func (s *DjangoSuite) TestModel(c *C) {
	for _, model := range s.gen.Models {
		switch model.Name {
		case "Event":
			checkEvent(c, model)

		case "Category":
			checkCategory(c, model)

		default:
			c.Errorf("unexpected model (%s)", model.Name)
		}
	}
}

func checkEvent(c *C, model Model) {
	fmt.Println(model)
	c.Assert(model.Name, Equals, "Event")

	/*
		c.Assert(len(model.Fields), Equals, 20)

		check4Field(c, model, "min_ticket_price_amount", "float64", false)
		check4Field(c, model, "min_ticket_price_currency_code", "string", false)
		check4Field(c, model, "min_ticket_price_display", "string", false)
		check4Field(c, model, "genre", "ForeignKey", false)
		check4Field(c, model, "venue", "ForeignKey", false)
	*/
}

func checkCategory(c *C, model Model) {
	fmt.Println(model)
	c.Assert(model.Name, Equals, "Category")
}

func check4Field(c *C, model *Model, name, typ string, nullable bool) {
	for _, field := range model.Fields {
		if field.Name == name {
			c.Assert(typ, Equals, field.Type)
			return
		}
	}
	c.Errorf("missing field %s", name)
}

func structGenerator(c *C) <-chan gorm.Struct {
	out := make(chan gorm.Struct)

	go func() {
		defer close(out)

		var strts []gorm.Struct
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
