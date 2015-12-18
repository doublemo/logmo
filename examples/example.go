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
package main

import (
	"github.com/doublemo/logmo"
	"time"
)

func main() {
	// logmo默认支持控制台输出
	log := logmo.New()
	// 增加文件写入
	fw := logmo.NewAdapterFile(10000)
	fw.Filename = "async.log"
	fw.MaxLine = 100
	fw.MaxSize = 1 << 30
	fw.MaxDays = 2
	fw.Rotation = 10
	// 增加日志等级过滤
	fw.AddHook("level", &logmo.HookLevel{logmo.ERROR})
	go fw.Run()

	log.AddAdapter("asyncfile", fw)
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
