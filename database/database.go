package database

import (
	"context"
	"errors"
	"fmt"
	"github.com/jcobhams/asari/builder"
	"github.com/jcobhams/asari/document"
	"github.com/jcobhams/asari/operator"
	"github.com/jcobhams/asari/queryfilter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"reflect"
	"time"
)

type client struct {
	Connection *mongo.Database
}

var (
	Instance *client
)

func Init(mongoDSN, databaseName string) *client {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	mClient, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoDSN))

	if err != nil {
		log.Panic(fmt.Sprintf("Could Not Connect To Server - %v", mongoDSN))
	}

	err = mClient.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Panic("Could Not Connect To Primary shard")
	}
	database := mClient.Database(databaseName)

	return &client{Connection: database}
}

func (c *client) findOne(collection string, filters []bson.E, target interface{}, findOneOptions ...*options.FindOneOptions) error {
	if err := c.validateDocumentKind(target); err != nil {
		return err
	}

	filters = c.applyIsDeletedFilter(filters)
	if err := c.validateFilters(filters); err != nil {
		return err
	}

	if preFindOne, ok := target.(document.PreFindOne); ok {
		if err := preFindOne.PreFindOne(c.Connection); err != nil {
			return errors.New(fmt.Sprintf("asari: PreOneUpdate Hook Error: %v", err))
		}
	}

	err := c.Connection.Collection(collection).FindOne(context.TODO(), filters, findOneOptions...).Decode(target)

	if err == nil {
		if postFindOne, ok := target.(document.PostFindOne); ok {
			if err := postFindOne.PostFindOne(c.Connection); err != nil {
				return errors.New(fmt.Sprintf("asari: PostOneUpdate Hook Error: %v", err))
			}
		}
	}

	return err
}

// FindOne searches for a single document that matches the provided filters.
// If projection is nil, all fields are returned.
// To specify only select fields, use a bson.M - eg: bson.M{"email":1, "phone":1}
// Target has to be a pointer to a struct where the document will be unmarshalled into.
func (c *client) FindOne(collection string, filters []bson.E, projection, target interface{}, findOneOptions ...*options.FindOneOptions) error {
	if err := c.validateProjection(projection); err != nil {
		return err
	}
	opts := &options.FindOneOptions{
		Projection: projection,
	}
	findOneOptions = append(findOneOptions, opts)

	return c.findOne(collection, filters, target, findOneOptions...)
}

// FindOneByID finds a document that matches the provided ID in the collection.
// If projection is nil, all fields are returned.
// To specify only select fields, use a bson.M - eg: bson.M{"email":1, "phone":1}
// Target has to be a pointer to a struct where the document will be unmarshalled into.
func (c *client) FindOneByID(collection string, id primitive.ObjectID, projection, target interface{}) error {
	if err := c.validateProjection(projection); err != nil {
		return err
	}
	filters := queryfilter.New().AddFilter(bson.E{Key: "_id", Value: id}).GetFilters()
	opts := &options.FindOneOptions{
		Projection: projection,
	}
	return c.findOne(collection, filters, target, opts)
}

// FindOneByField finds a document that matches the provided field and value pair.
// If projection is nil, all fields are returned.
// To specify only select fields, use a bson.M - eg: bson.M{"email":1, "phone":1}
// Target has to be a pointer to a struct where the document will be unmarshalled into.
func (c *client) FindOneByField(collection, field string, value, projection, target interface{}) error {
	if err := c.validateProjection(projection); err != nil {
		return err
	}

	filters := queryfilter.New().AddFilter(bson.E{Key: field, Value: value}).GetFilters()
	opts := &options.FindOneOptions{
		Projection: projection,
	}
	return c.findOne(collection, filters, target, opts)
}

