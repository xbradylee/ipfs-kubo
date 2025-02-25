package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"testing"

	"github.com/xbradylee/ipfs-kubo/thirdparty/unit"

	config "github.com/xbradylee/ipfs-kubo/config"
	random "github.com/jbenet/go-random"
)

func main() {
	if err := compareResults(); err != nil {
		log.Fatal(err)
	}
}

func compareResults() error {
	var amount unit.Information
	for amount = 10 * unit.MB; amount > 0; amount = amount * 2 {
		if results, err := benchmarkAdd(int64(amount)); err != nil { // TODO compare
			return err
		} else {
			log.Println(amount, "\t", results)
		}
	}
	return nil
}

func benchmarkAdd(amount int64) (*testing.BenchmarkResult, error) {
	results := testing.Benchmark(func(b *testing.B) {
		b.SetBytes(amount)
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			tmpDir, err := os.MkdirTemp("", "")
			if err != nil {
				b.Fatal(err)
			}
			defer os.RemoveAll(tmpDir)

			env := append(os.Environ(), fmt.Sprintf("%s=%s", config.EnvDir, path.Join(tmpDir, config.DefaultPathName)))
			setupCmd := func(cmd *exec.Cmd) {
				cmd.Env = env
			}

			cmd := exec.Command("ipfs", "init", "-b=2048")
			setupCmd(cmd)
			if err := cmd.Run(); err != nil {
				b.Fatal(err)
			}

			const seed = 1
			f, err := os.CreateTemp("", "")
			if err != nil {
				b.Fatal(err)
			}
			defer os.Remove(f.Name())

			err = random.WritePseudoRandomBytes(amount, f, seed)
			if err != nil {
				b.Fatal(err)
			}
			if err := f.Close(); err != nil {
				b.Fatal(err)
			}

			b.StartTimer()
			cmd = exec.Command("ipfs", "add", f.Name())
			setupCmd(cmd)
			if err := cmd.Run(); err != nil {
				b.Fatal(err)
			}
			b.StopTimer()
		}
	})
	return &results, nil
}
