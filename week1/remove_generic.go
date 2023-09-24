// Copyright@daidai53 2023
package week1

// Remove
// 切片删除，支持泛型
func Remove[T comparable](index int, slice []T) []T {
	if slice == nil || len(slice) == 0 {
		return slice
	}
	if index < 0 || index >= len(slice) {
		return slice
	}
	return append(slice[:index], slice[index+1:]...)
}

var shrinkRatio = 0.8

// ShrinkRemove
// 带缩容的切片删除，支持泛型
func ShrinkRemove[T comparable](index int, slice []T) []T {
	if slice == nil {
		return nil
	}
	if len(slice) == 0 && cap(slice) > 0 {
		return make([]T, 0)
	}
	if index >= 0 && index < len(slice) {
		slice = append(slice[:index], slice[index+1:]...)
		if float64(len(slice))/float64(cap(slice)) <= shrinkRatio {
			newSlice := make([]T, len(slice))
			copy(newSlice, slice)
			slice = newSlice
		}
	}
	return slice
}
