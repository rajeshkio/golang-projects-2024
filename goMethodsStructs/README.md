# Go Structs and Methods with Testing

This repository demonstrates how to create and use structs in Go, along with proper testing techniques. It focuses on Go's encapsulation patterns, methods, and table-driven tests.

## Project Structure

```
.
├── go.mod
├── main.go
└── task1
    ├── task1.go
    └── task1_test.go
```

## Task Overview

The main task involves creating a `Person` struct with private fields, getters, setters, and a `String()` method for formatting. The project also includes comprehensive testing using Go's testing package.

## Key Concepts

### 1. Struct Definition and Encapsulation

Go uses field name capitalization to control visibility:
- Lowercase field names are private (unexported)
- Uppercase field names are public (exported)

We use private fields with public methods to enforce encapsulation:

```go
type Person struct {
    name    string  // Private field
    age     int     // Private field
    address string  // Private field
}
```

### 2. Constructor Pattern

Since we can't directly initialize private fields from outside the package, we use a constructor:

```go
// NewPerson creates a new Person instance
func NewPerson(name string, age int, address string) Person {
    return Person{
        name:    name,
        age:     age,
        address: address,
    }
}
```

### 3. Getters and Setters

Methods to access and modify private fields:

```go
// Getter example
func (p Person) GetName() string {
    return p.name
}

// Setter example
func (p *Person) SetName(name string) {
    p.name = name
}
```

### 4. String Method

Implementing the `fmt.Stringer` interface for nice string representation:

```go
// String implements the Stringer interface
func (p Person) String() string {
    return fmt.Sprintf("Person{name: %s, age: %d, address: %s}", 
                      p.name, p.age, p.address)
}
```

## Testing in Go

### Basic Tests

Each test function:
- Starts with `Test`
- Takes `*testing.T` as a parameter
- Checks expected values against actual results

```go
func TestNewPerson(t *testing.T) {
    person := NewPerson("John", 30, "New York")
    
    if person.name != "John" {
        t.Errorf("Expected name to be 'John', got '%s'", person.name)
    }
    
    // More assertions...
}
```

### Table-Driven Tests

A powerful pattern for testing multiple scenarios:

```go
func TestPersonGettersWithTableDriven(t *testing.T) {
    testCases := []struct {
        name            string
        person          Person
        expectedName    string
        expectedAge     int
        expectedAddress string
    }{
        {
            name:            "Regular person",
            person:          NewPerson("Alice", 30, "Wonderland"),
            expectedName:    "Alice",
            expectedAge:     30,
            expectedAddress: "Wonderland",
        },
        // More test cases...
    }
    
    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            // Assertions...
        })
    }
}
```

## Running the Code

### Setup

1. Clone the repository
2. Make sure Go is installed (version 1.16+ recommended)

### Running the Main Program

```bash
go run main.go
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./task1

# Run a specific test
go test -v ./task1 -run TestNewPerson
```

## Common Go Patterns Demonstrated

1. **Private fields with public methods** - For encapsulation
2. **Constructor functions** - For proper initialization
3. **Pointer receivers for setters** - To modify the actual struct
4. **Value receivers for getters** - For immutable operations
5. **Table-driven tests** - For comprehensive test coverage

## Avoiding Common Mistakes

1. **Accessibility**: Remember that lowercase field names cannot be accessed directly from outside their package
2. **Pointer vs. Value Receivers**: Use pointer receivers (`*Person`) when methods need to modify the struct
3. **Test Coverage**: Aim to test all methods and edge cases

## Next Steps

- Add validation logic to the setter methods
- Implement custom JSON marshaling/unmarshaling
- Add benchmarks to measure performance
