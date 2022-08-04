package man

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/alecthomas/kong"
	mangokong "github.com/alecthomas/mango-kong"
	"github.com/muesli/roff"
)

// Man is a gum sub-command that generates man pages.
type Man struct {
	Output string `short:"o" type:"path" default:"/usr/share/man/man1/projects.1"`
}

const startedIn = 2022

func until(t time.Time) string {
	y := t.Year()
	if startedIn != y {
		return fmt.Sprintf("%d-%d", startedIn, y)
	} else {
		return fmt.Sprintf("%d", startedIn)
	}
}

// BeforeApply implements Kong BeforeApply hook.
func (m Man) BeforeApply(app *kong.Kong) error {
	man := mangokong.NewManPage(1, app.Model)
	man = man.WithSection(
		"Copyright",
		fmt.Sprintf("(c) %s Tim Ploß\nReleased under MIT license.", until(time.Now())),
	)
	out, err := m.outFile()
	if err != nil {
		return err
	}
	defer out.Close()
	fmt.Fprint(out, man.Build(roff.NewDocument()))
	app.Exit(0)
	return nil
}

func (m Man) outFile() (io.WriteCloser, error) {
	f, err := os.Create(m.Output)
	if err != nil {
		return nil, err
	}
	return f, nil
}
