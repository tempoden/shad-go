package commands

import (
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "export source code to public repo",
	Run:   exportCode,
}

var (
	flagPush bool
)

func init() {
	rootCmd.AddCommand(exportCmd)

	exportCmd.Flags().BoolVar(&flagPush, "push", false, "push to public repo")
}

func git(args ...string) {
	log.Println("git", strings.Join(args, " "))

	cmd := exec.Command("git", args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

func exportCode(cmd *cobra.Command, args []string) {
	git("checkout", "-b", "temp")
	git("reset", "public")

	privateFiles := listPrivateFiles(".")
	for _, f := range privateFiles {
		log.Printf("rm %s", f)
		if err := os.Remove(f); err != nil {
			log.Fatal(err)
		}
	}

	git("checkout", "public")
	git("branch", "-D", "temp")

	git("add", "-A")
	git("commit", "-m", "export public files", "--allow-empty")

	if flagPush {
		git("push", "public", "public:master")
	}

	git("checkout", "master")
}