// FindPaginated searches for document that matches the provided filters.
// PageOpts control CurrentPage and PerPage value.
// If projection is nil, all fields are returned.
// To specify only select fields, use a bson.M - eg: bson.M{"email":1, "phone":1}
// sort should be a bson.D - eg: bson.D{bson.E{Key: "_id", Value: -1}, bson.E{Key: "another, Value: "value"}}
// FindPaginated will return the Mongo Cursor in the PaginatedResult struct.
// REMEMBER TO CALL Cursor.Close(ctx) WHEN DONE READING
func (c *client) FindPaginated(collection string, pageOptions PageOpts, filters []bson.E, projection interface{}, sort bson.D) (*PaginatedResult, error) {
	if sort == nil {
		sort = bson.D{bson.E{Key: "_id", Value: -1}}
	}

	filters = c.applyIsDeletedFilter(filters)
	if err := c.validateFilters(filters); err != nil {
		return nil, err
	}

	if err := c.validateProjection(projection); err != nil {
		return nil, err
	}

	paginator := NewPaginator(pageOptions)
	paginator.SetOffset()
	opts := &options.FindOptions{
		Projection: projection,
		Skip:       &paginator.Offset,
		Limit:      &paginator.PerPage,
		Sort:       sort,
	}

	totalRows, err := c.Connection.Collection(collection).CountDocuments(nil, filters)
	if err != nil {
		return nil, err
	}
	paginator.TotalRows = totalRows

	cur, err := c.Connection.Collection(collection).Find(context.TODO(), filters, opts)
	if err != nil {
		return nil, err
	}

	paginator.SetTotalPages()
	paginator.SetPrevPage()
	paginator.SetNextPage()
	return &PaginatedResult{Cursor: cur, Paginator: *paginator}, nil
}

// FindLast returns the most recent document in the collection that matches the provided filters.
// It sorts based on the mongo objectId
func (c *client) FindLast(collection string, filters []bson.E, projection, target interface{}) error {
	if err := c.validateDocumentKind(target); err != nil {
		return err
	}

	cur, err := c.findLast(collection, 1, filters, projection)
	if err != nil {
		return err
	}

	hasResults := false
	for cur.Next(nil) {
		cur.Decode(target)
		hasResults = true
	}

	if hasResults {
		return nil
	}
	return mongo.ErrNoDocuments
}

// FindLastN returns the N (limit) most recent documents in the collection that matches the provided filters.
// It sorts based on provided mongo objectId
func (c *client) FindLastN(collection string, limit int, filters []bson.E, projection interface{}) (*mongo.Cursor, error) {
	return c.findLast(collection, limit, filters, projection)
}

func (c *client) findLast(collection string, limit int, filters []bson.E, projection interface{}) (*mongo.Cursor, error) {
	sort := bson.D{bson.E{Key: "_id", Value: -1}}

	if err := c.validateProjection(projection); err != nil {
		return nil, err
	}

	filters = c.applyIsDeletedFilter(filters)
	if err := c.validateFilters(filters); err != nil {
		return nil, err
	}

	opts := &options.FindOptions{
		Projection: projection,
		Sort:       sort,
	}
	opts.SetLimit(int64(limit))

	return c.Connection.Collection(collection).Find(context.TODO(), filters, opts)
}

func (c *client) applyIsDeletedFilter(filters []bson.E) []bson.E {
	//Extra Redundancy Incase is_deleted is not provided, default to false

	for _, v := range filters {
		if v.Key == "is_deleted" {
			return filters
		}
	}

	filters = append(filters, bson.E{Key: "is_deleted", Value: false})
	return filters
}

// FindAll - returns a list of all the document that match the filter or returns an error.
// To be used with care as a lot of document could be returned and use up a lot of memory.
func (c *client) FindAll(collection string, filters []bson.E, projection interface{}, sort bson.D) (*mongo.Cursor, error) {
	if sort == nil {
		sort = bson.D{bson.E{Key: "_id", Value: -1}}
	}

	if err := c.validateProjection(projection); err != nil {
		return nil, err
	}

	filters = c.applyIsDeletedFilter(filters)
	if err := c.validateFilters(filters); err != nil {
		return nil, err
	}

	opts := &options.FindOptions{
		Projection: projection,
		Sort:       sort,
	}

	return c.Connection.Collection(collection).Find(context.TODO(), filters, opts)
}

func (c *client) updateDocument(collection string, filters []bson.E, doc interface{}) (*mongo.SingleResult, error) {
	if err := c.validateFilters(filters); err != nil {
		return nil, err
	}

	doc.(document.Document).BeforeUpdate()
	result := c.Connection.Collection(collection).FindOneAndReplace(context.TODO(), filters, doc)
	if result.Err() != nil {
		return nil, errors.New(result.Err().Error())
	}
	return result, nil
}

