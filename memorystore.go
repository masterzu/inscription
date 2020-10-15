package main

// stub type for devel store
type inMemoryStorage struct{}

func (s *inMemoryStorage) TemplateFromURL(url string) string {
	return "all is LOVE"
}
func (s *inMemoryStorage) GetModel() FormModel {
	return FormModel{
		Nom:    "Cao",
		Prenom: "Patrick",
	}
}

func (s *inMemoryStorage) RecordModel(model FormModel) error {
	return nil
}

func (s *inMemoryStorage) GetHashs() map[string]FormModel {
	return map[string]FormModel{}
}

func NewMemoryStore() *inMemoryStorage {
	return &inMemoryStorage{}
}
