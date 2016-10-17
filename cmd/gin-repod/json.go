package main

import (
	"fmt"
	"io"
)

const (
	elmObject = "object"
	elmArray  = "array"
)

type element struct {
	elm   string
	count int
}

type JSONWriter struct {
	W   io.Writer
	err error
	elm []element
}

func (w *JSONWriter) writeString(data string) {
	if w.err != nil {
		return
	}
}

func (w *JSONWriter) pop(elmType string) {
	if w.err != nil {
		return
	}

	l := len(w.elm)
	if l < 1 {
		w.err = fmt.Errorf("trying to close %s on empty stack", elmType)
	} else if e := w.elm[l-1]; e.elm != elmType {
		w.err = fmt.Errorf("type mismatch: trying to close %s with %s", e.elm, elmType)
	}

	w.elm = w.elm[:l-1]
}

func (w *JSONWriter) push(elmType string) {
	if w.err != nil {
		return
	}
	w.elm = append(w.elm, element{elm: elmType, count: 0})
}

//Err returns the error during the writer operations, if any
//otherwise nil
func (w *JSONWriter) Err() error {
	return w.err
}

//BeginArray starts an array
func (w *JSONWriter) BeginArray() {
	w.push(elmArray)
	w.writeString("[")
}

//EndArray finishes an array
func (w *JSONWriter) EndArray() {
	w.pop(elmArray)
	w.writeString("]")
}

//BeginObject starts an object
func (w *JSONWriter) BeginObject() {
	w.push(elmObject)
	w.writeString("{")
}

//EndObject finishes an object
func (w *JSONWriter) EndObject() {
	w.pop(elmObject)
	w.writeString("}")
}

func (w *JSONWriter) checkMember(key string) bool {
	if w.err != nil {
		return false
	}

	if len(w.elm) < 1 {
		w.err = fmt.Errorf("Trying to add member %q outside array/object", key)
		return false
	}
	return true
}

func (w *JSONWriter) maybeSeperator() {
	if w.err != nil {
		return
	}

	l := len(w.elm)
	if w.elm[l-1].count > 0 {
		w.writeString(", ")
	}
}

//MemberString adds a string member to an array or object
func (w *JSONWriter) MemberString(key string, format string, a ...interface{}) {
	if !w.checkMember(key) {
		return
	}

	w.maybeSeperator()
	w.writeString(fmt.Sprintf("%q: ", key))
	w.writeString(fmt.Sprintf(format, a...))
}

//MemberArray adds an array to an array or object
func (w *JSONWriters) MemberArray(key string) {

}

//MemberObject adds an object to an array or object
func (w *JSONWriter) MemberObject(key string) {

}
