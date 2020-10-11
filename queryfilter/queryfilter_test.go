package queryfilter

import (
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
)

func TestNew(t *testing.T) {
	qf := New()

	assert.Equal(t, "is_deleted", qf.GetFilters()[0].Key)
	assert.Equal(t, false, qf.GetFilters()[0].Value)
}

func TestNewWithDeleted(t *testing.T) {
	qf := NewWithDeleted()

	assert.Equal(t, "is_deleted", qf.GetFilters()[0].Key)
	assert.Equal(t, true, qf.GetFilters()[0].Value)
}

func TestQueryFilter_AddFilter(t *testing.T) {
	qf := New().AddFilter(bson.E{Key: "someKey", Value: "someValue"})

	assert.Equal(t, "someKey", qf.GetFilters()[1].Key)
	assert.Equal(t, "someValue", qf.GetFilters()[1].Value)

	//Test Blank Keys Are Not Added
	qf.AddFilter(bson.E{})
	assert.Equal(t, 2, len(qf.GetFilters()))
}

func TestQueryFilter_GetFilters(t *testing.T) {
	qf := New().
		AddFilter(bson.E{Key: "someKey", Value: "someValue"}).
		AddFilter(bson.E{Key: "someKey", Value: "someValue"}).
		AddFilter(bson.E{Key: "someKey", Value: "someValue"})

	assert.Equal(t, 4, len(qf.GetFilters()))
}
