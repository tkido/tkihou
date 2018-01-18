package myarr

import (
	"reflect"
	"regexp"
	"testing"
)

func TestConcat(t *testing.T) {
	want := NewMyArr("1", "2", "3", "4")
	got := NewMyArr("1", "2")
	got.Concat(NewMyArr("3", "4"))
	if !reflect.DeepEqual(want, got) {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestFirst(t *testing.T) {
	arr := NewMyArr("1", "2", "3")
	if got, want := arr.First(), "1"; got != want {
		t.Errorf("First: got %v want %v", got, want)
	}
}

func TestMap(t *testing.T) {
	want := NewMyArr("11", "22", "33")
	got := NewMyArr("1", "2", "3")
	got.Map(repeat)
	if !reflect.DeepEqual(want, got) {
		t.Errorf("got %v want %v", got, want)
	}
}
func repeat(s string) string {
	return s + s
}
func TestPop(t *testing.T) {
	arr := NewMyArr("1", "2", "3")
	want := "1"
	got := arr.Pop()
	if got != want {
		t.Errorf("got %v want %v", got, want)
	}

	wantArr := NewMyArr("2", "3")
	gotArr := arr
	if !reflect.DeepEqual(wantArr, gotArr) {
		t.Errorf("got %v want %v", gotArr, wantArr)
	}
}

func TestPush(t *testing.T) {
	want := NewMyArr("1", "2", "3", "4", "5")
	got := NewMyArr("1")
	got.Push("2").Push("3").Push("4", "5")
	if !reflect.DeepEqual(want, got) {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestSize(t *testing.T) {
	arr := NewMyArr("1", "2", "3")
	want := 3
	got := arr.Size()
	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestTakeBlock(t *testing.T) {
	re := regexp.MustCompile(`^#`)
	arr := NewMyArr("#1", "#2", "3", "4")
	want := NewMyArr("1", "2")
	wantRest := NewMyArr("3", "4")
	got := arr.TakeBlock(re)
	if !reflect.DeepEqual(want, got) {
		t.Errorf("got %v want %v", got, want)
	}
	if !reflect.DeepEqual(wantRest, arr) {
		t.Errorf("got %v want %v", arr, wantRest)
	}
}
func TestUnshift(t *testing.T) {
	want := NewMyArr("1", "2", "3")
	got := NewMyArr("3")
	got.Unshift("2").Unshift("1")
	if !reflect.DeepEqual(want, got) {
		t.Errorf("Push: got %v want %v", got, want)
	}
}
