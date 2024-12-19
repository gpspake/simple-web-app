package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitializeApp(t *testing.T) {
	db, e := InitializeApp()

	assert.NotNil(t, db, "Database should not be nil")
	assert.NotNil(t, e, "Echo instance should not be nil")

	// Additional checks can be added for routes, templates, etc.
	defer db.Close()
}
