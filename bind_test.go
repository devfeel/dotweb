package dotweb

import (
	"github.com/devfeel/dotweb/test"

	"testing"
)

type Person struct {
	Hair     string
	HasGlass bool
	Age      int
	Legs     []string
}


//json
func TestBinder_Bind_json(t *testing.T) {

	binder := newBinder()

	if binder == nil {
		t.Error("binder can not be nil!")
	}

	//init DotServer
	app := New()

	if app == nil {
		t.Error("app can not be nil!")
	}

	//expected
	expected := &Person{
		Hair:     "Brown",
		HasGlass: true,
		Age:      10,
		Legs:     []string{"Left", "Right"},
	}

	//init param
	param := &InitContextParam{
		t,
		expected,
		"application/json",
		test.ToJson,
	}

	//init param
	context := initContext(param)
	//actual
	person := &Person{}

	err := binder.Bind(person, context)

	//check error must nil
	test.Nil(t, err)

	//check expected
	test.Equal(t, expected, person)

	t.Log("person:", person)
	t.Log("expected:", expected)

}

//json
func TestBinder_Bind_json_error(t *testing.T) {

	binder := newBinder()

	if binder == nil {
		t.Error("binder can not be nil!")
	}

	//init DotServer
	app := New()

	if app == nil {
		t.Error("app can not be nil!")
	}

	//expected
	expected := &Person{
		Hair:     "Brown",
		HasGlass: true,
		Age:      10,
		Legs:     []string{"Left", "Right"},
	}

	//init param
	param := &InitContextParam{
		t,
		expected,
		"application/xml",
		test.ToJson,
	}

	//init param
	context := initContext(param)
	//actual
	person := &Person{}

	err := binder.Bind(person, context)

	//check error must not nil
	test.NotNil(t, err)
}

//xml
func TestBinder_Bind_xml(t *testing.T) {

	binder := newBinder()

	if binder == nil {
		t.Error("binder can not be nil!")
	}

	//init DotServer
	app := New()

	if app == nil {
		t.Error("app can not be nil!")
	}

	//expected
	expected := &Person{
		Hair:     "Brown",
		HasGlass: true,
		Age:      10,
		Legs:     []string{"Left", "Right"},
	}
	param := &InitContextParam{
		t,
		expected,
		"application/xml",
		test.ToXML,
	}

	//init param
	context := initContext(param)
	//actual
	person := &Person{}

	err := binder.Bind(person, context)

	//check error must nil
	test.Nil(t, err)

	//check expected
	test.Equal(t, expected, person)

	t.Log("person:", person)
	t.Log("expected:", expected)

}

//xml
func TestBinder_Bind_xml_error(t *testing.T) {

	binder := newBinder()

	if binder == nil {
		t.Error("binder can not be nil!")
	}

	//init DotServer
	app := New()

	if app == nil {
		t.Error("app can not be nil!")
	}

	//expected
	expected := &Person{
		Hair:     "Brown",
		HasGlass: true,
		Age:      10,
		Legs:     []string{"Left", "Right"},
	}
	param := &InitContextParam{
		t,
		expected,
		"application/json",
		test.ToXML,
	}

	//init param
	context := initContext(param)
	//actual
	person := &Person{}

	err := binder.Bind(person, context)

	//check error must not nil
	test.NotNil(t, err)
}

//else
func TestBinder_Bind_default(t *testing.T) {

	binder := newBinder()

	if binder == nil {
		t.Error("binder can not be nil!")
	}

	//init DotServer
	app := New()

	if app == nil {
		t.Error("app can not be nil!")
	}

	//expected
	expected := &Person{
		Hair:     "Brown",
		HasGlass: true,
		Age:      10,
		Legs:     []string{"Left", "Right"},
	}
	param := &InitContextParam{
		t,
		expected,
		"",
		test.ToDefault,
	}

	//init param
	context := initContext(param)

	form := make(map[string][]string)
	form["Hair"] = []string{"Brown"}
	form["HasGlass"] = []string{"true"}
	form["Age"] = []string{"10"}
	form["Legs"] = []string{"Left", "Right"}

	context.request.Form = form
	//actual
	person := &Person{}

	err := binder.Bind(person, context)

	//check error must nil
	test.Nil(t, err)

	//check expected
	test.Equal(t, expected, person)

	t.Log("person:", person)
	t.Log("expected:", expected)

}

//else
func TestBinder_Bind_default_error(t *testing.T) {

	binder := newBinder()

	if binder == nil {
		t.Error("binder can not be nil!")
	}

	//init DotServer
	app := New()

	if app == nil {
		t.Error("app can not be nil!")
	}

	//expected
	expected := &Person{
		Hair:     "Brown",
		HasGlass: true,
		Age:      10,
		Legs:     []string{"Left", "Right"},
	}
	param := &InitContextParam{
		t,
		expected,
		"application/xml",
		test.ToDefault,
	}

	//init param
	context := initContext(param)

	form := make(map[string][]string)
	form["Hair"] = []string{"Brown"}
	form["HasGlass"] = []string{"true"}
	form["Age"] = []string{"10"}
	form["Legs"] = []string{"Left", "Right"}

	context.request.Form = form
	//actual
	person := &Person{}

	err := binder.Bind(person, context)

	//check error must not nil
	test.NotNil(t, err)

}

//default
//TODO:content type is null but body not null,is it right??
func TestBinder_Bind_ContentTypeNull(t *testing.T) {

	binder := newBinder()

	if binder == nil {
		t.Error("binder can not be nil!")
	}

	//init DotServer
	app := New()

	if app == nil {
		t.Error("app can not be nil!")
	}

	//expected
	expected := &Person{
		Hair:     "Brown",
		HasGlass: true,
		Age:      10,
		Legs:     []string{"Left", "Right"},
	}
	param := &InitContextParam{
		t,
		expected,
		"",
		test.ToXML,
	}

	//init param
	context := initContext(param)
	//actual
	person := &Person{}

	err := binder.Bind(person, context)

	//check error must nil?
	test.Nil(t, err)
}
