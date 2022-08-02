// Copyright © 2021 sealos.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sshutil

import (
	"testing"

	"github.com/fanux/sealos/pkg/logger"
)

func TestSSHCopyLocalToRemote(t *testing.T) {
	type args struct {
		host       string
		localPath  string
		remotePath string
	}
	var (
		host = "127.0.0.1"
		ssh  = SSH{
			User:       "lllrrrccc",
			Password:   "root888",
			PkFile:     "",
			PkPassword: "",
			Timeout:    nil,
		}
	)
	tests := []struct {
		name   string
		fields SSH
		args   args
	}{
		{"test for copy file to remote server", ssh, args{
			host,
			"/root/test.done",
			"/data/test.done",
		}},
		{"test copy for file when local file is not exist", ssh, args{
			host,
			// local file  001 is not exist.
			"/root/test.done01",
			"/data/test.done01",
		}},
		{"test copy dir to remote server", ssh, args{
			host,
			"/home/tmp",
			"/data/tmp",
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ss := &SSH{
				User:       tt.fields.User,
				Password:   tt.fields.Password,
				PkFile:     tt.fields.PkFile,
				PkPassword: tt.fields.PkPassword,
				Timeout:    tt.fields.Timeout,
			}

			if !fileExist(tt.args.localPath) {
				logger.Error("local filepath is not exist")
				return
			}
			if ss.IsFileExist(host, tt.args.remotePath) {
				logger.Error("remote filepath is exist")
				return
			}
			// test copy dir
			ss.CopyLocalToRemote(tt.args.host, tt.args.localPath, tt.args.remotePath)

			// test the copy result
			ss.Cmd(tt.args.host, "ls -lh "+tt.args.remotePath)

			// rm remote file
			ss.Cmd(tt.args.host, "rm -rf "+tt.args.remotePath)
		})
	}
}
