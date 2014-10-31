package ca

import (
	"testing"

	"github.com/jgcarvalho/zeca/rules"
)

type testdata struct {
	id       string
	begin    []byte
	expected []byte
	rule     *rules.Rule
}

var tests = []testdata{
	{"ca_create_ok", []byte{0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0}, []byte{1, 0, 0, 1, 0, 0, 0, 1, 0, 0, 1}, nil},
	{"ca_create_difflen", []byte{0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0}, []byte{1, 0, 0, 1, 0, 0, 0, 1, 0, 0, 1, 2, 3, 4}, nil},
}

func TestCreate1D(t *testing.T) {
	// teste caso 1
	id := tests[0].id
	begin := tests[0].begin
	end := tests[0].expected
	rule, _ := rules.Create(begin, end, true, 3)
	ca, err := Create1D(id, begin, end, rule, 10)
	if err != nil {
		t.Error("Caso 1. Erro na criação do autômato celular")
	}
	if ca == nil {
		t.Error("Caso 1. Erro no retorno do autômato celular")
	}

	// teste caso 2
	id = tests[1].id
	begin = tests[1].begin
	end = tests[1].expected
	rule, _ = rules.Create(begin, end, true, 3)
	ca, err = Create1D(id, begin, end, rule, 19)
	if err == nil {
		t.Error("Caso 2. Erro na criação do autômato celular que não capturou o erro. ")
	}
	if ca != nil {
		t.Error("Caso 2. Erro no retorno do autômato celular que retornou um autômato inválido")
	}

}

func Testenconde(t *testing.T) {

}

// func TestConfusionMatrix(t *testing.T) {
// 	id := tests[1].id
// 	begin := tests[1].begin
// 	end := tests[1].expected
// 	rule, _ := rules.Create(begin, end, true, 3)
// 	ca, err := Create1D(id, begin, end, rule, 19)
// 	if err == nil {
// 		t.Error("Caso 2. Erro na criação do autômato celular que não capturou o erro. ")
// 	}
// 	if ca != nil {
// 		t.Error("Caso 2. Erro no retorno do autômato celular que retornou um autômato inválido")
// 	}
// 	fmt.Println(ca.ConfusionMatrix())
// }
