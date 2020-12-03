package document

import "go.mongodb.org/mongo-driver/mongo"

type (
	// PreCreator Interface
	PreCreator interface {
		// PreCreate runs the concrete implementation before a new document is saved.
		PreCreate(dbConnection *mongo.Database) error
	}

	// PostCreator
	PostCreator interface {
		// PostCreate runs the concrete implementation after a new document is saved successfully.
		PostCreate(dbConnection *mongo.Database) error
	}

	// PreUpdater
	PreUpdater interface {
		// PreUpdate runs the concrete implementation before an existing document updated is saved.
		PreUpdate(dbConnection *mongo.Database) error
	}

	// PostUpdater
	PostUpdater interface {
		// PostUpdate runs the concrete implementation after an existing document is updated successfully.
		PostUpdate(dbConnection *mongo.Database) error
	}

	// PreSoftDeleter
	PreSoftDeleter interface {
		// PreSoftDelete runs the concrete implementation before a document is soft deleted.
		PreSoftDelete(dbConnection *mongo.Database) error
	}

	// PostSoftDeleter
	PostSoftDeleter interface {
		// PostSoftDelete runs the concrete implementation after a document is soft deleted.
		PostSoftDelete(dbConnection *mongo.Database) error
	}

	// PreHardDeleter
	PreHardDeleter interface {
		// PreHardDelete runs the concrete implementation before a document is hard deleted.
		PreHardDelete(dbConnection *mongo.Database) error
	}

	// PostHardDeleter
	PostHardDeleter interface {
		// PostHardDelete runs the concrete implementation after a document is hard deleted.
		PostHardDelete(dbConnection *mongo.Database) error
	}

	// PreFindOne
	PreFindOne interface {
		//PreFindOne runs the concrete implementation before any FindOne*() is called
		PreFindOne(dbConnection *mongo.Database) error
	}

	// PostFindOne - Runs after a document is successfully found and marshalled into a target
	PostFindOne interface {
		//PostFindOne runs the concrete implementation after any FindOne*() is called
		PostFindOne(dbConnection *mongo.Database) error
	}
)
