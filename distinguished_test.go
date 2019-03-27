package dn

import (
	"strings"
	"testing"
)

// Test basic functions
type basetest struct {
	in       string
	shouldbe string
}

func Test1DN(t *testing.T) {

	tvals := []basetest{
		basetest{in: "CN=Ethel the Aardvark,OU=Australia,O=Marsupials,C=OZ",
			shouldbe: "CN=Ethel the Aardvark|OU=Australia|O=Marsupials|C=OZ"},
		basetest{in: "  CN=Ethel the Aardvark;  OU=Australia  ,O=Marsupials;C=OZ",
			shouldbe: "CN=Ethel the Aardvark|OU=Australia|O=Marsupials|C=OZ"},
		basetest{in: "Ethel the Aardvark,OU=Australia,O=Marsupials,C=OZ",
			shouldbe: "Ethel the Aardvark|OU=Australia|O=Marsupials|C=OZ"},
		basetest{in: ",OU=Australia,O= ,C=OZ",
			shouldbe: "OU=Australia|O=|C=OZ"},
		basetest{in: ",;,;,,",
			shouldbe: ""},
	}

	for i, test := range tvals {
		out := strings.Join(Canonical2Strings(test.in), "|")
		if out != test.shouldbe {
			t.Errorf("%d: Basic split of %s failed\nProduced: %s\nShould have been: %s",
				i, test.in, out, test.shouldbe)
		}
	}
}

func Test2DN(t *testing.T) {

	tvals := []basetest{
		basetest{in: "CN=Ethel the Aardvark,OU=Australia,O=Marsupials,C=OZ",
			shouldbe: "CN=Ethel the Aardvark; OU=Australia; O=Marsupials; C=OZ"},
		basetest{in: "  CN=  Ethel the Aardvark ;  OU  = Australia  ,O= Marsupials,C=OZ",
			shouldbe: "CN=Ethel the Aardvark; OU=Australia; O=Marsupials; C=OZ"},
		basetest{in: "Ethel the Aardvark,OU=Australia,O=Marsupials,C=OZ",
			shouldbe: "Ethel the Aardvark; OU=Australia; O=Marsupials; C=OZ"},
		basetest{in: ",OU=Australia,O= ,C=OZ,=Sunny",
			shouldbe: "OU=Australia; C=OZ"},
		basetest{in: ",;,;,,",
			shouldbe: ""},
		basetest{in: "",
			shouldbe: ""},
	}

	// This time test complete
	for i, test := range tvals {
		dn := New(test.in)
		out := dn.String()
		if out != test.shouldbe {
			t.Errorf("%d: Complete parse of %s failed\nProduced: %s\nShould have been: %s",
				i, test.in, out, test.shouldbe)
		}
	}
}

// Test GetValues
type gvtest struct {
	in       string
	key      string
	shouldbe string
}

func TestGetValues(t *testing.T) {

	tvals := []gvtest{
		gvtest{in: "CN=Ethel the Aardvark,OU=Australia,O=Marsupials,C=OZ", key: "OU",
			shouldbe: "Australia"},
		gvtest{in: "CN=Ethel the Aardvark; OU=Four Feets; OU=Australia; OU=Cuteness; C=OZ", key: "OU",
			shouldbe: "Four Feets|Australia|Marsupials|Cuteness"},
		gvtest{in: "CN=Ethel the Aardvark,job=Quantity Surveying", key: "job",
			shouldbe: "Quantity Surveying"},
		gvtest{in: "CN=Ethel the Aardvark,job=Quantity Surveying,email=ethelaard@gmail.com", key: "email",
			shouldbe: "ethelaard@gmail.com"},
	}

	for i, test := range tvals {
		dn := New(test.in)
		out := strings.Join(dn.GetValues(test.key), "|")
		if out != test.shouldbe {
			t.Errorf("%d: Complete parse of %s failed\nProduced: %s\nShould have been: %s",
				i, test.in, out, test.shouldbe)
		}
	}

}

// Test CommonName
type cntest struct {
	in          string
	shouldbe    string
	shouldbeerr bool
}

func TestCommonName(t *testing.T) {

	tvals := []cntest{
		cntest{in: "CN=Ethel the Aardvark; OU=Australia; O=Marsupials; C=OZ",
			shouldbe: "Ethel the Aardvark", shouldbeerr: false},
		cntest{in: "CN=Ethel the Aardvark; OU=Australia; O=Marsupials; C=OZ",
			shouldbe: "Ethel the Aardvark", shouldbeerr: false},
	}

	for i, test := range tvals {
		dn := New(test.in)
		out, err := dn.CommonName()
		if err != nil {
			if !test.shouldbeerr {
				t.Errorf("%d: CommonName of %s produced false error\nError: %s\n",
					i, test.in, err.Error())
			} else {
				continue // Correctly produced an error
			}
		} else {
			if test.shouldbeerr {
				t.Errorf("%d: CommonName of %s didn't produce an error, when it should\nProduced: %s\n",
					i, test.in, out)
			} else {
				if out != test.shouldbe {
					t.Errorf("%d: CommonName of %s failed\nProduced: %s\nShould have been: %s",
						i, test.in, out, test.shouldbe)
				}
			}
		}
	}
}
