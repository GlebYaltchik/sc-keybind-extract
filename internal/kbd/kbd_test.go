package kbd

import "testing"

func TestNormalize(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		args string
		want string
	}{
		{
			name: "empty",
			args: "",
			want: "",
		},
		{
			name: "letter",
			args: "k",
			want: "K",
		},
		{
			name: "word",
			args: "upper",
			want: "Upper",
		},
		{
			name: "save case",
			args: "uppEr",
			want: "UppEr",
		},
		{
			name: "two letter word",
			args: "lctrl",
			want: "LCtrl",
		},
		{
			name: "plus combo",
			args: "ctrl+y",
			want: "Ctrl+Y",
		},
		{
			name: "underline combo",
			args: "wheel_up",
			want: "Wheel_Up",
		},
		{
			name: "mixed combo",
			args: "lctrl+wheel_up",
			want: "LCtrl+Wheel_Up",
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := Normalize(tt.args); got != tt.want {
				t.Errorf("Normalize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_normalize(t *testing.T) {
	type args struct {
		v string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := normalize(tt.args.v); got != tt.want {
				t.Errorf("normalize() = %v, want %v", got, tt.want)
			}
		})
	}
}
