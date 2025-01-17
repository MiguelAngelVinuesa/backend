package main

import (
	"bufio"
	"os"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/prng.git/rng"
)

func main() {
	r := rng.NewRNG()
	buf := make([]uint8, 1<<12)
	w := bufio.NewWriterSize(os.Stdout, 64*1024)
	for {
		r.Read(buf)
		w.Write(buf)
	}
}
