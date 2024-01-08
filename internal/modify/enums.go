package modify

//go:generate go-enum -f=$GOFILE -t ../enums.tmpl --marshal --names --nocomments

// Determines what type of modification to perform.
//
// The append/prepend behavior differs slightly depending on
// the destination content type. Strings are concatenated,
// numbers are added, and objects are recursively merged.
// Arrays are concatenated by default, but that behavior can
// be customized via the 'merge' enum.
//
// Replace and delete behave consistently across all types.
/*
	ENUM(
		append   // Append to the destination content.
		prepend  // Prepend to the destination content.
		replace  // Replace the destination.
		delete   // Delete the destination content.
	).
*/
type Action string

// Determines merge behavior for arrays - either when modifying them directly
// or when recursively merging objects containing arrays.
/*
	ENUM(
		concat   // Concatenate source and destination arrays.
		upsert   // Add source array items if not present in the destination.
		replace  // Replace the destination with the source.
	).
*/
type MergeType string
