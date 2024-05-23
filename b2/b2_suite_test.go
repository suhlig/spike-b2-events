package b2_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestB2(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "B2 Suite")
}
