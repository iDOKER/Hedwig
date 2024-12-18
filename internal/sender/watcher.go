package sender

import (
	"Alarm2File/internal/common"
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func WatchAndForward(dir, filePrefix string, fileSuffix string, targetURL string, timeout int, retryCount int, retryInterval int, sslVerify bool, headerKey string, headerValue string, encryptToken string, dataBackupEnable bool, dataBackupDir string) error {
	for {
		files, err := ioutil.ReadDir(dir)
		if err != nil {
			common.Logger.Errorf("Failed to read directory: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		for _, file := range files {
			fileName := file.Name()
			if strings.HasSuffix(fileName, fileSuffix) && strings.HasPrefix(fileName, filePrefix) {
				filePath := filepath.Join(dir, file.Name())

				// 解密和转发
				err := processFile(filePath, targetURL, timeout, retryCount, retryInterval, sslVerify, headerKey, headerValue, encryptToken)
				if err != nil {
					common.Logger.Warnf("Failed to process file %s: %v", filePath, err)
					err := os.Rename(filePath, filepath.Join(dataBackupDir, file.Name()+".http400"))
					if err != nil {
						return err
					}
				} else {
					if dataBackupEnable {
						common.Logger.Debugf("Backing up file %s", filePath)
						err := os.Rename(filePath, filepath.Join(dataBackupDir, file.Name()))
						if err != nil {
							common.Logger.Errorf("Failed to rename file: %v", err)
							return err
						}
					} else {
						common.Logger.Debugf("Deleteing file %s", filePath)
						err := os.Remove(filePath)
						if err != nil {
							common.Logger.Errorf("Failed to delete file: %v", err)
							return err
						} // 删除已处理的文件
					}
				}
			}
		}
		time.Sleep(5 * time.Second)
	}
}

func processFile(filePath, targetURL string, timeout int, retryCount int, retryInterval int, sslVerify bool, headerKey string, headerValue string, encryptToken string) error {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		common.Logger.Errorf("failed to read encrypted file: %v", err)
		return fmt.Errorf("failed to read encrypted file: %v", err)
	}

	common.Logger.Debugf("Decrypting data with token: [%s]", encryptToken)
	key := common.GenerateKey(encryptToken)
	decryptedData, err := common.UuidDecrypt(string(data), []byte(key))
	if err != nil {
		common.Logger.Warnf("failed to decrypt data: %v", err)
		return fmt.Errorf("failed to decrypt data: %v", err)
	}

	common.Logger.Debugf("Forwarding data: [%+v] to [%s]", decryptedData, targetURL)

	// 非 Json 数据解密逻辑
	/*	var decodeData map[string]interface{}
		err = json.Unmarshal(data, &decodeData)
		if err != nil {
			return err
		}*/

	// 创建 HTTP 客户端
	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: !sslVerify},
		},
	}

	for attempt := 0; attempt < retryCount; attempt++ {

		// 创建 HTTP 请求
		// 第二次请求时候请求体会被消耗，需要重新创建请求体和请求对象
		// 在 http.NewRequest 中使用的 bytes.Buffer 或其他类型的 io.Reader 在发送请求时会被消耗掉。这意味着如果你在多次发送请求时使用同一个 bytes.Buffer，第二次发送时它将为空。

		req, err := http.NewRequest("POST", targetURL, bytes.NewBuffer([]byte(decryptedData)))
		if err != nil {
			return err
		}

		// 设置请求头
		req.Header.Set(headerKey, headerValue)
		req.Header.Set("Content-Type", "application/json")

		// 发送请求
		resp, err := client.Do(req)
		if err != nil {
			common.Logger.Warnf("Attempt %d: Request failed: %v", attempt+1, err)
			if attempt < retryCount-1 {
				time.Sleep(time.Duration(retryInterval) * time.Second)
				continue
			}
			return err
		}

		// 立即关闭 resp.Body
		if resp != nil {
			err := resp.Body.Close()
			if err != nil {
				common.Logger.Warnf("Failed to close response body: %v", err)
			}
		}

		// 处理响应
		if resp.StatusCode != http.StatusOK {
			body, _ := ioutil.ReadAll(resp.Body)
			common.Logger.Warnf("Response status: %s", resp.Status)
			common.Logger.Debugf("Failed to forward data: %s", body)
			if attempt < retryCount-1 {
				time.Sleep(time.Duration(retryInterval) * time.Second)
				continue
			}
			return fmt.Errorf("failed to forward data: %s", body)
		}

		// 如果请求成功，退出循环
		break
	}

	/*
		// 发送请求
		defer 在函数返回之前它不会执行
		https://stackoverflow.com/questions/45617758/proper-way-to-release-resources-with-defer-in-a-loop
		The whole point of defer is that it does not execute until the function returns, so the appropriate place to put it would be immediately after the resource you want to close is opened. However, since you're creating the resource inside the loop, you should not use defer at all - otherwise, you're not going to close any of the resources created inside the loop until the function exits, so they'll pile up until then. Instead, you should close them at the end of each loop iteration, without defer:
	*/
	/*resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			common.Logger.Errorf("Failed to close response body: %v", err)
		}
	}(resp.Body)*/

	common.Logger.Debugf("Successfully forwarded data to %s", targetURL)
	return nil
}

func decryptData(data []byte, token string) ([]byte, error) {
	// Simulated logic; can be replaced with real encryption
	key := common.GenerateKey(token)
	decryptedStr, err := common.UuidDecrypt(string(data), []byte(key))
	if err != nil {
		return nil, err
	}
	return []byte(decryptedStr), nil
}
