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
    "sync"
    "os"
    "time"
    "path/filepath"
    "strings"
    "log"
    "io"
    "bytes"
)

type AdapterFile struct {
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
    
    //锁
    lock sync.Mutex
    fwg sync.WaitGroup
    dwg sync.WaitGroup
    
    // 写入
    mutexWriter *fileMutex
    
    // 定义输出
    out *log.Logger
    
    // 文件名称
    Filename string
    
    // 文件最大行
    MaxLine int
    
    // 文件最大
    MaxSize int
    
    // 最大保存天数
    MaxDays int64
    
    //  循环滚动次数
    Rotation int
    
    // 记录当前文件line
    line int
    
    // 记录当前文件size
    size int
    
    lastDay int
}


func (adapter *AdapterFile) write( message Message ) error {
    adapter.lock.Lock()
    defer adapter.lock.Unlock()
    
    // 格式化
    msg, err := adapter.formatter.Format( message )
    if err != nil {
        return err
    }
    
    size := len(msg)
    adapter.check(size)
    adapter.out.Println(string(msg))
    return nil
}

func (adapter *AdapterFile) AsyncWrite( message Message ) error {
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

func (adapter *AdapterFile) SyncWrite( message Message ) error {
    // 执行hook
    for _, hook := range adapter.hooks {
       err :=  hook.Fire(message)
       if err != nil {
           return nil
       }
    }
    
    return adapter.write( message )
}

func (adapter *AdapterFile) SetFormatter( formatter Formatter ) error {
    adapter.formatter = formatter
    return nil
}

func (adapter *AdapterFile) AddHook( name string, hook Hook ) error {
    adapter.lock.Lock()
    defer adapter.lock.Unlock()
    
    if _, ok := adapter.hooks[name]; ok {
        return nil
    }
    
    adapter.hooks[name] = hook
    return nil
}

func (adapter *AdapterFile) DeleteHook( name string ) error {
    adapter.lock.Lock()
    defer adapter.lock.Unlock()
    
    if _, ok := adapter.hooks[name]; !ok {
        return nil
    }
    
    delete(adapter.hooks, name)
    return nil
}

func (adapter *AdapterFile) Async( b bool ) {
    adapter.async = b
}

func (adapter *AdapterFile) IsAsync() bool {
    return adapter.async
}

func (adapter *AdapterFile) Destroy() {
    adapter.dwg.Add(1)
    adapter.event <- ADAPTER_EVENT_DESTORY
    adapter.dwg.Wait()
}

func (adapter *AdapterFile) Flush() {
    adapter.fwg.Add(1)
    adapter.event <- ADAPTER_EVENT_FLUSH
    adapter.fwg.Wait()
}

func (adapter *AdapterFile) Run() {
    // 初始化
    adapter.Initialize()
    for{
        select {
            case message := <-adapter.channel:
              err := adapter.write( message )
              if err != nil {
                  fmt.Println("Run:",err)
              }
              
            case e := <-adapter.event:
              switch e {
                  case ADAPTER_EVENT_DESTORY:
                    close(adapter.channel)
                    close(adapter.event)
                    adapter.mutexWriter.Close()
                    adapter.dwg.Done()
                    return
                    
                  case ADAPTER_EVENT_FLUSH:
                    if len(adapter.channel) > 0 {
                        adapter.write( <-adapter.channel )
                    }
                    adapter.mutexWriter.Flush()
                    adapter.fwg.Done()
                    return
              }
        }
    }
}

// 检查文件是否满足条件,进行日志分割
func (adapter *AdapterFile) check( size int ) error {
    if (adapter.size > adapter.MaxSize && adapter.size > 0) || 
       (adapter.line > adapter.MaxLine && adapter.line > 0) ||
       adapter.lastDay != time.Now().Day() {
        if err := adapter.rotate(); err != nil {
            fmt.Fprintf(os.Stderr, "AdapterFile(%q): %s\n", adapter.Filename, err)
            return err
        }
    }
    
    adapter.line ++
    adapter.size += size
    return nil
}

// 滚动日志分割
func (adapter *AdapterFile) rotate() error {
    adapter.mutexWriter.Lock()
    defer adapter.mutexWriter.Unlock()
    adapter.mutexWriter.Close()
    _, e := os.Lstat(adapter.Filename)
    if e != nil {
        return nil
    }
    
    for n := adapter.Rotation; n > 0 ; n -- {
        fname  := adapter.Filename + fmt.Sprintf(".%s.%04d", time.Now().Format("2006-01-02"), n)
        _, err := os.Lstat(fname)
        if err != nil  {
            continue
        }
        
        if n >= adapter.Rotation {
           err = os.Remove(fname)
           if err != nil {
               return err
           }
           continue
        } 
        
        tname := adapter.Filename + fmt.Sprintf(".%s.%04d", time.Now().Format("2006-01-02"), n + 1)
        err = os.Rename(fname, tname)
        if err != nil {
            return err
        }
    }
    
    
    tname := adapter.Filename + fmt.Sprintf(".%s.%04d", time.Now().Format("2006-01-02"), 1)
    err := os.Rename(adapter.Filename, tname)
    if err != nil {
        return err
    }
    
    go adapter.deleteExpiredLog()
    adapter.Initialize()
    return nil
}

// 删除过期日志
func (adapter *AdapterFile) deleteExpiredLog() {
    dir := filepath.Dir(adapter.Filename)
    filepath.Walk(dir, func(path string, info os.FileInfo, err error) (returnErr error) {
        defer func() {
            if r := recover(); r != nil {
                returnErr = fmt.Errorf("Unable to delete old log '%s', error: %+v", path, r)
                fmt.Println(returnErr)
            }
        }()
        
        if !info.IsDir() && info.ModTime().Unix() < (time.Now().Unix()-60*60*24*adapter.MaxDays) {
            if strings.HasPrefix(filepath.Base(path), filepath.Base(adapter.Filename)) {
                os.Remove(path)
            }
        }
        return
    })
}

func (adapter *AdapterFile) Initialize() {
    size, err := adapter.mutexWriter.Open(adapter.Filename)
    if err != nil {
        fmt.Fprintf(os.Stderr, "AdapterFile - Initialize-(%q): %s\n", adapter.Filename, err)
        return
    }
    
    adapter.size     = size 
    adapter.lastDay  = time.Now().Day()
    adapter.line    = 0
    
    if adapter.size > 0 {
        num, err := adapter.lines()
         if err != nil {
             return
         }
         
        adapter.line = num
    }
}

func (adapter *AdapterFile) lines() (int, error) {
    fd, err := os.Open(adapter.Filename)
    if err != nil {
        return 0, err
    }
    
    defer fd.Close()
    buf   := make([]byte, 32768) // 32k
    count := 0
    lineSep := []byte{'\n'}
    
    for {
        c, err := fd.Read(buf)
        if err != nil && err != io.EOF {
            return count, err
        }
        
        count += bytes.Count(buf[:c], lineSep)
        if err == io.EOF {break}
    }
    
    return count, nil
}

func NewAdapterFile( channelLen int ) *AdapterFile {
    adapter:= &AdapterFile{
        channel : make(chan Message, channelLen),
        event   : make(chan AdapterEvent),
        formatter : new(FormatterText),
        hooks   : make(map[string]Hook),
        async   : true,
    }
    
    adapter.mutexWriter = new(fileMutex)
    adapter.out         = log.New(adapter.mutexWriter, "", log.Ldate|log.Ltime)
    
    // 默认配置
    adapter.Filename    = "log"
    adapter.MaxLine     = 10000
    adapter.MaxSize     = 1 << 28
    adapter.MaxDays     = 30
    adapter.Rotation    = 50
    return adapter
}