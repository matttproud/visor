// Copyright (c) 2012, SoundCloud Ltd.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Source code and contact info at http://github.com/soundcloud/visor

package visor

import (
	"testing"
)

func appSetup(name string) (app *App) {
	s, err := Dial(DEFAULT_ADDR, DEFAULT_ROOT)
	if err != nil {
		panic(err)
	}

	r, _ := s.conn.Rev()
	err = s.conn.Del("apps", r)
	rev, err := Init(s)
	if err != nil {
		panic(err)
	}

	app = NewApp(name, "git://cat.git", "whiskers", s)
	app = app.FastForward(rev)

	return
}

func TestAppRegistration(t *testing.T) {
	app := appSetup("lolcatapp")

	check, _, err := app.conn.Exists(app.Path.Dir)
	if err != nil {
		t.Error(err)
		return
	}
	if check {
		t.Error("App already registered")
		return
	}

	app2, err := app.Register()
	if err != nil {
		t.Error(err)
		return
	}
	check, _, err = app2.conn.Exists(app.Path.Dir)
	if err != nil {
		t.Error(err)
		return
	}
	if !check {
		t.Error("App registration failed")
		return
	}
	_, err = app.Register()
	if err == nil {
		t.Error("App allowed to be registered twice")
	}
	_, err = app2.Register()
	if err == nil {
		t.Error("App allowed to be registered twice")
	}
}

func TestEnvPersistenceOnRegister(t *testing.T) {
	app := appSetup("envyapp")

	app.Env["VAR1"] = "VAL1"
	app.Env["VAR2"] = "VAL2"

	app, err := app.Register()
	if err != nil {
		t.Error(err)
		return
	}

	env, err := app.EnvironmentVars()
	if err != nil {
		t.Error(err)
		return
	}
	for key, val := range app.Env {
		if env[key] != val {
			t.Errorf("%s should be '%s', got '%s'", key, val, env[key])
		}
	}
}

func TestAppUnregistration(t *testing.T) {
	app := appSetup("dog")

	app, err := app.Register()
	if err != nil {
		t.Error(err)
		return
	}

	err = app.Unregister()
	if err != nil {
		t.Error(err)
		return
	}

	check, _, err := app.conn.Exists(app.Path.Dir)
	if err != nil {
		t.Error(err)
	}
	if check {
		t.Error("App still registered")
	}
}

func TestAppUnregistrationFailure(t *testing.T) {
	app := appSetup("dog-fail")

	app2, err := app.Register()
	if err != nil {
		t.Error(err)
		return
	}

	err = app.Unregister()
	if err == nil {
		t.Error("App allowed to be unregistered with old revision")
		return
	}

	err = app2.Unregister()
	if err != nil {
		t.Error(err)
		return
	}

	app3 := app2.FastForward(-1)
	_, err = app3.Register()
	if err != nil {
		t.Error(err)
		return
	}
}

func TestSetAndGetEnvironmentVar(t *testing.T) {
	app := appSetup("lolcatapp")

	app, err := app.SetEnvironmentVar("meow", "w00t")
	if err != nil {
		t.Error(err)
		return
	}
	if app.Env["meow"] != "w00t" {
		t.Error("app.Env should be updated")
	}

	value, err := app.GetEnvironmentVar("meow")
	if err != nil {
		t.Error(err)
		return
	}

	if value != "w00t" {
		t.Errorf("EnvironmentVar 'meow' expected %s got %s", "w00t", value)
	}
}

func TestSetAndDelEnvironmentVar(t *testing.T) {
	app := appSetup("catalolna")

	app, err := app.SetEnvironmentVar("wuff", "lulz")
	if err != nil {
		t.Error(err)
	}

	app, err = app.DelEnvironmentVar("wuff")
	if err != nil {
		t.Error(err)
		return
	}

	v, err := app.GetEnvironmentVar("wuff")
	if err == nil {
		t.Errorf("EnvironmentVar wasn't deleted: %#v", v)
		return
	}
}

func TestEnvironmentVars(t *testing.T) {
	app := appSetup("cat-A-log")

	_, err := app.SetEnvironmentVar("whiskers", "purr")
	if err != nil {
		t.Error(err)
	}
	app, err = app.SetEnvironmentVar("lasers", "pew pew")
	if err != nil {
		t.Error(err)
	}

	vars, err := app.EnvironmentVars()
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
	app := appSetup("apps-test")
	names := []string{"cat", "dog", "lol"}

	for i := range names {
		a := NewApp(names[i], "zebra", "joke", app.Path.Snapshot)
		_, err := a.Register()
		if err != nil {
			t.Error(err)
		}
	}
	app = app.FastForward(-1)

	s, _ := Dial(DEFAULT_ADDR, DEFAULT_ROOT)

	apps, err := Apps(s)
	if err != nil {
		t.Error(err)
	}
	if len(apps) != len(names) {
		t.Errorf("expected length %d returned length %d", len(names), len(apps))
	} else {
		for i := range apps {
			if apps[i].Name != names[i] {
				t.Errorf("expected %s got %s", names[i], apps[i].Name)
			}
		}
	}
}
