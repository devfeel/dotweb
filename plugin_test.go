package dotweb

import "fmt"

type testPlugin struct {
}

func (p *testPlugin) Name() string {
	return "test"
}
func (p *testPlugin) Run() error {
	fmt.Println(p.Name(), "runing")
	//panic("error test run")
	return nil
}
func (p *testPlugin) IsValidate() bool {
	return true
}
