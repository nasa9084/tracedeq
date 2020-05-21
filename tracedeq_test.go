package tracedeq_test

import (
	"strconv"
	"testing"

	"github.com/nasa9084/tracedeq"
)

func compareResults(t *testing.T, got, want tracedeq.Result) {
	t.Helper()

	if got.IsEqual != want.IsEqual {
		t.Errorf("unexpected Result.IsEqual: %t != %t", got.IsEqual, want.IsEqual)
		return
	}
	if len(got.Trace) != len(want.Trace) {
		t.Errorf("unexpected length of Result.Trace: %d != %d", len(got.Trace), len(want.Trace))
		return
	}
	for i := 0; i < len(got.Trace); i++ {
		if got.Trace[i] != want.Trace[i] {
			t.Errorf("unexpected Result.Trace[%d]: %s != %s", i, got.Trace[i], want.Trace[i])
			return
		}
	}
	t.Logf("X: %v, %v", got.X, want.Y)
	t.Logf("Y: %v, %v", got.Y, want.Y)
}

func TestDeepEqualNil(t *testing.T) {
	got := tracedeq.DeepEqual(nil, nil)
	want := tracedeq.Result{
		IsEqual: true,
	}
	compareResults(t, got, want)
}

func TestDeepEqualDifferentTypes(t *testing.T) {
	got := tracedeq.DeepEqual("foo", 3)
	want := tracedeq.Result{
		IsEqual: false,
		Trace:   []string{"TYPE"},
		X:       "foo",
		Y:       3,
	}
	compareResults(t, got, want)
}

func TestDeepEqualString(t *testing.T) {
	tests := []struct {
		x, y string
		want tracedeq.Result
	}{
		{
			x: "",
			y: "",
			want: tracedeq.Result{
				IsEqual: true,
			},
		},
		{
			x: "foo",
			y: "foo",
			want: tracedeq.Result{
				IsEqual: true,
			},
		},
		{
			x: "foo",
			y: "bar",
			want: tracedeq.Result{
				IsEqual: false,
				X:       "foo",
				Y:       "bar",
			},
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := tracedeq.DeepEqual(tt.x, tt.y)
			compareResults(t, got, tt.want)
		})
	}
}

func TestDeepEqualBool(t *testing.T) {
	tests := []struct {
		x, y bool
		want tracedeq.Result
	}{
		{
			x: true,
			y: true,
			want: tracedeq.Result{
				IsEqual: true,
			},
		},
		{
			x: false,
			y: false,
			want: tracedeq.Result{
				IsEqual: true,
			},
		},
		{
			x: true,
			y: false,
			want: tracedeq.Result{
				IsEqual: false,
				X:       true,
				Y:       false,
			},
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := tracedeq.DeepEqual(tt.x, tt.y)
			compareResults(t, got, tt.want)
		})
	}
}

func TestDeepEqualInt(t *testing.T) {
	tests := []struct {
		x, y int
		want tracedeq.Result
	}{
		{
			x: 0,
			y: 0,
			want: tracedeq.Result{
				IsEqual: true,
			},
		},
		{
			x: 20,
			y: 20,
			want: tracedeq.Result{
				IsEqual: true,
			},
		},
		{
			x: 20,
			y: 25,
			want: tracedeq.Result{
				IsEqual: false,
				X:       20,
				Y:       25,
			},
		},
		{
			x: -20,
			y: 20,
			want: tracedeq.Result{
				IsEqual: false,
				X:       -20,
				Y:       20,
			},
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := tracedeq.DeepEqual(tt.x, tt.y)
			compareResults(t, got, tt.want)
		})
	}
}

func TestDeepEqualUint(t *testing.T) {
	tests := []struct {
		x, y uint
		want tracedeq.Result
	}{
		{
			x: 0,
			y: 0,
			want: tracedeq.Result{
				IsEqual: true,
			},
		},
		{
			x: 20,
			y: 20,
			want: tracedeq.Result{
				IsEqual: true,
			},
		},
		{
			x: 20,
			y: 25,
			want: tracedeq.Result{
				IsEqual: false,
				X:       20,
				Y:       25,
			},
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := tracedeq.DeepEqual(tt.x, tt.y)
			compareResults(t, got, tt.want)
		})
	}
}