// SaveDocument will create a new document or update an existing document if doc is not new.
// if doc is new and implements the PreCreator interface, the PreCreate hook will fire or return appropriate error.
// if doc is new and implements the PostCreator interface, the PostCreate hook will fire or return appropriate error.
// if doc is existing and implements the PreUpdater interface, the PreUpdate hook will fire or return appropriate error.
// if doc is new and implements the PostUpdater interface, the PostUpdate hook will fire or return appropriate error.
func (c *client) SaveDocument(collection string, doc interface{}) (interface{}, error) {
	if err := c.validateDocumentKind(doc); err != nil {
		return nil, err
	}

	if !doc.(document.Document).CanSave() {
		return nil, errors.New("asari: cannot save new document. call document.Setup() before calling SaveDocument()")
	}

	if doc.(document.Document).IsNew() {

		if preCreator, ok := doc.(document.PreCreator); ok {
			if err := preCreator.PreCreate(c.Connection); err != nil {
				return nil, errors.New(fmt.Sprintf("asari: PreCreate Hook Error: %v", err))
			}
		}

		_, err := c.Connection.Collection(collection).InsertOne(context.TODO(), doc)
		if err == nil {
			doc.(document.Document).SetIsNew(false)

			if postCreator, ok := doc.(document.PostCreator); ok {
				if err := postCreator.PostCreate(c.Connection); err != nil {
					return nil, errors.New(fmt.Sprintf("asari: PostCreate Hook Error: %v", err))
				}
			}
		}
		return doc, err
	} else {

		if preUpdater, ok := doc.(document.PreUpdater); ok {
			if err := preUpdater.PreUpdate(c.Connection); err != nil {
				return nil, errors.New(fmt.Sprintf("asari: PreUpdate Hook Error: %v", err))
			}
		}

		qf := queryfilter.New().AddFilter(bson.E{Key: "_id", Value: doc.(document.Document).GetID()}).GetFilters()
		_, err := c.updateDocument(collection, qf, doc)

		if err == nil {
			if postUpdater, ok := doc.(document.PostUpdater); ok {
				if err := postUpdater.PostUpdate(c.Connection); err != nil {
					return nil, errors.New(fmt.Sprintf("asari: PostUpdate Hook Error: %v", err))
				}
			}
		}

		return doc, err
	}
}

// UpdateMany finds the documents that match the filter and update them based on the operators configured in the UpdateManyBuilder
func (c *client) UpdateMany(collection string, filters []bson.E, updateBuilder *builder.UpdateManyBuilder, updateOptions *options.UpdateOptions) (*mongo.UpdateResult, error) {
	if updateBuilder.HasValues() {
		return c.Connection.Collection(collection).UpdateMany(context.TODO(), filters, updateBuilder.Get(), updateOptions)
	}
	return nil, errors.New("empty UpdateManyBuilder provided")
}

// CountDocuments returns a count of all the documents that match the provided filters or error otherwise
func (c *client) CountDocuments(collection string, filters interface{}) (int, error) {
	count, err := c.Connection.Collection(collection).CountDocuments(context.TODO(), filters)
	return int(count), err
}

// SoftDeleteDocument marks a document as deleted and sets the deleted timestamp. This does not remove the item from the
// DB but it hides it from future queries except deleted records is added to the filters
func (c *client) SoftDeleteDocument(collection string, doc interface{}) (*mongo.SingleResult, error) {
	if err := c.validateDocumentKind(doc); err != nil {
		return nil, err
	}

	d := doc.(document.Document)
	d.BeforeSoftDelete()
	id := d.GetID()

	qf := queryfilter.New().AddFilter(bson.E{Key: "_id", Value: id}).GetFilters()

	if preSoftDeleter, ok := doc.(document.PreSoftDeleter); ok {
		if err := preSoftDeleter.PreSoftDelete(c.Connection); err != nil {
			return nil, errors.New(fmt.Sprintf("asari: PreSoftDeleter Hook Error: %v", err))
		}
	}

	result, err := c.updateDocument(collection, qf, doc)

	if err == nil {
		if postSoftDeleter, ok := doc.(document.PostSoftDeleter); ok {
			if err := postSoftDeleter.PostSoftDelete(c.Connection); err != nil {
				return nil, errors.New(fmt.Sprintf("asari: PostSoftDeleter Hook Error: %v", err))
			}
		}
	}

	return result, err
}

