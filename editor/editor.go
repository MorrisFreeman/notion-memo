package editor

import (
	"io/ioutil"
	"os"
	"os/exec"
)

func ReadEditor(name string) ([]byte, error) {
	tf, err := ioutil.TempFile("", name)
	if err != nil {
		return nil, err
	}
	defer os.Remove(tf.Name())

	cmd := exec.Command("sh", "-c", "editor "+tf.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadFile(tf.Name())
	if err != nil {
		return nil, err
	}

	return b, nil
}
