package visor

import (
	"testing"
)

func appSetup(name string) (c *Client, app *App) {
	app, err := NewApp(name, "git://cat.git", "whiskers")
	if err != nil {
		panic(err)
	}
	c, err = Dial(DEFAULT_ADDR, DEFAULT_ROOT)
	if err != nil {
		panic(err)
	}

	c.Del("apps")

	return
}

func TestAppRegistration(t *testing.T) {
	c, app := appSetup("lolcatapp")

	check, err := c.Exists(app.Path())
	if err != nil {
		t.Error(err)
	}
	if check {
		t.Error("App already registered")
	}

	err = app.Register(c)
	if err != nil {
		t.Error(err)
	}

	check, err = c.Exists(app.Path())
	if err != nil {
		t.Error(err)
	}
	if !check {
		t.Error("App registration failed")
	}

	err = app.Register(c)
	if err == nil {
		t.Error("App allowed to be registered twice")
	}
}

func TestAppUnregistration(t *testing.T) {
	c, app := appSetup("dog")

	err := app.Register(c)
	if err != nil {
		t.Error(err)
	}

	err = app.Unregister(c)
	if err != nil {
		t.Error(err)
	}

	check, err := c.Exists(app.Path())
	if err != nil {
		t.Error(err)
	}
	if check {
		t.Error("App still registered")
	}
}

func TestSetAndGetEnvironmentVar(t *testing.T) {
	c, app := appSetup("lolcatapp")

	err := app.SetEnvironmentVar(c, "meow", "w00t")
	if err != nil {
		t.Error(err)
	}

	value, err := app.GetEnvironmentVar(c, "meow")
	if err != nil {
		t.Error(err)
	}

	if value != "w00t" {
		t.Errorf("EnvironmentVar 'meow' expected %s got %s", "w00t", value)
	}
}

func TestSetAndDelEnvironmentVar(t *testing.T) {
	c, app := appSetup("catalolna")

	err := app.SetEnvironmentVar(c, "wuff", "lulz")
	if err != nil {
		t.Error(err)
	}

	err = app.DelEnvironmentVar(c, "wuff")
	if err != nil {
		t.Error(err)
	}

	_, err = app.GetEnvironmentVar(c, "wuff")
	if err == nil {
		t.Error(err)
		t.Error("EnvironmentVar wasn't deleted")
	}

	if err != ErrKeyNotFound {
		t.Error(err)
	}
}

func TestEnvironmentVars(t *testing.T) {
	c, app := appSetup("cat-A-log")

	err := app.SetEnvironmentVar(c, "whiskers", "purr")
	if err != nil {
		t.Error(err)
	}
	err = app.SetEnvironmentVar(c, "lasers", "pew pew")
	if err != nil {
		t.Error(err)
	}

	vars, err := app.EnvironmentVars(c)
	if err != nil {
		t.Error(err)
	}
	if vars["whiskers"] != "purr" {
		t.Error("Var not set")
	}
	if vars["lasers"] != "pew pew" {
		t.Error("Var not set")
	}
}

func TestApps(t *testing.T) {
	c, _ := appSetup("apps-test")
	names := []string{"cat", "dog", "lol"}

	for i := range names {
		a, err := NewApp(names[i], "zebra", "joke")
		if err != nil {
			t.Error(err)
		}
		err = a.Register(c)
		if err != nil {
			t.Error(err)
		}
	}

	apps, err := Apps(c)
	if err != nil {
		t.Error(err)
	}
	if len(apps) != len(names) {
		t.Errorf("expected length %d returned length %d", len(names), len(apps))
	} else {
		for i := range apps {
			if apps[i].Name != names[i] {
				t.Error("expected %s got %s", names[i], apps[i].Name)
			}
		}
	}
}
