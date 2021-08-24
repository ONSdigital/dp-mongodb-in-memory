package main

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestDefaultOptions(t *testing.T) {
	versions := []string{"4.4.8", "5.0.2"}

	for _, version := range versions {
		t.Run(version, func(t *testing.T) {
			server, err := Start(version)
			require.NoError(t, err)
			defer server.Stop()

			client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(server.URI()))
			require.NoError(t, err)

			require.NoError(t, client.Ping(context.Background(), nil))
		})
	}
}
