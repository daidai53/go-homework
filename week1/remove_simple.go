// Copyright@daidai53 2023
package week1

// SimpleRemove
// 实现对int切片的删除操作
func SimpleRemove(index int, slice []int) []int {
	if slice == nil {
		return nil
	}
	length := len(slice)
	if index < 0 || index >= length {
		return slice
	}
	return append(slice[:index], slice[index+1:]...)
}
