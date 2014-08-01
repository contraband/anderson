package integration_test

import (
	"fmt"
	"os/exec"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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
		andersonCommand.Dir = filepath.Join("src", "github.com", "xoebus", "prime")
		andersonCommand.Env = append(andersonCommand.Env, fmt.Sprintf("GOPATH=%s", filepath.Join("integration")))
	})

	It("runs", func() {
		session, err := gexec.Start(
			andersonCommand,
			gexec.NewPrefixedWriter("\x1b[32m[o]\x1b[95m[anderson]\x1b[0m ", GinkgoWriter),
			gexec.NewPrefixedWriter("\x1b[91m[e]\x1b[95m[anderson]\x1b[0m ", GinkgoWriter),
		)
		Ω(err).ShouldNot(HaveOccurred())

		Eventually(session).Should(gexec.Exit(0))
	})
})
