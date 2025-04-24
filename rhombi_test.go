package rhombifer

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/racg0092/rhombifer/pkg/models"
)

func TestRhombi(t *testing.T) {
	//NOTE: imitating os args
	input := "watch -i 2 ls -all"
	os.Args = strings.Split(input, " ")

	root := Root()
	root.Run = func(args ...string) error {
		f, _ := FindFlag("i")
		fmt.Println(f.Values)
		fmt.Println(args)
		return nil
	}

	root.Name = "watch"
	c := GetConfig()
	c.AllowFlagsInRoot = true
	root.AddFlags(&models.Flag{Name: "interval", ShortFormat: "i", SingleValue: false})

	if e := Start(); e != nil {
		t.Error(e)
	}

}
