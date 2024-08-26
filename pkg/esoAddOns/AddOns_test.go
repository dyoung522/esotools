package esoAddOns_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/dyoung522/esotools/pkg/esoAddOns"
)

var _ = Describe("AddOns", func() {
	var myAddOn esoAddOns.AddOn
	var myAddOns esoAddOns.AddOns

	BeforeEach(func() {
		myAddOns = esoAddOns.AddOns{}

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

	Describe("ToKey()", func() {
		Context("when a string is provided", func() {
			It("should return a key", func() {
				Expect(esoAddOns.ToKey("Test AddOn ")).To(Equal("Test-AddOn"))
			})
		})
	})

	Describe("Add()", func() {
		Context("when a new AddOn is provided", func() {
			It("should add the AddOn to the list", func() {
				myAddOns.Add(myAddOn)
				Expect(myAddOns[myAddOn.Key()]).To(Equal(myAddOn))
			})
		})

		Context("when an existing AddOn is provided", func() {
			It("should not add the AddOn to the list", func() {
				myAddOns.Add(myAddOn)

				Expect(func() { myAddOns.Add(myAddOn) }).To(Panic())
				Expect(myAddOns).To(HaveLen(1))
			})
		})
	})

	Describe("Update()", func() {
		Context("when a new AddOn is provided", func() {
			It("should not add the AddOn to the list", func() {
				Expect(func() { myAddOns.Update(myAddOn) }).To(Panic())
				Expect(myAddOns).To(HaveLen(0))
			})
		})

		Context("when an existing AddOn is provided", func() {
			It("should update the AddOn already in the list", func() {
				myAddOns.Add(myAddOn)

				myAddOn.Title = "Updated Title"

				myAddOns.Update(myAddOn)
				Expect(myAddOns[myAddOn.Key()].Title).To(Equal("Updated Title"))
				Expect(myAddOns).To(HaveLen(1))
			})
		})
	})

})
