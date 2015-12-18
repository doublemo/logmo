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

// 日志处理
package logmo

import (
	"errors"
	"fmt"
	"os"
	"path"
	"runtime"
	"sync"
	"time"
)

// 定义日志等级
const (
	EMERGENCY = iota
	ALERT
	CRITICAL
	ERROR
	WARNING
	NOTICE
	INFO
	DEBUG
)

type Logger struct {
	// 适配器
	adapters map[string]Adapter

	// 定义错误深度
	ExtraCalldepth int

	// 信息记数器
	counter int64

	lock sync.Mutex
}

// 增加适配器
func (log *Logger) AddAdapter(name string, adapter Adapter) error {
	log.lock.Lock()
	defer log.lock.Unlock()
	if _, ok := log.adapters[name]; ok {
		return nil
	}

	log.adapters[name] = adapter
	return nil
}

// 删除适配器
func (log *Logger) DeleteAdapter(name string) error {
	log.lock.Lock()
	defer log.lock.Unlock()
	if _, ok := log.adapters[name]; !ok {
		return nil
	}

	delete(log.adapters, name)
	return nil
}

// 获取适配器
func (log *Logger) GetAdapter(name string) (Adapter, error) {
	if _, ok := log.adapters[name]; ok {
		return log.adapters[name], nil
	}

	return nil, errors.New(fmt.Sprintf("Adapter:%s not found", name))
}

// 输入信息
func (log *Logger) Write(level byte, prefix string, msg string, data interface{}, sync bool) error {
	log.counter++
	message := new(DefaultMessage)
	message.Level = level
	message.Message = msg
	message.Prefix = prefix
	message.Time = time.Now()
	message.Data = data
	message.Pid = os.Getpid()
	message.Id = time.Now().UnixNano() + log.counter
	_, file, line, ok := runtime.Caller(2 + log.ExtraCalldepth)
	if ok {
		_, filename := path.Split(file)
		message.Line = line
		message.File = filename
	}
	errs := []error{}
	for _, adapter := range log.adapters {
		if sync {
			err := adapter.SyncWrite(message)
			if err != nil {
				errs = append(errs, err)
			}

			continue
		}

		if adapter.IsAsync() {
			err := adapter.AsyncWrite(message)
			if err != nil {
				errs = append(errs, err)
			}
		} else {
			err := adapter.SyncWrite(message)
			if err != nil {
				errs = append(errs, err)
			}
		}
	}

	if len(errs) > 0 {
		return errs[0]
	}
	return nil
}

func (log *Logger) Flush() {
	for _, adapter := range log.adapters {
		adapter.Flush()
	}
}

func (log *Logger) Close() {
	for _, adapter := range log.adapters {
		adapter.Destroy()
	}

	log.adapters = make(map[string]Adapter)
}

// 紧急
func (log *Logger) Emerg(format string, v ...interface{}) {
	log.Write(EMERGENCY, "M", fmt.Sprintf(format, v...), nil, false)
}

// 报警
func (log *Logger) Alert(format string, v ...interface{}) {
	log.Write(ALERT, "A", fmt.Sprintf(format, v...), nil, false)
}

// 严重
func (log *Logger) Crit(format string, v ...interface{}) {
	log.Write(CRITICAL, "C", fmt.Sprintf(format, v...), nil, false)
}

// 错误
func (log *Logger) Err(format string, v ...interface{}) {
	log.Write(ERROR, "E", fmt.Sprintf(format, v...), nil, false)
}

// 警告
func (log *Logger) Warn(format string, v ...interface{}) {
	log.Write(WARNING, "W", fmt.Sprintf(format, v...), nil, false)
}

// 提示
func (log *Logger) Notice(format string, v ...interface{}) {
	log.Write(NOTICE, "N", fmt.Sprintf(format, v...), nil, false)
}

// 信息
func (log *Logger) Info(format string, v ...interface{}) {
	log.Write(INFO, "I", fmt.Sprintf(format, v...), nil, false)
}

// 调试
func (log *Logger) Debug(format string, v ...interface{}) {
	log.Write(DEBUG, "D", fmt.Sprintf(format, v...), nil, false)
}

// 紧急
func (log *Logger) SyncEmerg(format string, v ...interface{}) {
	log.Write(EMERGENCY, "M", fmt.Sprintf(format, v...), nil, true)
}

// 报警
func (log *Logger) SyncAlert(format string, v ...interface{}) {
	log.Write(ALERT, "A", fmt.Sprintf(format, v...), nil, true)
}

// 严重
func (log *Logger) SyncCrit(format string, v ...interface{}) {
	log.Write(CRITICAL, "C", fmt.Sprintf(format, v...), nil, true)
}

// 错误
func (log *Logger) SyncErr(format string, v ...interface{}) {
	log.Write(ERROR, "E", fmt.Sprintf(format, v...), nil, true)
}

// 警告
func (log *Logger) SyncWarn(format string, v ...interface{}) {
	log.Write(WARNING, "W", fmt.Sprintf(format, v...), nil, true)
}

// 提示
func (log *Logger) SyncNotice(format string, v ...interface{}) {
	log.Write(NOTICE, "N", fmt.Sprintf(format, v...), nil, true)
}

// 信息
func (log *Logger) SyncInfo(format string, v ...interface{}) {
	log.Write(INFO, "I", fmt.Sprintf(format, v...), nil, true)
}

// 调试
func (log *Logger) SyncDebug(format string, v ...interface{}) {
	log.Write(DEBUG, "D", fmt.Sprintf(format, v...), nil, true)
}

func New() *Logger {
	logger := &Logger{adapters: make(map[string]Adapter)}
	console := NewAdapterConsole(10000)
	go console.Run()
	logger.AddAdapter("default", console)
	return logger
}

// 创建默认日志方便调用
var logmo *Logger

func init() {
	logmo = New()
}

// 紧急
func Emerg(format string, v ...interface{}) {
	logmo.Emerg(format, v)
}

// 报警
func Alert(format string, v ...interface{}) {
	logmo.Alert(format, v)
}

// 严重
func Crit(format string, v ...interface{}) {
	logmo.Crit(format, v)
}

// 错误
func Err(format string, v ...interface{}) {
	logmo.Err(format, v)
}

// 警告
func Warn(format string, v ...interface{}) {
	logmo.Warn(format, v)
}

// 提示
func Notice(format string, v ...interface{}) {
	logmo.Notice(format, v)
}

// 信息
func Info(format string, v ...interface{}) {
	logmo.Info(format, v)
}

// 调试
func Debug(format string, v ...interface{}) {
	logmo.Info(format, v)
}

func AddAdapter(name string, adapter Adapter) error {
	return logmo.AddAdapter(name, adapter)
}

func DeleteAdapter(name string) error {
	return logmo.DeleteAdapter(name)
}

func GetAdapter(name string) (Adapter, error) {
	return logmo.GetAdapter(name)
}
