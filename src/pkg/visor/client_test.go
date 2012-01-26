package visor

import (
	"github.com/soundcloud/doozer"
	"testing"
)

func setup(path string) (c *Client, conn *doozer.Conn) {
	c, err := Dial(DEFAULT_ADDR, DEFAULT_ROOT)
	if err != nil {
		panic(err)
	}

	c.Del(path)

	return c, c.conn
}

func TestDel(t *testing.T) {
	path := "del-test"
	c, conn := setup(path)

	_, err := conn.Set(DEFAULT_ROOT+"/del-test/deep/blue", c.rev, []byte{})
	if err != nil {
		t.Error(err)
	}

	err = c.Del(path)
	if err != nil {
		t.Error(err)
	}

	_, rev, err := conn.Stat(DEFAULT_ROOT+"/"+path, nil)
	if err != nil {
		t.Error(err)
	}

	if rev != 0 {
		t.Error("path isn't deleted")
	}
}

func TestExists(t *testing.T) {
	path := "exists-test"
	c, conn := setup(path)

	exists, err := c.Exists(path)
	if err != nil {
		t.Error(err)
	}
	if exists {
		t.Error("path shouldn't exist")
	}

	_, err = conn.Set(DEFAULT_ROOT+"/"+path, 0, []byte{})
	if err != nil {
		t.Error(err)
	}

	exists, err = c.Exists(path)
	if err != nil {
		t.Error(err)
	}
	if !exists {
		t.Error("path doesn't exist")
	}
}

func TestGet(t *testing.T) {
	path := "get-test"
	body := "aloha"
	c, conn := setup(path)

	_, err := conn.Set(DEFAULT_ROOT+"/"+path, 0, []byte(body))
	if err != nil {
		t.Error(err)
	}

	b, err := c.Get(path)
	if err != nil {
		t.Error(err)
	}
	if string(b) != body {
		t.Errorf("expected %s got %s", body, string(b))
	}
}

func TestKeys(t *testing.T) {
	path := "keys-test"
	keys := []string{"bar", "baz", "foo"}
	c, conn := setup(path)

	for i := range keys {
		_, err := conn.Set(DEFAULT_ROOT+"/"+path+"/"+keys[i], 0, []byte{})
		if err != nil {
			t.Error(err)
		}
	}

	k, err := c.Keys(path)
	if err != nil {
		t.Error(err)
	}
	if len(k) != len(keys) {
		t.Errorf("expected length %d returned lenght %d", len(keys), len(k))
	} else {
		for i := range keys {
			if keys[i] != k[i] {
				t.Errorf("expected %s got %s", keys[i], k[i])
			}
		}
	}
}

func TestSet(t *testing.T) {
	path := "set-test"
	body := "hola"
	c, conn := setup(path)

	err := c.Set(path, body)
	if err != nil {
		t.Error(err)
	}

	b, _, err := conn.Get(DEFAULT_ROOT+"/"+path, nil)
	if err != nil {
		t.Error(err)
	}
	if string(b) != body {
		t.Errorf("Expected %s got %s", body, string(b))
	}
}

func TestDifferentRoot(t *testing.T) {
	path := "visor"
	body := "test"
	c, conn := setup(path)

	client := &Client{Addr: c.Addr, conn: conn, Root: "/notvisor", rev: c.rev}
	err := client.Set("root", body)
	if err != nil {
		t.Error(err)
	}

	b, _, err := conn.Get("/notvisor/root", nil)
	if err != nil {
		t.Error(err)
	}
	if string(b) != body {
		t.Errorf("Expected %s got %s", body, string(b))
	}
}