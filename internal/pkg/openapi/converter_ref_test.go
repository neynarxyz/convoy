package openapi

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConverter_ExtractWebhooks_OpenAPI31Refs(t *testing.T) {
	spec := []byte(`{
  "openapi": "3.1.0",
  "info": {"title": "refs", "version": "1"},
  "webhooks": {"user.created": {"$ref": "#/components/pathItems/UserCreated"}},
  "components": {
    "pathItems": {"UserCreated": {"post": {"description": "user created",
      "requestBody": {"$ref": "#/components/requestBodies/User"}}}},
    "requestBodies": {"User": {"content": {"application/json": {
      "schema": {"$ref": "#/components/schemas/User"}}}}},
    "schemas": {"User": {"type": "object", "required": ["id"],
      "properties": {"id": {"type": "string", "minLength": 3}}}}
  }
}`)

	c, err := NewFromBytes(spec)
	require.NoError(t, err)
	collection, err := c.ExtractWebhooks()
	require.NoError(t, err)

	webhook := collection.Webhooks["user.created"]
	require.NotNil(t, webhook)
	require.Equal(t, "user created", webhook.Description)
	require.Equal(t, []string{"id"}, webhook.Schema.Required)
	require.Equal(t, uint64(3), webhook.Schema.Properties["id"].Value.MinLength)
}
