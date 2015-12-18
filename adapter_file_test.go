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

// 文件日志测试
package logmo

import(
    "testing"
    "time"
)

var logap = createLog()

func TestAFile( t *testing.T ) {

    logap.Emerg("specific language governing permissions")
    logap.Alert("specific language governing permissions")
    logap.Crit("specific language governing permissions")
    logap.Err("specific language governing permissions")
    logap.Warn("specific language governing permissions")
    logap.Notice("specific language governing permissions")
    logap.Info("specific language governing permissions")
    logap.Debug("specific language governing permissions")
    
    time.Sleep(time.Second * 1)
}

func BenchmarkAsyncFile(b *testing.B) {
    b.StopTimer() 
    b.StartTimer()
	for i := 0; i < b.N; i++ {
		logap.Debug("specific language governing permissions")
    }
    
    b.StopTimer() 
    time.Sleep(time.Second * 1)
}

func createLog() *Logger{
    log := New()
    fw  := NewAdapterFile( 10000 )
    fw.Filename = "async.log"
    fw.MaxLine  = 100
    fw.MaxSize  = 1 << 30
    fw.MaxDays  = 2
    fw.Rotation = 10
    go fw.Run()
    log.AddAdapter("asyncfile", fw)
    return log
}