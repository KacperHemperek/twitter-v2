package dbmap

import (
	"encoding/json"
	"time"
)

type NilTime struct {
	Val   time.Time
	Empty bool
}

func (nt NilTime) MarshalJSON() ([]byte, error) {
	if nt.Empty {
		return json.Marshal(nil)
	}
	return json.Marshal(nt.Val)
}

func (nt *NilTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		nt.Val = time.Time{}
		nt.Empty = true
		return nil
	}
	nt.Empty = false
	return json.Unmarshal(data, &nt.Val)
}

type NilInt struct {
	Val   int
	Empty bool
}

func (ni NilInt) MarshalJSON() ([]byte, error) {
	if ni.Empty {
		return json.Marshal(nil)
	}
	return json.Marshal(ni.Val)
}

func (ni *NilInt) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		ni.Val = 0
		ni.Empty = true
		return nil
	}
	ni.Empty = false
	return json.Unmarshal(data, &ni.Val)
}

type NilString struct {
	Val   string
	Empty bool
}

func (ns *NilString) MarshalJSON() ([]byte, error) {
	if ns.Empty {
		return json.Marshal(nil)
	}
	return json.Marshal(ns.Val)
}

func (ns *NilString) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		ns.Val = ""
		ns.Empty = true
		return nil
	}
	ns.Empty = false
	return json.Unmarshal(b, &ns.Val)
}
