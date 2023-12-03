package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/otiai10/gosseract/v2"
)

var (
	Selector      = flag.String("s", "slurp", "screenshot area selection tool [slurp, slop, hacksaw, samurai-select, etc...]")
	ScreenShotter = flag.String("c", "grim", "screenshot tool [grim]")
)

func System(cmd string) error {
	fmt.Printf(cmd)
	c := exec.Command("sh", "-c", cmd)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	err := c.Run()

	if err != nil {
		return err
	}
	return nil
}

func Shot(dimensions string, output string) error {
	grimPath, err := exec.LookPath("grim")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Could not find grim")
		return err
	}

	cmd := fmt.Sprintf("%s -g %s %s", grimPath, dimensions, output)

	return System(cmd)
}

func Sel(sel string) (string, error) {
	if sel == "" {
		sel = "slurp"
	}

	slurp := exec.Command(sel, "-o")
	b, err := slurp.Output()
	if err != nil {
		return "", err
	}

	out := string(b)
	out = out[:len(out)-1]
	return out, nil
}

func main() {
	dim, err := Sel("slurp")
	if err != nil {
		fmt.Println(err)
	}
	dim = fmt.Sprintf("'%s'", dim)

	out := time.Now()
	output := fmt.Sprintf("%v-%v-%v_%v-%v.jpg",
		out.Month(),
		out.Day(),
		out.Year(),
		out.Hour(),
		out.Minute(),
	)

	err = Shot(dim, output)
	if err != nil {
		fmt.Println(err)
	}

	client := gosseract.NewClient()
	defer client.Close()
	client.SetImage(output)
	text, _ := client.Text()
	fmt.Println(text)
}
