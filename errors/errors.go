// Package errors implements functions to manipulate errors.
//
// The New function creates errors whose contents consist of a textual
// description of the error and contextual information. This contextual
// information can be used to trace back the source of the error. The Trace
// function can be used to write traceback information to an io.Writer in a
// human-friendly format.
//
// Errors are frequently caused by other errors. To account for this, the New
// function takes two parameters: a textual description of the error and its
// parent error.
//
// Package developers may choose to predefine certain errors, such as the io.EOF
// error:
//
//	var EOF = errors.New("EOF")
//
// However, the context of this error is the line on which it is declared. To
// update the contextual information, the Here method can be invoked when
// intending to raise a predefined error from a current context, like so:
//
//	io.EOF.Here()
//
// A function which handles errors from multiple concurrently executing
// processes should, in most cases, return the first error it receives:
//
//	func f() error {
//		errs := make(chan error)
//		for i := 0; i < 5; i++ {
//			go g(errs)
//		}
//		for err := range errs {
//			if err != nil {
//				return err
//			}
//		}
//		return nil
//	}
//
// Alternatively, if the developer's intent is to return all received errors,
// the Join function is provided to return all errors as one.
package errors

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
)

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
	// Return a copy of the error with its context set to the current context.
	Here() Error
	// Offending line number in error source file.
	Line() int
	// Parent error which caused the current error.
	Parent() Error
	// Short textual error description.
	// Does not include textual error descriptions of all ancestor errors.
	Text() string
}

type errtype struct {
	address uintptr
	file    string
	fn      string
	line    int
	parent  Error
	text    string
}

func (e errtype) Address() uintptr {
	return e.address
}

func (e errtype) Error() string {
	var text string
	var err Error
	for err = e; err.Parent() != nil; err = err.Parent() {
		text += err.Text() + ": "
	}
	text += err.Text()
	return text
}

func (e errtype) File() string {
	return e.file
}

func (e errtype) Func() string {
	return e.fn
}

func (e errtype) Here() Error {
	err := e
	addr, file, line, _ := runtime.Caller(1)
	f := runtime.FuncForPC(addr)
	fn := f.Name()
	err.address = addr
	err.file = file
	err.fn = fn
	err.line = line
	return err
}

func (e errtype) Line() int {
	return e.line
}

func (e errtype) Parent() Error {
	return e.parent
}

func (e errtype) Text() string {
	return e.text
}

// Has reports whether the textual error description of err or any of its parent
// errors match the textual error description of target.
func Has(err, target Error) bool {
	e := err
	for ; e != nil; e = e.Parent() {
		if e.Text() == target.Text() {
			return true
		}
	}
	return false
}

// Is reports whether the textual error description of err matches the textual
// error description of target.
func Is(err, target Error) bool {
	return err.Text() == target.Text()
}

// New returns an error whose textual error description is given by text, and 
// whose parent error is err. If the new error has no parent, err should be
// given as nil.
//
// The current filename, line, program counter, and parent function name are
// stored within the error interface. Each call to New returns a distinct error
// value even if text is identical.
func New(text string, err Error) Error {
	addr, file, line, _ := runtime.Caller(1)
	f := runtime.FuncForPC(addr)
	fn := f.Name()
	return errtype{
		address: addr,
		file:    file,
		fn:      fn,
		line:    line,
		parent:  err,
		text:    text,
	}
}

// Join returns a Error that combines the given errs. Any nil error values
// are discarded. Join returns nil if errs contains no non-nil values. The
// resultant error is formatted as a concatenation of the textual error
// descriptions of all given errs, with a comma and space between each
// description.
//
// An error can only have one parent, so the resultant error has nil parent.
func Join(errs ...Error) Error {
	var text string
	for _, err := range errs {
		if err != nil {
			text += err.Text() + ", "
		}
	}
	strings.TrimSuffix(text, ", ")
	addr, file, line, _ := runtime.Caller(1)
	f := runtime.FuncForPC(addr)
	fn := f.Name()
	return errtype{
		address: addr,
		file:    file,
		fn:      fn,
		line:    line,
		parent:  nil,
		text:    text,
	}
}

// Trace writes human-friendly error traceback information from err to w. If w
// is nil, Trace writes to the standard error stream.
func Trace(w io.Writer, err Error) {
	if w == nil {
		w = os.Stderr
	}
	fmt.Fprintf(w, "\nerror: %s\n", e.Error())
	defer fmt.Fprint(w, "\n")
	e := err
	for ; e != nil; e = e.Parent() {
		defer fmt.Fprintf(w, "\t%s:%d\n", e.File(), e.Line())
		defer fmt.Fprintf(w, "%s(...)\n", e.Func())
	}
}
