package reflections

import (
	"reflect"
	"testing"
	"time"
)

func TestNewValue(t *testing.T) {
	tests := []struct {
		name    string
		value   any
		kind    reflect.Value
		want    reflect.Value
		wantErr bool
	}{
		{
			name:  "string to bool",
			value: "true",

			want:    reflect.ValueOf(true),
			wantErr: false,
		},
		{
			name:    "string to time",
			value:   "2025-07-30 20:00:00",
			want:    reflect.ValueOf(time.Date(2025, 7, 30, 20, 0, 0, 0, time.UTC)),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := NewValue(tt.value, tt.want)

			if (err != nil) != tt.wantErr {
				t.Errorf("newValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newValue() = %v, want %v", got, tt.want)
			}
		})
	}
}
