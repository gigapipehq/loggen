package schemas

import (
	"github.com/brianvoe/gofakeit/v6"
)

//easyjson:json
type CategoryList []string

//easyjson:json
type GeneratorFunctionList []GeneratorFunction
type GeneratorFunction struct {
	Display     string           `json:"display"`
	Category    string           `json:"category"`
	Description string           `json:"description"`
	Example     string           `json:"example"`
	Params      []gofakeit.Param `json:"params,omitempty"`
}
