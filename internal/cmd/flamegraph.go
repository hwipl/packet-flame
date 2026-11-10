package cmd

import (
	"fmt"
	"io"
	"log"
	"os/exec"
	"strings"
)

// createFlamegraph creates a flamegraph from all collected counters and writes
// it to the Writer w
func createFlamegraph(w io.Writer) {
	// generate input string for flamegraph tool
	var builder strings.Builder
	for l, c := range counter.getAll() {
		fmt.Fprintln(&builder, l, c)
	}
	input := builder.String()

	// generate flamegraph
	cmd := exec.Command(flamegraphExe)
	if flamegraphArgs != nil {
		cmd = exec.Command(flamegraphExe, flamegraphArgs...)
	}
	cmd.Stdin = strings.NewReader(input)
	cmd.Stdout = w
	err := cmd.Run()
	if err != nil {
		log.Println(err)
		return
	}
}
