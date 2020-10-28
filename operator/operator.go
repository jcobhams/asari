package operator

// Credits: https://github.com/Kamva/mgm/blob/master/operator

const (
	AddFields      = "$addFields"
	Bucket         = "$bucket"
	BucketAuto     = "$bucketAuto"
	CollStats      = "$collStats"
	Count          = "$count"
	Facet          = "$facet"
	GeoNear        = "$geoNear"
	GraphLookup    = "$graphLookup"
	Group          = "$group"
	IndexStats     = "$indexStats"
	Limit          = "$limit"
	ListSessions   = "$listSessions"
	Lookup         = "$lookup"
	Match          = "$match"
	Merge          = "$merge"
	Out            = "$out"
	PlanCacheStats = "$planCacheStats"
	Project        = "$project"
	Redact         = "$redact"
	ReplaceRoot    = "$replaceRoot"
	ReplaceWith    = "$replaceWith"
	Sample         = "$sample"
	Set            = "$set"
	Skip           = "$skip"
	Sort           = "$sort"
	SortByCount    = "$sortByCount"
	Unset          = "$unset"
	Mul            = "$mul"
	Unwind         = "$unwind"

	Avg        = "$avg"
	Max        = "$max"
	Min        = "$min"
	StdDevPop  = "$stdDevPop"
	StdDevSamp = "$stdDevSamp"
	Sum        = "$sum"

	AddToSet = "$addToSet"
	Pop      = "$pop"
	Pull     = "$pull"
	Push     = "$push"
	PullAll  = "$pullAll"

	// Comparison
	Eq  = "$eq"
	Gt  = "$gt"
	Gte = "$gte"
	In  = "$in"
	Lt  = "$lt"
	Lte = "$lte"
	Ne  = "$ne"
	Nin = "$nin"

	// Logical
	And = "$and"
	Not = "$not"
	Nor = "$nor"
	Or  = "$or"

	// Element
	Exists = "$exists"
	Type   = "$type"

	// Evaluation
	Expr       = "$expr"
	JSONSchema = "$jsonSchema"
	Mod        = "$mod"
	Regex      = "$regex"
	Text       = "$text"
	Where      = "$where"

	// Geo spatial
	GeoIntersects = "$geoIntersects"
	GeoWithin     = "$geoWithin"
	Near          = "$near"
	NearSphere    = "$nearSphere"

	// Array
	All       = "$all"
	ElemMatch = "$elemMatch"
	Size      = "$size"

	// Bitwise
	BitsAllClear = "$bitsAllClear"
	BitsAllSet   = "$bitsAllSet"
	BitsAnyClear = "$bitsAnyClear"
	BitsAnySet   = "$bitsAnySet"

	// Comments
	Comment = "$comment"

	// Projection operators

	Dollar = "$"
	Meta   = "$meta"
	Slice  = "$slice"

	//LookUpKeys
	LookupFrom         = "from"
	LookupLocalField   = "localField"
	LookupForeignField = "foreignField"
	LookupAs           = "as"
	LookupLet          = "let"
	LookupPipeline     = "pipeline"
)
