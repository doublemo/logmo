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

// +build !windows

// 控制台文字颜色控制
package logmo

import(
    "io"
    "fmt"
)

// 定义日志等级色彩
var logColors = map[byte]string{
    EMERGENCY : "1;34",
    ALERT     : "1;36",
    CRITICAL  : "1;35",
    ERROR     : "1;31",
    WARNING   : "1;33",
    NOTICE    : "1;32",
    INFO      : "1;37",
    DEBUG     : "1;37",
}



// 向控制台写入颜色信息
func consoleWriteColor(out io.Writer, level byte, msg []byte) {
    msg = "\033[" + logColors[Level] + "m" + msg + "\033[0m"
    _, err := fmt.Fprintln(out, string(msg))
    return err
}