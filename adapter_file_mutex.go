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

// 文件日志同步写入
package logmo

import(
    "os"
    "sync"
)

type fileMutex struct{
    sync.Mutex
    fd *os.File
}

// 同步写入
func (mutex *fileMutex) Write(b []byte) (int, error) {
    mutex.Lock()
    defer mutex.Unlock()
    
    return mutex.fd.Write(b)
}

// 打开文件并返回文件大小
func (mutex *fileMutex) Open( file string ) (int, error) {
    if mutex.fd != nil {
        mutex.fd.Close()
    }
    
    fd, err := os.OpenFile(file, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0660)
    if err != nil {
        return 0, err
    }
    
    finfo, err := fd.Stat()
    if err != nil {
        fd.Close()
        return 0, err
    }
    mutex.fd = fd
    return int(finfo.Size()), nil
}

func (mutex *fileMutex) Close() {
    mutex.fd.Close()
}

func (mutex *fileMutex) Flush() {
    mutex.fd.Sync()
}