func TestDeepEqualFloat(t *testing.T) {
	tests := []struct {
		x, y float64
		want tracedeq.Result
	}{
		{
			x: 0,
			y: 0,
			want: tracedeq.Result{
				IsEqual: true,
			},
		},
		{
			x: 20.3,
			y: 20.3,
			want: tracedeq.Result{
				IsEqual: true,
			},
		},
		{
			x: 20.5,
			y: 20.3,
			want: tracedeq.Result{
				IsEqual: false,
				X:       20.5,
				Y:       20.3,
			},
		},
		{
			x: -20,
			y: 20,
			want: tracedeq.Result{
				IsEqual: false,
				X:       -20,
				Y:       20,
			},
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := tracedeq.DeepEqual(tt.x, tt.y)
			compareResults(t, got, tt.want)
		})
	}
}

func TestDeepEqualArray(t *testing.T) {
	tests := []struct {
		x, y [3]string
		want tracedeq.Result
	}{
		{
			x: [3]string{"foo", "bar", "baz"},
			y: [3]string{"foo", "bar", "baz"},
			want: tracedeq.Result{
				IsEqual: true,
			},
		},
		{
			x: [3]string{"foo", "bar", ""},
			y: [3]string{"foo", "bar", ""},
			want: tracedeq.Result{
				IsEqual: true,
			},
		},
		{
			x: [3]string{"bar", "foo", "baz"},
			y: [3]string{"foo", "bar", "baz"},
			want: tracedeq.Result{
				IsEqual: false,
				Trace:   []string{"0"},
				X:       [3]string{"bar", "foo", "baz"},
				Y:       [3]string{"foo", "bar", "baz"},
			},
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := tracedeq.DeepEqual(tt.x, tt.y)
			compareResults(t, got, tt.want)
		})
	}
}

func TestDeepEqualSlice(t *testing.T) {
	same := []string{"foo", "bar"}
	tests := []struct {
		x, y []string
		want tracedeq.Result
	}{
		{
			x: []string{"foo", "bar", "baz"},
			y: []string{"foo", "bar", "baz"},
			want: tracedeq.Result{
				IsEqual: true,
			},
		},
		{
			x: same,
			y: same,
			want: tracedeq.Result{
				IsEqual: true,
			},
		},
		{
			x: []string{"bar", "foo", "baz"},
			y: []string{"foo", "bar", "baz"},
			want: tracedeq.Result{
				IsEqual: false,
				Trace:   []string{"0"},
				X:       []string{"bar", "foo", "baz"},
				Y:       []string{"foo", "bar", "baz"},
			},
		},
		{
			x: []string{"foo", "bar", "baz"},
			y: []string{"foo", "bar", "baz", "qux"},
			want: tracedeq.Result{
				IsEqual: false,
				Trace:   []string{"LENGTH"},
				X:       []string{"foo", "bar", "baz"},
				Y:       []string{"foo", "bar", "baz", "qux"},
			},
		},
		{
			x: []string{"foo", "bar", "baz"},
			y: nil,
			want: tracedeq.Result{
				IsEqual: false,
				X:       []string{"foo", "bar", "baz"},
				Y:       nil,
			},
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := tracedeq.DeepEqual(tt.x, tt.y)
			compareResults(t, got, tt.want)
		})
	}
}

func stringPtr(s string) *string { return &s }

func TestDeepEqualPointer(t *testing.T) {
	same := stringPtr("foo")
	tests := []struct {
		x, y *string
		want tracedeq.Result
	}{
		{
			x: stringPtr(""),
			y: stringPtr(""),
			want: tracedeq.Result{
				IsEqual: true,
			},
		},
		{
			x: stringPtr("foo"),
			y: stringPtr("foo"),
			want: tracedeq.Result{
				IsEqual: true,
			},
		},
		{
			x: same,
			y: same,
			want: tracedeq.Result{
				IsEqual: true,
			},
		},
		{
			x: stringPtr("foo"),
			y: stringPtr("bar"),
			want: tracedeq.Result{
				IsEqual: false,
				X:       "foo",
				Y:       "bar",
			},
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := tracedeq.DeepEqual(tt.x, tt.y)
			compareResults(t, got, tt.want)
		})
	}
}

