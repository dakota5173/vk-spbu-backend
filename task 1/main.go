package main

import "fmt"

func twoSum(nums []int, target int) []int {
	seen := make(map[int]int, len(nums))

	for i, a := range nums {
		b := target - a
		if j, exists := seen[b]; exists {
			return []int{j, i}
		}
		seen[a] = i
	}
	return nil
}

func main() {
	nums := []int{-3, 1, 18, -5, -16, 11}
	target := -19

	fmt.Println(twoSum(nums, target))
}
