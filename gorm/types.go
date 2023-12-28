package gorm

import (
	"encoding/json"
	"strings"

	"github.com/stoewer/go-strcase"
)

type UpdateType int

const (
	AddField UpdateType = iota
	DeleteField
	AddStruct
)

type Update struct {
	Struct   string
	Metadata map[string]string
	Type     UpdateType
	Fields   []Field
}

func (t Update) String() string {
	b, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return string(b)
	}
	return string(b)
}

type Struct struct {
	IsModel  bool
	Name     string
	Metadata map[string]string
	Fields   []Field
}

func (t Struct) SnakeName() string {
	// return strcase.ToSnake(t.Name)
	return strcase.SnakeCase(t.Name)
}

func (t Struct) TableName() string {
	return t.Metadata["tablename"]
}

func (t Struct) HasMany2Many() ([]Field, bool) {
	fields := make([]Field, 0, 5)
	for _, f := range t.Fields {
		if _, yes := f.IsMany2Many(); yes {
			fields = append(fields, f)
		}
	}
	return fields, len(fields) > 0
}

func (t Struct) String() string {
	b, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return string(b)
	}
	return string(b)
}

func (t Struct) HasRelation(key string) bool {
	for _, field := range t.Fields {
		if strings.HasPrefix(field.Name, key[:len(key)-2]) {
			return true
		}
	}
	return false
}

type Field struct {
	Name string
	Type string
	Tags map[string]Tag
}

func (t Field) IsID() bool {
	return t.Name == "ID"
}

func (t Field) IsPrimaryKey() bool {
	for _, tag := range t.Tags {
		if tag.IsPrimaryKey() {
			return t.IsID()
		}
	}
	return false
}

func (t Field) IsForeignKey() (string, bool) {
	for _, tag := range t.Tags {
		if tag.IsForeignKey() {
			return tag.Value[:len(tag.Value)-2], true
		}
	}
	return "", false
}

func (t Field) IsMany2Many() (string, bool) {
	for _, tag := range t.Tags {
		if tag.IsMany2Many() {
			return tag.Value, true
		}
	}
	return "", false
}

func (t Field) IsRelationshipID() bool {
	return strings.HasSuffix(t.Name, "ID") && !t.IsID()
}

func (t Field) HasConstraints() (map[string]string, bool) {
	for _, tag := range t.Tags {
		if m, yes := tag.HasConstraints(); yes {
			return m, true
		}
	}
	return nil, false
}

func (t Field) SnakeName() string {
	// return strcase.ToSnake(t.Name)
	return strcase.SnakeCase(t.Name)
}

func (t Field) GetType() (string, bool) {
	if t.IsArray() {
		if strings.HasPrefix(t.Type[2:], "*") {
			return t.Type[3:], true
		}
		return t.Type[2:], true
	}

	if strings.HasPrefix(t.Type, "*") {
		return t.Type[1:], true
	}

	return t.Type, false
}

func (t Field) GetTags(source string) map[string]string {
	m := make(map[string]string)
	for _, tag := range t.Tags {
		if source == tag.Source {
			m[tag.Name] = tag.Value
		}
	}
	return m
}

func (t Field) IsInt64() bool {
	return t.Type == "int64"
}

func (t Field) IsArray() bool {
	return strings.HasPrefix(t.Type, "[]")
}

func (t Field) IsMap() bool {
	return strings.HasPrefix(t.Type, "map[")
}

func (t Field) IsStruct() bool {
	goType, _ := t.GetType()

	switch GoType(goType) {
	case IntType, Int8Type, Int16Type, Int32Type, Int64Type:
		return false
	case UIntType, UInt8Type, UInt16Type, UInt32Type, UInt64Type:
		return false
	case Float32Type, Float64Type:
		return false
	case StringType:
		return false
	case BoolType:
		return false
	case ByteType, RuneType:
		return false
	case Complex64Type, Complex128Type:
		return false
	case TimeType:
		return false
	}
	return !t.IsArray() && !t.IsMap()
}

func (t Field) IsEmbedded() bool {
	for _, tag := range t.Tags {
		if tag.IsEmbedded() {
			return true
		}
	}
	return false
}

func (t Field) String() string {
	b, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return string(b)
	}
	return string(b)
}

type Tag struct {
	Name   string
	Value  string
	Source string
}

func (t Tag) IsPrimaryKey() bool {
	return t.Name == "primaryKey"
}

func (t Tag) IsAutoIncrement() bool {
	return t.Name == "autoIncrement"
}

func (t Tag) IsNotNull() bool {
	return t.Name == "not null"
}

func (t Tag) IsEmbedded() bool {
	return t.Name == "embedded"
}

func (t Tag) IsForeignKey() bool {
	return t.Name == "foreignKey"
}

func (t Tag) IsMany2Many() bool {
	return t.Name == "many2many"
}

func (t Tag) IsAutoCreateTime() bool {
	return t.Name == "autoCreateTime"
}

func (t Tag) IsAutoUpdateTime() bool {
	return t.Name == "autoUpdateTime"
}

func (t Tag) HasConstraints() (map[string]string, bool) {
	if t.Name == "constraint" {
		return parseConstraints(t.Value), true
	}
	return nil, false
}

func parseConstraints(optionsString string) map[string]string {
	options := make(map[string]string)

	// Split the string into records separated by commas
	records := strings.Split(optionsString, ",")

	// Iterate over each record
	for _, record := range records {
		// Split the record into key and value separated by colon
		parts := strings.SplitN(record, ":", 2)
		if len(parts) == 2 {
			// key := strcase.ToSnake(strings.TrimSpace(parts[0]))
			key := strcase.SnakeCase(strings.TrimSpace(parts[0]))
			value := strings.TrimSpace(parts[1])
			options[key] = value
		}
	}
	return options
}

type GoType string

const (
	IntType        GoType = "int"
	Int8Type       GoType = "int8"
	Int16Type      GoType = "int16"
	Int32Type      GoType = "int32"
	Int64Type      GoType = "int64"
	UIntType       GoType = "uint"
	UInt8Type      GoType = "uint8"
	UInt16Type     GoType = "uint16"
	UInt32Type     GoType = "uint32"
	UInt64Type     GoType = "uint64"
	Float32Type    GoType = "float32"
	Float64Type    GoType = "float64"
	StringType     GoType = "string"
	BoolType       GoType = "bool"
	ByteType       GoType = "byte"
	RuneType       GoType = "rune"
	Complex64Type  GoType = "complex64"
	Complex128Type GoType = "complex128"
	TimeType       GoType = "time.Time"
)