func TestDeepEqualMap(t *testing.T) {
	tests := []struct {
		x, y map[string]string
		want tracedeq.Result
	}{
		{
			x: nil,
			y: nil,
			want: tracedeq.Result{
				IsEqual: true,
			},
		},
		{
			x: map[string]string{},
			y: map[string]string{},
			want: tracedeq.Result{
				IsEqual: true,
			},
		},
		{
			x: map[string]string{
				"foo": "hoge",
			},
			y: map[string]string{
				"foo": "hoge",
			},
			want: tracedeq.Result{
				IsEqual: true,
			},
		},
		{
			x: map[string]string{},
			y: nil,
			want: tracedeq.Result{
				IsEqual: false,
				X:       map[string]string{},
				Y:       nil,
			},
		},
		{
			x: map[string]string{
				"foo": "hoge",
			},
			y: map[string]string{
				"foo": "fuga",
			},
			want: tracedeq.Result{
				IsEqual: false,
				Trace:   []string{"foo"},
				X: map[string]string{
					"foo": "hoge",
				},
				Y: map[string]string{
					"foo": "fuga",
				},
			},
		},
		{
			x: map[string]string{
				"foo": "hoge",
			},
			y: map[string]string{
				"foo": "hoge",
				"bar": "fuga",
			},
			want: tracedeq.Result{
				IsEqual: false,
				Trace:   []string{"LENGTH"},
				X: map[string]string{
					"foo": "hoge",
				},
				Y: map[string]string{
					"foo": "hoge",
					"bar": "fuga",
				},
			},
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := tracedeq.DeepEqual(tt.x, tt.y)
			compareResults(t, got, tt.want)
		})
	}
}

func TestDeepEqualSimpleStruct(t *testing.T) {
	type ExportedStruct struct {
		ExportedField   string
		unexportedField string
	}
	tests := []struct {
		x, y ExportedStruct
		want tracedeq.Result
	}{
		{
			x: ExportedStruct{},
			y: ExportedStruct{},
			want: tracedeq.Result{
				IsEqual: true,
			},
		},
		{
			x: ExportedStruct{
				ExportedField: "foo",
			},
			y: ExportedStruct{
				ExportedField: "foo",
			},
			want: tracedeq.Result{
				IsEqual: true,
			},
		},
		{
			x: ExportedStruct{
				ExportedField: "foo",
			},
			y: ExportedStruct{},
			want: tracedeq.Result{
				IsEqual: false,
				Trace:   []string{"ExportedField"},
				X:       "foo",
				Y:       "",
			},
		},
		{
			x: ExportedStruct{
				unexportedField: "foo",
			},
			y: ExportedStruct{
				unexportedField: "foo",
			},
			want: tracedeq.Result{
				IsEqual: true,
			},
		},
		{
			x: ExportedStruct{
				unexportedField: "foo",
			},
			y: ExportedStruct{},
			want: tracedeq.Result{
				IsEqual: false,
				Trace:   []string{"unexportedField"},
				X:       "foo",
				Y:       "",
			},
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := tracedeq.DeepEqual(tt.x, tt.y)
			compareResults(t, got, tt.want)
		})
	}
}

type A struct {
	B *B
}
type B struct {
	C *C
}
type C struct {
	A *A
}

func TestDeepEqualCyclic(t *testing.T) {
	t.Run("each", func(t *testing.T) {
		x := A{
			B: &B{
				C: &C{},
			},
		}
		x.B.C.A = &x
		y := A{
			B: &B{
				C: &C{},
			},
		}
		y.B.C.A = &y
		got := tracedeq.DeepEqual(x, y)
		want := tracedeq.Result{
			IsEqual: true,
		}
		compareResults(t, got, want)
	})

	t.Run("same", func(t *testing.T) {
		x := A{
			B: &B{
				C: &C{},
			},
		}
		y := A{
			B: &B{
				C: &C{},
			},
		}
		x.B.C.A = &y
		y.B.C.A = &y

		got := tracedeq.DeepEqual(x, y)
		want := tracedeq.Result{
			IsEqual: true,
		}
		compareResults(t, got, want)
	})

	t.Run("large", func(t *testing.T) {
		x := A{
			B: &B{
				C: &C{
					A: &A{
						B: &B{
							C: &C{},
						},
					},
				},
			},
		}
		x.B.C.A.B.C.A = &x

		got := tracedeq.DeepEqual(x, *(x.B.C.A))
		want := tracedeq.Result{
			IsEqual: true,
		}
		compareResults(t, got, want)
	})
}

func TestDeepEqualFunc(t *testing.T) {
	t.Run("different", func(t *testing.T) {
		x := func() {}
		y := func() {}

		got := tracedeq.DeepEqual(x, y)
		want := tracedeq.Result{
			IsEqual: false,
			Trace:   []string{"FUNC"},
			X:       x,
			Y:       y,
		}
		compareResults(t, got, want)
	})

	t.Run("nil", func(t *testing.T) {
		var x func() = nil
		var y func() = nil

		got := tracedeq.DeepEqual(x, y)
		want := tracedeq.Result{
			IsEqual: true,
		}
		compareResults(t, got, want)
	})
}
