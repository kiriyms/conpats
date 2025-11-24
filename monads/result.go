package monads

type Result struct {
	value int
	err   error
}

// Map applies a function to the value inside the Result if there is no error.
func (r Result) Map(f func(int) int) Result {}

// Pure creates a Result containing a value with no error.
func Pure(value int) Result {}

// Apply applies a Result containing a function to a Result containing a value.
func (r Result) Apply(rf Result) Result {}

// Bind applies a function that returns a Result to the value inside the Result if there is no error.
func (r Result) Bind(f func(int) Result) Result {}
