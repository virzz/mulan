package maintain

import "github.com/virzz/mulan/db"

var models = []db.Modeler{}

func Register(model ...db.Modeler) { models = append(models, model...) }

func Models() []any {
	var result = make([]any, len(models))
	for i, model := range models {
		result[i] = model
	}
	return result
}
