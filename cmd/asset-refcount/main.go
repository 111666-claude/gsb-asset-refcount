// Command asset-refcount 跑资源引用计数样例。
package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"example.com/assets"
)

// Run 执行一次命令行调用，返回退出码。
func Run(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("asset-refcount", flag.ContinueOnError)
	flags.SetOutput(stderr)
	sample := flags.String("sample", "load", "load / release / budget / lru / scan")
	if err := flags.Parse(args); err != nil {
		return 2
	}
	switch *sample {
	case "load":
		store := assets.New(2)
		store.Load("hero")
		store.Load("hero")
		fmt.Fprintf(stdout, "loads=%d\n", store.Loads())
	case "release":
		store := assets.New(2)
		first := store.Load("hero")
		store.Load("hero")
		store.Release(first)
		store.Tick()
		fmt.Fprintf(stdout, "loaded=%v\n", store.Loaded("hero"))
	case "budget":
		store := assets.New(2)
		for index := 0; index < 5; index++ {
			handle := store.Load(fmt.Sprintf("a-%d", index))
			store.Release(handle)
		}
		unloaded := store.Tick()
		fmt.Fprintf(stdout, "unloaded=%d\n", unloaded)
	case "lru":
		store := assets.New(1)
		handleA := store.Load("a")
		handleB := store.Load("b")
		handleC := store.Load("c")
		store.Use("c", 10)
		store.Release(handleC)
		store.Release(handleA)
		store.Release(handleB)
		store.Tick()
		fmt.Fprintf(stdout, "loaded_a=%v loaded_c=%v\n", store.Loaded("a"), store.Loaded("c"))
	case "scan":
		store := assets.New(100)
		for index := 0; index < 3000; index++ {
			store.Load(fmt.Sprintf("a-%d", index))
		}
		for index := 0; index < 3000; index++ {
			store.RefCount("a-2999")
		}
		fmt.Fprintf(stdout, "scanned=%d\n", store.Scanned())
	default:
		fmt.Fprintln(stderr, "需要 --sample load|release|budget|lru|scan")
		return 2
	}
	return 0
}

func main() {
	os.Exit(Run(os.Args[1:], os.Stdout, os.Stderr))
}
