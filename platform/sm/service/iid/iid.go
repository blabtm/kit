package iid

import (
	"fmt"
	"strings"
	"time"
)

const Layout = "20060102_150405"

type IID struct {
	Name string
	Time time.Time
}

func New(name string) *IID {
	return &IID{
		Name: name,
		Time: time.Now().Local(),
	}
}

func Parse(name string, short string) *IID {
	t, err := time.Parse(Layout, short)

	if err != nil {
		panic(err)
	}

	return &IID{
		Name: name,
		Time: t,
	}
}

func (iid *IID) Short() string {
	return iid.Time.Format(Layout)
}

func (iid *IID) Full() string {
	ns := strings.ReplaceAll(iid.Name, ".", "_")
	return fmt.Sprintf("%s_%s", ns, iid.Time.Format(Layout))
}
