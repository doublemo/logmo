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

// +build windows

// 控制台文字颜色控制

package logmo

import(
    "io"
    "fmt"
    "syscall"
)

var(
    kernel32DLL                 = syscall.NewLazyDLL("kernel32.dll")
	setConsoleTextAttributeProc = kernel32DLL.NewProc("SetConsoleTextAttribute")
)

var logColors = map[byte]uint16{
    EMERGENCY : 0x0004,
    ALERT     : 0x0008,
    CRITICAL  : 0x0005,
    ERROR     : 0x0004,
    WARNING   : 0x0006,
    NOTICE    : 0x0002,
    INFO      : 0x0007,
    DEBUG     : 0x0003,
}

type fileInterface interface {
	Fd() uintptr
}

// 向控制台写入颜色信息
func consoleWriteColor(out io.Writer, level byte, msg []byte) error {
    if f, ok := out.(fileInterface); ok {
        setConsoleTextAttribute(f, logColors[level])
        _, err := fmt.Fprintln(out, string(msg))
        setConsoleTextAttribute(f, 0x0007)
        
        return err
    }
    
    _, err := fmt.Fprintln(out, string(msg))
    return err
}

// setConsoleTextAttribute sets the attributes of characters written to the
// console screen buffer by the WriteFile or WriteConsole function.
// See http://msdn.microsoft.com/en-us/library/windows/desktop/ms686047(v=vs.85).aspx.
func setConsoleTextAttribute( f fileInterface, attribute uint16) bool {
    ok, _, _ := setConsoleTextAttributeProc.Call(f.Fd(), uintptr(attribute), 0)
    return ok != 0
}
