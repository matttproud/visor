// Copyright (c) 2012, SoundCloud Ltd.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Source code and contact info at http://github.com/soundcloud/visor

package visor

import (
	"testing"
)

func snapshotSetup() (s Snapshot) {
	s, err := Dial(DEFAULT_ADDR, "/snapshot-test")
	if err != nil {
		panic(err)
	}
	err = s.Del("/")
	Init(s)
	s = s.FastForward(-1)

	return
}

func TestSnapshotExists(t *testing.T) {
	s := snapshotSetup()
	c := s.conn
	k := "key"
	v := "value"

	rev, err := c.Set(k, -1, []byte(v))
	if err != nil {
		panic(err)
	}

	exists, _, err := s.Exists(k)
	if exists {
		t.Errorf("path %s shouldn't exist yet", k)
	}

	exists, rev1, err := s.FastForward(rev).Exists(k)
	if !exists {
		t.Errorf("path %s should exist", k)
	}
	if rev1 != rev {
		t.Error("snapshot rev should match file rev")
	}
}

func TestSnapshotSetGet(t *testing.T) {
	s := snapshotSetup()
	k := "key"
	v := "value"

	s1, err := s.Set(k, v)
	if err != nil {
		t.Error(err)
	}

	_, _, err = s.Get(k)
	if !IsErrNoEnt(err) {
		t.Error("expected NoEnt error")
	}

	val, rev, err := s1.Get(k)
	if err != nil {
		t.Error(err)
	}
	if val != v {
		t.Errorf("expected value '%s', got '%s'", v, val)
	}
	if rev != s1.Rev {
		t.Errorf("unexpected rev")
	}

	// REV_MISMATCH

	s2, err := s.Set(k, v)
	if err == nil {
		t.Error("expected REV_MISMATCH")
	}
	if s2 != s {
		t.Error("expected return snapshot to be same on error")
	}
}

func TestSnapshotIsNoEnt(t *testing.T) {
	s := snapshotSetup()

	_, _, err := s.Get("/does-not-exist")
	if !IsErrNoEnt(err) {
		t.Error("expected NoEnt error")
	}
}

func TestSnapshotUpdate(t *testing.T) {
	s := snapshotSetup()
	k := "key"
	v := "value"

	s1, err := s.Set(k, v)
	if err != nil {
		panic(err)
	}

	_, err = s.Update(k, "#")
	if !IsErrNoEnt(err) {
		t.Error("expected NoEnt error")
	}

	s2, err := s1.Update(k, "#")
	if err != nil {
		t.Error(err)
	}

	s3, err := s2.Update(k, "*")
	if err != nil {
		t.Error(err)
	}

	val, _, _ := s.conn.Get(k, &s2.Rev)
	if string(val) != "#" {
		t.Error("expected '#' value")
	}

	val, _, _ = s.conn.Get(k, &s3.Rev)
	if string(val) != "*" {
		t.Error("expected '*' value")
	}

	if !(s3.Rev > s2.Rev && s2.Rev > s1.Rev) {
		t.Error("incorrect revision ordering")
	}
}

func TestSnapshotGetAndSetScale(t *testing.T) {
	s := snapshotSetup()

	app := NewApp("scale-app", "git://scale.git", "scale-stack", s)
	pty := NewProcType(app, "scale-proc", s)
	rev := NewRevision(app, "scale-ref", s)

	if _, err := app.Register(); err != nil {
		panic(err)
	}
	if _, err := rev.Register(); err != nil {
		panic(err)
	}
	if _, err := pty.Register(); err != nil {
		panic(err)
	}
	s = s.FastForward(-1)

	scale, _, err := s.GetScale(app.Name, rev.Ref, string(pty.Name))
	if err != nil {
		t.Error(err)
	}
	if scale != 0 {
		t.Error("expected initial scale of 0")
	}

	s1, err := s.SetScale(app.Name, rev.Ref, string(pty.Name), 9)
	if err != nil {
		t.Error(err)
	}

	scale, _, err = s1.GetScale(app.Name, rev.Ref, string(pty.Name))
	if err != nil {
		t.Error(err)
	}
	if scale != 9 {
		t.Errorf("expected scale of 9, got %d", scale)
	}

	scale, _, err = s1.GetScale("invalid-app", rev.Ref, string(pty.Name))
	if scale != 0 {
		t.Errorf("expected scale to be 0")
	}
}
