package integration_test

import (
	"fmt"
	"os/exec"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"

	"testing"
)

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Integration Suite")
}

var andersonPath string

var _ = BeforeSuite(func() {
	var err error
	andersonPath, err = gexec.Build("github.com/xoebus/anderson")
	Ω(err).ShouldNot(HaveOccurred())
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
})

var _ = Describe("Anderson", func() {
	var andersonCommand *exec.Cmd

	BeforeEach(func() {
		andersonCommand = exec.Command(andersonPath)
		andersonCommand.Dir = filepath.Join("_ignore", "src", "github.com", "xoebus", "prime")

		gopath, err := filepath.Abs(".")
		Ω(err).ShouldNot(HaveOccurred())
		andersonCommand.Env = append(andersonCommand.Env, fmt.Sprintf("GOPATH=%s", gopath))
	})

	It("does some cheesy dredd scene-setting", func() {
		session, err := gexec.Start(
			andersonCommand,
			gexec.NewPrefixedWriter("\x1b[32m[o]\x1b[95m[anderson]\x1b[0m ", GinkgoWriter),
			gexec.NewPrefixedWriter("\x1b[91m[e]\x1b[95m[anderson]\x1b[0m ", GinkgoWriter),
		)
		Ω(err).ShouldNot(HaveOccurred())

		Eventually(session).Should(gbytes.Say("Hold still citizen, scanning dependencies for contraband..."))
		Eventually(session).Should(gexec.Exit(0))
	})

	It("shows whitelisted licenses as 'CHECKS OUT'", func() {
		session, err := gexec.Start(
			andersonCommand,
			gexec.NewPrefixedWriter("\x1b[32m[o]\x1b[95m[anderson]\x1b[0m ", GinkgoWriter),
			gexec.NewPrefixedWriter("\x1b[91m[e]\x1b[95m[anderson]\x1b[0m ", GinkgoWriter),
		)
		Ω(err).ShouldNot(HaveOccurred())

		Eventually(session).Should(gbytes.Say("github.com/xoebus/whitelist.*CHECKS OUT"))
		Eventually(session).Should(gexec.Exit(0))
	})

	It("shows blacklisted licenses as 'CONTRABAND'", func() {
		session, err := gexec.Start(
			andersonCommand,
			gexec.NewPrefixedWriter("\x1b[32m[o]\x1b[95m[anderson]\x1b[0m ", GinkgoWriter),
			gexec.NewPrefixedWriter("\x1b[91m[e]\x1b[95m[anderson]\x1b[0m ", GinkgoWriter),
		)
		Ω(err).ShouldNot(HaveOccurred())

		Eventually(session).Should(gbytes.Say("github.com/xoebus/blacklist.*CONTRABAND"))
		Eventually(session).Should(gexec.Exit(0))
	})

	It("shows projects with no license as 'NO LICENSE'", func() {
		session, err := gexec.Start(
			andersonCommand,
			gexec.NewPrefixedWriter("\x1b[32m[o]\x1b[95m[anderson]\x1b[0m ", GinkgoWriter),
			gexec.NewPrefixedWriter("\x1b[91m[e]\x1b[95m[anderson]\x1b[0m ", GinkgoWriter),
		)
		Ω(err).ShouldNot(HaveOccurred())

		Eventually(session).Should(gbytes.Say("github.com/xoebus/no-license.*NO LICENSE"))
		Eventually(session).Should(gexec.Exit(0))
	})
})
