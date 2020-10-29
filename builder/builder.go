package builder

import (
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"sync"
)

type UpdateManyBuilder struct {
	m                sync.Mutex
	updateOperations bson.D
}

func NewUpdateManyBuilder() *UpdateManyBuilder {
	return &UpdateManyBuilder{}
}

// Get returns all the updateOperations.
// Test for empty value or use the HasValues() method to check if empty.
func (u *UpdateManyBuilder) Get() bson.D {
	return u.updateOperations
}

//HasValues checks if there are any values in updateOperations.
func (u *UpdateManyBuilder) HasValues() bool {
	if len(u.updateOperations) < 1 {
		return false
	}
	return true
}

// Add creates a structure used for an UpdateMany command. The order in which command are provided are preserved.
// Example:
// u.Add(operator.Set, bson.E{Key: "name", "Asari"}).
//		Add(operator.Mul, bson.E{Key: "count", Value: 2}).
// 		Add(operator.Set, bson.E{Key: "email", "asari@gmail.com"})
// Will result in {$set: {name: "Asari", email: "asari@gmail.com"}, $mul: {count: 2}}
func (u *UpdateManyBuilder) Add(operator string, values ...bson.E) *UpdateManyBuilder {
	if len(values) < 1 {
		log.Println("Call to Add() with zero fields....Ignoring this :) ")
		return u
	}

	u.m.Lock()
	defer u.m.Unlock()

	for key, updateOperator := range u.updateOperations {
		if updateOperator.Key == operator {
			tmp := append(updateOperator.Value.([]bson.E), values...)
			u.updateOperations[key].Value = tmp
			return u
		}
	}

	u.updateOperations = append(u.updateOperations, bson.E{Key: operator, Value: values})
	return u
}
