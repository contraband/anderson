package integration_test

import (
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/types"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var _ = Describe("Subcommand", func() {
	Describe("ginkgo bootstrap", func() {
		It("should generate a bootstrap file, as long as one does not exist", func() {
			pkgPath := tmpPath("foo")
			os.Mkdir(pkgPath, 0777)
			session := startGinkgo(pkgPath, "bootstrap")
			Eventually(session).Should(gexec.Exit(0))
			output := session.Out.Contents()

			Ω(output).Should(ContainSubstring("foo_suite_test.go"))

			content, err := ioutil.ReadFile(filepath.Join(pkgPath, "foo_suite_test.go"))
			Ω(err).ShouldNot(HaveOccurred())
			Ω(content).Should(ContainSubstring("func TestFoo(t *testing.T) {"))
			Ω(content).Should(ContainSubstring("RegisterFailHandler"))
			Ω(content).Should(ContainSubstring("RunSpecs"))

			Ω(content).Should(ContainSubstring("\t" + `. "github.com/onsi/ginkgo"`))
			Ω(content).Should(ContainSubstring("\t" + `. "github.com/onsi/gomega"`))

			session = startGinkgo(pkgPath, "bootstrap")
			Eventually(session).Should(gexec.Exit(1))
			output = session.Out.Contents()
			Ω(output).Should(ContainSubstring("foo_suite_test.go already exists"))
		})

		It("should import nodot declarations when told to", func() {
			pkgPath := tmpPath("foo")
			os.Mkdir(pkgPath, 0777)
			session := startGinkgo(pkgPath, "bootstrap", "--nodot")
			Eventually(session).Should(gexec.Exit(0))
			output := session.Out.Contents()

			Ω(output).Should(ContainSubstring("foo_suite_test.go"))

			content, err := ioutil.ReadFile(filepath.Join(pkgPath, "foo_suite_test.go"))
			Ω(err).ShouldNot(HaveOccurred())
			Ω(content).Should(ContainSubstring("func TestFoo(t *testing.T) {"))
			Ω(content).Should(ContainSubstring("RegisterFailHandler"))
			Ω(content).Should(ContainSubstring("RunSpecs"))

			Ω(content).Should(ContainSubstring("var It = ginkgo.It"))
			Ω(content).Should(ContainSubstring("var Ω = gomega.Ω"))

			Ω(content).Should(ContainSubstring("\t" + `"github.com/onsi/ginkgo"`))
			Ω(content).Should(ContainSubstring("\t" + `"github.com/onsi/gomega"`))
		})
	})

	Describe("nodot", func() {
		It("should update the declarations in the bootstrap file", func() {
			pkgPath := tmpPath("foo")
			os.Mkdir(pkgPath, 0777)

			session := startGinkgo(pkgPath, "bootstrap", "--nodot")
			Eventually(session).Should(gexec.Exit(0))

			byteContent, err := ioutil.ReadFile(filepath.Join(pkgPath, "foo_suite_test.go"))
			Ω(err).ShouldNot(HaveOccurred())

			content := string(byteContent)
			content = strings.Replace(content, "var It =", "var MyIt =", -1)
			content = strings.Replace(content, "var Ω = gomega.Ω\n", "", -1)

			err = ioutil.WriteFile(filepath.Join(pkgPath, "foo_suite_test.go"), []byte(content), os.ModePerm)
			Ω(err).ShouldNot(HaveOccurred())

			session = startGinkgo(pkgPath, "nodot")
			Eventually(session).Should(gexec.Exit(0))

			byteContent, err = ioutil.ReadFile(filepath.Join(pkgPath, "foo_suite_test.go"))
			Ω(err).ShouldNot(HaveOccurred())

			Ω(byteContent).Should(ContainSubstring("var MyIt = ginkgo.It"))
			Ω(byteContent).ShouldNot(ContainSubstring("var It = ginkgo.It"))
			Ω(byteContent).Should(ContainSubstring("var Ω = gomega.Ω"))
		})
	})

	Describe("ginkgo generate", func() {
		var pkgPath string

		BeforeEach(func() {
			pkgPath = tmpPath("foo_bar")
			os.Mkdir(pkgPath, 0777)
		})

		Context("with no arguments", func() {
			It("should generate a test file named after the package", func() {
				session := startGinkgo(pkgPath, "generate")
				Eventually(session).Should(gexec.Exit(0))
				output := session.Out.Contents()

				Ω(output).Should(ContainSubstring("foo_bar_test.go"))

				content, err := ioutil.ReadFile(filepath.Join(pkgPath, "foo_bar_test.go"))
				Ω(err).ShouldNot(HaveOccurred())
				Ω(content).Should(ContainSubstring(`var _ = Describe("FooBar", func() {`))
				Ω(content).Should(ContainSubstring("\t" + `. "github.com/onsi/ginkgo"`))
				Ω(content).Should(ContainSubstring("\t" + `. "github.com/onsi/gomega"`))

				session = startGinkgo(pkgPath, "generate")
				Eventually(session).Should(gexec.Exit(1))
				output = session.Out.Contents()

				Ω(output).Should(ContainSubstring("foo_bar_test.go already exists"))
			})
		})

		Context("with an argument of the form: foo", func() {
			It("should generate a test file named after the argument", func() {
				session := startGinkgo(pkgPath, "generate", "baz_buzz")
				Eventually(session).Should(gexec.Exit(0))
				output := session.Out.Contents()

				Ω(output).Should(ContainSubstring("baz_buzz_test.go"))

				content, err := ioutil.ReadFile(filepath.Join(pkgPath, "baz_buzz_test.go"))
				Ω(err).ShouldNot(HaveOccurred())
				Ω(content).Should(ContainSubstring(`var _ = Describe("BazBuzz", func() {`))
			})
		})

		Context("with an argument of the form: foo.go", func() {
			It("should generate a test file named after the argument", func() {
				session := startGinkgo(pkgPath, "generate", "baz_buzz.go")
				Eventually(session).Should(gexec.Exit(0))
				output := session.Out.Contents()

				Ω(output).Should(ContainSubstring("baz_buzz_test.go"))

				content, err := ioutil.ReadFile(filepath.Join(pkgPath, "baz_buzz_test.go"))
				Ω(err).ShouldNot(HaveOccurred())
				Ω(content).Should(ContainSubstring(`var _ = Describe("BazBuzz", func() {`))

			})
		})

		Context("with an argument of the form: foo_test", func() {
			It("should generate a test file named after the argument", func() {
				session := startGinkgo(pkgPath, "generate", "baz_buzz_test")
				Eventually(session).Should(gexec.Exit(0))
				output := session.Out.Contents()

				Ω(output).Should(ContainSubstring("baz_buzz_test.go"))

				content, err := ioutil.ReadFile(filepath.Join(pkgPath, "baz_buzz_test.go"))
				Ω(err).ShouldNot(HaveOccurred())
				Ω(content).Should(ContainSubstring(`var _ = Describe("BazBuzz", func() {`))
			})
		})

		Context("with an argument of the form: foo_test.go", func() {
			It("should generate a test file named after the argument", func() {
				session := startGinkgo(pkgPath, "generate", "baz_buzz_test.go")
				Eventually(session).Should(gexec.Exit(0))
				output := session.Out.Contents()

				Ω(output).Should(ContainSubstring("baz_buzz_test.go"))

				content, err := ioutil.ReadFile(filepath.Join(pkgPath, "baz_buzz_test.go"))
				Ω(err).ShouldNot(HaveOccurred())
				Ω(content).Should(ContainSubstring(`var _ = Describe("BazBuzz", func() {`))
			})
		})

		Context("with multiple arguments", func() {
			It("should generate a test file named after the argument", func() {
				session := startGinkgo(pkgPath, "generate", "baz", "buzz")
				Eventually(session).Should(gexec.Exit(0))
				output := session.Out.Contents()

				Ω(output).Should(ContainSubstring("baz_test.go"))
				Ω(output).Should(ContainSubstring("buzz_test.go"))

				content, err := ioutil.ReadFile(filepath.Join(pkgPath, "baz_test.go"))
				Ω(err).ShouldNot(HaveOccurred())
				Ω(content).Should(ContainSubstring(`var _ = Describe("Baz", func() {`))

				content, err = ioutil.ReadFile(filepath.Join(pkgPath, "buzz_test.go"))
				Ω(err).ShouldNot(HaveOccurred())
				Ω(content).Should(ContainSubstring(`var _ = Describe("Buzz", func() {`))
			})
		})

		Context("with nodot", func() {
			It("should not import ginkgo or gomega", func() {
				session := startGinkgo(pkgPath, "generate", "--nodot")
				Eventually(session).Should(gexec.Exit(0))
				output := session.Out.Contents()

				Ω(output).Should(ContainSubstring("foo_bar_test.go"))

				content, err := ioutil.ReadFile(filepath.Join(pkgPath, "foo_bar_test.go"))
				Ω(err).ShouldNot(HaveOccurred())
				Ω(content).ShouldNot(ContainSubstring("\t" + `. "github.com/onsi/ginkgo"`))
				Ω(content).ShouldNot(ContainSubstring("\t" + `. "github.com/onsi/gomega"`))
			})
		})
	})

	Describe("ginkgo blur", func() {
		It("should unfocus tests", func() {
			pathToTest := tmpPath("focused")
			copyIn("focused_fixture", pathToTest)

			session := startGinkgo(pathToTest, "--noColor")
			Eventually(session).Should(gexec.Exit(types.GINKGO_FOCUS_EXIT_CODE))
			output := session.Out.Contents()

			Ω(output).Should(ContainSubstring("3 Passed"))
			Ω(output).Should(ContainSubstring("3 Skipped"))

			session = startGinkgo(pathToTest, "blur")
			Eventually(session).Should(gexec.Exit(0))

			session = startGinkgo(pathToTest, "--noColor")
			Eventually(session).Should(gexec.Exit(0))
			output = session.Out.Contents()
			Ω(output).Should(ContainSubstring("6 Passed"))
			Ω(output).Should(ContainSubstring("0 Skipped"))
		})
	})

	Describe("ginkgo version", func() {
		It("should print out the version info", func() {
			session := startGinkgo("", "version")
			Eventually(session).Should(gexec.Exit(0))
			output := session.Out.Contents()

			Ω(output).Should(MatchRegexp(`Ginkgo Version \d+\.\d+\.\d+`))
		})
	})

	Describe("ginkgo help", func() {
		It("should print out usage information", func() {
			session := startGinkgo("", "help")
			Eventually(session).Should(gexec.Exit(0))
			output := string(session.Err.Contents())

			Ω(output).Should(MatchRegexp(`Ginkgo Version \d+\.\d+\.\d+`))
			Ω(output).Should(ContainSubstring("ginkgo watch"))
			Ω(output).Should(ContainSubstring("-succinct"))
			Ω(output).Should(ContainSubstring("-nodes"))
			Ω(output).Should(ContainSubstring("ginkgo generate"))
		})
	})
})
