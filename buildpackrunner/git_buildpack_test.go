package buildpackrunner_test

import (
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	. "github.com/cloudfoundry-incubator/buildpack_app_lifecycle/buildpackrunner"

	. "github.com/cloudfoundry-incubator/buildpack_app_lifecycle/Godeps/_workspace/src/github.com/onsi/ginkgo"
	. "github.com/cloudfoundry-incubator/buildpack_app_lifecycle/Godeps/_workspace/src/github.com/onsi/gomega"
)

var _ = Describe("GitBuildpack", func() {

	Describe("Clone", func() {
		var cloneTarget string
		BeforeEach(func() {
			var err error
			cloneTarget, err = ioutil.TempDir(tmpDir, "clone")
			Ω(err).ShouldNot(HaveOccurred())
		})

		AfterEach(func() {
			os.RemoveAll(cloneTarget)
		})

		It("clones a URL", func() {
			err := GitClone(gitUrl, cloneTarget)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(currentBranch(cloneTarget)).Should(Equal("master"))
		})

		It("clones a URL with a branch", func() {
			branchUrl := gitUrl
			branchUrl.Fragment = "a_branch"
			err := GitClone(branchUrl, cloneTarget)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(currentBranch(cloneTarget)).Should(Equal("a_branch"))
		})

		It("clones a URL with a lightweight tag", func() {
			branchUrl := gitUrl
			branchUrl.Fragment = "a_lightweight_tag"
			err := GitClone(branchUrl, cloneTarget)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(currentBranch(cloneTarget)).Should(Equal("a_lightweight_tag"))
		})

		Context("with bogus git URLs", func() {
			It("returns an error", func() {
				By("passing an invalid path", func() {
					badUrl := gitUrl
					badUrl.Path = "/a/bad/path"
					err := GitClone(badUrl, cloneTarget)
					Ω(err).Should(HaveOccurred())
				})

				By("passing a bad tag/branch", func() {
					badUrl := gitUrl
					badUrl.Fragment = "notfound"
					err := GitClone(badUrl, cloneTarget)
					Ω(err).Should(HaveOccurred())
				})
			})
		})

	})
})

func currentBranch(gitDir string) string {
	cmd := exec.Command("git", "symbolic-ref", "--short", "-q", "HEAD")
	cmd.Dir = gitDir
	bytes, err := cmd.Output()
	if err != nil {
		// try the tag
		cmd := exec.Command("git", "name-rev", "--name-only", "--tags", "HEAD")
		cmd.Dir = gitDir
		bytes, err = cmd.Output()
	}
	Ω(err).ShouldNot(HaveOccurred())
	return strings.TrimSpace(string(bytes))
}