// HardDeleteDocument deletes a record from the DB. Careful with this as the document is irrecoverable.
// Use SoftDeleteDocument() instead except you want the document truly gone.
func (c *client) HardDeleteDocument(collection string, doc interface{}) (*mongo.DeleteResult, error) {
	if err := c.validateDocumentKind(doc); err != nil {
		return nil, err
	}

	qf := queryfilter.New().
		AddFilter(bson.E{Key: "_id", Value: doc.(document.Document).GetID()}).
		GetFilters()

	if preHardDeleter, ok := doc.(document.PreHardDeleter); ok {
		if err := preHardDeleter.PreHardDelete(c.Connection); err != nil {
			return nil, errors.New(fmt.Sprintf("asari: PreHardDeleter Hook Error: %v", err))
		}
	}

	result, err := c.Connection.Collection(collection).DeleteOne(context.TODO(), qf)

	if err == nil {
		if postHardDeleter, ok := doc.(document.PostHardDeleter); ok {
			if err := postHardDeleter.PostHardDelete(c.Connection); err != nil {
				return nil, errors.New(fmt.Sprintf("asari: PostHardDeleter Hook Error: %v", err))
			}
		}
	}
	return result, err
}

func (c *client) aggregate(collection string, pipeline mongo.Pipeline, aggregateOptions *options.AggregateOptions) (*mongo.Cursor, error) {
	return c.Connection.Collection(collection).Aggregate(nil, pipeline, aggregateOptions)
}

// Aggregate runs a simple aggregation pipeline and returns a cursor if successful or error if any.
// If no aggregation options are provided, allowDiskUse is set to true by default.
func (c *client) Aggregate(collection string, pipeline mongo.Pipeline, aggregateOptions *options.AggregateOptions) (*mongo.Cursor, error) {

	if aggregateOptions == nil {
		aggregateOptions = &options.AggregateOptions{}
		aggregateOptions.SetAllowDiskUse(true)
	}

	return c.aggregate(collection, pipeline, aggregateOptions)
}

// AggregatePaginated runs an aggregation that will have paginated results.
func (c *client) AggregatePaginated(collection string, pageOptions PageOpts, pipeline mongo.Pipeline, aggregateOptions *options.AggregateOptions) (*AggregationPaginatedResult, error) {

	if aggregateOptions == nil {
		aggregateOptions = &options.AggregateOptions{}
		aggregateOptions.SetAllowDiskUse(true)
	}

	paginator := NewPaginator(pageOptions)
	paginator.SetOffset()

	facetPipeline := bson.D{
		{
			operator.Facet,
			bson.D{
				{Key: "meta", Value: bson.A{bson.M{operator.Count: "total"}}},
				{Key: "data", Value: bson.A{bson.M{operator.Skip: paginator.Offset}, bson.M{operator.Limit: paginator.PerPage}}},
			},
		},
	}

	pipeline = append(pipeline, facetPipeline)

	cur, err := c.aggregate(collection, pipeline, aggregateOptions)
	if err != nil {
		return nil, err
	}
	defer cur.Close(nil)

	type (
		meta struct {
			Total int64 `json:"total" bson:"total"`
		}
		result struct {
			Meta []meta     `json:"meta" bson:"meta"`
			Data []bson.Raw `json:"data" bson:"data"`
		}
	)

	var r result
	for cur.Next(nil) {
		var tmpR result
		if err := cur.Decode(&tmpR); err != nil {
			return nil, err
		}
		r = tmpR
	}

	if len(r.Data) > 0 && len(r.Meta) > 0 {

		paginator.TotalRows = r.Meta[0].Total
		paginator.SetTotalPages()
		paginator.SetPrevPage()
		paginator.SetNextPage()

		return &AggregationPaginatedResult{Paginator: *paginator, Data: r.Data}, nil
	}

	return &AggregationPaginatedResult{Paginator: *paginator, Data: r.Data}, nil
}

func (c *client) validateDocumentKind(obj interface{}) error {
	if reflect.TypeOf(obj).Kind() != reflect.Ptr {
		return errors.New("asari: doc must be a pointer to a document")
	}
	return nil
}

func (c *client) validateProjection(projection interface{}) error {
	if projection != nil {
		if _, ok := projection.(bson.M); !ok {
			return errors.New("asari: projections can only be bson.M types")
		}
	}
	return nil
}

func (c *client) validateFilters(filters []bson.E) error {
	for _, f := range filters {
		if f.Key == "" {
			return errors.New("asari: document field names in filters cannot be empty. Key required")
		}
	}
	return nil
}
