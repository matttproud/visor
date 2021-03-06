// Copyright (c) 2012, SoundCloud Ltd.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Source code and contact info at http://github.com/soundcloud/visor

// Visor is a library which provides an abstract interface
// over a global process state.
//
// This process state is referred to as the registry.
//
// Example usage:
//
//     package main
//
//     import "soundcloud/visor"
//
//     func main() {
//         client, err := visor.Dial("coordinator:8046", "/", new(visor.StringCodec))
//         if err != nil {
//           panic(err)
//         }
//
//         l := make(chan *visor.Event)
//
//         // Watch for changes in the global process state
//         go visor.WatchEvent(client.Snapshot, l)
//
//         for {
//             fmt.Println(<-l)
//         }
//     }
//
package visor

import (
	"errors"
	"fmt"
	"path"
	"strconv"
	"time"
)

const DEFAULT_URI string = "doozer:?ca=localhost:8046"
const DEFAULT_ADDR string = "localhost:8046"
const DEFAULT_ROOT string = "/visor"
const SCALE_PATH string = "scale"
const START_PORT int = 8000
const START_PORT_PATH string = "/next-port"
const UID_PATH string = "/uid"

type ProcessName string
type Stack string
type State string

func Init(s Snapshot) (rev int64, err error) {
	var s1 Snapshot

	exists, _, err := s.conn.Exists(START_PORT_PATH)
	if err != nil {
		return
	}

	if !exists {
		s1, err = s.Set(START_PORT_PATH, strconv.Itoa(START_PORT))
		if err != nil {
			return
		}

		return s1.Rev, err
	}
	return s.conn.Rev()
}

func ClaimNextPort(s Snapshot) (port int, err error) {
	for {
		f, err := GetLatest(s, START_PORT_PATH, new(IntCodec))
		if err == nil {
			port = f.Value.(int)

			f, err = f.Set(port + 1)
			if err == nil {
				break
			} else {
				s = f.Snapshot
				time.Sleep(time.Second / 10)
			}
		} else {
			return -1, err
		}
	}

	return
}

func Scale(app string, revision string, processName string, factor int, s Snapshot) (err error) {
	if factor < 0 {
		return errors.New("scaling factor needs to be a positive integer")
	}

	exists, _, err := s.conn.Exists(path.Join(APPS_PATH, app, REVS_PATH, revision))
	if !exists || err != nil {
		return fmt.Errorf("%s@%s not found", app, revision)
	}
	exists, _, err = s.conn.Exists(path.Join(APPS_PATH, app, PROCS_PATH, processName))
	if !exists || err != nil {
		return fmt.Errorf("proc '%s' doesn't exist", processName)
	}

	op := OpStart
	tickets := factor

	current, _, err := s.GetScale(app, revision, processName)
	if err != nil {
		return
	}

	// Check response isn't empty
	if current > 0 {
		tickets = factor - current

		if tickets < 0 {
			op = OpStop
			tickets = -tickets
		}
	}

	s, err = s.SetScale(app, revision, processName, factor)
	if err != nil {
		return
	}

	for i := 0; i < tickets; i++ {
		var ticket *Ticket

		ticket, err = CreateTicket(app, revision, ProcessName(processName), op, s)
		if err != nil {
			return
		}

		s = s.FastForward(ticket.Rev)
	}

	return
}
