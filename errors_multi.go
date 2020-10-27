package bonusly

/*
Source: github.com/mongodb/grip

Error Catcher

The MutiCatcher type makes it possible to collect from a group of
operations and then aggregate them as a single error.
*/

import (
	"fmt"
	"strings"
	"sync"

	"github.com/pkg/errors"
)

// CheckFunction are functions which take no arguments and return an
// error.
type CheckFunction func() error

// Catcher is an interface for an error collector for use when
// implementing continue-on-error semantics in concurrent
// operations. There are three different Catcher implementations
// provided by this package that differ *only* in terms of the
// string format returned by String() (and also the format of the
// error returned by Resolve().)
//
// If you do not use github.com/pkg/errors to attach
// errors, the implementations are usually functionally
// equivalent. The Extended variant formats the errors using the "%+v"
// (which returns a full stack trace with pkg/errors,) the Simple
// variant uses %s (which includes all the wrapped context,) and the
// basic catcher calls error.Error() (which should be equvalent to %s
// for most error implementations.)
type catcher interface {
	Add(error)
	AddWhen(bool, error)
	Extend([]error)
	ExtendWhen(bool, []error)
	Len() int
	HasErrors() bool
	String() string
	Resolve() error
	Errors() []error

	New(string)
	NewWhen(bool, string)
	Errorf(string, ...interface{})
	ErrorfWhen(bool, string, ...interface{})

	Wrap(error, string)
	Wrapf(error, string, ...interface{})

	Check(CheckFunction)
	CheckExtend([]CheckFunction)
	CheckWhen(bool, CheckFunction)
}

// multiCatcher provides an interface to collect and coalesse error
// messages within a function or other sequence of operations. Used to
// implement a kind of "continue on error"-style operations. The
// methods on MultiCatatcher are thread-safe.
type baseCatcher struct {
	errs  []error
	mutex sync.RWMutex
	fmt.Stringer
}

// Add takes an error object and, if it's non-nil, adds it to the
// internal collection of errors.
func (c *baseCatcher) Add(err error) {
	if err == nil {
		return
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.errs = append(c.errs, err)
}

// Len returns the number of errors stored in the collector.
func (c *baseCatcher) Len() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return len(c.errs)
}

// HasErrors returns true if the collector has ingested errors, and
// false otherwise.
func (c *baseCatcher) HasErrors() bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return len(c.errs) > 0
}

// Extend adds all non-nil errors, passed as arguments to the catcher.
func (c *baseCatcher) Extend(errs []error) {
	if len(errs) == 0 {
		return
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	for _, err := range errs {
		if err == nil {
			continue
		}

		c.errs = append(c.errs, err)
	}
}

func (c *baseCatcher) Errorf(form string, args ...interface{}) {
	if form == "" {
		return
	} else if len(args) == 0 {
		c.New(form)
		return
	}
	c.Add(errors.Errorf(form, args...))
}

func (c *baseCatcher) New(e string) {
	if e == "" {
		return
	}
	c.Add(errors.New(e))
}

func (c *baseCatcher) Wrap(err error, m string) { c.Add(errors.Wrap(err, m)) }

func (c *baseCatcher) Wrapf(err error, f string, args ...interface{}) {
	c.Add(errors.Wrapf(err, f, args...))
}

func (c *baseCatcher) AddWhen(cond bool, err error) {
	if !cond {
		return
	}

	c.Add(err)
}

func (c *baseCatcher) ExtendWhen(cond bool, errs []error) {
	if !cond {
		return
	}

	c.Extend(errs)
}

func (c *baseCatcher) ErrorfWhen(cond bool, form string, args ...interface{}) {
	if !cond {
		return
	}

	c.Errorf(form, args...)
}

func (c *baseCatcher) NewWhen(cond bool, e string) {
	if !cond {
		return
	}

	c.New(e)
}

func (c *baseCatcher) Check(fn CheckFunction) { c.Add(fn()) }

func (c *baseCatcher) CheckWhen(cond bool, fn CheckFunction) {
	if !cond {
		return
	}

	c.Add(fn())
}

func (c *baseCatcher) CheckExtend(fns []CheckFunction) {
	for _, fn := range fns {
		c.Add(fn())
	}
}

func (c *baseCatcher) Errors() []error {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	out := make([]error, len(c.errs))

	copy(out, c.errs)

	return out
}

// Resolve returns a final error object for the Catcher. If there are
// no errors, it returns nil, and returns an error object with the
// string form of all error objects in the collector.
func (c *baseCatcher) Resolve() error {
	if !c.HasErrors() {
		return nil
	}

	return errors.New(c.String())
}

type basicCatcher struct{ *baseCatcher }

// newBasicCatcher collects error messages and formats them using a
// new-line separated string of the output of error.Error()
func newBasicCatcher() catcher {
	c := &baseCatcher{}
	c.Stringer = &basicCatcher{c}
	return c
}

func (c *basicCatcher) String() string {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	output := make([]string, len(c.errs))

	for idx, err := range c.errs {
		output[idx] = err.Error()
	}

	return strings.Join(output, "\n")
}
