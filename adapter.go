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

// 适配器接口
package logmo

type Adapter interface {
    // 同步信息写入
    SyncWrite( message Message ) error
    
    // 导步信息写入
    AsyncWrite( message Message ) error 
    
    // 设置处理方式
    Async( b bool )
    
    // 是否异步处理
    IsAsync() bool
    
    // 信息格式化
    SetFormatter( formatter Formatter ) error
    
    // Hook
    AddHook( name string, hook Hook ) error
    
    // Delete Hook
    DeleteHook( name string ) error
    
    // Close
    Destroy()
    
    // Run
    Run()
    
    // Flush
    Flush()
}

// 定义事件驱动
const(
    // 销毁事件
    ADAPTER_EVENT_DESTORY = iota
    
    ADAPTER_EVENT_FLUSH
)

// 定义事件类型
type AdapterEvent byte