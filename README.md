# builder-gen

<p align="center"><img src="https://user-images.githubusercontent.com/2807589/104822381-fac68500-580f-11eb-9e26-0dc9a0ed2776.png" width="300"></p>

builder gen is a go tool that helps you 
generate optional functions for a struct 
via a builder type

### Command Line Usage
```text
  -dir string
        the directory to run builder-gen on. Will use working directory if not provided
  -recursive
        set to true to recursively iterate directories; auto excludes any directories starting with .
  -recursive-exclusions string
        a comma separated list of directories to exclude when recursively iterating

```


### Setup Example

1. Add a comment with `::builder-gen` inside of the struct
you want to generate a builder for
   - ```go
     // struct is saved inside of named my_opts.go
     type MyOpts struct {
       // ::builder-gen
       ExampleShown bool
       Cool *string
       Steps []string
     }
     ```
1. add the following inside of a generate.go file within the same package
   - ```go
     //go:generate go run github.com/darwayne/builder-gen`
     ```
1. Run `go generate ./...`
1. A file called `build_gen_my_opts.go` gets generated with useful builder / optional funcs helpers for your struct
   - note: a file gets generated for each struct type within the package with the `::builder-gen` comment

### Features

***Null Check Helper Methods***
```go
 type MyOpts struct {
   // ::builder-gen
   ExampleShown bool
   Cool *string
   Steps []string
 }

var opt MyOpts

// Because Cool is a pointer we generate a HasCool helper function
// to help you check if it was set
if opt.HasCool() {...}
```

***Builder Generated Just for your Struct***

```go
 type MyOpts struct {
   // ::builder-gen
   ExampleShown bool
   Cool *string
   Steps []string
 }

optFns := NewMyOptsBuilder().
    ExampleShown(true).
    Cool("beans").
    Steps("1", "2", "3").
    Build()
```

***Generate Original Struct from OptionalFunctions***
```go
 type MyOpts struct {
   // ::builder-gen
   ExampleShown bool
   Cool *string
   Steps []string
 }

opts := ToMyOpts(NewMyOptsBuilder().
    ExampleShown(true).
    Cool("beans").
    Steps("1", "2", "3").
    Build()...)

// opt will be equivalent to the following
cool := "beans"
opts2 := MyOpts{
    Cool: &cool, 
    ExampleShown: true, 
    Steps: []string{"1", "2", "3"},
}
```

***Configure generation for a particular struct***
```
Usage of ::builder-gen flags:
  -no-builder
        set this flag if you want to exclude creating the builder object
  -prefix string
        if set this will be the prefix of your global functions. Note: with-globals option required
  -suffix string
        if set this will be the suffix of your global functions. Note: with-globals option required
  -with-globals
        set this flag if you want to generate global functions as well

```

Configuration is set directly in the comment 
```go
 type MyOpts struct {
   // ::builder-gen -with-globals -prefix=With -suffix=Opt -no-builder
   ExampleShown bool
   Cool *string
   Steps []string
 }

opts := ToMyOpts(
    WithExampleShownOpt(true),
    WithCoolOpt("beans"),
    WithStepsOpt("1", "2", "3"))

// opt will be equivalent to the following
cool := "beans"
opts2 := MyOpts{
    Cool: &cool, 
    ExampleShown: true, 
    Steps: []string{"1", "2", "3"},
}
```

***Set Defaults on Your Struct And Have Optional Functions Override Them***
```go
cool := "beans"
opts := MyOpts{
    Cool: &cool, 
    ExampleShown: true, 
    Steps: []string{"1", "2", "3"},
}

ToMyOptsWithDefault(&opts, NewMyOptsBuilder().
                               ExampleShown(false).
                               Cool("son").
                               Build()...)

// opts now contains
cool := "son"
MyOpts{
    Cool: &cool, 
    ExampleShown: false, 
    Steps: []string{"1", "2", "3"},
}
```


