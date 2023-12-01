package gorm

import (
	"encoding/json"
	"strings"

	"github.com/iancoleman/strcase"
)

type UpdateType int

const (
	Add UpdateType = iota
	Delete
)

type Update struct {
	Struct string
	Type   UpdateType
	Field  Field
}

type Struct struct {
	IsModel  bool
	Name     string
	Metadata map[string]string
	Fields   []Field
}

func (t Struct) SnakeName() string {
	return strcase.ToSnake(t.Name)
}

func (t Struct) TableName() string {
	return t.Metadata["tablename"]
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
	return strcase.ToSnake(t.Name)
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
	switch t.Type {
	case "int", "int32", "int64", "float32", "float64", "bool":
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
	Name  string
	Value string
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

func (t Tag) IsAutoCreateTime() bool {
	return t.Name == "autoCreateTime"
}

func (t Tag) IsAutoUpdateTime() bool {
	return t.Name == "autoUpdateTime"
}

func (t Tag) HasConstraints() (map[string]string, bool) {
	if t.Name != "constraints" {
		return nil, false
	}
	return parseConstraints(t.Value), true
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
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			options[key] = value
		}
	}

	return options
}
