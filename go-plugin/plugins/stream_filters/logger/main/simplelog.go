package main

import (
	"bytes"
	"log"
	"time"

	"github.com/natefinch/lumberjack"
)

type logger struct {
	chanTags chan [TAG_END]string
	log      *log.Logger
}

func NewLogger(path string, config map[string]string) (*logger, error) {
	ll := log.New(&lumberjack.Logger{
		Filename:  path,
		MaxSize:   configIntValue("max_size", config, maxSize), // 单位为MB,默认为1MB
		MaxAge:    configIntValue("max_age", config, maxAge),   // 文件最多保存3天
		LocalTime: true,                                        // 采用本地时间
		Compress:  false,                                       // 是否压缩日志
	}, "", log.Lmsgprefix)

	l := &logger{
		chanTags: make(chan [TAG_END]string, 1024),
		log:      ll,
	}
	go l.loop()
	return l, nil
}

func (l *logger) Print(tags [TAG_END]string) {
	l.chanTags <- tags
}

func (l *logger) Close() {
}

func (l *logger) loop() {
	for tags := range l.chanTags {
		l.print(tags)
	}
}

func (l *logger) print(tags [TAG_END]string) {
	currentTime := time.Now()
	var buffer bytes.Buffer
	buffer.WriteString(currentTime.Format(glogTmFmtWithMS))
	buffer.WriteString(("||"))
	buffer.WriteString(`{`)
	for index, value := range tags {
		if key, ok := tagsName[index]; ok && len(value) != 0 {
			if index != 0 {
				buffer.WriteString(`,"`)
			} else {
				buffer.WriteString(`"`)
			}
			buffer.WriteString(key)
			buffer.WriteString(`":"`)
			buffer.WriteString(value)
			buffer.WriteString(`"`)
		}
	}
	buffer.WriteString(`}`)
	l.log.Println(buffer.String())
}
