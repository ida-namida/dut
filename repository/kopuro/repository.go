package kopuro

import (
	"dut/model"
)

type Repository struct {
	kc kopuroClient
}

func NewRepository(kopuroBaseURL string) Repository {
	kc := newKopuroClient(kopuroBaseURL)
	return Repository{kc}
}

func (repo Repository) GetDecision(name string) (model.Decision, error) {
	var d model.Decision
	err := repo.kc.readJSONFile(name, &d)
	if err != nil {
		return model.Decision{}, err
	}

	return d, nil
}