package templates

import (
	"embed"
	"html/template"
	"io"
	"io/fs"
	"log"
	"strings"
)

//go:embed *
var TemplateDir embed.FS

type TplData struct {
	NavData
	ContentData interface{}
}

type tplData struct {
	TplData
	RootData
}

type TemplateRepo struct {
	templates map[string]*template.Template
	rootData  RootData
}

func New(templatesDir fs.FS, data RootData, serviceCode string) *TemplateRepo {
	t := &TemplateRepo{
		templates: map[string]*template.Template{},
		rootData:  data,
	}

	funcs := template.FuncMap{
		"displayName": getDisplayNameFunc(serviceCode),
	}

	err := fs.WalkDir(templatesDir, ".", func(path string, d fs.DirEntry, err error) error {
		if parts := strings.Split(d.Name(), "."); parts[len(parts)-1] != "gohtml" {
			return nil
		}
		if d.Name() == "root.gohtml" {
			t.templates["root"], err = template.New("root").Funcs(funcs).ParseFS(templatesDir, "root.gohtml")
			if err != nil {
				return err
			}
		}
		// Skip directories
		if d.IsDir() {
			return nil
		}
		// Skip components directory as they will be parsed with each template
		if strings.Contains(path, "components/") {
			return nil
		}
		name := strings.Split(d.Name(), ".")[0]
		t.templates[name], err = template.New(name).Funcs(funcs).ParseFS(templatesDir, "root.gohtml", "nav.gohtml", path, "components/*")
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return t
}

func (t *TemplateRepo) Render(w io.Writer, name string, data *TplData) error {
	var d tplData
	if data == nil {
		d = tplData{
			TplData: TplData{
				NavData:     NavData{},
				ContentData: nil,
			},
			RootData: t.rootData,
		}
	} else {
		d = tplData{
			TplData:  *data,
			RootData: t.rootData,
		}
	}
	return t.templates[name].ExecuteTemplate(w, "root", d)
}

func (t *TemplateRepo) RenderFragment(w io.Writer, name, fragment string, data interface{}) error {
	return t.templates[name].ExecuteTemplate(w, fragment, data)
}
