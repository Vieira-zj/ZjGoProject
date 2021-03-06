package sort

import (
	"fmt"
)

// ------------------------------
// #1. 二分查找 有序数组 O(logN) 递归
// ------------------------------

func binarySearch01(arr []int, start, end, val int) int {
	if start > end {
		return -1
	}

	mid := start + (end-start)/2
	if val > arr[mid] {
		return binarySearch01(arr, mid+1, end, val)
	} else if val < arr[mid] {
		return binarySearch01(arr, start, mid-1, val)
	} else {
		return mid
	}
}

// ------------------------------
// #2. 二分查找 有序数组 O(logN) 非递归
// ------------------------------

func binarySearch02(arr []int, val int) int {
	start := 0
	end := len(arr) - 1
	var mid int

	for start <= end {
		mid = start + (end-start)/2
		if val > arr[mid] {
			start = mid + 1
		} else if val < arr[mid] {
			end = mid - 1
		} else {
			return mid
		}
	}
	return -1
}

// TestSearchAlgorithms test for search algorithms.
func TestSearchAlgorithms() {
	if false {
		fmt.Println("\n#1. 二分查找 有序数组 O(logN) 递归")
		arr := []int{1, 3, 4, 6, 8, 9, 10, 12, 13, 77}
		for _, val := range []int{1, 12, 77} {
			fmt.Println("binary search results by index:", binarySearch01(arr, 0, len(arr)-1, val))
		}

		fmt.Println("\n#2. 二分查找 有序数组 O(logN) 非递归")
		for _, val := range []int{1, 12, 77} {
			fmt.Println("binary search results by index:", binarySearch02(arr, val))
		}
	}

	fmt.Println("search algorithms done.")
}
