package main

import (
	"testing"
)

// TestLinkedListInsert tests various scenarios of the Insert method
func TestLinkedListInsert(t *testing.T) {
	l := New()
	// Test inserting into an empty list
	if !l.Insert(0, 10) {
		t.Error("Insert failed to insert the first element")
	}
	if l.head.Value != 10 || l.length != 1 {
		t.Errorf("Expected head value of 10 and length 1, got value %d and length %d", l.head.Value, l.length)
	}

	// Test inserting at the end of the list
	if !l.Insert(1, 20) {
		t.Error("Insert failed to append second element")
	}
	if l.head.Next.Value != 20 || l.length != 2 {
		t.Errorf("Expected second value of 20 and length 2, got value %d and length %d", l.head.Next.Value, l.length)
	}

	// Test inserting at the middle of the list
	if !l.Insert(1, 15) {
		t.Error("Insert failed to insert element in the middle")
	}
	if l.head.Next.Value != 15 || l.head.Next.Next.Value != 20 || l.length != 3 {
		t.Errorf("Expected middle value of 15, got %d, expected third value of 20, got %d, expected length 3, got %d", l.head.Next.Value, l.head.Next.Next.Value, l.length)
	}

	// Test inserting out of bounds
	if l.Insert(5, 30) {
		t.Error("Insert did not fail when trying to insert out of bounds")
	}
}

// TestLinkedListRemove tests various scenarios of the Remove method
func TestLinkedListRemove(t *testing.T) {
	l := New()
	l.Insert(0, 10)
	l.Insert(1, 20)
	l.Insert(2, 30)

	// Test removing the first element
	if !l.Remove(0) {
		t.Error("Remove failed to remove the first element")
	}
	if l.head.Value != 20 || l.length != 2 {
		t.Errorf("Expected new head value of 20 and length 2, got value %d and length %d", l.head.Value, l.length)
	}

	// Test removing a middle element
	if !l.Remove(1) { // Now this is index 1, which should be 30
		t.Error("Remove failed to remove the middle element")
	}
	if l.head.Next != nil || l.length != 1 {
		t.Errorf("Expected final element to be nil and length 1, got next value %v and length %d", l.head.Next, l.length)
	}

	// Test removing the last element
	if !l.Remove(0) {
		t.Error("Remove failed to remove the last element")
	}
	if l.head != nil || l.length != 0 {
		t.Errorf("Expected empty list, got head %v and length %d", l.head, l.length)
	}

	// Test removing from an empty list
	if l.Remove(0) {
		t.Error("Remove did not fail when trying to remove from an empty list")
	}
}

// TestLinkedListGet tests various scenarios of the Get method
func TestLinkedListGet(t *testing.T) {
	l := New()
	l.Insert(0, 10)
	l.Insert(1, 20)
	l.Insert(2, 30)
	// Test getting each element
	tests := []struct {
		index    uint
		expected int
		ok       bool
	}{
		{0, 10, true},
		{1, 20, true},
		{2, 30, true},
		{3, 0, false}, // out of bounds
	}
	for _, tt := range tests {
		value, ok := l.Get(tt.index)
		if ok != tt.ok || (ok && value != tt.expected) {
			t.Errorf("Get(%d): expected %d, ok %t, got %d, ok %t", tt.index, tt.expected, tt.ok, value, ok)
		}
	}
}

// TestLinkedListFind tests various scenarios of the Find method
func TestLinkedListFind(t *testing.T) {
	l := New()
	l.Insert(0, 10)
	l.Insert(1, 20)
	l.Insert(2, 30)
	// Test finding each element
	tests := []struct {
		val      int
		expected uint
		found    bool
	}{
		{10, 0, true},
		{20, 1, true},
		{30, 2, true},
		{40, 0, false}, // not present
	}
	for _, tt := range tests {
		index, found := l.Find(tt.val)
		if found != tt.found || (found && index != tt.expected) {
			t.Errorf("Find(%d): expected index %d, found %t, got index %d, found %t", tt.val, tt.expected, tt.found, index, found)
		}
	}
}
