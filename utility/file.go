/*
 * @Author: SpenserCai
 * @Date: 2023-02-06 14:57:35
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-02-06 15:41:32
 * @Description: file content
 */
package utility

import (
	"archive/zip"
	"io"
	"os"
)

type FileLGS struct {
}

func (f *FileLGS) CopyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	// 创建目标文件
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()
	// 拷贝文件
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}
	return nil
}

// 将所有文件解压到dst目录下，如果遇到文件夹则递归解压，将文件路径存入fileList
func (f *FileLGS) UnZip(src string, dst string, fileList []string) error {
	// 判断dst的值后面是否有"/"，如果没有则添加
	if dst[len(dst)-1:] != "/" {
		dst = dst + "/"
	}
	zipReader, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer zipReader.Close()
	for _, file := range zipReader.File {
		if file.FileInfo().IsDir() {
			continue
		}
		// 解压文件到flingTrainerPath
		fileReader, err := file.Open()
		if err != nil {
			return err
		}
		defer fileReader.Close()
		fileWriter, err := os.Create(dst + file.Name)
		if err != nil {
			return err
		}
		defer fileWriter.Close()
		_, err = io.Copy(fileWriter, fileReader)
		if err != nil {
			return err
		}
		fileList = append(fileList, dst+file.Name)
	}
	return nil
}
