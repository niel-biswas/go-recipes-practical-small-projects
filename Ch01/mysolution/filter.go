package main

import (
	"fmt"
	"reflect"
)

// filter function defines the implementation logic of predicate i.e. return a slice with only values that predicate returns true
func filter(pred func(any) bool, values []any) []any {
	result := make([]any, 0)
	for _, val := range values {
		if pred(val) {
			result = append(result, val)
		}
	}
	return result
}

func isOdd(n any) bool {
	// check if the entered parameter is of Integer type
	if reflect.TypeOf(n).String() == "int" {
		return n.(int)%2 == 1 // Check modulo operation(i.e. division with 2 remainder value). Result(Non-Zero remainder): TRUE -> Odd, FALSE -> Even
	}
	return false
}

func isEven(n any) bool {
	// check if the entered parameter is of Integer type
	if reflect.TypeOf(n).String() == "int" {
		return n.(int)%2 == 0 // Check modulo operation(i.e. division with 2 remainder value). Result(Zero remainder): TRUE -> Even, FALSE -> Odd
	}
	return false
}

func NoApplicableFilter(n any) bool {
	// check if any filter criteria is applicable to the input Type
	if isOdd(n) == false && isEven(n) == false && isPrime(n) == false && isPalindromeString(n) == false {
		return true
	}
	return false
}

func isPrime(n any) bool {
	var count int = 0
	// check if the entered parameter is of Integer type
	if reflect.TypeOf(n).String() == "int" {
		for i := 2; i <= n.(int)/2; i++ {
			if n.(int)%i == 0 {
				count++
				break
			}
		}
		return count == 0 && n.(int) != 1 // Check if the number has exactly 2 factors i.e. self & 1. Result: TRUE -> Prime, FALSE -> Non-Prime
	}
	return false
}

func isPalindromeString(n any) bool {
	// check if the entered parameter is of String type
	if reflect.TypeOf(n).String() == "string" {
		inputSlice := fmt.Sprintf("%v", n)
		for i := 0; i < len(inputSlice)/2; i++ {

			if inputSlice[i] != inputSlice[len(inputSlice)-1-i] {
				return false
			}

		}

		return true
	}

	return false
}

func main() {
	anyVals := []any{1, 2, 3, 4, 5, "madam", 6, 7, 8, "kayak", 45.25, 9, 10, 11, 12, 13, 123321, "NotPalindrome", 452345}
	fmt.Printf("Entered any values: %v\n", anyVals)
	fmt.Printf("Filter(isOdd): %v\n", filter(isOdd, anyVals))
	fmt.Printf("Filter(isEven): %v\n", filter(isEven, anyVals))
	fmt.Printf("Filter(isPrime): %v\n", filter(isPrime, anyVals))
	fmt.Printf("Filter(isPalindromeString): %v\n", filter(isPalindromeString, anyVals))
	fmt.Printf("NoApplicableFilter: %v\n", filter(NoApplicableFilter, anyVals))
}
