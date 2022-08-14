package rendering

import (
	"os/exec"
)

func GraphVizCmd(sourcePath string, targetPath string, targetTyp *OutputFormat) *exec.Cmd {
	// dot -Tsvg -otarget.svg source.dot
	return exec.Command(
		"dot",
		"-T"+targetTyp.FlagValue,
		"-o"+targetPath,
		sourcePath)
}
