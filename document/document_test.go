package document

import (
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
	"time"
)

type testDoc struct {
	Base `bson:",inline"`
}

func TestBase_Setup(t *testing.T) {
	b := testDoc{}
	assert.Equal(t, primitive.NilObjectID, b.ID)
	assert.True(t, b.CreatedAt.IsZero())
	assert.True(t, b.UpdatedAt.IsZero())
	assert.Equal(t, false, b.isNew)

	err := b.Setup()
	assert.Nil(t, err)
	assert.NotEqual(t, primitive.NilObjectID, b.ID)
	assert.False(t, b.CreatedAt.IsZero())
	assert.False(t, b.UpdatedAt.IsZero())
	assert.Equal(t, true, b.isNew)

	err = b.Setup()
	assert.NotNil(t, err)
}

func TestBase_CanSave(t *testing.T) {
	b := testDoc{}
	assert.False(t, b.CanSave())

	b.Setup()
	assert.True(t, b.CanSave())
}

func TestBase_GetID(t *testing.T) {
	b := testDoc{}
	b.Setup()
	assert.NotEqual(t, primitive.NilObjectID, b.GetID())
}

func TestBase_SetIsNew(t *testing.T) {
	b := testDoc{}
	assert.False(t, b.IsNew())

	b.SetIsNew(true)
	assert.True(t, b.IsNew())
}

func TestBase_BeforeSoftDelete(t *testing.T) {
	b := testDoc{}
	assert.False(t, b.IsDeleted)
	assert.True(t, b.DeletedAt.IsZero())

	b.BeforeSoftDelete()
	assert.True(t, b.IsDeleted)
	assert.False(t, b.DeletedAt.IsZero())
}

func TestBase_GetCreatedAt(t *testing.T) {
	b := testDoc{}
	assert.True(t, b.GetCreatedAt().IsZero())

	b.Setup()
	assert.False(t, b.GetCreatedAt().IsZero())
}

func TestBase_GetUpdatedAt(t *testing.T) {
	b := testDoc{}
	assert.True(t, b.GetUpdatedAt().IsZero())

	b.Setup()
	assert.False(t, b.GetUpdatedAt().IsZero())
}

func TestBase_BeforeUpdate(t *testing.T) {
	b := testDoc{}
	b.Setup()
	oldTime := b.GetUpdatedAt()

	//Sleep for 1 Second
	d := 1 * time.Second
	time.Sleep(d)

	b.BeforeUpdate()
	newTime := b.GetUpdatedAt()
	assert.True(t, newTime.After(oldTime))
}

func TestBase_FormatDate(t *testing.T) {
	b := testDoc{}
	b.Setup()

	assert.Equal(t, b.FormatDate(b.CreatedAt, "Jan 02, 2006"), b.GetCreatedAt().Format("Jan 02, 2006"))
}

func TestBase_FormatDateShort(t *testing.T) {
	b := testDoc{}
	b.Setup()

	assert.Equal(t, b.FormatDateShort(b.CreatedAt), b.GetCreatedAt().Format("Jan 02, 2006"))
}

func TestBase_FormatDateTimeShort(t *testing.T) {
	b := testDoc{}
	b.Setup()

	assert.Equal(t, b.FormatDateTimeShort(b.CreatedAt), b.GetCreatedAt().Format("Jan 02, 2006 - 15:04"))
}

func TestBase_GetFormattedCreatedAt(t *testing.T) {
	b := testDoc{}

	assert.Nil(t, b.GetFormattedCreatedAt())

	b.Setup()

	tms := b.GetFormattedCreatedAt()
	assert.Equal(t, tms.DateShort, b.GetCreatedAt().Format("Jan 02, 2006"))
	assert.Equal(t, tms.DateTimeShort, b.GetCreatedAt().Format("Jan 02, 2006 - 15:04"))
	assert.Equal(t, tms.ISO, b.GetCreatedAt().Format(time.RFC3339))
}

func TestBase_GetFormattedUpdatedAt(t *testing.T) {
	b := testDoc{}

	assert.Nil(t, b.GetFormattedUpdatedAt())

	b.Setup()

	tms := b.GetFormattedUpdatedAt()
	assert.Equal(t, tms.DateShort, b.GetUpdatedAt().Format("Jan 02, 2006"))
	assert.Equal(t, tms.DateTimeShort, b.GetUpdatedAt().Format("Jan 02, 2006 - 15:04"))
	assert.Equal(t, tms.ISO, b.GetUpdatedAt().Format(time.RFC3339))
}

func TestBase_GetFormattedDeletedAt(t *testing.T) {
	b := testDoc{}

	assert.Nil(t, b.GetFormattedDeletedAt())

	b.Setup()
	b.BeforeSoftDelete()

	tms := b.GetFormattedDeletedAt()
	assert.Equal(t, tms.DateShort, b.DeletedAt.Format("Jan 02, 2006"))
	assert.Equal(t, tms.DateTimeShort, b.DeletedAt.Format("Jan 02, 2006 - 15:04"))
	assert.Equal(t, tms.ISO, b.DeletedAt.Format(time.RFC3339))
}

func TestBase_GetAllTimeStamps(t *testing.T) {
	b := testDoc{}
	b.Setup()
	b.BeforeSoftDelete()

	tms := b.GetAllTimeStamps()
	assert.Equal(t, 3, len(tms))
	assert.NotNil(t, tms["createdAt"])
	assert.NotNil(t, tms["updatedAt"])
	assert.NotNil(t, tms["deletedAt"])
}
