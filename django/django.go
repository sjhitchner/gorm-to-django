package django

import (
	"encoding/json"
	"fmt"
)

type DjangoFieldType string

const (
	AutoField                 DjangoFieldType = "AutoField"
	BigAutoField              DjangoFieldType = "BigAutoField"
	BigIntegerField           DjangoFieldType = "BigIntegerField"
	BinaryField               DjangoFieldType = "BinaryField"
	BooleanField              DjangoFieldType = "BooleanField"
	CharField                 DjangoFieldType = "CharField"
	DateField                 DjangoFieldType = "DateField"
	DateTimeField             DjangoFieldType = "DateTimeField"
	DecimalField              DjangoFieldType = "DecimalField"
	DurationField             DjangoFieldType = "DurationField"
	EmailField                DjangoFieldType = "EmailField"
	FileField                 DjangoFieldType = "FileField"
	FilePathField             DjangoFieldType = "FilePathField"
	ForeignKey                DjangoFieldType = "ForeignKey"
	GenericIPAddressField     DjangoFieldType = "GenericIPAddressField"
	ImageField                DjangoFieldType = "ImageField"
	IntegerField              DjangoFieldType = "IntegerField"
	FloatField                DjangoFieldType = "FloatField"
	JSONField                 DjangoFieldType = "JSONField"
	PositiveBigIntegerField   DjangoFieldType = "PositiveBigIntegerField"
	PositiveIntegerField      DjangoFieldType = "PositiveIntegerField"
	PositiveSmallIntegerField DjangoFieldType = "PositiveSmallIntegerField"
	SlugField                 DjangoFieldType = "SlugField"
	SmallIntegerField         DjangoFieldType = "SmallIntegerField"
	SmallAutoField            DjangoFieldType = "SmallAutoField"
	TextField                 DjangoFieldType = "TextField"
	TimeField                 DjangoFieldType = "TimeField"
	URLField                  DjangoFieldType = "URLField"
	UUIDField                 DjangoFieldType = "UUIDField"
)

type Field struct {
	Name           string
	Type           string
	IsNullable     bool
	IsPrimaryKey   bool
	IsRelationship bool
	Contraints     map[string]string
}

type Model struct {
	Name      string
	TableName string
	Fields    Fields
}

func (t Model) String() string {
	b, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return string(b)
	}
	return string(b)
}

type Fields []Field

func (t Fields) Len() int {
	return len(t)
}

func (t Fields) Less(i, j int) bool {
	if t[i].Name == "ID" {
		return true

	} else if t[j].Name == "ID" {
		return false

	}
	return t[i].Name < t[j].Name
}

func (t Fields) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t Field) DjangoType() (DjangoFieldType, error) {

	if t.IsPrimaryKey {
		switch t.Type {
		case "uint", "uint32", "uint8", "uint16", "int", "int32":
			return AutoField, nil
		case "int64", "uint64":
			return BigAutoField, nil
		}
		return "", fmt.Errorf("Invalid type for AutoField (%s)", t.Type)
	}

	if t.IsRelationship {
		return ForeignKey, nil
	}

	switch t.Type {
	case "map":
		return JSONField, nil
	case "uint", "uint32":
		return PositiveIntegerField, nil
	case "uint8", "uint16":
		return PositiveSmallIntegerField, nil
	case "uint64":
		return PositiveBigIntegerField, nil
	case "int", "int32":
		return IntegerField, nil
	case "int64":
		return BigIntegerField, nil
	case "float32", "float64":
		return FloatField, nil
	case "string":
		return TextField, nil
	case "bool":
		return BooleanField, nil
	case "time.Time":
		return DateTimeField, nil
	case "time.Duration":
		return DurationField, nil
	}

	return "", fmt.Errorf("Unhandled gotype (%s)", t.Type)
}

/*
AutoField
BigAutoField
BigIntegerField
BooleanField
CharField
DateField
DateTimeField
DecimalField
DurationField
EmailField
FileField
FileField and FieldFile
FilePathField
FloatField
GenericIPAddressField
ImageField
IntegerField
JSONField
PositiveBigIntegerField
PositiveIntegerField
PositiveSmallIntegerField
SlugField
SmallAutoField
SmallIntegerField
TextField
TimeField
URLField
UUIDField
*/

