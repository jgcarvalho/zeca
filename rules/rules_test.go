package rules

import (
	"fmt"
	"testing"
)

func TestCreate(t *testing.T) {
	rule, _ := Create(prmDefault)
	fmt.Println(rule)
}
