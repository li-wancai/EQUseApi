/*
Created on Fri Sep 16 17:04:36 2024
@author:liwancai

	QQ:248411282
	Tel:13199701121
*/
package EQUseApi

import (
	"github.com/li-wancai/logger"
	"github.com/li-wancai/wxwebhook/pkg/wechatsend"
)

func SendTxT(msg string,
	sendtogrouplist []string, //群webhook的key列表
	Mentioned []string, //@列表,或者@all
	Mobiles []string, //@的手机号列表
) {
	for _, towxid := range sendtogrouplist {
		client := wechatsend.NewClient(towxid)
		TextMessage := wechatsend.TextMessage{
			Content:   msg,
			Mentioned: Mentioned,
			Mobiles:   Mobiles,
		}
		if err := client.Send(TextMessage); err != nil {
			log.Errorf("发送消息到 %s 失败： %v", towxid, err)
		}
	}
}

var log *logger.LogN

func SetLogger(l *logger.LogN) {
	log = l //配置log信息
}

// SendLogTxT 是一个用于发送日志事务的函数。
//
// 它接受两个参数：
// - msg: 一个字符串，表示信息。
// - towxid: 一个字符串，表示目标微信 ID。
// - v: 一个可变参数，可以是 INFO, WARN, DEBUG, RECORD, FATAL, ERROR 等。
// 它不返回任何值。
func SendLogTxT(msg string, sendtogrouplist []string, v ...interface{}) {
	SendTxT(msg, sendtogrouplist, []string{}, []string{})
	level := "ERROR" // default log level
	if len(v) > 0 {
		if str, ok := v[0].(string); ok {
			level = str
		}
	}
	switch level {
	case "INFO":
		log.Info(msg)
	case "WARN":
		log.Warn(msg)
	case "DEBUG":
		log.Debug(msg)
	case "RECORD":
		log.Record(msg)
	case "FATAL":
		log.Fatal(msg)
	case "ERROR":
		log.Error(msg)
	default:
		log.Error(msg)
	}
}
