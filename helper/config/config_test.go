package config

import "testing"

func TestGetConfig(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			"get config",
			args{filename: "domain"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GetDomainConfig(tt.args.filename)
		})
	}
}
