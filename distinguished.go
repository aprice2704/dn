package dn

import (
	"fmt"
	"strings"
)

// Utilities for dealing with OSI X.500 distinguished names in the simplest version of the
// Canonical format defined here https://tools.ietf.org/html/rfc1779; that is to say,
// things like: CN=Ethel the Aardvark, OU=AntBGone Dept
// #TODO more complete support of syntax with https://github.com/alecthomas/participle or some such.

// AJP for T-Mobile via UST Global 25-Mar-2019

// DN represents a canonical distinguished name in a way that facilitates examining and manipulating it
type DN struct {
	orig string
	//	parsed  bool
	parts   []Part
	partMap map[string][]string
}

// Part is a name component pair e.g. "CN", "Ethel the Aardvark"
type Part struct {
	Key   string
	Value string
}

func (p Part) String() string {
	s := ""
	if len(p.Key) > 0 {
		s = s + p.Key + "=" + p.Value
		return s
	}
	s = s + p.Value
	return s
}

// Parts is a DN split into components
type Parts []Part

// PartSeparator is put between parts during reconstruction of a string
// Change this to ", ", or ";" etc. as you wish.
var PartSeparator = "; "

// Standardised Keywords
const (
	CommonNameKey          = "CN"
	LocalityNameKey        = "L"
	StateOrProvinceNameKey = "ST"
	OrganizationNameKey    = "O"
	OrganizationalUnitKey  = "OU"
	CountryNameKey         = "C"
	StreetAddressKey       = "STREET"
)

// String rebuilds a DN to a string from parts
func (d DN) String() string {
	s := ""
	for _, p := range d.parts {
		s = s + p.String() + PartSeparator
	}
	if len(s) > 2 {
		return s[0 : len(s)-2]
	}
	return s
}

// GetValues retrieves all the values for a given key, no errors
func (d DN) GetValues(key string) []string {
	return d.partMap[key]
}

// CommonName fetches CommonName value, with errors
func (d DN) CommonName() (string, error) {
	vals := d.GetValues(CommonNameKey)
	if len(vals) < 1 {
		return "", fmt.Errorf("no CommonName elements found")
	}
	if len(vals) > 1 {
		return "", fmt.Errorf("more than one CommonName element found")
	}
	return vals[0], nil
}

// New makes a new DN object from a string
func New(s string) (d *DN) {
	d = &DN{orig: s}
	d.parts = Canonical2Parts(s)
	d.partMap = make(map[string][]string)
	for _, v := range d.parts {
		_, there := d.partMap[v.Key]
		if there {
			d.partMap[v.Key] = append(d.partMap[v.Key], v.Value)
		} else {
			d.partMap[v.Key] = make([]string, 0)
		}
	}
	return d
}

// Canonical2Parts converts a full DN into its key, value pairs
func Canonical2Parts(can string) (parts Parts) {
	ps := Canonical2Strings(can)
	newp := Part{}
	for _, v := range ps {
		sides := strings.Split(v, "=")
		if len(sides) > 0 {
			for i, s := range sides {
				sides[i] = strings.TrimSpace(s)
			}
			if len(sides) == 1 { // there is only a value
				newp = Part{Value: sides[0]}
			} else { // both key and value exist
				if len(sides[0]) > 0 { // is the key non-blank?
					newp = Part{Key: sides[0], Value: sides[1]} // success! we have both
				} else {
					newp = Part{} // blank, but present, key -> nada
				}
			}
			if len(newp.Value) > 0 { // ignore those with blank values
				parts = append(parts, newp)
			}
		}
	}
	return parts
}

// Canonical2Strings converts a canonical format DN to a sequence of string pieces as separated by , or ;
//   with leading and trailing space trimmed off. e.g. "CN=Ethel the Aardvark"
func Canonical2Strings(can string) (parts []string) {
	if len(can) == 0 {
		return parts
	}
	f := func(c rune) bool {
		return (c == ',' || c == ';')
	}
	parts = strings.FieldsFunc(can, f)
	for i, s := range parts {
		parts[i] = strings.TrimSpace(s)
	}
	return parts
}
