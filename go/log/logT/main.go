package main

import (
	"github.com/mkideal/log"
)

func main() {
	// Init and defer Uninit
	defer log.Uninit(log.InitFile("./log/app.log"))

	// 默认日志等级是 INFO, 可以按以下方式修改等级:
	//
	//	log.SetLevel(log.LvTRACE)
	// 	log.SetLevel(log.LvDEBUG)
	// 	log.SetLevel(log.LvINFO)
	// 	log.SetLevel(log.LvWARN)
	// 	log.SetLevel(log.LvERROR)
	// 	log.SetLevel(log.LvFATAL)

	log.Trace("%s cannot be printed", "TRACE")
	log.Debug("%s cannot be printed", "DEBUG")

	log.Info("%s should be printed", "INFO")
	log.Warn("%s should be printed", "WARN")
	log.Error("%s should be printed", "ERROR")

	log.If(true).Info("%v should be printed", true)

	// 这个 if, else if, else 的特性有时候可以简化代码
	iq := 250
	log.If(iq < 250).Info("IQ less than 250").
		ElseIf(iq > 250).Info("IQ greater than 250").
		Else().Info("IQ equal to 250")

	log.With("hello").Info("With a string field")
	log.With(1).Info("With an int field")
	log.With(true).Info("With a bool field")
	log.With(1, "2", false).Info("With %d fields", 3)

	// log.M 其实就是 map[string]interface{}
	log.With(log.M{
		"a": 1,
		"b": "hello",
		"c": true,
	}).Info("With a map")

	// 用 JSON 格式输出 With 的数据
	log.WithJSON(log.M{
		"a": 1,
		"b": "hello",
		"c": true,
	}).Info("With a map and using JSONFormatter")
	// 还可以使用别的格式:
	// log.With(data).SetFormatter(formatter).Info("...")
	// formatter 是一个实现了 log.Formatter 接口的任意对象
	//	type Formatter interface {
	//		Format(v interface{}) []byte
	//	}

	// 携带上下文数据反复使用
	ctxlogger := log.With(log.M{"module": "hello"})
	ctxlogger.Info("I am in module hello")
	ctxlogger.Warn("I am in module hello")

	// NoHeader 函数将禁止输出日志头，就是时间啊，文件名，行号之类的东西
	log.NoHeader()

	log.Info("This message have no header")

	log.Fatal("%s should be printed and exit program with status code 1", "FATAL")

	log.Info("You cannot see me")
}