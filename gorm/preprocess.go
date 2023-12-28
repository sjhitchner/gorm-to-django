package gorm

import (
	"fmt"
	"sync"

	"github.com/stoewer/go-strcase"
)

type PreprocessFunc func(out chan<- Update, field Field, st Struct, structMap map[string]Struct) error

var preprocessFuncs = []PreprocessFunc{
	ignoreRelationshipID,
	makeForeignKey,
	many2Many,
}

func Preprocess(in <-chan Struct) (<-chan Struct, <-chan error) {
	out := make(chan Struct)
	errCh := make(chan error)

	go func() {
		defer close(out)
		defer close(errCh)

		structMap := make(map[string]Struct)
		for s := range in {
			structMap[s.Name] = s
		}

		updateCh := make(chan Update)
		var updates []Update

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()

			for u := range updateCh {
				updates = append(updates, u)
			}
		}()

		go func() {
			defer wg.Done()
			defer close(updateCh)

			for _, s := range structMap {
				if s.IsModel {
					err := preprocessStruct(updateCh, s, structMap)
					if err != nil {
						errCh <- err
					}
				}
			}
		}()

		wg.Wait()

		for _, u := range updates {
			s := structMap[u.Struct]

			switch u.Type {
			case AddField:
				addFields(&s, u.Fields...)
				structMap[u.Struct] = s

			case DeleteField:
				removeFields(&s, u.Fields...)
				structMap[u.Struct] = s

			case AddStruct:
				structMap[u.Struct] = Struct{
					IsModel:  true,
					Name:     u.Struct,
					Metadata: u.Metadata,
					Fields:   u.Fields,
				}
			}
		}

		for _, s := range structMap {
			out <- s
		}
	}()

	return out, errCh
}

func preprocessStruct(out chan<- Update, s Struct, structMap map[string]Struct) error {
	for _, field := range s.Fields {
		for _, fn := range preprocessFuncs {
			if err := fn(out, field, s, structMap); err != nil {
				return err
			}
		}
	}
	return nil
}

func ignoreRelationshipID(out chan<- Update, field Field, st Struct, structMap map[string]Struct) error {
	if field.IsRelationshipID() {
		if !st.HasRelation(field.Name) {
			return fmt.Errorf("Model has RelationshipID (%s) but no corresponding relation", field.Name)
		}

		out <- Update{
			Struct: st.Name,
			Type:   DeleteField,
			Fields: []Field{
				Field{
					Name: field.Name,
				},
			},
		}
		return nil
	}
	return nil
}

func makeForeignKey(out chan<- Update, field Field, st Struct, structMap map[string]Struct) error {
	if relationName, yes := field.IsForeignKey(); yes {
		modelName, _ := field.GetType()

		_, found := structMap[modelName]
		if !found {
			return fmt.Errorf("Relationship defined but model doens't exist (%s)", modelName)
		}

		out <- Update{
			Struct: modelName,
			Type:   AddField,
			Fields: []Field{
				Field{
					Name: relationName,
					Type: fmt.Sprintf("*%s", st.Name),
					Tags: field.Tags,
				},
			},
		}

		out <- Update{
			Struct: st.Name,
			Type:   DeleteField,
			Fields: []Field{
				Field{
					Name: field.Name,
				},
			},
		}

		return nil
	}
	return nil
}

func many2Many(out chan<- Update, field Field, st Struct, structMap map[string]Struct) error {
	if tableName, yes := field.IsMany2Many(); yes {
		fieldType, _ := field.GetType()

		out <- Update{
			Struct: strcase.UpperCamelCase(tableName),
			Metadata: map[string]string{
				"tablename": tableName,
			},
			Type: AddStruct,
			Fields: []Field{
				Field{
					Name: "ID",
					Type: "int64",
					Tags: map[string]Tag{
						"primaryKey": Tag{
							Name:   "primaryKey",
							Source: "gorm",
						},
					},
				},
				Field{
					Name: st.Name,
					Type: "*" + st.Name,
					Tags: map[string]Tag{
						"constraint": Tag{
							Name:   "constraint",
							Value:  "OnDelete:CASCADE",
							Source: "gorm",
						},
					},
				},
				Field{
					Name: fieldType,
					Type: "*" + fieldType,
					Tags: map[string]Tag{
						"constraint": Tag{
							Name:   "constraint",
							Value:  "OnDelete:PROTECT",
							Source: "gorm",
						},
					},
				},
			},
		}

		out <- Update{
			Struct: st.Name,
			Type:   DeleteField,
			Fields: []Field{
				Field{
					Name: field.Name,
				},
			},
		}
	}

	return nil
}

func addFields(s *Struct, fields ...Field) {
	removeFields(s, fields...)
	s.Fields = append(s.Fields, fields...)
}

func removeFields(s *Struct, fields ...Field) {
	out := make([]Field, 0, len(s.Fields))
	for _, f := range s.Fields {
		for _, field := range fields {
			if field.Name != f.Name {
				out = append(out, f)
			}
		}
	}
	s.Fields = out
}
