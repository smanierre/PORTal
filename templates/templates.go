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

func New(templatesDir fs.FS, data RootData) *TemplateRepo {
	t := &TemplateRepo{
		templates: map[string]*template.Template{},
		rootData:  data,
	}
	err := fs.WalkDir(templatesDir, ".", func(path string, d fs.DirEntry, err error) error {
		if parts := strings.Split(d.Name(), "."); parts[len(parts)-1] != "gohtml" {
			return nil
		}
		if d.Name() == "root.gohtml" {
			t.templates["root"], err = template.ParseFS(templatesDir, "root.gohtml")
			if err != nil {
				return err
			}
		}
		if d.IsDir() {
			return nil
		}
		t.templates[strings.Split(d.Name(), ".")[0]], err = template.ParseFS(templatesDir, "root.gohtml", "nav.gohtml", path)
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
