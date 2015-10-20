package cga

import (
	"fmt"
	"testing"

	"bitbucket.org/jgcarvalho/zeca/rules"
)

func TestNewProbs(t *testing.T) {
	// p := NewProbs(rules.PrmDefault)
	// fmt.Println(p)
}

func TestGenRule(t *testing.T) {
	rule, _ := rules.Create(rules.PrmDefault)
	p := NewProbs(rule.Prm)
	ruleNew1 := p.GenRule()
	ruleNew2 := p.GenRule()
	fmt.Println("@@@Nova regra (1)", ruleNew1)
	fmt.Println("@@@Nova regra (2)", ruleNew2)
}

func TestAdjustByDuel(t *testing.T) {
	rule, _ := rules.Create(rules.PrmDefault)
	p := NewProbs(rule.Prm)
	fmt.Println(p)
	r1 := p.GenRule()
	r2 := p.GenRule()
	p.AdjustByDuel(r1, r2)
	for i := 0; i < 1000; i++ {
		r1 = p.GenRule()
		r2 = p.GenRule()
		p.AdjustByDuel(r1, r2)
		fmt.Println("Geração ", i)
		fmt.Println(p)
	}
}
