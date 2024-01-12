package django

import (
	"sync"
	"testing"

	"github.com/sjhitchner/gorm-to-django/gorm"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type DjangoSuite struct {
}

var _ = Suite(&DjangoSuite{})

func (s *DjangoSuite) SetUpSuite(c *C) {
}

func (s *DjangoSuite) TestModelGenerate(c *C) {
	// c.Assert(s.gen.GenerateModels(os.Stdout), IsNil)
}

func (s *DjangoSuite) TestAdminGenerate(c *C) {
	// c.Assert(s.gen.GenerateAdmin(os.Stdout), IsNil)
}

func (s *DjangoSuite) TestModel(c *C) {
	var wg sync.WaitGroup
	wg.Add(1)

	ch, err := gorm.Parse("../gorm/test")
	c.Assert(err, IsNil)

	out, errCh := gorm.Preprocess(ch)

	go func() {
		defer wg.Done()
		for err := range errCh {
			c.Assert(err, IsNil)
		}
	}()

	gen, err := New()
	c.Assert(err, IsNil)
	c.Assert(gen.Build(out), IsNil)

	wg.Wait()

	c.Assert(len(gen.Models), Equals, 8)

	for _, model := range gen.Models {
		switch model.Name {
		case "Event":
			checkEvent(c, &model)

		case "Country":
			checkCountry(c, &model)

		case "Genre":
			checkGenre(c, &model)

		case "Venue":
			checkVenue(c, &model)

		case "State":
			checkState(c, &model)

		case "City":
			checkCity(c, &model)

		case "Category":
			checkCategory(c, &model)

		case "EventCategories":
			checkEventCategories(c, &model)

		default:
			c.Errorf("unexpected model (%s)", model.Name)
		}
	}
}

func checkEvent(c *C, model *Model) {
	c.Assert(model.Name, Equals, "Event")
	c.Assert(len(model.Fields), Equals, 21)

	check4Field(c, model, "id", "int64", false, BigAutoField)
	check4Field(c, model, "name", "string", false, TextField)
	check4Field(c, model, "categories", "Category", false, ManyToMany)
	check4Field(c, model, "date_confirmed", "bool", false, BooleanField)
	check4Field(c, model, "end_date", "time.Time", true, DateTimeField)
	check4Field(c, model, "event_api", "string", false, TextField)
	check4Field(c, model, "external_id", "int64", false, BigIntegerField)
	check4Field(c, model, "genre", "Genre", false, ForeignKey)
	check4Field(c, model, "link", "string", false, TextField)
	check4Field(c, model, "min_ticket_price_amount", "float64", false, FloatField)
	check4Field(c, model, "min_ticket_price_currency_code", "string", false, TextField)
	check4Field(c, model, "min_ticket_price_display", "string", false, TextField)
	check4Field(c, model, "on_sale_date", "time.Time", true, DateTimeField)
	check4Field(c, model, "start_date", "time.Time", false, DateTimeField)
	check4Field(c, model, "status", "string", false, TextField)
	check4Field(c, model, "time_confirmed", "bool", false, BooleanField)
	check4Field(c, model, "type", "string", false, TextField)
	check4Field(c, model, "venue", "Venue", false, ForeignKey)
	check4Field(c, model, "created_at", "time.Time", true, DateTimeField)
	check4Field(c, model, "deleted_at", "time.Time", true, DateTimeField)
	check4Field(c, model, "updated_at", "time.Time", true, DateTimeField)
}

func checkCategory(c *C, model *Model) {
	c.Assert(model.Name, Equals, "Category")
	c.Assert(len(model.Fields), Equals, 3)
	check4Field(c, model, "id", "int64", false, BigAutoField)
	check4Field(c, model, "name", "string", false, TextField)
	check4Field(c, model, "role", "string", false, TextField)
}

func checkEventCategories(c *C, model *Model) {
	c.Assert(model.Name, Equals, "EventCategories")
	c.Assert(len(model.Fields), Equals, 3)
	check4Field(c, model, "id", "int64", false, BigAutoField)
	check4Field(c, model, "event", "Event", false, ForeignKey)
	check4Field(c, model, "category", "Category", false, ForeignKey)
}

func checkCountry(c *C, model *Model) {
	c.Assert(model.Name, Equals, "Country")
	c.Assert(len(model.Fields), Equals, 23)
}

func checkCity(c *C, model *Model) {
	c.Assert(model.Name, Equals, "City")
	c.Assert(len(model.Fields), Equals, 12)
}

func checkState(c *C, model *Model) {
	c.Assert(model.Name, Equals, "State")
	c.Assert(len(model.Fields), Equals, 11)
}

func checkGenre(c *C, model *Model) {
	c.Assert(model.Name, Equals, "Genre")
	c.Assert(len(model.Fields), Equals, 5)
}

func checkVenue(c *C, model *Model) {
	c.Assert(model.Name, Equals, "Venue")
	c.Assert(len(model.Fields), Equals, 13)
}

func check4Field(c *C, model *Model, name, typ string, nullable bool, djangoType DjangoFieldType) {
	for _, field := range model.Fields {
		if field.Name == name {
			c.Assert(field.Type, Equals, typ, Commentf(field.Name))
			c.Assert(field.IsNullable, Equals, nullable, Commentf(field.Name))

			djt, err := field.DjangoField()
			c.Assert(err, IsNil, Commentf(field.Name))
			c.Assert(djt, Equals, djangoType, Commentf(field.Name))

			return
		}
	}
	c.Errorf("missing field %s", name)
}
