package queryfilter

import "go.mongodb.org/mongo-driver/bson"

type queryFilter struct {
	filters []bson.E
}

func newQF() *queryFilter {
	return &queryFilter{}
}

//New returns a pointer to a new queryFilter struct. By default, New will set is_deleted to false as part of the filters
//Each call to to New() will empty the filters slice and restart.
func New() *queryFilter {
	qf := newQF()
	qf.filters = append(qf.filters, bson.E{Key: "is_deleted", Value: false})
	return qf
}

// NewWithDeleted returns a pointer to a queryFilter struct. NewWithDeleted will set is_deleted to true as part of the filters
//Each call to to NewWithDeleted() will empty the filters slice and restart.
func NewWithDeleted() *queryFilter {
	qf := newQF()
	qf.filters = append(qf.filters, bson.E{Key: "is_deleted", Value: true})
	return qf
}

//AddFilter appends the provided filter to the list and returns the queryFilter pointer so calls to AddFilter() can be
//Chained.
func (qf *queryFilter) AddFilter(filter bson.E) *queryFilter {
	if filter.Key == "" {
		return qf
	}
	qf.filters = append(qf.filters, filter)
	return qf
}

//GetFilters returns all the added filters. Filters can no longer be added after this method is called.
func (qf *queryFilter) GetFilters() []bson.E {
	return qf.filters
}
