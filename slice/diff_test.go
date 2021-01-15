package slice

import (
	"reflect"
	"testing"
)

func TestDiff(t *testing.T) {
	type args struct {
		a []string
		b []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "It should return a single element when it's not present",
			args: args{
				a: []string{"a", "b", "c"},
				b: []string{"b", "a"},
			},
			want: []string{"c"},
		},
		{
			name: "It should return the diff when the two slices have the same size but different elements",
			args: args{
				a: []string{"_DSC2179.JPG", "999.JPG"},
				b: []string{"_DSC2179.JPG", "_DSC2180.JPG"},
			},
			want: []string{"999.JPG"},
		},
		{
			name: "It should return no diff when the two slices have different sizes but the same elements",
			args: args{
				a: []string{"_DSC2179.JPG"},
				b: []string{"_DSC2179.JPG", "_DSC2180.JPG"},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Diff(tt.args.a, tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Diff() = %v, want %v", got, tt.want)
			}
		})
	}
}
