// pmm-agent
// Copyright (C) 2018 Percona LLC
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.

package supervisor

import (
	"context"
	"syscall"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/percona/pmm/api/agent"
)

func TestRaceCondition(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	ctx, cancel := context.WithCancel(context.Background())
	m := New(ctx, &agent.SetStateRequest_AgentProcess{
		Type: agent.Type_MYSQLD_EXPORTER,
		Args: []string{"-web.listen-address=127.0.0.1:11111"},
		Env: []string{
			`DATA_SOURCE_NAME="pmm:pmm@(127.0.0.1:3306)/pmm-managed-dev"`,
		},
	})
	go func() {
		time.Sleep(1 * time.Second)
		cancel()
	}()
	assert.NotPanics(t, func() {
		for {
			state := <-m.Changes()
			if state == STOPPED {
				break
			}
		}
	})
}

func TestStates(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	m := New(ctx, &agent.SetStateRequest_AgentProcess{
		Type: agent.Type_MYSQLD_EXPORTER,
		Args: []string{"-web.listen-address=127.0.0.1:11112"},
		Env: []string{
			`DATA_SOURCE_NAME="pmm:pmm@(127.0.0.1:3306)/pmm-managed-dev"`,
		},
	})

	time.Sleep(3 * time.Second)
	err := syscall.Kill(m.pid(), syscall.SIGKILL)
	assert.NoError(t, err)
	assert.Equal(t, STARTING, <-m.Changes())
	assert.Equal(t, RUNNING, <-m.Changes())
	assert.Equal(t, BACKOFF, <-m.Changes())
	assert.Equal(t, STARTING, <-m.Changes())
	assert.Equal(t, RUNNING, <-m.Changes())
	m.Stop()
	assert.Equal(t, STOPPING, <-m.Changes())
	assert.Equal(t, STOPPED, <-m.Changes())
}
