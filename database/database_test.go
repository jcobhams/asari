package database

import (
	"github.com/jcobhams/asari/builder"
	"github.com/jcobhams/asari/document"
	"github.com/jcobhams/asari/operator"
	"github.com/jcobhams/asari/queryfilter"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"testing"
)

var TestClient *client

const UserCollection string = "users"

type User struct {
	document.Base `bson:",inline"`
	FirstName     string `bson:"first_name"`
	LastName      string `bson:"last_name"`
	Email         string `bson:"email"`
	Level         int    `bson:"level"`
}

func TestMain(m *testing.M) {
	mongoDSN, ok := os.LookupEnv("MONGO_DSN")
	if !ok {
		panic("MONGO_DSN Environment Variable Required")
	}

	databaseName, ok := os.LookupEnv("DATABASE_NAME")
	if !ok {
		panic("DATABASE_NAME Environment Variable Required")
	}

	TestClient = Init(mongoDSN, databaseName)

	code := m.Run()
	os.Exit(code)
}

func TestInit(t *testing.T) {
	mongoDSN, _ := os.LookupEnv("MONGO_DSN")
	databaseName, _ := os.LookupEnv("DATABASE_NAME")

	var c *client
	if assert.NotPanics(t, func() { c = Init(mongoDSN, databaseName) }) {
		assert.NotNil(t, c)
		assert.Equal(t, databaseName, c.Connection.Name())
	}

	assert.Panics(t, func() { Init("", "") })
}

func TestClient_FindOneInternal(t *testing.T) {
	user := &User{
		FirstName: "Joseph",
		LastName:  "Cobhams",
		Level:     1,
	}
	user.Setup()
	TestClient.Connection.Collection(UserCollection).InsertOne(nil, user)

	var u User
	qf := queryfilter.New().AddFilter(bson.E{Key: "first_name", Value: "Joseph"})
	err := TestClient.findOne(nil, UserCollection, qf.GetFilters(), &u)
	assert.Nil(t, err)
	assert.Equal(t, user.FirstName, u.FirstName)
	assert.Equal(t, user.LastName, u.LastName)
	assert.Equal(t, user.Level, u.Level)

	//Test Error is returned if doc is not a pointer
	assert.Error(t, TestClient.findOne(nil, UserCollection, qf.GetFilters(), u))

	//Test Error is returned if queryfilter contains empty field names
	qf2 := qf
	qf2f := append(qf2.GetFilters(), bson.E{Key: "", Value: "some"})
	assert.Error(t, TestClient.findOne(nil, UserCollection, qf2f, &u, nil))

	//Test Invalid Projections
	opts := &options.FindOneOptions{
		Projection: map[string]interface{}{"test": 1},
	}
	assert.Error(t, TestClient.findOne(nil, UserCollection, qf2f, &u, opts))

	//Test Projection
	var u2 User

	opts = &options.FindOneOptions{
		Projection: bson.M{"last_name": 1},
	}
	err = TestClient.findOne(nil, UserCollection, qf.GetFilters(), &u2, opts)
	assert.Nil(t, err)
	assert.Equal(t, user.LastName, u2.LastName)
	assert.NotEqual(t, user.FirstName, u2.FirstName)

	tearDown()
}

func TestClient_FindOne(t *testing.T) {
	user := &User{
		FirstName: "Joseph",
		LastName:  "Cobhams",
		Level:     1,
	}
	user.Setup()
	TestClient.Connection.Collection(UserCollection).InsertOne(nil, user)

	var u User
	qf := queryfilter.New().AddFilter(bson.E{Key: "first_name", Value: "Joseph"}).GetFilters()
	err := TestClient.FindOne(nil, UserCollection, qf, nil, &u)

	assert.Nil(t, err)
	assert.Equal(t, user.FirstName, u.FirstName)
	assert.Equal(t, user.LastName, u.LastName)
	assert.Equal(t, user.Level, u.Level)

	tearDown()
}

