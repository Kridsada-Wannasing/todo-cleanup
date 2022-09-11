package todo

import (
	"testing"
)

func TestCreateTodoNotAllowSleepTask(t *testing.T) {
	handler := NewTodoHandler(&TestDB{})
	c := &TestContext{}

	handler.NewTask(c)

	want := "not allowed"

	if want != c.v["error"] {
		t.Errorf("want %s but get %s\n", want, c.v["error"])
	}
}

type TestDB struct{}

func (TestDB) New(*Todo) error {
	return nil
}

type TestContext struct {
	v map[string]interface{}
}

func (TestContext) Bind(v interface{}) error {
	*v.(*Todo) = Todo{Title: "sleep"}
	return nil
}

func (c *TestContext) JSON(code int, v interface{}) {
	c.v = v.(map[string]interface{})
}

func (TestContext) TransactionID() string {
	return "TestTransactionID"
}

func (TestContext) Audience() string {
	return "Unit Test"
}
