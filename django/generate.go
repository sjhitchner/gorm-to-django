package django

import (
	"fmt"
	"io"
	"sort"
	"sync"
	"text/template"

	"github.com/sjhitchner/gorm-to-django/gorm"

	"github.com/iancoleman/strcase"
)

type Generator struct {
	models []Model
	tmpl   *template.Template
}

func New(templateDir string) (*Generator, error) {

	tmpl := template.New("django")

	var err error
	if templateDir == "" {
		tmpl, err = tmpl.New("model").Parse(ModelTemplate)
		if err != nil {
			return nil, err
		}

		tmpl, err = tmpl.New("admin").Parse(AdminTemplate)
		if err != nil {
			return nil, err
		}

	} else {

	}

	return &Generator{
		models: make([]Model, 0, 10),
		tmpl:   tmpl,
	}, nil
}

func (t *Generator) Build(in <-chan gorm.Struct) error {

	structMap := make(map[string]*gorm.Struct)
	for gm := range in {
		fmt.Println(gm)
		structMap[gm.Name] = &gm
	}

	// Preprocessing
	for _, s := range structMap {
		if s.IsModel {
			err := preprocessModel(s, structMap)
			if err != nil {
				return err
			}
		}
	}

	// Final Wrap up
	for _, s := range structMap {
		if s.IsModel {
			model, err := makeModel(s, structMap)
			if err != nil {
				return err
			}

			sort.Sort(model.Fields)

			t.models = append(t.models, *model)
		}
	}

	return nil
}

func preprocessModel(s *gorm.Struct, structMap map[string]*gorm.Struct) error {

	for _, field := range s.Fields {
	}

	return
}

func makeModel(s *gorm.Struct, structMap map[string]*gorm.Struct) (*Model, error) {

	//	models := chan

	var wg sync.WaitGroup

	fieldCh := make(chan Field)
	errCh := make(chan error)
	for _, field := range s.Fields {
		for _, fn := range fieldFuncs {
			wg.Add(1)
			go func() {
				defer wg.Done()
				fn(fieldCh, errCh, field, s, structMap)
			}()
		}
	}

	wg.Wait()
	close(fieldCh)
	close(errCh)

	go func() {
		var fields []Field
		for field := range fieldCh {
			fields = append(fields, field)
		}
	}()

	return &Model{
		Name:      s.Name,
		TableName: strcase.ToSnake(s.TableName()),
		Fields:    fields,
	}, nil
}

func (t *Generator) Models(w io.Writer) error {
	if err := t.tmpl.ExecuteTemplate(w, "model", t.models); err != nil {
		return err
	}
	return nil
}

func (t *Generator) Admin(w io.Writer) error {
	if err := t.tmpl.ExecuteTemplate(w, "admin", t.models); err != nil {
		return err
	}
	return nil
}
