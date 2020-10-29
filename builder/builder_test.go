package builder

import (
	"github.com/jcobhams/asari/operator"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
)

func TestUpdateManyBuilder_AddGetHasValues(t *testing.T) {
	b := NewUpdateManyBuilder()
	assert.False(t, b.HasValues())

	//Test No Fields
	b.Add(operator.Set)
	assert.Equal(t, 0, len(b.Get()))

	b.Add(operator.Set, bson.E{Key: "name", Value: "asari"}, bson.E{Key: "score", Value: 500}).
		Add(operator.Unset, bson.E{Key: "email", Value: "asari@gmail.com"}).
		Add(operator.Mul, bson.E{Key: "count", Value: 2})

	assert.Equal(t, 3, len(b.Get()))
	assert.Equal(t, 2, len(b.Get()[0].Value.([]bson.E)))
	assert.Equal(t, 1, len(b.Get()[1].Value.([]bson.E)))
	assert.Equal(t, 1, len(b.Get()[2].Value.([]bson.E)))
}
