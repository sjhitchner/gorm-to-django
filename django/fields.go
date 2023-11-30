package django

import (
	"fmt"

	"github.com/sjhitchner/gorm-to-django/gorm"
)

type PreprocessFunc func(field gorm.Field, st *gorm.Struct, structMap map[string]*gorm.Struct) (bool, error)

var preprocessFuncs = []PreprocessFunc{
	ignoreForeignKeyIDs,
}

type FieldFunc func(chan<- Field, errCh chan<- error, gorm.Field,  *gorm.Struct,  map[string]*gorm.Struct) 

var fieldFuncs = []FieldFunc{
	makeEmbedded,
	makeForeignKeys,
	makeID,
	makeArray,
	makeMap,
	makeField,
}

func makeEmbedded(field gorm.Field, st *gorm.Struct, structMap map[string]*gorm.Struct) ([]Field, bool, error) {

	if field.IsEmbedded() {

		embeddedType, _ := field.GetType()
		embeddedPrefix, found := field.Tags["embeddedPrefix"]
		if !found {
			return nil, false, fmt.Errorf("Embedded model (%s) has no prefix", embeddedType)
		}

		embeddedModel, found := structMap[embeddedType]
		if !found {
			return nil, false, fmt.Errorf("Embedded model (%s) not found", embeddedType)
		}

		var fields []Field
		for _, modelField := range embeddedModel.Fields {

			modelType, nullable := modelField.GetType()
			embeddedField := Field{
				Name:       fmt.Sprintf("%s%s", embeddedPrefix.Value, modelField.SnakeName()),
				Type:       modelType,
				IsNullable: nullable,
			}

			fields = append(fields, embeddedField)
		}
		return fields, true, nil
	}

	return nil, false, nil
}

// genre = models.ForeignKey('Genre', on_delete=models.DO_NOTHING, blank=False, null=False)
//  venue = models.ForeignKey('Venue', on_delete=models.DO_NOTHING, blank=False, null=False)
//  status = models.ForeignKey('Status', on_delete=models.DO_NOTHING, blank=False, null=False)
func makeForeignKeys(field gorm.Field, st *gorm.Struct, structMap map[string]*gorm.Struct) ([]Field, bool, error) {

	if field.IsStruct() && !field.IsEmbedded() {
		modelType, nullable := field.GetType()
		return []Field{
			Field{
				Name:           field.SnakeName(),
				Type:           modelType,
				IsNullable:     nullable,
				IsRelationship: true,
			},
		}, true, nil

	}
	return nil, false, nil
}



// id = models.AutoField(primary_key=True)
func makeID(field gorm.Field, st *gorm.Struct, structMap map[string]*gorm.Struct) ([]Field, bool, error) {
	if field.IsPrimaryKey() {
		return []Field{
			Field{
				Name:         field.SnakeName(),
				Type:         field.Type,
				IsNullable:   false,
				IsPrimaryKey: true,
			},
		}, true, nil
	}
	return nil, false, nil
}

func makeArray(field gorm.Field, st *gorm.Struct, structMap map[string]*gorm.Struct) ([]Field, bool, error) {
	if field.IsArray() {
		fmt.Println("array", field.Name)

		// TODO figure out a way to add reference to associated model
		/*
			modelType, nullable := field.GetType()
				[]Field{
						Field{
							Name:       field.SnakeName(),
							Type:       modelType,
							IsNullable: nullable,
						},
					}
		*/
		return nil, true, nil
	}
	return nil, false, nil
}

func makeMap(field gorm.Field, st *gorm.Struct, structMap map[string]*gorm.Struct) ([]Field, bool, error) {

	if field.IsMap() {
		modelType, nullable := field.GetType()
		return []Field{
			Field{
				Name:       field.SnakeName(),
				Type:       modelType,
				IsNullable: nullable,
			},
		}, true, nil
	}
	return nil, false, nil
}

func makeField(field gorm.Field, st *gorm.Struct, structMap map[string]*gorm.Struct) ([]Field, bool, error) {
	modelType, nullable := field.GetType()
	return []Field{
		Field{
			Name:       field.SnakeName(),
			Type:       modelType,
			IsNullable: nullable,
		},
	}, true, nil
}


func ignoreForeignKeyIDs(field gorm.Field, st *gorm.Struct, structMap map[string]*gorm.Struct) (bool, error) {

	if field.IsForeignKeyID() {
		if !st.HasRelation(field.Name) {
			return false, fmt.Errorf("Model has ForeignKeyID (%s) but no corresponding relation", field.Name)
		}
		return true, nil
	}

	return  false, nil
}
