package feature

import (
	"testing"

	"github.com/atotto/gobdd/_testdata/calc"
)

var cal *calc.Calc

func Given_IHaveEntered_val1_IntoTheCalculator(t *testing.T, v int64) {
	cal = calc.NewCalc()
	cal.Push(v)
}

func And_IHaveEntered_val1_IntoTheCalculator(t *testing.T, v int64) {
	cal.Push(v)
}

func When_IPress_val1_(t *testing.T, button string) {
	switch button {
	case "add":
		cal.Add()
	}
}

func Then_TheResultShouldBe_val1_OnTheScreen(t *testing.T, v int64) {
	if v != cal.Result() {
		t.Errorf("want %d, got %d", v, cal.Result())
	}
}
