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
// update the contextual information, the Raise function can be invoked when
// intending to raise a predefined error from a current context, like so:
//
//	errors.Raise(EOF)
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
// the Join function is provided to return all errors as one. However, this
// removes the context behind the combined errors.
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
type Error struct {
	addr   uintptr
	file   string
	fn     string
	line   int
	parent error
	text   string
}

func (e Error) Addr() uintptr {
	return e.addr
}

func (e Error) Error() string {
	next := true
	var text string
	var err error = e
	for next {
		switch t := err.(type) {
		case Error:
			if t.Parent() == nil {
				next = false
				text += t.Text()
			} else {
				text += t.Text() + ": "
				err = t.Parent()
			}
		case error:
			text += t.Error()
			next = false
		}
	}
	return text
}

func (e Error) File() string {
	return e.file
}

func (e Error) Func() string {
	return e.fn
}

func (e Error) raise() error {
	err := e
	addr, file, line, _ := runtime.Caller(2)
	f := runtime.FuncForPC(addr)
	fn := f.Name()
	err.addr = addr
	err.file = file
	err.fn = fn
	err.line = line
	return err
}

func (e Error) Line() int {
	return e.line
}

func (e Error) Parent() error {
	return e.parent
}

func (e Error) Text() string {
	return e.text
}

// Has reports whether the textual error description of err or any of its parent
// errors match the textual error description of target.
func Has(err, target error) bool {
	text := ""
	if err == nil && target == nil {
		return true
	} else if target == nil {
		return false
	}
	switch t := target.(type) {
	case Error:
		text = t.Text()
	case error:
		text, _, _ = strings.Cut(t.Error(), ": ")
	}
	for {
		switch t := err.(type) {
		case Error:
			if t.Text() == text {
				return true
			}
			err = t.Parent()
			if err == nil {
				return false
			}
		case error:
			utext, _, _ := strings.Cut(err.Error(), ": ")
			if utext == text {
				return true
			} else {
				return false
			}
		}
	}
}

// Is reports whether the textual error description of err matches the textual
// error description of target.
func Is(err, target error) bool {
	if err == nil && target == nil {
		return true
	} else if err == nil || target == nil {
		return false
	}
	switch t := err.(type) {
	case Error:
		switch u := target.(type) {
		case Error:
			return t.Text() == u.Text()
		case error:
			text, _, _ := strings.Cut(u.Error(), ": ")
			return t.Text() == text
		}
	case error:
		switch u := target.(type) {
		case Error:
			text, _, _ := strings.Cut(t.Error(), ": ")
			return text == u.Text()
		}
	}
	utext, _, _ := strings.Cut(err.Error(), ": ")
	vtext, _, _ := strings.Cut(target.Error(), ": ")
	return utext == vtext
}

// New returns an error whose textual error description is given by text, and
// whose parent error is err. If the new error has no parent, err should be
// given as nil.
//
// The current filename, line, program counter, and parent function name are
// stored within the error interface. Each call to New returns a distinct error
// value even if text is identical.
func New(text string, err error) error {
	addr, file, line, _ := runtime.Caller(1)
	f := runtime.FuncForPC(addr)
	fn := f.Name()
	return Error{
		addr:   addr,
		file:   file,
		fn:     fn,
		line:   line,
		parent: err,
		text:   text,
	}
}

// Join returns an error that combines the given errs. Any nil error values
// are discarded. Join returns nil if errs contains no non-nil values. The
// resultant error is formatted as a concatenation of the textual error
// descriptions of all given errs, with a comma and space between each
// description.
//
// An error can only have one parent, so the resultant error has nil parent.
func Join(errs ...error) error {
	first := true
	var text string
	for _, err := range errs {
		if err != nil {
			if first {
				first = false
			} else {
				text += ", "
			}
			switch t := err.(type) {
			case Error:
				text += t.Text()
			case error:
				text += t.Error()
			}
		}
	}
	if text == "" {
		return nil
	}
	addr, file, line, _ := runtime.Caller(1)
	f := runtime.FuncForPC(addr)
	fn := f.Name()
	return Error{
		addr:   addr,
		file:   file,
		fn:     fn,
		line:   line,
		parent: nil,
		text:   text,
	}
}

func Raise(err error) error {
	switch e := err.(type) {
	case Error:
		return e.raise()
	}
	return err
}

// Trace writes human-friendly error traceback information from err to w. If w
// is nil, Trace writes to the standard error stream.
func Trace(w io.Writer, err error) {
	if w == nil {
		w = os.Stderr
	}
	fmt.Fprintln(w, "Traceback (most recent call first):")
	for {
		switch t := err.(type) {
		case Error:
			defer fmt.Fprintf(w, "\t%s:%d\n", t.File(), t.Line())
			defer fmt.Fprintf(w, "\t%s\n", t.Text())
			defer fmt.Fprintf(w, "%s(...)\n", t.Func())
			err = t.Parent()
			if err == nil {
				return
			}
		case error:
			defer fmt.Fprintf(w, "\t%s\n", t.Error())
			defer fmt.Fprintf(w, "no-context error:\n")
			return
		}
	}
}
