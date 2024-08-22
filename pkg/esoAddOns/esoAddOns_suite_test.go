package esoAddOns_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestEsoAddOns(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "EsoAddOns Suite")
}
