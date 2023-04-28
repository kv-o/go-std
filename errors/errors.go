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

	"codeberg.org/kvo/std"
)

type errtype struct {
	address uintptr
	file    string
	fn      string
	line    int
	parent  std.Error
	text    string
}

func (e errtype) Address() uintptr {
	return e.address
}

func (e errtype) Error() string {
	var text string
	var err std.Error
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

func (e errtype) Here() {
	addr, file, line, _ := runtime.Caller(1)
	f := runtime.FuncForPC(addr)
	fn := f.Name()
	e.address = addr
	e.file = file
	e.fn = fn
	e.line = line
}

func (e errtype) Line() int {
	return e.line
}

func (e errtype) Parent() std.Error {
	return e.parent
}

func (e errtype) Text() string {
	return e.text
}

// New returns an error whose textual error description is given by text, and 
// whose parent error is err. If the new error has no parent, err should be
// given as nil.
//
// The current filename, line, program counter, and parent function name are
// stored within the error interface. Each call to New returns a distinct error
// value even if text is identical.
func New(text string, err std.Error) std.Error {
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

// Join returns a std.Error that combines the given errs. Any nil error values
// are discarded. Join returns nil if errs contains no non-nil values. The
// resultant error is formatted as a concatenation of the textual error
// descriptions of all given errs, with a comma and space between each
// description.
//
// An error can only have one parent, so the resultant error has nil parent.
func Join(errs ...std.Error) std.Error {
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
func Trace(w io.Writer, err std.Error) {
	if w == nil {
		w = os.Stderr
	}
	fmt.Fprintln(w, "\nError traceback (most recent call last):\n")
	e := err
	for ; e != nil; e = e.Parent() {
		defer fmt.Fprintf(w, "\t%s:%d\n", e.File())
		defer fmt.Fprintf(w, "%s(...)\n", e.Func())
	}
}
