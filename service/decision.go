package service

import (
	"bytes"
	"dut/model"
	templatePkg "dut/pkg/template"
	"dut/pkg/dictionary"
	"dut/repository"
	"dut/view"
	"encoding/json"
	"fmt"
	"text/template"
)

type DecisionService struct {
	repository repository.Repository
}

func NewDecisionService(repository repository.Repository) DecisionService {
	return DecisionService{repository}
}

func (ds DecisionService) GetDecisionInputForm(request view.DecisionFormRequest) (view.DecisionFormResponse, error) {
	decision, err := ds.repository.GetDecision(request.Name)
	if err != nil {
		return view.DecisionFormResponse{}, err
	}

	formFields := []view.FormField{}
	for _, field := range decision.InputForm {
		chainLink := field.ChainLink
		if chainLink != nil {
			chainFormFields, err := ds.getChainedFormFields(*chainLink)
			if err != nil {
				return view.DecisionFormResponse{}, err
			}
			formFields = append(formFields, chainFormFields...)
			continue
		}

		formFields = append(formFields, getFormFieldView(field))
	}

	return view.DecisionFormResponse{
		DecisionName: decision.Name,
		FormFields:   formFields,
	}, nil
}

func (ds DecisionService) EvaluateDecision(name string, form map[string]interface{}) (map[string]interface{}, error) {
	decision, err := ds.repository.GetDecision(name)
	if err != nil {
		return nil, err
	}

	templateFuncMap := map[string]interface{}{
		"sum": templatePkg.Sum,
		"greater_than": templatePkg.Gt,
		"greater_than_or_equals": templatePkg.Gte,
		"less_than": templatePkg.Lt,
		"less_than_or_equals": templatePkg.Lte,
	}

	input := make(map[string]interface{})
	for _, field := range decision.InputForm {
		chainLink := field.ChainLink
		if chainLink != nil {
			chainFormFields, err := ds.evaluateChainedFormFields(*chainLink, form, templateFuncMap)
			if err != nil {
				return nil, err
			}
			input = dictionary.MergeStringToInterfaceMaps(input, chainFormFields)
			continue
		}

		value := form[field.Key]
		err = validateValue(value, field.Type)
		if err != nil {
			return nil, err
		}
		input[field.Key] = value
	}

	output := make(map[string]interface{})
	for _, field := range decision.OutputForm {
		rule := field.Rule
		if rule == nil {
			continue
		}
		tmpl, err := template.New(decision.Name).Funcs(template.FuncMap(templateFuncMap)).Parse(*rule)
		if err != nil {
			return nil, err
		}

		var b bytes.Buffer
		err = tmpl.Execute(&b, input)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		m := make(map[string]interface{})
		err = json.Unmarshal(b.Bytes(), &m)
		if err != nil {
			return nil, err
		}

		outputValue, ok := m[field.Key]
		if !ok {
			return nil, fmt.Errorf("key %s not found in the output of decision %s", field.Key, decision.Name)
		}

		output[field.Key] = outputValue
	}

	return output, nil
}

func (ds DecisionService) evaluateChainedFormFields(chainLink model.DecisionChainLink, form, templateFuncMap map[string]interface{}) (map[string]interface{}, error) {
	decision, err := ds.repository.GetDecision(chainLink.DecisionName)
	if err != nil {
		return nil, err
	}

	input := make(map[string]interface{})
	for _, field := range decision.InputForm {
		nextChainLink := field.ChainLink
		if nextChainLink != nil {
			chainFormFields, err := ds.evaluateChainedFormFields(*nextChainLink, form, templateFuncMap)
			if err != nil {
				return nil, err
			}
			input = dictionary.MergeStringToInterfaceMaps(input, chainFormFields)
			continue
		}

		value := form[field.Key]
		err = validateValue(value, field.Type)
		if err != nil {
			return nil, err
		}
		input[field.Key] = value
	}

	output := make(map[string]interface{})
	for _, field := range decision.OutputForm {
		rule := field.Rule
		if rule == nil {
			continue
		}
		tmpl, err := template.New(decision.Name).Funcs(template.FuncMap(templateFuncMap)).Parse(*rule)
		if err != nil {
			return nil, err
		}

		var b bytes.Buffer
		err = tmpl.Execute(&b, input)
		if err != nil {
			return nil, err
		}

		m := make(map[string]interface{})
		err = json.Unmarshal(b.Bytes(), &m)
		if err != nil {
			return nil, err
		}

		outputValue, ok := m[chainLink.SourceKey]
		if !ok {
			return nil, fmt.Errorf("key %s not found in the output of decision %s", chainLink.SourceKey, decision.Name)
		}

		output[chainLink.DestinationKey] = outputValue
	}

	return output, nil
}

func validateValue(value interface{}, valueType string) error {
	switch valueType {
	case "Number":
		_, ok := value.(float64)
		if !ok {
			return fmt.Errorf("%v is not a valid %s", value, valueType)
		}
		return nil
	case "Bool":
		_, ok := value.(bool)
		if !ok {
			return fmt.Errorf("%v is not a valid %s", value, valueType)
		}
		return nil
	case "String":
		_, ok := value.(string)
		if !ok {
			return fmt.Errorf("%v is not a valid %s", value, valueType)
		}
		return nil
	default:
		return fmt.Errorf("invalid value type %s", valueType)
	}
}

func (ds DecisionService) getChainedFormFields(chainLink model.DecisionChainLink) ([]view.FormField, error) {
	decision, err := ds.repository.GetDecision(chainLink.DecisionName)
	if err != nil {
		return []view.FormField{}, err
	}

	formFields := []view.FormField{}
	for _, field := range decision.InputForm {
		nextChainLink := field.ChainLink
		if nextChainLink != nil {
			chainFormFields, err := ds.getChainedFormFields(*nextChainLink)
			if err != nil {
				return []view.FormField{}, err
			}
			formFields = append(formFields, chainFormFields...)
			continue
		}
		formFields = append(formFields, getFormFieldView(field))
	}

	return formFields, nil
}

func getFormFieldView(field model.FormField) view.FormField {
	return view.FormField{
		Key:         field.Key,
		Title:       field.Title,
		Description: field.Description,
		Unit:        field.Unit,
		Type:        field.Type,
	}
}