/*
type Event struct {
	ID             int64        `json:"id" gorm:"primaryKey"`
	EventAPI       EventAPIName `json:"event_api" gorm:"index"`
	Name           string       `json:"name"`
	Link           string       `json:"link"`
	StartDate      time.Time    `json:"start_date" gorm:"index"`
	EndDate        *time.Time   `json:"end_date,omitempty"`
	OnSaleDate     *time.Time   `json:"on_sale_date,omitempty" gorm:"index"`
	DateConfirmed  bool         `json:"date_confirmed"`
	TimeConfirmed  bool         `json:"time_confirmed"`
	Type           string       `json:"type" gorm:"index"`
	MinTicketPrice *Money       `json:"min_ticket_price,omitempty" gorm:"embedded;embeddedPrefix:min_ticket_price_"`
	Status         EventStatus  `json:"status"`
	GenreID        int64        `json:"genre_id,omitempty"`
	Genre          *Genre       `json:"genre,omitempty"`
	VenueID        int64        `json:"venue_id,omitempty"`
	Venue          *Venue       `json:"venue,omitempty"`
	Categories     []*Category  `json:"categories" gorm:"foreignKey:EventID" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt      *time.Time   `json:"created_at,omitempty" gorm:"autoCreateTime"`
	UpdatedAt      *time.Time   `json:"updated_at,omitempty" gorm:"autoUpdateTime"`
	DeletedAt      *time.Time   `json:"deleted_at,omitempty"`
}


# Create your models here.
class Country(models.Model):
    id = models.AutoField(primary_key=True)
    name = models.CharField(max_length=255)
    iso3 = models.CharField(max_length=10)
    iso2 = models.CharField(max_length=10, unique=True)
    numeric_code = models.CharField(max_length=10)
    phone_code = models.CharField(max_length=10)
    capital = models.CharField(max_length=255)
    currency = models.CharField(max_length=255)
    currency_name = models.CharField(max_length=255)
    currency_symbol = models.CharField(max_length=10)
    tld = models.CharField(max_length=10)
    native = models.CharField(max_length=255)
    region = models.CharField(max_length=255)
    subregion = models.CharField(max_length=255)
    timezones = models.TextField()
    latitude = models.FloatField()
    longitude = models.FloatField()
    emoji = models.CharField(max_length=10)
    emoji_u = models.CharField(max_length=20)
    active = models.BooleanField()
    created_at = models.DateTimeField(blank=True, null=True)
    updated_at = models.DateTimeField(blank=True, null=True)
    deleted_at = models.DateTimeField(blank=True, null=True)

    def __str__(self):
        return self.name


*/

/*
type Field interface {
	Name() string
	PythonType() string
	GolangType() string
	FieldArgs() []string
}

// BaseField represents common fields across Django field types
type BaseField struct {
	name string
}

// Name returns the field name
func (f *BaseField) Name() string {
	return f.name
}

// StringField represents Django CharField
type StringField struct {
	BaseField
	MaxLen int
}

// PythonType returns the Python type of the field
func (f *StringField) PythonType() string {
	return "str"
}

// GolangType returns the Golang type of the field
func (f *StringField) GolangType() string {
	return "string"
}

// FieldArgs returns the field arguments as strings
func (f *StringField) FieldArgs() []string {
	return []string{fmt.Sprintf("maxLen: %d", f.MaxLen)}
}

// IntegerField represents Django IntegerField
type IntegerField struct {
	BaseField
}

// PythonType returns the Python type of the field
func (f *IntegerField) PythonType() string {
	return "int"
}

// GolangType returns the Golang type of the field
func (f *IntegerField) GolangType() string {
	return "int"
}

// FieldArgs returns the field arguments as strings
func (f *IntegerField) FieldArgs() []string {
	return []string{}
}


type BaseField struct {
	name string
}

// Name returns the field name
func (f *BaseField) Name() string {
	return f.name
}

// CharField represents Django CharField
type CharField struct {
	BaseField
	MaxLen int
}

// PythonType returns the Python type of the field
func (f *CharField) PythonType() string {
	return "str"
}

// GolangType returns the Golang type of the field
func (f *CharField) GolangType() string {
	return "string"
}

// FieldArgs returns the field arguments as strings
func (f *CharField) FieldArgs() []string
	return []string{fmt.Sprintf("maxLen: %d", f.MaxLen)}
}

// IntegerField represents Django IntegerField
type IntegerField struct {
	BaseField
}

// PythonType returns the Python type of the field
func (f *IntegerField) PythonType() string {
	return "int"
}

// GolangType returns the Golang type of the field
func (f *IntegerField) GolangType() string {
	return "int"
}

// FieldArgs returns the field arguments as strings
func (f *IntegerField) FieldArgs() []string {
	return []string{}
}

// BooleanField represents Django BooleanField
type BooleanField struct {
	BaseField
}

// PythonType returns the Python type of the field
func (f *BooleanField) PythonType() string {
	return "bool"
}

// GolangType returns the Golang type of the field
func (f *BooleanField) GolangType() string {
	return "bool"
}

// FieldArgs returns the field arguments as strings
func (f *BooleanField) FieldArgs() []string {
	return []string{}
}
*/

