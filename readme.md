# multierr

[![Go Report Card](https://goreportcard.com/badge/github.com/maja42/multierr)](https://goreportcard.com/report/github.com/maja42/multierr)
[![GoDoc](https://godoc.org/github.com/maja42/multierr?status.svg)](https://godoc.org/github.com/maja42/multierr)

`multierr` is a package that allows combining multiple errors into a single `error` type.
This allows functions to return multiple errors at once.

Callers can either use the returned multi-error as a conventional error (which is printed as a nice human-readable string), 
or continue working with it by appending or unwrapping individual errors.

## Use cases


### Validation
When validating configurations or user-input, it's always a great user-experience to see all problems at once.
Just append all individual errors to a multi-error and return it to the user.

### APIs
When implementing APIs (be it a WebServer's REST-API, a protobuf RPC API or any other interface), 
validating incoming data and **returning all problems at once** greatly improves a developer's quality of life.

No longer do API-users need to call an endpoint just to receive the next error they need to fix. 

### Collecting go-routine errors

Sometimes, multiple concurrently-running go-routines can each return an error.
Which error should you report? The first one? What about the others, log them or ignore them?
A multi-error can simply collect all those errors and return them at once.

## Usage

### Simple validation
```go
type Input struct {
	Name string
	Age  int
}

func (i *Input) Validate() error {
	var valErr error

	if i.Name == "" {
		valErr = multierr.Append(valErr, errors.New("missing name"))
	}
	if i.Age < 18 {
		valErr = multierr.Append(valErr, errors.New("too young"))
	}
	return valErr
}
```

This prints the following output:
```
2 errors occurred:
  - missing name
  - too young
```

### Custom title

If you instead return `multierr.Titled(valErr, "Invalid input:")`, you can get the following output:

```
Invalid input:
  - missing name
  - too young
```

### Custom prefix

Alternatively, you can prefix each error via `multierr.Prefixed(valErr, "Invalid input: ")` to get the following output:

```
Invalid input: missing name
Invalid input: too young
```

### Combining multiple multi-errors

When validating nested structures, you often receive errors from sub-validators. 
The same can happen when calling functions.

These cases can be handled in 4 different ways, all of them producing great error messages:


#### Option 1: `multiErr.Append(valErr, err)`

```go
type Input struct {
	Name    string
	Age     int
	Address Address
}

type Address struct {
	City   string
	Street string
}

func (i *Input) Validate() error {
	var valErr error

	if i.Name == "" {
		valErr = multierr.Append(valErr, errors.New("missing name"))
	}
	if i.Age < 18 {
		valErr = multierr.Append(valErr, errors.New("too young"))
	}
	valErr = multierr.Append(valErr, i.Address.Validate())

	return multierr.Titled(valErr, "invalid input:")
}

func (a *Address) Validate() error {
	var valErr error

	if a.City == "" {
		valErr = multierr.Append(valErr, errors.New("missing city"))
	}
	if a.Street == "" {
		valErr = multierr.Append(valErr, errors.New("missing street"))
	}
	return valErr
}
```

This is the simplest version. \
And you just got rid of those nasty if-error-checks. \
You don't need to check for nil-errors when validating the address.
If there is no error, `Append()` will simply do nothing.

You get the following error message:
```
invalid input:
  - missing name
  - too young
  - 2 errors occurred:
      - missing city
      - missing street
```



#### Option 2: `multiErr.Append(valErr, multierr.Titled(...))`

You can get a slightly better error message by choosing your own title. \
Replace the address validation with this piece of code:

```go
err := i.Address.Validate()
valErr = multierr.Append(valErr, multierr.Titled(err, "invalid address:"))
```

Again - you don't need to check for errors. The `Titled`-function simply returns `nil` if there was no error. \
You get the following output:

```
invalid input:
  - missing name
  - too young
  - invalid address:
      - missing city
      - missing street
```


#### Option 3: `multierr.Merge(valErr, err)` 

Now, what if you don't want nested error messages? Just merge them! \
Replace the address validation with this:
```go
valErr = multierr.Merge(valErr, i.Address.Validate())
```

You will get the following:

```
invalid input:
  - missing name
  - too young
  - missing city
  - missing street
```


#### Option 4: `multierr.MergePrefixed(valErr, err)`

In the above example, you do not see that city and street are subfields of the address.
You can keep that information by using prefixes.
Replace the address validation with this:
```go
valErr = multierr.MergePrefixed(valErr, "invalid adress: ", i.Address.Validate())
```

You will get the following:

```
invalid input:
  - missing name
  - too young
  - invalid adress: missing city
  - invalid adress: missing street
```



#### Option 5: `multiErr.Append(err, fmt.Errorf(...))`

And, of course, calling `fmt.Errorf()` instead of `multierr.Append()` also yields great results.

Perform the address validation as follows:
```go
if err := i.Address.Validate(); err != nil {
	valErr = multierr.Append(valErr, fmt.Errorf("invalid address: %s", err))
}
```

You will get:

```
invalid input:
  - missing name
  - too young
  - invalid address: 2 errors occurred:
      - missing city
      - missing street
```

### Custom error format

Sometimes, you just want to format errors differently. And that's entirely possible:

```go
err := multierr.Append(
	errors.New("error 1"),
	errors.New("error 2"),
)
err.Formatter = func(errs []error) string {
	return fmt.Sprintf("there are %d errors", len(errs))
}
```

This is not feasible if you want to have a different error format **globally** though. \
In that case, you can overwrite the default formatter:

```go
multierr.DefaultFormatter = func(errs []error) string {
	return fmt.Sprintf("there are %d errors", len(errs))
}
```

## Accessing the list of errors

You can access a list with all sub-errors by simply calling 
```go
errList := multierr.Inspect(multiErr)
```

This also works if the provided argument is not actually a multi-error. \
If it's a normal `error`, the returned list will have the error as a single element.

## Unwrapping specific sub-errors

Multi-errors support the standard library's `errors.Unwrap()`, `errors.As()` and `errors.Is()` methods. \
It's therefore possible to inspect certain root-causes of an error.


