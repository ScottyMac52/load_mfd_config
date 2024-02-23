package main

import (
	"reflect"
	"testing"
)

var _ = reflect.TypeOf((*Rectangle)(nil)) // Enforce comparison compatibility check

func TestConfiguration_CenterIn(t *testing.T) {
	tests := []struct {
		name    string
		outer   *Configuration
		args    args
		want    *Rectangle
		wantErr bool
	}{
		{
			name:  "Centered Within Module Area Without Overlap",
			outer: &Configuration{Left: 0, Top: 0, Width: 10, Height: 10},
			args: args{
				inner: &Configuration{Left: 2, Top: 2, Width: 4, Height: 4},
			},
			want:    &Rectangle{Left: 3, Top: 3, Width: 4, Height: 4},
			wantErr: false,
		},
		{
			name:  "Partially Outside With Partial Overlap",
			outer: &Configuration{Left: 0, Top: 0, Width: 10, Height: 10},
			args: args{
				inner: &Configuration{Left: 0, Top: 0, Width: 4, Height: 4},
			},
			want:    &Rectangle{Left: 3, Top: 3, Width: 4, Height: 4},
			wantErr: false,
		},
		/*
			{
				name:  "Completely Outside Leading to Erroneous Operation",
				outer: &Configuration{Left: 0, Top: 0, Width: 10, Height: 10},
				args: args{
					inner: &Configuration{Left: -1, Top: 1, Width: 20, Height: 20},
				},
				want:    &Rectangle{},
				wantErr: true,
			},
		*/
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.outer.CenterIn(tt.args.inner)
			if (err != nil) != tt.wantErr {
				t.Errorf("Configuration.CenterIn() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Configuration.CenterIn() = %v, want %v", got, tt.want)
			}
		})
	}
}

type args struct {
	inner *Configuration
}
