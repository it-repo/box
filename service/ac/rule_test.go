package ac

import (
	"testing"
)

func TestRule_checkAND(t *testing.T) {
	type fields struct {
		rules []string
		both  bool
	}
	type args struct {
		rules []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			"and example1",
			fields{[]string{"a", "b"}, true},
			args{[]string{"a", "b", "c"}},
			true,
		},
		{
			"and example2",
			fields{[]string{"a", "b"}, true},
			args{[]string{"a", "c"}},
			false,
		},
		{
			"and example3",
			fields{[]string{"a", "b"}, true},
			args{[]string{"b", "c"}},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Rule{
				rules: tt.fields.rules,
				both:  tt.fields.both,
			}
			if got := r.checkAND(tt.args.rules); got != tt.want {
				t.Errorf("Rule.checkAND() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRule_checkOR(t *testing.T) {
	type fields struct {
		rules []string
		both  bool
	}
	type args struct {
		rules []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			"or example1",
			fields{[]string{"a", "b"}, false},
			args{[]string{"a", "c"}},
			true,
		},
		{
			"or example2",
			fields{[]string{"a", "b"}, false},
			args{[]string{"b", "c"}},
			true,
		},
		{
			"or example3",
			fields{[]string{"a", "b"}, false},
			args{[]string{"c"}},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Rule{
				rules: tt.fields.rules,
				both:  tt.fields.both,
			}
			if got := r.checkOR(tt.args.rules); got != tt.want {
				t.Errorf("Rule.checkOR() = %v, want %v", got, tt.want)
			}
		})
	}
}
