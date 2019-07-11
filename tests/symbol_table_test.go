package tests

import (
	"testing"
	"azula/compiler"
)

func TestDefine(t *testing.T) {
	expected := map[string]compiler.Symbol{
		"a": compiler.Symbol{Name: "a", Scope: compiler.GlobalScope, Index: 0},
		"b": compiler.Symbol{Name: "b", Scope: compiler.GlobalScope, Index: 1},
		"c": compiler.Symbol{Name: "c", Scope: compiler.LocalScope, Index: 0},
		"d": compiler.Symbol{Name: "d", Scope: compiler.LocalScope, Index: 1},
	}

	global := compiler.NewSymbolTable()

	a := global.Define("a")
	if a != expected["a"] {
		t.Errorf("expected a=%+v, got=%+v", expected["a"], a)
	}

	b := global.Define("b")
	if b != expected["b"] {
		t.Errorf("expected b=%+v, got=%+v", expected["b"], b)
	}

	local := compiler.NewEnclosedSymbolTable(global)

	c := local.Define("c")
	if c != expected["c"] {
		t.Errorf("expected c=%+v, got=%+v", expected["c"], c)
	}

	d := local.Define("d")
	if d != expected["d"] {
		t.Errorf("expected d=%+v, got=%+v", expected["d"], d)
	}
}

func TestResolveGlobal(t *testing.T) {
	global := compiler.NewSymbolTable()
	global.Define("a")
	global.Define("b")

	expected := []compiler.Symbol{
		compiler.Symbol{Name: "a", Scope: compiler.GlobalScope, Index: 0},
		compiler.Symbol{Name: "b", Scope: compiler.GlobalScope, Index: 1},
	}

	for _, sym := range expected {
		result, ok := global.Resolve(sym.Name)
		if !ok {
			t.Errorf("name %s not resolvable", sym.Name)
			continue
		}

		if result != sym {
			t.Errorf("expected %s to resolve to %+v, got=%+v", sym.Name, sym, result)
		}
	}
}

func TestResolveLocal(t *testing.T) {
	global := compiler.NewSymbolTable()
	global.Define("a")
	global.Define("b")

	local := compiler.NewEnclosedSymbolTable(global)
	local.Define("c")
	local.Define("d")

	expected := []compiler.Symbol{
		compiler.Symbol{Name: "a", Scope: compiler.GlobalScope, Index: 0},
		compiler.Symbol{Name: "b", Scope: compiler.GlobalScope, Index: 1},
		compiler.Symbol{Name: "c", Scope: compiler.LocalScope, Index: 0},
		compiler.Symbol{Name: "d", Scope: compiler.LocalScope, Index: 1},
	}

	for _, sym := range expected {
		result, ok := local.Resolve(sym.Name)
		if !ok {
			t.Errorf("name %s not resolvable", sym.Name)
			continue
		}

		if result != sym {
			t.Errorf("expected %s to resolve to %+v, got=%+v", sym.Name, sym, result)
		}
	}
}