func TestClient_FindOneByID(t *testing.T) {
	user := &User{
		FirstName: "Joseph",
		LastName:  "Cobhams",
		Level:     1,
	}
	user.Setup()
	TestClient.Connection.Collection(UserCollection).InsertOne(nil, user)

	var u User
	err := TestClient.FindOneByID(nil, UserCollection, user.ID, nil, &u)

	assert.Nil(t, err)
	assert.Equal(t, user.FirstName, u.FirstName)
	assert.Equal(t, user.LastName, u.LastName)
	assert.Equal(t, user.Level, u.Level)

	tearDown()
}

func TestClient_FindOneByField(t *testing.T) {
	user := &User{
		FirstName: "Joseph",
		LastName:  "Cobhams",
		Level:     1,
	}
	user.Setup()
	TestClient.Connection.Collection(UserCollection).InsertOne(nil, user)

	var u User
	err := TestClient.FindOneByField(nil, UserCollection, "last_name", user.LastName, nil, &u)

	assert.Nil(t, err)
	assert.Equal(t, user.FirstName, u.FirstName)
	assert.Equal(t, user.LastName, u.LastName)
	assert.Equal(t, user.Level, u.Level)

	tearDown()
}

func TestClient_FindPaginated(t *testing.T) {
	user1 := &User{
		FirstName: "Joseph",
		LastName:  "Cobhams",
		Level:     1,
	}
	user1.Setup()
	TestClient.Connection.Collection(UserCollection).InsertOne(nil, user1)

	user2 := &User{
		FirstName: "Asari",
		LastName:  "Cobhams",
		Level:     2,
	}
	user2.Setup()
	TestClient.Connection.Collection(UserCollection).InsertOne(nil, user2)

	user3 := &User{
		FirstName: "Ivy",
		LastName:  "Cobhams",
		Level:     3,
	}
	user3.Setup()
	TestClient.Connection.Collection(UserCollection).InsertOne(nil, user3)

	pageOpts := PageOpts{
		Page:    1,
		PerPage: 1,
	}

	qf := queryfilter.New().GetFilters()
	users, err := TestClient.FindPaginated(nil, UserCollection, pageOpts, qf, nil, nil)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), users.Paginator.CurrentPage)
	assert.Equal(t, int64(2), users.Paginator.NextPage)
	assert.Equal(t, int64(3), users.Paginator.TotalPages)
	assert.Equal(t, int64(3), users.Paginator.TotalRows)

	for users.Cursor.Next(nil) {
		var u User
		assert.Nil(t, users.Cursor.Decode(&u))
		assert.Equal(t, "Ivy", u.FirstName)
	}
	users.Cursor.Close(nil)

	//Load Page 2
	pageOpts.Page = 2
	users, err = TestClient.FindPaginated(nil, UserCollection, pageOpts, qf, nil, nil)
	assert.Nil(t, err)
	assert.Equal(t, int64(2), users.Paginator.CurrentPage)
	assert.Equal(t, int64(3), users.Paginator.NextPage)
	assert.Equal(t, int64(3), users.Paginator.TotalPages)
	assert.Equal(t, int64(3), users.Paginator.TotalRows)

	for users.Cursor.Next(nil) {
		var u User
		assert.Nil(t, users.Cursor.Decode(&u))
		assert.Equal(t, "Asari", u.FirstName)
	}
	users.Cursor.Close(nil)

	//bson.Unmarshal(users.Results[0], &u)
	//assert.Equal(t, "Asari", u.FirstName)

	tearDown()
}

func TestClient_FindLast(t *testing.T) {
	user1 := &User{
		FirstName: "Joseph",
		LastName:  "Cobhams",
		Level:     1,
	}
	user1.Setup()
	TestClient.Connection.Collection(UserCollection).InsertOne(nil, user1)

	user2 := &User{
		FirstName: "Asari",
		LastName:  "Cobhams",
		Level:     2,
	}
	user2.Setup()
	TestClient.Connection.Collection(UserCollection).InsertOne(nil, user2)

	var u User
	qf := queryfilter.New().AddFilter(bson.E{Key: "last_name", Value: "Cobhams"}).GetFilters()
	err := TestClient.FindLast(nil, UserCollection, qf, nil, &u)

	assert.Nil(t, err)
	assert.Equal(t, user2.FirstName, u.FirstName)
	assert.Equal(t, user2.LastName, u.LastName)
	assert.Equal(t, user2.Level, u.Level)

	tearDown()
}

