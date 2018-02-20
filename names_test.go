package main

import (
	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type NamesSuite struct{}

func (s *NamesSuite) TestX(t sweet.T) {
	pkg, err := parseDir("./testing/names")
	Expect(err).To(BeNil())

	v := newNameExtractor()
	walk(pkg, v)
	Expect(v.names).To(ConsistOf([]string{
		"OuterType",
		"StructType",
		"IntefaceType",
		"SimpleType",
	}))
}