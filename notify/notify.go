package notify

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const (
	win         = "windows"
	linux       = "linux"
	minLoopTime = 500
	exitErr     = "exit status 3"
)

type Notify struct {
	Root     string
	LoopTime int
	ModTimes map[string]time.Time
}

func Start(root string, loopTime int) error {
	if loopTime < minLoopTime {
		loopTime = minLoopTime
	}
	notify := &Notify{Root: root, LoopTime: loopTime}
	notify.ModTimes = make(map[string]time.Time)
	return notify.start()
}
func (n *Notify) reloadLoop() {
	for {
		filepath.Walk(n.Root, n.visit)
		time.Sleep(time.Duration(n.LoopTime) * time.Millisecond)
	}
}
func (n *Notify) visit(path string, fileinfo os.FileInfo, err error) error {
	if err != nil {
		return fmt.Errorf("访问文件失败%s", err)
	}
	if !fileinfo.IsDir() && strings.HasSuffix(path, "go") {
		modTime := fileinfo.ModTime()
		if oldModTime, ok := n.ModTimes[path]; !ok {
			n.ModTimes[path] = modTime
		} else {
			if oldModTime.Before(modTime) {
				fmt.Printf("%s文件发生变化，重新加载\n", filepath.Base(path))
				os.Exit(3)
			}
		}
	}
	return nil
}

func (n *Notify) start() error {
	if reloadEnv := os.Getenv("DOTWEB_RELOAD"); reloadEnv != "true" {
		stdErr := make(chan string)
		for {
			read, write, _ := os.Pipe()
			go func() {
				var buf bytes.Buffer
				io.Copy(&buf, read)
				stdErr <- buf.String()
			}()
			arg := []string{"run"}
			_, file := filepath.Split(os.Args[0])
			if runtime.GOOS == win {
				file = filepath.Base(file)
				ext := filepath.Ext(file)
				file = strings.TrimSuffix(file, ext)
			}
			file = file + ".go"
			arg = append(arg, file)
			arg = append(arg, os.Args[1:]...)
			command := exec.Command("go", arg...)
			command.Env = append(command.Env, "DOTWEB_RELOAD=true")
			command.Env = append(command.Env, os.Environ()...)
			command.Stdout = os.Stdout
			command.Stderr = write
			if err := command.Run(); err != nil {
				write.Close()
				if !strings.Contains(<-stdErr, exitErr) {
					fmt.Println(stdErr)
					return err
				}
			} else {
				return nil
			}

		}
	} else {
		go func() {
			n.reloadLoop()
		}()
	}
	return nil
}
