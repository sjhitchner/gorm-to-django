package django

import (
	"fmt"

	"github.com/sjhitchner/gorm-to-django/gorm"
)

type FieldFunc func(out chan<- Field, field gorm.Field, st gorm.Struct, structMap map[string]gorm.Struct) (bool, error)

var fieldFuncs = []FieldFunc{
	makeEmbedded,
	makeRelationship,
	makeID,
	makeMap,
	makeField,
}

func makeEmbedded(out chan<- Field, field gorm.Field, st gorm.Struct, structMap map[string]gorm.Struct) (bool, error) {
	if field.IsEmbedded() {
		embeddedType, _ := field.GetType()
		embeddedPrefix, found := field.Tags["embeddedPrefix"]
		if !found {
			return false, fmt.Errorf("Embedded model (%s) has no prefix", embeddedType)
		}

		embeddedModel, found := structMap[embeddedType]
		if !found {
			return false, fmt.Errorf("Embedded model (%s) not found", embeddedType)
		}

		for _, modelField := range embeddedModel.Fields {
			modelType, nullable := modelField.GetType()
			out <- Field{
				Name:       fmt.Sprintf("%s%s", embeddedPrefix.Value, modelField.SnakeName()),
				Type:       modelType,
				IsNullable: nullable,
			}
		}
		return true, nil
	}

	return false, nil
}

func makeRelationship(out chan<- Field, field gorm.Field, st gorm.Struct, structMap map[string]gorm.Struct) (bool, error) {
	if field.IsStruct() && !field.IsEmbedded() {

		constraints, _ := field.HasConstraints()

		modelType, _ := field.GetType()
		out <- Field{
			Name:           field.SnakeName(),
			Type:           modelType,
			IsNullable:     false,
			IsRelationship: true,
			Constraints:    constraints,
		}
		return true, nil
	}
	return false, nil
}

func makeID(out chan<- Field, field gorm.Field, st gorm.Struct, structMap map[string]gorm.Struct) (bool, error) {
	if field.IsPrimaryKey() {
		out <- Field{
			Name:         field.SnakeName(),
			Type:         field.Type,
			IsNullable:   false,
			IsPrimaryKey: true,
		}
		return true, nil
	}

	return false, nil
}

func makeMap(out chan<- Field, field gorm.Field, st gorm.Struct, structMap map[string]gorm.Struct) (bool, error) {

	if field.IsMap() {
		modelType, nullable := field.GetType()
		out <- Field{
			Name:       field.SnakeName(),
			Type:       modelType,
			IsNullable: nullable,
		}
		return true, nil
	}
	return false, nil
}

func makeField(out chan<- Field, field gorm.Field, st gorm.Struct, structMap map[string]gorm.Struct) (bool, error) {
	modelType, nullable := field.GetType()
	out <- Field{
		Name:       field.SnakeName(),
		Type:       modelType,
		IsNullable: nullable,
		Tags:       field.GetTags("django"),
	}
	return true, nil
}
