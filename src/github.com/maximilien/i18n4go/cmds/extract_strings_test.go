package cmds

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("extract_strings cmd", func() {

	Describe("NewExtractStrings", func() {
		It("create new ExtractStrings object with params", func() {
			//test
		})
	})
	Describe("InspectFile", func() {
		It("should extract strings from file", func() {
			//test
		})
	})
	Describe("InspectDir", func() {
		Context("when recursive flag is passed (true)", func() {
			It("should extract strings from directory and recursively in subdirs", func() {
				//test
			})
		})
		Context("when recursive flag is false or not passed", func() {
			It("should extract strings from directory only", func() {
				//test
			})
		})
		Context("when ignore dir regex is passed", func() {
			It("should extract strings from directory that are not ignored", func() {
				//test
			})
		})
	})
	Describe("Println & Printf", func() {
		Context("when -v is set", func() {
			It("should print to stdout if verbose is set", func() {
				//test
			})
		})
		Context("when -v is NOT set", func() {
			It("should ignore all strings passed", func() {
				//test
			})
		})
	})
})
