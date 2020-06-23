package haro

import (
	"reflect"
	"testing"
)

func Test_registry_Exists(t *testing.T) {
	type args struct {
		topicName string
	}
	tests := []struct {
		name string
		t    registry
		args args
		want bool
	}{
		{
			name: "Topic exists",
			t: map[string]*topic{
				"test": &topic{},
			},
			args: args{
				topicName: "test",
			},
			want: true,
		},
		{
			name: "Topic not exist",
			t: map[string]*topic{
				"test": &topic{},
			},
			args: args{
				topicName: "lol",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.t.Exists(tt.args.topicName); got != tt.want {
				t.Errorf("Exists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_registry_Get(t *testing.T) {
	type args struct {
		topicName string
	}
	tests := []struct {
		name    string
		t       registry
		args    args
		want    *topic
		wantErr bool
	}{
		{
			name: "Topic retrieved",
			t: map[string]*topic{
				"test": &topic{},
			},
			args: args{
				topicName: "test",
			},
			want:    &topic{},
			wantErr: false,
		},
		{
			name: "Topic not found",
			t: map[string]*topic{
				"test": &topic{},
			},
			args: args{
				topicName: "lol",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.t.Get(tt.args.topicName)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}
