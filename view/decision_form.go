package view

type DecisionFormRequest struct {
	Name string
}

type DecisionFormResponse struct {
	DecisionName string      `json:"decision_name"`
	FormFields   []FormField `json:"form_fields"`
}

type FormField struct {
	Key         string `json:"key"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Unit        string `json:"unit,omitempty"`
	Type        string `json:"type"`
}