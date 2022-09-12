package algorithms

/*
#include "basic.h"
#include "sort.h"
#include "stdlib.h"
*/
import "C"
import (
	"encoding/json"
	"fmt"
	"time"
	"unsafe"
)

//Sorting package struct
type Sorting struct{}

//Method parses input and executes the the requested algorithm
func (s Sorting) run(request Parameters) (string, string, error) {
	//Parsing
	var numArray []int32
	var start time.Time
	var executionTime time.Duration
	err := json.Unmarshal([]byte("["+request.Input+"]"), &numArray)
	if err != nil {
		return "Incorrect input, please refer to documentation", "", nil
	}
	if len(numArray) == 0 {
		return "Empty array", "", nil
	}
	arrayPointer := unsafe.Pointer(&numArray[0])
	//Algorithm execution
	if len(request.Parameters) < 1 || request.Parameters[0] == "ascending" {
		switch request.Algorithm {
		case "InsertionSort":
			start = time.Now()
			C.insertionSort((*C.int)(arrayPointer), C.int(len(numArray)))
			executionTime = time.Since(start)
		case "BubbleSort":
			start = time.Now()
			C.bubbleSort((*C.int)(arrayPointer), C.int(len(numArray)))
			executionTime = time.Since(start)
		case "QuickSort":
			start = time.Now()
			C.quickSort((*C.int)(arrayPointer), C.int(0), C.int(len(numArray)-1))
			executionTime = time.Since(start)
		case "MergeSort":
			start = time.Now()
			C.mergeSort((*C.int)(arrayPointer), C.int(0), C.int(len(numArray)-1))
			executionTime = time.Since(start)
		case "SelectionSort":
			start = time.Now()
			C.selectionSort((*C.int)(arrayPointer), C.int(len(numArray)))
			executionTime = time.Since(start)
		}
	} else if request.Parameters[0] == "descending" {
		switch request.Algorithm {
		case "InsertionSort":
			start = time.Now()
			C.insertionSortReverse((*C.int)(arrayPointer), C.int(len(numArray)))
			executionTime = time.Since(start)
		case "BubbleSort":
			start = time.Now()
			C.bubbleSortReverse((*C.int)(arrayPointer), C.int(len(numArray)))
			executionTime = time.Since(start)
		case "QuickSort":
			start = time.Now()
			C.quickSortReverse((*C.int)(arrayPointer), C.int(0), C.int(len(numArray)-1))
			executionTime = time.Since(start)
		case "MergeSort":
			start = time.Now()
			C.mergeSortReverse((*C.int)(arrayPointer), C.int(0), C.int(len(numArray)-1))
			executionTime = time.Since(start)
		case "SelectionSort":
			start = time.Now()
			C.selectionSortReverse((*C.int)(arrayPointer), C.int(len(numArray)))
			executionTime = time.Since(start)
		}
	} else {
		return "Order option unavailable, please refer to documentation", "", nil
	}
	str := fmt.Sprint(numArray)
	return str[1 : len(str)-1], executionTime.String(), nil
}
