package printer_test

import (
	"testing"

	"github.com/minodisk/thriftfmt/formatter/printer"
)

func TestIndent_ByReference(t *testing.T) {
	for _, c := range []struct {
		base int
		ref  int
		want int
	}{
		{
			0,
			1,
			0,
		},
		{
			1,
			2,
			1,
		},
		{
			2,
			3,
			2,
		},
	} {
		i := printer.NewIndent(c.base)
		i.Increment()
		func(i *printer.Indent) {
			if int(*i) != c.ref {
				t.Errorf("can't receive by reference")
			}
			i.Decrement()
		}(i)
		if int(*i) != c.want {
			t.Errorf("can't pass by reference")
		}
	}
}

func TestIndent_Increment(t *testing.T) {
	for _, c := range []struct {
		base   int
		repeat int
		want   int
	}{
		{
			0,
			0,
			0,
		},
		{
			0,
			1,
			1,
		},
		{
			0,
			2,
			2,
		},
		{
			1,
			0,
			1,
		},
		{
			1,
			1,
			2,
		},
		{
			1,
			2,
			3,
		},
		{
			2,
			0,
			2,
		},
		{
			2,
			1,
			3,
		},
		{
			2,
			2,
			4,
		},
	} {
		indent := printer.NewIndent(c.base)
		for i := 0; i < c.repeat; i++ {
			indent.Increment()
		}
		got := int(*indent)
		if got != c.want {
			t.Errorf("want %d, but got %d", c.want, got)
		}
	}
}

func TestIndent_Decrement(t *testing.T) {
	for _, c := range []struct {
		base   int
		repeat int
		want   int
	}{
		{
			0,
			0,
			0,
		},
		{
			0,
			1,
			-1,
		},
		{
			0,
			2,
			-2,
		},
		{
			1,
			0,
			1,
		},
		{
			1,
			1,
			0,
		},
		{
			1,
			2,
			-1,
		},
		{
			2,
			0,
			2,
		},
		{
			2,
			1,
			1,
		},
		{
			2,
			2,
			0,
		},
	} {
		indent := printer.NewIndent(c.base)
		for i := 0; i < c.repeat; i++ {
			indent.Decrement()
		}
		got := int(*indent)
		if got != c.want {
			t.Errorf("want %d, but got %d", c.want, got)
		}
	}
}

func TestIndent_String(t *testing.T) {
	for _, c := range []struct {
		indent int
		want   string
	}{
		{
			0,
			"",
		},
		{
			1,
			"  ",
		},
		{
			2,
			"    ",
		},
	} {
		i := printer.NewIndent(c.indent)
		got := i.String()
		if got != c.want {
			t.Errorf("\nwant: '%s'\n got: '%s'", c.want, got)
		}
	}
}
