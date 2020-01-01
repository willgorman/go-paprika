package paprika_test

import (
	"os"
	"reflect"
	"testing"

	"github.com/willgorman/go-paprika"
)

var (
	User, Password string
)

func init() {
	var ok bool
	User, ok = os.LookupEnv("PAPRIKA_USER")
	if !ok {
		panic("PAPRIKA_USER must be set")
	}

	Password, ok = os.LookupEnv("PAPRIKA_PASSWORD")
	if !ok {
		panic("PAPRIKA_PASSWORD must be set")
	}
}

func TestClient_Recipes(t *testing.T) {
	type fields struct {
		username string
		password string
	}
	tests := []struct {
		name    string
		fields  fields
		want    []paprika.RecipeItem
		wantErr bool
	}{
		{
			name: "test",
			fields: fields{
				username: User,
				password: Password,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := paprika.NewClient(tt.fields.username, tt.fields.password)
			got, err := c.Recipes()
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.Recipes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.Recipes() = %v, want %v", got, tt.want)
			}
		})
	}
}
