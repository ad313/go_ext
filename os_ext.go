package GoExt

import (
	"fmt"
	"os"
)

func deleteTmpFilesByPath(path string) error {
	err := os.RemoveAll(path)
	if err != nil {
		fmt.Println("删除文件夹失败：", err)
		return err
	} else {
		//fmt.Println("成功删除文件夹及其内容")
	}

	return nil
}
