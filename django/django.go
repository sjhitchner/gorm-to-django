package django

import (
	"encoding/json"
	"fmt"
	"strings"
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
	Constraints    map[string]string
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

type Models []Model

func (t Models) Len() int {
	return len(t)
}

func (t Models) Less(i, j int) bool {
	return t[i].Name < t[j].Name
}

func (t Models) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

type Fields []Field

func (t Fields) Len() int {
	return len(t)
}

func (t Fields) Less(i, j int) bool {
	if t[i].Name == "id" {
		return true

	} else if t[j].Name == "id" {
		return false

	}

	if t[i].Name == "name" {
		return true

	} else if t[j].Name == "name" {
		return false

	}

	if strings.HasSuffix(t[i].Name, "_at") &&
		strings.HasSuffix(t[j].Name, "_at") {
		return t[i].Name < t[j].Name
	} else if strings.HasSuffix(t[i].Name, "_at") {
		return false
	} else if strings.HasSuffix(t[i].Name, "_at") {
		return true
	}

	return t[i].Name < t[j].Name
}

func (t Fields) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t Field) DjangoField() (DjangoFieldType, error) {

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

func (t Field) DjangoArgs() string {

	var args []string

	if t.IsRelationship {
		args = append(args, fmt.Sprintf("'%s'", t.Type))

		onDelete, _ := t.OnDelete()
		args = append(args, fmt.Sprintf("on_delete=models.%s", onDelete))
	}

	if t.IsPrimaryKey {
		args = append(args, "primary_key=True")
	}

	if t.IsNullable {
		args = append(args, "null=True")
	}

	args = append(args, "blank=True")

	return strings.Join(args, ", ")
}

func (t Field) OnDelete() (string, bool) {
	onDelete, found := t.Constraints["on_delete"]
	if !found {
		return "DO_NOTHING", false
	}
	return onDelete, true
}

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
