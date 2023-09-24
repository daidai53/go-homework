// Copyright@daidai53 2023
package week1

import (
	"testing"
)

func compareIntSlice(sli1, sli2 []int) bool {
	if sli1 == nil && sli2 == nil {
		return true
	} else if sli1 == nil || sli2 == nil {
		return false
	}

	if len(sli1) != len(sli2) {
		return false
	}

	for i := range sli1 {
		if sli1[i] != sli2[i] {
			return false
		}
	}
	return true
}

func compareSlice[T comparable](sli1, sli2 []T) bool {
	if sli1 == nil && sli2 == nil {
		return true
	} else if sli1 == nil || sli2 == nil {
		return false
	}

	if len(sli1) != len(sli2) {
		return false
	}

	for i := range sli1 {
		if sli1[i] != sli2[i] {
			return false
		}
	}
	return true
}

func Test_SimpleRemove(t *testing.T) {
	slice := []int{
		1, 2, 3, 4, 5,
	}
	slice = SimpleRemove(-1, slice)
	if !compareIntSlice(slice, []int{1, 2, 3, 4, 5}) {
		t.Fatalf("使用非法下标删除成功")
	}

	slice = SimpleRemove(2, slice)
	if !compareIntSlice(slice, []int{1, 2, 4, 5}) {
		t.Fatalf("使用下标删除失败")
	}
}

func Test_Remove(t *testing.T) {
	sliceInt := []int{1, 2, 3, 4, 5}
	sliceInt = Remove(-1, sliceInt)
	if !compareSlice(sliceInt, []int{1, 2, 3, 4, 5}) {
		t.Fatalf("使用非法下标删除Int切片成功")
	}
	sliceInt = Remove(2, sliceInt)
	if !compareSlice(sliceInt, []int{1, 2, 4, 5}) {
		t.Fatalf("使用下标删除Int切片失败")
	}

	sliceString := []string{"1", "2", "3", "4", "5"}
	sliceString = Remove(5, sliceString)
	if !compareSlice(sliceString, []string{"1", "2", "3", "4", "5"}) {
		t.Fatalf("使用非法下标删除string切片成功")
	}
	sliceString = Remove(2, sliceString)
	if !compareSlice(sliceString, []string{"1", "2", "4", "5"}) {
		t.Fatalf("使用下标删除Int切片失败")
	}
}

func Test_Remove_WithShrink(t *testing.T) {
	sliceInt64 := []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	sliceInt64 = ShrinkRemove(2, sliceInt64)
	// cap10,len9,不触发缩容
	if !compareSlice(sliceInt64, []int64{1, 2, 4, 5, 6, 7, 8, 9, 10}) ||
		cap(sliceInt64) != 10 {
		t.Fatalf("使用带缩容的下标删除Int切片成功，但是触发缩容")
	}

	sliceInt64 = ShrinkRemove(2, sliceInt64)
	// cap10,len8,触发缩容
	if !compareSlice(sliceInt64, []int64{1, 2, 5, 6, 7, 8, 9, 10}) ||
		cap(sliceInt64) != 8 {
		t.Fatalf("使用带缩容的下标删除Int切片成功，未触发缩容")
	}
}
