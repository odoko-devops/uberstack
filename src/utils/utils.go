package utils

import (
  "log"
  "fmt"
  "io"
  "bufio"
  "os/exec"
)


/***********************************************************************
 * Check and quit on errors
 */
func Check(err error) {
  if err!=nil {
    log.Fatal(err);
    panic(err)
  }
}

func Required(value string, message string) {
  if value == "" {
    panic(message)
  }
}

/***********************************************************************
 * Execute a command, with streamed output for slow running commands
 */
func watchOutputStream(typ string, r bufio.Reader) {
    for {
      line, _, err := r.ReadLine()
      if err == io.EOF {
        break
      }
      fmt.Printf("%s: %s\n", typ, line)
    }
}
func Execute(cmd *exec.Cmd) {
  stdout, _ := cmd.StdoutPipe()
  stderr, _ := cmd.StderrPipe()
  cmd.Start()
  stdoutReader := bufio.NewReader(stdout)
  stderrReader := bufio.NewReader(stderr)
  go watchOutputStream("stdout", *stdoutReader)
  go watchOutputStream("stderr", *stderrReader)
  cmd.Wait()
}

func ExecuteAndRetrieve(cmd *exec.Cmd) string {
  output, err := cmd.Output()
  Check(err)
  return string(output)
}