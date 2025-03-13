# Task 3: Value Receivers vs Pointer Receivers in Go

This task demonstrates the important differences between value receivers and pointer receivers in Go methods. Understanding these differences is crucial for writing effective Go code and making appropriate design decisions.

## Objective

- Create identical methods with both value receivers and pointer receivers
- Test modification behavior with each type
- Analyze when to use each approach

## Code Explanation

### Value Receiver Methods

```go
// Value receiver - creates a copy, doesn't modify the original
func UpdateAgeValue(p Person, newAge int) Person {
    p.Age = newAge // Changes only affect the local copy
    return p       // Return the modified copy
}
```

When using a value receiver, Go creates a **copy** of the original struct. Any modifications only affect this copy, not the original. To see the changes, you must use the returned value.

### Pointer Receiver Methods

```go
// Pointer receiver - operates on the original object
func UpdateAgePointer(p *Person, newAge int) {
    p.Age = newAge // Changes affect the original
}
```

When using a pointer receiver, Go passes a reference to the original struct. Any modifications directly affect the original struct.

## Key Observations

Our demonstrations show:

1. Value receivers:
   - Don't modify the original struct
   - Return a modified copy that must be captured
   - Example: `updatedPerson := UpdateAgeValue(person, 35)`

2. Pointer receivers:
   - Directly modify the original struct
   - No need to capture a return value
   - Example: `UpdateAgePointer(&person, 40)`

## When to Use Each

### Use Value Receivers When:

- You don't need to modify the receiver
- The struct is small (like bool, int)
- You want immutability (method won't modify original)
- Method logic treats the receiver as a value (conceptually)

### Use Pointer Receivers When:

- You need to modify the receiver
- The struct is large (to avoid costly copying)
- The receiver shouldn't be copied (e.g., contains a mutex)
- For consistency (if one method needs a pointer, use pointers for all)

## Memory and Performance Considerations

- Value receivers create a copy, which uses more memory for large structs
- Pointer receivers only pass a memory address, more efficient for large structs
- For tiny structs (like a single int), the difference may be negligible

## Go's Convenience Feature

Go allows you to call pointer methods on values and value methods on pointers. The compiler automatically takes the address of a value when needed, making the code more readable.

## Testing Approach

Testing value vs pointer receivers involves:
1. Creating an initial struct
2. Calling methods with different receiver types
3. Verifying if/how the original was modified
4. Checking the returned values when applicable