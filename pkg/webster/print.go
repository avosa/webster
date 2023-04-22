// pkg/webster/print.go

package webster

import (
	"fmt"
	"reflect"
	"strings"
)

// PrintValue prints the value of a single object.
func PrintValue(val interface{}) {
	// Check if the value is an array or slice.
	valType := reflect.TypeOf(val)
	if valType.Kind() == reflect.Array || valType.Kind() == reflect.Slice {
		// If it is an array or slice, convert it to a slice of interface{} values.
		arr := reflect.ValueOf(val)
		slice := make([]interface{}, arr.Len())
		for i := 0; i < arr.Len(); i++ {
			slice[i] = arr.Index(i).Interface()
		}
		// Print the slice using fmt.Sprintf and strings.TrimPrefix to remove the square brackets.
		fmt.Print(strings.TrimPrefix(fmt.Sprintf("%v", slice), "["))
		fmt.Print(strings.TrimSuffix(fmt.Sprintf("%v", slice), "]"))
	} else {
		// If it is not an array or slice, simply print the value using fmt.Println.
		fmt.Println(val)
	}
}

// PrintStatement represents a print statement in the program.
type PrintStatement struct {
	Value interface{}
}

// Execute prints the value of the expression to the console.
func (ps *PrintStatement) Execute() error {
	PrintValue(ps.Value)
	return nil
}
