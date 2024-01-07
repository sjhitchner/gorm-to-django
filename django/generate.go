package django

import (
	"fmt"
	"io"
	"sort"
	"sync"
	"text/template"

	"github.com/sjhitchner/gorm-to-django/gorm"

	"github.com/stoewer/go-strcase"
)

type Generator struct {
	Models Models
	tmpl   *template.Template
}

func NewWithCustomTemplates(templateDir string) (*Generator, error) {
	return nil, nil
}

func New() (*Generator, error) {
	tmpl := template.New("django")

	var err error
	tmpl, err = tmpl.New("model").Parse(ModelTemplate)
	if err != nil {
		return nil, err
	}

	tmpl, err = tmpl.New("admin").Parse(AdminTemplate)
	if err != nil {
		return nil, err
	}

	return &Generator{
		Models: make([]Model, 0, 10),
		tmpl:   tmpl,
	}, nil
}

func (t *Generator) Build(in <-chan gorm.Struct) error {
	out, errCh := gorm.Preprocess(in)

	errCh2 := make(chan error)
	go func() {
		defer close(errCh2)

		structMap := make(map[string]gorm.Struct)
		for gm := range out {
			fmt.Println("Q", gm)
			structMap[gm.Name] = gm
		}

		// Final Wrap up
		for _, s := range structMap {
			if s.IsModel {
				model, err := makeModel(s, structMap)
				if err != nil {
					errCh2 <- err
				}

				sort.Sort(model.Fields)

				t.Models = append(t.Models, *model)
			}
		}

		sort.Sort(t.Models)
	}()

	for err := range mergeErrors(errCh, errCh2) {
		return err
	}

	return nil
}

func makeModel(s gorm.Struct, structMap map[string]gorm.Struct) (*Model, error) {
	var fields []Field
	var wg sync.WaitGroup
	fieldCh := make(chan Field)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for field := range fieldCh {
			fields = append(fields, field)
		}
	}()

	for _, field := range s.Fields {
		for _, fn := range fieldFuncs {
			stop, err := fn(fieldCh, field, s, structMap)
			if err != nil {
				return nil, err
			}

			if stop {
				break
			}
		}
	}

	close(fieldCh)
	wg.Wait()

	return &Model{
		Name: s.Name,
		Metadata: map[string]string{
			"db_table": fmt.Sprintf("'%s'", strcase.SnakeCase(s.TableName())),
		},
		Fields: fields,
	}, nil
}

func (t *Generator) GenerateModels(w io.Writer) error {
	if err := t.tmpl.ExecuteTemplate(w, "model", t.Models); err != nil {
		return err
	}
	return nil
}

func (t *Generator) GenerateAdmin(w io.Writer) error {
	if err := t.tmpl.ExecuteTemplate(w, "admin", t.Models); err != nil {
		return err
	}
	return nil
}

func mergeErrors(channels ...<-chan error) <-chan error {
	var wg sync.WaitGroup
	merged := make(chan error)

	// Start a goroutine for each input channel
	for _, ch := range channels {
		wg.Add(1)
		go func(c <-chan error) {
			defer wg.Done()
			for err := range c {
				merged <- err
			}
		}(ch)
	}

	// Close the merged channel when all goroutines are done
	go func() {
		wg.Wait()
		close(merged)
	}()

	return merged
}
