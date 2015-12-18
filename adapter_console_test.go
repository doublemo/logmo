// Copyright 2015 doublemo. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// 控制台日志测试
package logmo

import(
    "os"
    "testing"
    "time"
)

func TestAConsole( t *testing.T ) {
    log := New()
    c,err := log.GetAdapter("default")
    if err != nil {
        t.Fatal(err)
    }
    
    c.AddHook("level", &HookLevel{ERROR})
   
    log.Emerg("specific language governing permissions")
    log.Alert("specific language governing permissions")
    log.Crit("specific language governing permissions")
    
    log.Err("specific language governing permissions")
    log.Warn("specific language governing permissions")
    log.Notice("specific language governing permissions")
    log.Info("specific language governing permissions")
    log.Debug("specific language governing permissions")
    
    time.Sleep(time.Second * 1)
}


func BenchmarkSyncConsole(b *testing.B) {
	 f := &FormatterText{}
     c := &AdapterConsole{
        channel:make(chan Message, 10000),
        formatter:f,
        out: os.Stdout,
     }
     
     m := &DefaultMessage{
        Level:NOTICE,
        File :"t.go",
        Line : 5,
        Message : "specific language governing permissions",
        Data :nil,
        Time : time.Now(),
        Prefix :"ERR",
    }
	for i := 0; i < b.N; i++ {
		c.SyncWrite(m)
    }
}

func BenchmarkAsyncConsole(b *testing.B) {
	 f := &FormatterText{}
     c := &AdapterConsole{
        channel:make(chan Message, 1000),
        formatter:f,
        out: os.Stdout,
     }
     
     m := &DefaultMessage{
        Level:NOTICE,
        File :"t.go",
        Line : 5,
        Message : "specific language governing permissions",
        Data :nil,
        Time : time.Now(),
        Prefix :"ERR",
    }
    
    go c.Run()
	for i := 0; i < b.N; i++ {
		c.AsyncWrite(m)
    }
}