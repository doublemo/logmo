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
    "io"
    "sync"
    "os"
)

type AdapterConsole struct { 
    // 定义通道
    channel chan Message
    
    // 定义事件通道
    event chan AdapterEvent
    
    // 格式化
    formatter Formatter
    
    // hooks
    hooks map[string]Hook
    
    // 处理模式
    async bool
    
    // out
    out io.Writer
    
    lock sync.Mutex
    
    fwg sync.WaitGroup
    dwg sync.WaitGroup
}

func (adapter *AdapterConsole) write( message Message ) error {
    adapter.lock.Lock()
    defer adapter.lock.Unlock()
    
    // 格式化
    msg, err := adapter.formatter.Format( message )
    if err != nil {
        return err
    }
    
    return consoleWriteColor(adapter.out, message.GetLevel(), msg)
}

func (adapter *AdapterConsole) SyncWrite( message Message ) error {
     // 执行hook
    for _, hook := range adapter.hooks {
       err :=  hook.Fire(message)
       if err != nil {
           return nil
       }
    }
    
    return adapter.write( message )
}

func (adapter *AdapterConsole) AsyncWrite( message Message ) error {
    defer func() {
        if r := recover(); r != nil {
           fmt.Println(r)
        }
    }()
    
     // 执行hook
    for _, hook := range adapter.hooks {
       err :=  hook.Fire(message)
       if err != nil {
           return nil
       }
    }
    
    adapter.channel <- message
    return nil
}

func (adapter *AdapterConsole) SetFormatter( formatter Formatter ) error {
    adapter.formatter = formatter
    return nil
}

func (adapter *AdapterConsole) AddHook( name string, hook Hook ) error {
    adapter.lock.Lock()
    defer adapter.lock.Unlock()
    
    if _, ok := adapter.hooks[name]; ok {
        return nil
    }
    
    adapter.hooks[name] = hook
    return nil
}

func (adapter *AdapterConsole) DeleteHook( name string ) error {
    adapter.lock.Lock()
    defer adapter.lock.Unlock()
    
    if _, ok := adapter.hooks[name]; !ok {
        return nil
    }
    
    delete(adapter.hooks, name)
    return nil
}

func (adapter *AdapterConsole) Async( b bool ) {
    adapter.async = b
}

func (adapter *AdapterConsole) IsAsync() bool {
    return adapter.async
}

func (adapter *AdapterConsole) Destroy() {
    adapter.dwg.Add(1)
    adapter.event <- ADAPTER_EVENT_DESTORY
    adapter.dwg.Wait()
}

func (adapter *AdapterConsole) Flush() {
    adapter.fwg.Add(1)
    adapter.event <- ADAPTER_EVENT_FLUSH
    adapter.fwg.Wait()
}

func (adapter *AdapterConsole) Run() {
    for{
        select {
            case message := <-adapter.channel:
              err := adapter.write( message )
              if err != nil {
                  fmt.Println(err)
              }
              
            case e := <-adapter.event:
              switch e {
                  case ADAPTER_EVENT_DESTORY:
                    close(adapter.channel)
                    close(adapter.event)
                    adapter.dwg.Done()
                    return
                    
                  case ADAPTER_EVENT_FLUSH:
                    if len(adapter.channel) > 0 {
                        adapter.write( <-adapter.channel )
                    }
                    adapter.fwg.Done()
              }
        }
    }
}

func NewAdapterConsole( channelLen int ) *AdapterConsole{
    return &AdapterConsole{
        channel : make(chan Message, channelLen),
        event   : make(chan AdapterEvent),
        formatter : new(FormatterText),
        hooks   : make(map[string]Hook),
        async   : true,
        out     : os.Stdout,
    }
}