/*
// GORM model tag constants
const (
	gormTag        = "gorm"
	columnTag      = "column"
	typeTag        = "type"
	autoIncrement  = "AUTO_INCREMENT"
	primaryKeyTag  = "primary_key"
	uniqueTag      = "unique"
	notNullTag     = "not null"
	defaultTag     = "default"
	defaultExprTag = "default_expr"
)

// GolangToDjangoType converts Golang type to Django model field type
func GolangToDjangoType(golangType string) string {
	switch golangType {
	case "int", "int8", "int16", "int32", "uint", "uint8", "uint16", "uint32":
		return "IntegerField"
	case "int64", "uint64":
		return "BigIntegerField"
	case "float32", "float64":
		return "FloatField"
	case "string":
		return "CharField"
	case "bool":
		return "BooleanField"
	default:
		return "UnknownFieldType"
	}
}

// ConvertGORMToDjangoModel converts GORM model to Django model definition
func ConvertGORMToDjangoModel(gormModel interface{}) {
	modelType := reflect.TypeOf(gormModel)

	if modelType.Kind() != reflect.Struct {
		fmt.Println("Input is not a struct")
		return
	}

	modelName := modelType.Name()

	fmt.Printf("class %s(models.Model):\n", modelName)

	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		fieldName := field.Name
		gormTag := field.Tag.Get(gormTag)
		djangoType := GolangToDjangoType(field.Type.Name())

		fmt.Printf("    %s = models.%s(", fieldName, djangoType)

		// Parse GORM tags
		if gormTag != "" {
			tags := strings.Split(gormTag, ";")
			for _, tag := range tags {
				tagParts := strings.Split(tag, ":")
				tagName := tagParts[0]
				tagValue := ""
				if len(tagParts) > 1 {
					tagValue = tagParts[1]
				}

				switch tagName {
				case columnTag:
					fmt.Printf("db_column='%s', ", tagValue)
				case typeTag:
					// Handle type tag, if needed
				case primaryKeyTag:
					fmt.Print("primary_key=True, ")
				case uniqueTag:
					fmt.Print("unique=True, ")
				case notNullTag:
					fmt.Print("null=False, ")
				case defaultTag:
					fmt.Printf("default='%s', ", tagValue)
				case defaultExprTag:
					// Handle default_expr tag, if needed
				}
			}
		}

		fmt.Println(")")

		// Print any additional configuration for the field
		// Add your own logic based on your specific requirements

		fmt.Println()
	}
}

// Example GORM model
type GORMModel struct {
	ID        uint   `gorm:"column:id;primary_key;auto_increment" json:"id"`
	Name      string `gorm:"column:name;type:varchar(255);not null" json:"name"`
	Age       int    `gorm:"column:age;not null;default:0" json:"age"`
	IsStudent bool   `gorm:"column:is_student;not null;default:false" json:"is_student"`
}

func main() {
	// Example usage
	gormModel := &GORMModel{}
	ConvertGORMToDjangoModel(gormModel)
}
*/
