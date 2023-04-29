package std

// Error represents an error. The Error type holds a textual description of the
// error and contextual information describing the source of the error. It also
// holds a copy of its parent error, if one exists.
//
// Contextual information held in an Error can be used to trace back the source
// of an error, and can be returned to the user through functions such as
// errors.Trace
type Error interface {
	// Program counter of the error origin.
	Address() uintptr
	// Textual error description.
	// Includes textual error descriptions of all ancestor errors.
	Error() string
	// Error source filename.
	File() string
	// Name of the function in which the error occurred.
	Func() string
	// Update the error's source information.
	Here()
	// Offending line number in error source file.
	Line() int
	// Parent error which caused the current error.
	Parent() Error
	// Short textual error description.
	// Does not include textual error descriptions of all ancestor errors.
	Text() string
}
