package testpkg

import "testing"

func TestExample(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		h := MyOpts{Yo: "son"}
		ToMyOptsWithDefault(&h, NewMyOptsBuilder().
			Total(2).
			Build()...,
		)

		t.Logf("%+v", h)

		if h.Total != 2 {
			t.Fatal("expected total to be 2 but got", h.Total)
		}

		if h.Yo != "son" {
			t.Fatal("expected yo to be son but got", h.Yo)
		}
	})
}
