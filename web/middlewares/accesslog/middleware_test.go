package accesslog

import (
	"reflect"
	"testing"
	"web"
)

func TestMiddlewareBuilder_Build(t *testing.T) {
	tests := []struct {
		name string
		want web.Middleware
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mi := MiddlewareBuilder{}
			if got := mi.Build(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Build() = %v, want %v", got, tt.want)
			}
		})
	}
}
