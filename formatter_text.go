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

// 控制台日志支持
package logmo

import(
    "fmt"
)

type FormatterText struct {}

func (format *FormatterText) Format( message Message ) ( []byte, error ) {
    msg      := message.GetMessage()
    createat := message.GetTime().Format("2006/01/02 15:04:05")
    line     := message.GetLine()
    file     := message.GetFile()
    prefix   := message.GetPrefix()
    
    if line > 0 && file != "" {
        msg = fmt.Sprintf("%s [%s] [%s:%d] %s", createat, prefix, file, line, msg)
    } else {
        msg = fmt.Sprintf("%s [%s] %s", createat, prefix, msg)
    }
   
	return []byte(msg), nil
}