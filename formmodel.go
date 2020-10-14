package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
)

///////////////////////////////////////////////////////////
// Data Model

// FormModel is the struct for the first form filled by student
type FormModel struct {
	Nom    string
	Prenom string
}

func ReadFormModel(r io.Reader) (FormModel, error) {
	var model FormModel
	err := json.NewDecoder(r).Decode(&model)
	if err != nil {
		err = fmt.Errorf("problem reading JSON, %v", err)
	}

	return model, err
}

func (f *FormModel) String() string {
	return fmt.Sprintf("Nom: %s, Prenom: %s", f.Nom, f.Prenom)
}

// GetHash return the calculated hash from model
func (f *FormModel) GetHash() string {
	h := md5.New()
	h.Write([]byte(f.String()))
	resu := b2s(h.Sum(nil))
	return resu
}

func (f *FormModel) GetJSON() string {
	body, _ := json.Marshal(f)
	return string(body)
}

///////////////////////////////////////////////////////////
// interfaces
// type ReaderFormModel interface {
// 	Read(io.Reader) (FormModel, error)
// }

type Storage interface {
	// TemplateFromURL return a HTML template from a url
	TemplateFromURL(string) string

	// GetModel return a model from internal state
	GetModel() FormModel

	//RecordModel store model with a hash
	RecordModel(FormModel, string) error
}

///////////////////////////////////////////////////////////
// private tools

// b2s convert a []byte to a string in hex
func b2s(b []byte) string {
	return fmt.Sprintf("%x", b)
}