func TestClient_FindLastN(t *testing.T) {
	user1 := &User{
		FirstName: "Joseph",
		LastName:  "Cobhams",
		Level:     1,
	}
	user1.Setup()
	TestClient.Connection.Collection(UserCollection).InsertOne(nil, user1)

	user2 := &User{
		FirstName: "Asari",
		LastName:  "Cobhams",
		Level:     2,
	}
	user2.Setup()
	TestClient.Connection.Collection(UserCollection).InsertOne(nil, user2)

	user3 := &User{
		FirstName: "Ivy",
		LastName:  "Cobhams",
		Level:     3,
	}
	user3.Setup()
	TestClient.Connection.Collection(UserCollection).InsertOne(nil, user3)

	qf := queryfilter.New().AddFilter(bson.E{Key: "last_name", Value: "Cobhams"}).GetFilters()
	cur, err := TestClient.FindLastN(nil, UserCollection, 2, qf, nil)

	users := []User{}
	for cur.Next(nil) {
		var u User
		cur.Decode(&u)
		users = append(users, u)
	}

	assert.Nil(t, err)
	assert.Equal(t, users[0].FirstName, user3.FirstName)
	assert.Equal(t, users[1].FirstName, user2.FirstName)
	assert.Equal(t, 2, len(users))

	tearDown()
}

func TestClient_FindAll(t *testing.T) {
	user1 := &User{
		FirstName: "Joseph",
		LastName:  "Cobhams",
		Level:     1,
	}
	user1.Setup()
	TestClient.Connection.Collection(UserCollection).InsertOne(nil, user1)

	user2 := &User{
		FirstName: "Asari",
		LastName:  "Cobhams",
		Level:     2,
	}
	user2.Setup()
	TestClient.Connection.Collection(UserCollection).InsertOne(nil, user2)

	user3 := &User{
		FirstName: "Ivy",
		LastName:  "Cobhams",
		Level:     3,
	}
	user3.Setup()
	TestClient.Connection.Collection(UserCollection).InsertOne(nil, user3)

	qf := queryfilter.New().GetFilters()
	cur, err := TestClient.FindAll(nil, UserCollection, qf, nil, nil)

	users := []User{}
	for cur.Next(nil) {
		var u User
		cur.Decode(&u)
		users = append(users, u)
	}

	assert.Nil(t, err)
	assert.Equal(t, users[0].FirstName, user3.FirstName)
	assert.Equal(t, users[1].FirstName, user2.FirstName)
	assert.Equal(t, users[2].Level, user1.Level)
	assert.Equal(t, 3, len(users))

	tearDown()
}

func TestClient_SaveDocument(t *testing.T) {
	user := &User{
		FirstName: "Joseph",
		LastName:  "Cobhams",
		Level:     1,
	}

	//Test error returned if document is not a pointer
	_, err := TestClient.SaveDocument(nil, UserCollection, *user)
	assert.Error(t, err)

	//Test error is returned if document has never been setup
	_, err = TestClient.SaveDocument(nil, UserCollection, user)
	assert.Error(t, err)

	user.Setup()
	assert.True(t, user.IsNew())

	//Test Saving New Document
	_, err = TestClient.SaveDocument(nil, UserCollection, user)
	assert.Nil(t, err)
	assert.False(t, user.IsNew())

	var u User
	TestClient.FindOneByID(nil, UserCollection, user.ID, nil, &u)
	assert.Equal(t, user.FirstName, u.FirstName)

	//Test Updating Existing Document
	user.FirstName = "Asari"
	_, err = TestClient.SaveDocument(nil, UserCollection, user)

	TestClient.FindOneByID(nil, UserCollection, user.ID, nil, &u)
	assert.Equal(t, user.FirstName, u.FirstName)
	assert.NotEqual(t, "Joseph", u.FirstName)

	tearDown()
}

