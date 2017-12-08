package tests_test

import (
	"flag"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var myFlag string

func init() {
	// cmd: ginkgo -v src/tests/ -- -myFlag="flag text"
	flag.StringVar(&myFlag, "myFlag", "default", "myFlag is used to control my behavior")
}

// cmd: ginkgo -v --focus="demo01" src/tests/
var _ = Describe("TestDemo01", func() {
	var myText string

	BeforeSuite(func() {
		GinkgoWriter.Write([]byte("TEST: exec BeforeSuite\n"))
		By("my flag value: " + myFlag) // get external var
	})

	AfterSuite(func() {
		GinkgoWriter.Write([]byte("TEST: exec AfterSuite\n"))
	})

	BeforeEach(func() {
		GinkgoWriter.Write([]byte("TEST: exec BeforeEach\n"))
		myText = "test"
	})

	AfterEach(func() {
		GinkgoWriter.Write([]byte("TEST: exec AfterEach\n"))
	})

	JustBeforeEach(func() {
		GinkgoWriter.Write([]byte("TEST: exec JustBeforeEach\n"))
	})

	Describe("Test string", func() {
		Context("Test context", func() {
			It("[demo01] text is not null", func() {
				GinkgoWriter.Write([]byte("TEST: run test01\n"))
				By("sub step description")
				Expect(myText != "").Should(BeTrue(), "Failed, not null")
			})
		})

		Context("Test context", func() {
			It("[demo01] text length should be 4", func() {
				GinkgoWriter.Write([]byte("TEST: run test02\n"))
				Expect(len(myText)).To(Equal(4), "Failed, text length = 4")
			})
		})

		Context("Test context", func() {
			It("[demo01] Marking Specs as Failed", func() {
				By("TEST: run test03")
				Fail("Mark failed")
			})
		})
	})
})
