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
	ManyToMany                DjangoFieldType = "ManyToManyField"
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
	Name            string
	Type            string
	IsNullable      bool
	IsPrimaryKey    bool
	IsRelationship  bool
	ManyToManyTable string
	Constraints     map[string]string
	Tags            map[string]string
	Django          map[string]string
}

func (t Field) HasManyToMany() bool {
	return t.ManyToManyTable != ""
}

type Model struct {
	Name     string
	Metadata map[string]string
	Fields   Fields
}

func (t Model) DisplayList() string {
	var list []string
	for _, field := range t.Fields {
		if _, yes := field.Django["display_list"]; yes {
			list = append(list, field.Name)
		}
	}

	if len(list) > 0 {
		return "'" + strings.Join(list, "','") + "'"
	}

	for _, field := range t.Fields {
		if field.Name == "name" {
			return "'" + field.Name + "'"
		}
	}
	return "'id'"
}

func (t Model) ReadOnlyFields() string {
	var list []string
	for _, field := range t.Fields {
		if _, yes := field.Django["readonly_field"]; yes {
			list = append(list, field.Name)
		}
	}

	if len(list) > 0 {
		return "'" + strings.Join(list, "','") + "'"
	}

	for _, field := range t.Fields {
		list = append(list, field.Name)
	}

	if len(list) > 0 {
		return "'" + strings.Join(list, "','") + "'"
	}

	return ""
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

	if t.HasManyToMany() {
		return ManyToMany, nil
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
		if _, ok := t.Tags["size"]; ok {
			return CharField, nil
		}
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

	if t.HasManyToMany() {
		args = append(args, fmt.Sprintf("'%s'", t.Type))
		args = append(args, fmt.Sprintf("through='%s'", t.ManyToManyTable))
	}

	if t.IsPrimaryKey {
		args = append(args, "primary_key=True")
	}

	if t.IsNullable {
		args = append(args, "null=True")
	}

	if size, ok := t.Tags["size"]; ok {
		args = append(args, fmt.Sprintf("max_length=%s", size))
	}

	if !t.IsPrimaryKey && !t.IsRelationship {
		args = append(args, "blank=True")
	}

	return strings.Join(args, ", ")
}

func (t Field) OnDelete() (string, bool) {
	onDelete, found := t.Constraints["on_delete"]
	if !found {
		return "DO_NOTHING", false
	}
	return onDelete, true
}
