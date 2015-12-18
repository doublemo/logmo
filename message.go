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

// 信息结构
package logmo

import(
    "time"
)

type Message interface {
    // 获取日志信息
    GetMessage() string
    
    // 获取日志等级
    GetLevel() byte
    
    // 获取文件
    GetFile() string
    
    // 获取行号
    GetLine() int
    
    // 获取附加数据
    GetData() interface{} 
    
    // 获取时间
    GetTime() time.Time
    
    // 获取前缀
    GetPrefix() string
    
    // 获取进程号
    GetPID() int
    
    // 获取信息编号
    GetID() int64
}

type DefaultMessage struct {
    Level byte
    File  string
    Line  int
    Message string
    Data  interface{}
    Time  time.Time
    Prefix string
    Pid int
    Id  int64
}

func (msg *DefaultMessage) GetMessage() string {
    return msg.Message
}

func (msg *DefaultMessage) GetLevel() byte {
    return msg.Level
}

func (msg *DefaultMessage) GetFile() string {
    return msg.File
}

func (msg *DefaultMessage) GetLine() int {
    return msg.Line
}

func (msg *DefaultMessage) GetData() interface{} {
    return msg.Data
}

func (msg *DefaultMessage) GetTime() time.Time {
    return msg.Time
}

func (msg *DefaultMessage) GetPrefix() string {
    return msg.Prefix
}

func (msg *DefaultMessage) GetPID() int {
    return msg.Pid
} 

func (msg *DefaultMessage) GetID() int64 {
    return msg.Id
} 