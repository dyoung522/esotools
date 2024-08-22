package esoAddOns_test

import (
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/dyoung522/esotools/pkg/esoAddOns"
)

var _ = Describe("AddOn", func() {
	var myAddOn, myInvalidAddOn esoAddOns.AddOn

	BeforeEach(func() {
		myInvalidAddOn, _ = esoAddOns.NewAddOn("InvalidAddon")

		myAddOn, _ = esoAddOns.NewAddOn("ValidAddon")
		myAddOn.Title = "Test AddOn"
		myAddOn.Author = "Jimmy Tester"
		myAddOn.Contributors = "Test Contributors"
		myAddOn.Version = "0.0.1-rc1"
		myAddOn.Description = "This is a longer description of Test Addon."
		myAddOn.AddOnVersion = "867"
		myAddOn.APIVersion = "5309"
		myAddOn.SavedVariables = []string{"MyAddOnSavedVariables", "SomeOtherSavedVariables"}
		myAddOn.DependsOn = []string{"myOtherAddOn"}
		myAddOn.OptionalDependsOn = []string{"Test OptionalDependsOn"}
	})

	Describe("NewAddOn()", func() {
		Context("when a key is provided", func() {
			It("should return a new AddOn", func() {
				Expect(esoAddOns.NewAddOn("foo")).To(BeAssignableToTypeOf(esoAddOns.AddOn{}))
			})

			It("should return an AddOn with the correct key", func() {
				Expect(myAddOn.Key()).To(Equal("ValidAddon"))
				Expect(myInvalidAddOn.Key()).To(Equal("InvalidAddon"))
			})
		})

		Context("when a key is not provided", func() {
			It("should return an Error", func() {
				_, err := esoAddOns.NewAddOn("")

				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("TitleString()", func() {
		It("should return the correct title string", func() {
			Expect(myAddOn.TitleString()).To(Equal("Test AddOn (v0.0.1-rc1) by Jimmy Tester"))
		})
	})

	Describe("Dir() & IsSubmodule()", func() {
		var myOtherAddOn = myAddOn

		BeforeEach(func() {
			myAddOn.SetDir(filepath.Join("Test"))
			myAddOn.DependsOn = []string{"myOtherAddOn"}

			myOtherAddOn.SetDir(filepath.Join("MyAddOn", "Test"))
		})

		Context("when the directory is set", func() {
			It("should return the correct directory", func() {
				Expect(myAddOn.Dir()).To(Equal("Test"))
			})
		})

		Context("when the directory is not set", func() {
			It("should return an empty string", func() {
				Expect(myInvalidAddOn.Dir()).To(Equal(""))
			})
		})

		Context("when the directory has multiple paths", func() {
			It("should also be a sub-module", func() {
				Expect(myOtherAddOn.IsSubmodule()).To(BeTrue())
			})
		})

		Context("when the directory has a single path", func() {
			It("should not be a sub-module", func() {
				Expect(myAddOn.IsSubmodule()).To(BeFalse())
			})
		})
	})

	Describe("IsDependency()", func() {
		Context("when meta.dependency is set", func() {
			It("should return true", func() {
				testAddOn := myAddOn
				testAddOn.SetDependency(true)
				Expect(testAddOn.IsDependency()).To(BeTrue())
			})
		})

		Context("when meta.dependency is not set", func() {
			It("should return false", func() {
				Expect(myAddOn.IsDependency()).To(BeFalse())
			})
		})
	})

	Describe("IsLibrary()", func() {
		Context("when meta.library is set", func() {
			It("should return true", func() {
				testAddOn := myAddOn
				testAddOn.SetLibrary(true)
				Expect(testAddOn.IsLibrary()).To(BeTrue())
			})
		})

		Context("when meta.library is not set", func() {
			It("should return false", func() {
				Expect(myAddOn.IsLibrary()).To(BeFalse())
			})
		})
	})

	Describe("Validate()", func() {
		Context("when all required fields are set", func() {
			It("should return true", func() {
				Expect(myAddOn.Validate()).To(BeTrue())
			})
		})

		Context("when a required field is not set", func() {
			It("should fill in missing fields", func() {
				Expect(myInvalidAddOn.Validate()).To(BeTrue())
				Expect(myInvalidAddOn.Title).To(Equal("Invalid Addon"))
				Expect(myInvalidAddOn.Author).To(Equal("Unknown"))
				Expect(myInvalidAddOn.Version).To(Equal("0"))
			})
		})
	})
})
