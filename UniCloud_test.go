package unicloud

import (
	"fmt"
	"testing"
)

func TestUpload(t *testing.T) {
	url := Upload("test.png")
	fmt.Printf("返回的cdn地址为：%v", url)
}
