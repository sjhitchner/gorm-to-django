package gorm

import (
	"fmt"
	"sync"
)

type PreprocessFunc func(out chan<- Update, field Field, st Struct, structMap map[string]Struct) error

var preprocessFuncs = []PreprocessFunc{
	ignoreRelationshipID,
	makeForeignKey,
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

		fmt.Println(updates)
		for _, u := range updates {

			s := structMap[u.Struct]

			switch u.Type {
			case Add:
				addField(&s, u.Field)
			case Delete:
				removeField(&s, u.Field)
			}

			structMap[u.Struct] = s
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
			Type:   Delete,
			Field: Field{
				Name: field.Name,
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
			Type:   Add,
			Field: Field{
				Name: relationName,
				Type: fmt.Sprintf("*%s", st.Name),
				Tags: field.Tags,
			},
		}

		out <- Update{
			Struct: st.Name,
			Type:   Delete,
			Field: Field{
				Name: field.Name,
			},
		}

		return nil
	}
	return nil
}

func addField(s *Struct, field Field) {
	removeField(s, field)
	s.Fields = append(s.Fields, field)
}

func removeField(s *Struct, field Field) {
	fields := make([]Field, 0, len(s.Fields))
	for _, f := range s.Fields {
		if field.Name != f.Name {
			fields = append(fields, f)
		}
	}
	s.Fields = fields
}
