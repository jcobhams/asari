package document

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type (
	//Document is the base interface for all Documents. This helps all document to be initialized using a single Function.
	Document interface {
		Setup() error
		CanSave() bool
		BeforeUpdate()
		BeforeSoftDelete()
		GetID() primitive.ObjectID
		GetCreatedAt() time.Time
		GetUpdatedAt() time.Time
		SetIsNew(status bool)
		IsNew() bool
		FormatDate(d time.Time, layout string) string
		FormatDateShort(d time.Time) string
		FormatDateTimeShort(d time.Time) string
		GetTimeStamps() map[string]formattedTimestamp
		GetAllTimeStamps() map[string]formattedTimestamp
		GetFormattedCreatedAt() *formattedTimestamp
		GetFormattedUpdatedAt() *formattedTimestamp
		GetFormattedDeletedAt() *formattedTimestamp
	}

	//Base is the base document all documents must inherit. This ensure shared document properties can be set.
	Base struct {
		ID         primitive.ObjectID            `bson:"_id,omitempty" json:"_id"`
		CreatedAt  time.Time                     `bson:"created_at" json:"-"`
		UpdatedAt  time.Time                     `bson:"updated_at" json:"-"`
		DeletedAt  time.Time                     `bson:"deleted_at,omitempty" json:"-"`
		IsDeleted  bool                          `bson:"is_deleted" json:"-"`
		isNew      bool                          `json:"-" bson:"-"`
		Timestamps map[string]formattedTimestamp `bson:"-" json:"timestamps"`
	}

	formattedTimestamp struct {
		DateShort     string `json:"dateShort"`
		DateTimeShort string `json:"dateTimeShort"`
		ISO           string `json:"iso"`
	}
)

//Setup initializes a document with ID and Timestamps
func (d *Base) Setup() error {
	if d.ID == primitive.NilObjectID {
		d.ID = primitive.NewObjectID()
		d.CreatedAt = time.Now().UTC()
		d.UpdatedAt = time.Now().UTC()
		d.isNew = true
		return nil
	}
	return errors.New("cannot setup an already existing document. - Setup() only applies to new document")
}

func (d *Base) CanSave() bool {
	if d.ID == primitive.NilObjectID || d.CreatedAt.IsZero() || d.UpdatedAt.IsZero() {
		return false
	}
	return true
}

//GetID returns a document's ID
func (d *Base) GetID() primitive.ObjectID {
	return d.ID
}

//GetCreatedAt returns a document's created time
func (d *Base) GetCreatedAt() time.Time {
	return d.CreatedAt
}

//GetUpdatedAt returns a document's last update time
func (d *Base) GetUpdatedAt() time.Time {
	return d.UpdatedAt
}

//IsNew returns a document's initialization state.
func (d *Base) IsNew() bool {
	return d.isNew
}

//SetIsNew sets the isNew property to provided status.
//If status is true, This signals that the current instance has not been saved before and document.Save() will perform an insert.
//If status is false, This makes all calls to document.Save() perform an update operation not an insert.
func (d *Base) SetIsNew(status bool) {
	d.isNew = status
}

func (d *Base) BeforeUpdate() {
	d.UpdatedAt = time.Now().UTC()
}

func (d *Base) BeforeSoftDelete() {
	d.IsDeleted = true
	d.DeletedAt = time.Now().UTC()
}

//FormatDateShort returns a formatted time object in the format MMM DD, YYYY
func (d *Base) FormatDateShort(dt time.Time) string {
	return dt.Format("Jan 02, 2006")
}

//FormatDateTimeShort returns a formatted time object in the format MMM DD, YYYY - HH:MM
func (d *Base) FormatDateTimeShort(dt time.Time) string {
	return dt.Format("Jan 02, 2006 - 15:04")
}

//FormatDate returns a formatted time object based on the provided layout.
func (d *Base) FormatDate(dt time.Time, layout string) string {
	return dt.Format(layout)
}

func (d *Base) GetFormattedCreatedAt() *formattedTimestamp {
	if !d.CreatedAt.IsZero() {
		return &formattedTimestamp{
			DateShort:     d.FormatDateShort(d.CreatedAt),
			DateTimeShort: d.FormatDateTimeShort(d.CreatedAt),
			ISO:           d.CreatedAt.Format(time.RFC3339),
		}
	}
	return nil
}

func (d *Base) GetFormattedUpdatedAt() *formattedTimestamp {
	if !d.UpdatedAt.IsZero() {
		return &formattedTimestamp{
			DateShort:     d.FormatDateShort(d.UpdatedAt),
			DateTimeShort: d.FormatDateTimeShort(d.UpdatedAt),
			ISO:           d.UpdatedAt.Format(time.RFC3339),
		}
	}
	return nil
}

func (d *Base) GetFormattedDeletedAt() *formattedTimestamp {
	if !d.DeletedAt.IsZero() {
		return &formattedTimestamp{
			DateShort:     d.FormatDateShort(d.DeletedAt),
			DateTimeShort: d.FormatDateTimeShort(d.DeletedAt),
			ISO:           d.DeletedAt.Format(time.RFC3339),
		}
	}
	return nil
}

func (d *Base) GetTimeStamps() map[string]formattedTimestamp {
	t := map[string]formattedTimestamp{"createdAt": *d.GetFormattedCreatedAt()}
	if updatedAt := d.GetFormattedUpdatedAt(); updatedAt != nil {
		t["updatedAt"] = *d.GetFormattedUpdatedAt()
	}
	return t
}

func (d *Base) GetAllTimeStamps() map[string]formattedTimestamp {
	t := d.GetTimeStamps()
	if deletedAt := d.GetFormattedDeletedAt(); deletedAt != nil {
		t["deletedAt"] = *d.GetFormattedDeletedAt()
	}
	return t
}
