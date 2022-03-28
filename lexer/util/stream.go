package util

import (
	"bufio"
	"container/list"
	"io"
)

type Stream struct {
	scanner    *bufio.Scanner //扫描每行
	queueCache *list.List     // queue队列缓存
	endToken   string
	isEnd      bool
}

func NewStream(r io.Reader, et string) *Stream {
	scan := bufio.NewScanner(r) //根据换行来扫描每一行
	scan.Split(bufio.ScanRunes) // uft-8编码扫描
	return &Stream{
		scanner:    scan,
		queueCache: list.New(),
		endToken:   et,
		isEnd:      false,
	}
}
