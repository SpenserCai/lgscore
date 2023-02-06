/*
 * @Author: SpenserCai
 * @Date: 2023-02-06 14:57:05
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-02-06 15:05:51
 * @Description: file content
 */
package utility

import (
	"io"
	"net/http"
	"os"
)

type NewWorkLGS struct{}

func (n *NewWorkLGS) Download(downloadUrl, savePath string) error {
	resp, err := http.Get(downloadUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// 创建临时文件
	dlFile, err := os.Create(savePath)
	if err != nil {
		return err
	}
	defer dlFile.Close()
	// 将下载的内容写入临时文件
	_, err = io.Copy(dlFile, resp.Body)
	if err != nil {
		return err
	}
	return nil
}
