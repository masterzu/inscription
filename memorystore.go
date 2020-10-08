package main

// stub type for devel store
// FIXME: get a real datas storage
type stubWriter struct{}

func (s *stubWriter) GetForm(url string) string {
	return "all is LOVE"
}
func (s *stubWriter) GetModel() FormModel {
	return FormModel{
		Nom:    "Cao",
		Prenom: "Patrick",
	}
}

func InMemoryStore() *stubWriter {
	return &stubWriter{}
}
