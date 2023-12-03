package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/otiai10/gosseract/v2"
)

var (
	Selector      = flag.String("s", "slurp", "screenshot area selection tool [slurp, slop, hacksaw, samurai-select, etc...]")
	ScreenShotter = flag.String("c", "grim", "screenshot tool [grim]")
)

func System(cmd string) error {
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

func Select() (string, error) {
	sel, err := exec.LookPath("slurp")
	if err != nil {
		return "", err
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

func cleanup(img string) error {
	img, err := filepath.Abs(img)
	if err != nil {
		return err
	}

	err = os.Remove(img)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	dim, err := Select()
	if err != nil {
		log.Fatal(err)
	}

	// screen dimensions
	dim = fmt.Sprintf("'%s'", dim)

	// output jpg
	out := time.Now()
	output := fmt.Sprintf("/tmp/%v-%v-%v_%v-%v.jpg",
		out.Month(),
		out.Day(),
		out.Year(),
		out.Hour(),
		out.Minute(),
	)

	go notifyLinux("tessa", "tessa", "screenshot taken", "/usr/share/icons/Kanagawa/categories/symbolic/appimagekit-cacher-symbolic.svg")
	err = Shot(dim, output)
	if err != nil {
		fmt.Println(err)
	}

	// tesseract
	client := gosseract.NewClient()
	defer client.Close()
	client.SetImage(output)

	text, _ := client.Text()
	fmt.Println(text)

	if err := cleanup(output); err != nil {
		log.Fatal(err)
	}
}
