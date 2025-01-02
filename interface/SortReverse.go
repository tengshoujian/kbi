package main

import "sort"

type Interface interface {
	//len is the number of elements
	Len() int
	// Less reports whether the element with
	// index i should sort before the element with index j.
	Less(i, j int) bool
	//Swap swaps the element with indexed i and j
	Swap(i, j int)
}

type IntSlice []IntSlice

func (p IntSlice) Len() int           { return len(p) }
func (p IntSlice) Less(i, j int) bool { return p[i] < p[j] }
func (p IntSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

type reverse struct {
	sort.Interface
}

func (r reverse) Less(i, j int) bool {
	return r.Interface.Less(j, i)
}