func TestClient_SoftDeleteDocument(t *testing.T) {
	user := &User{
		FirstName: "Joseph",
		LastName:  "Cobhams",
		Level:     1,
	}
	user.Setup()
	TestClient.SaveDocument(nil, UserCollection, user)

	assert.False(t, user.IsDeleted)
	assert.True(t, user.DeletedAt.IsZero())

	_, err := TestClient.SoftDeleteDocument(nil, UserCollection, user)
	assert.Nil(t, err)

	var u User
	err = TestClient.FindOneByID(nil, UserCollection, user.ID, nil, &u)
	assert.Equal(t, mongo.ErrNoDocuments, err)

	//Manually Set is_deleted value to true to get document that has been marked as deleted
	qf := queryfilter.NewWithDeleted().
		AddFilter(bson.E{Key: "_id", Value: user.ID}).
		GetFilters()

	err = TestClient.FindOne(nil, UserCollection, qf, nil, &u)
	assert.Nil(t, err)
	assert.True(t, u.IsDeleted)
	assert.False(t, u.DeletedAt.IsZero())
	assert.Equal(t, user.FirstName, u.FirstName)

	tearDown()
}

func TestClient_CountDocuments(t *testing.T) {
	user := &User{
		FirstName: "Joseph",
		LastName:  "Cobhams",
		Level:     1,
	}
	user.Setup()
	TestClient.SaveDocument(nil, UserCollection, user)

	user2 := &User{
		FirstName: "Asari",
		LastName:  "Cobhams",
		Level:     1,
	}
	user2.Setup()
	TestClient.SaveDocument(nil, UserCollection, user2)

	qf := queryfilter.New().GetFilters()
	count, err := TestClient.CountDocuments(nil, UserCollection, qf)
	assert.Nil(t, err)
	assert.Equal(t, 2, count)

	tearDown()
}

func TestClient_HardDeleteDocument(t *testing.T) {
	user := &User{
		FirstName: "Joseph",
		LastName:  "Cobhams",
		Level:     1,
	}
	user.Setup()
	TestClient.SaveDocument(nil, UserCollection, user)

	_, err := TestClient.HardDeleteDocument(nil, UserCollection, user)
	assert.Nil(t, err)

	qf := queryfilter.New().AddFilter(bson.E{Key: "_id", Value: user.ID}).GetFilters()
	count, _ := TestClient.CountDocuments(nil, UserCollection, qf)
	assert.Equal(t, 0, count)

	qf = queryfilter.NewWithDeleted().AddFilter(bson.E{Key: "_id", Value: user.ID}).GetFilters()
	count, _ = TestClient.CountDocuments(nil, UserCollection, qf)
	assert.Equal(t, 0, count)

	tearDown()
}

func TestClient_UpdateMany(t *testing.T) {
	user1 := &User{
		FirstName: "Joseph",
		LastName:  "Dahryl",
	}
	user1.Setup()
	TestClient.Connection.Collection(UserCollection).InsertOne(nil, user1)

	user2 := &User{
		FirstName: "Asari",
		LastName:  "Dahryl",
	}
	user2.Setup()
	TestClient.Connection.Collection(UserCollection).InsertOne(nil, user2)

	user3 := &User{
		FirstName: "Ivy",
		LastName:  "Ariella",
	}
	user3.Setup()
	TestClient.Connection.Collection(UserCollection).InsertOne(nil, user3)

	qf := queryfilter.New().AddFilter(bson.E{Key: "last_name", Value: "Dahryl"}).GetFilters()

	ub := builder.NewUpdateManyBuilder().
		Add(operator.Set, bson.E{Key: "last_name", Value: "dahryl"}, bson.E{Key: "level", Value: 2})

	res, err := TestClient.UpdateMany(nil, UserCollection, qf, ub, nil)
	assert.Nil(t, err)
	assert.Equal(t, int64(2), res.MatchedCount)
	assert.Equal(t, int64(2), res.ModifiedCount)

	count, _ := TestClient.CountDocuments(nil, UserCollection, queryfilter.New().AddFilter(bson.E{Key: "last_name", Value: "dahryl"}).GetFilters())
	assert.Equal(t, 2, count)

	tearDown()
}

func tearDown() {
	TestClient.Connection.Collection(UserCollection).DeleteMany(nil, []bson.E{})
}
