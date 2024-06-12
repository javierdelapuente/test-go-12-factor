package config

import (
	"reflect"
	"testing"
)

func TestBuildCharmConfig(t *testing.T) {
	postgresqlUrl := "postgresql://test-username:test-password@test-postgresql:5432/test-database"
	tests := map[string]struct {
		environ  []string
		expected CharmConfig
	}{
		"bare configs": {
			environ: []string{"x=1", "y=2"},
			expected: CharmConfig{
				Configs: []Config{{Name: "x", Value: "1"}, {Name: "y", Value: "2"}},
			},
		},
		"postgresql": {
			environ: []string{"x=1", "POSTGRESQL_DB_CONNECT_STRING=" + postgresqlUrl},
			expected: CharmConfig{
				Configs: []Config{{Name: "x", Value: "1"}},
				Integrations: Integrations{
					PostgresqlUrl: &postgresqlUrl,
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actual := BuildCharmConfig(test.environ)

			if !reflect.DeepEqual(actual, test.expected) {
				t.Errorf("CharmConfig does not match: expected %#v, got %#v", test.expected, actual)
			}
		})
	}
}
