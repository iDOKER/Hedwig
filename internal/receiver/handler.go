package receiver

import (
	"Alarm2File/internal/common"
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
	"os"
	"path"
	"time"
)

// AlertEvent represents a single alert event
type AlertEvent struct {
	Content string   `json:"content"`
	WxPhone []string `json:"wxPhone"`
}

// AlertCallbackData represents the callback request body
type AlertCallbackData struct {
	Events []AlertEvent `json:"events"`
}

// StartServer starts the HTTP server to handle N9E callbacks
func StartServer(addr, outputDir string, filePrefix string, fileSuffix string, headerKey string, headerValue string, encryptToken string, dataBackupEnable bool, dataBackupDir string) error {
	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {

		// Check if the request contains the required headers
		if headerKey != "" && headerValue != "" {
			if r.Header.Get(headerKey) != headerValue {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
		}

		// Decode incoming JSON data
		var data AlertEvent
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			common.Logger.Errorf("Invalid JSON data: %v", err)
			http.Error(w, "Invalid JSON data", http.StatusBadRequest)
			return
		}

		common.Logger.Debugf("Received callback with Content length %d.", len(data.Content))
		common.Logger.Debugf("Received callback with Content: %s.", data.Content)
		common.Logger.Debugf("Received callback with WxPhone: %v.", data.WxPhone)

		// Save data to a JSON file
		fileName := path.Join(outputDir, time.Now().Format("20060102_150405_")+uuid.New().String()+".json")
		fileContent, _ := json.Marshal(data)
		err = saveEncryptedFile(fileName, filePrefix, fileSuffix, fileContent, encryptToken, dataBackupEnable, dataBackupDir)
		if err != nil {
			common.Logger.Errorf("Failed to save file: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		common.Logger.Debugf("Saved encrypted alert file: %s", fileName)
	})

	common.Logger.Infof("Receiver listening on %s/callback", addr)
	return http.ListenAndServe(addr, nil)
}

// saveEncryptedFile saves encrypted data to a file
func saveEncryptedFile(fileName string, filePrefix string, fileSuffix string, content []byte, encryptToken string, dataBackupEnable bool, dataBackupDir string) error {

	fileNameBackup := path.Join(dataBackupDir, path.Base(fileName))

	if dataBackupEnable {
		err := os.WriteFile(fileNameBackup+"_bak.json", content, 0644)
		if err != nil {
			common.Logger.Error("Failed to create backup file: %v", err)
		}
	}

	// Simulated encryption logic; replace with actual encryption if needed
	encryptedContent, err := encryptData(content, encryptToken)
	if err != nil {
		common.Logger.Error("Encryption failed: %v", err)
		return err
	}
	common.Logger.Debugf("Raw content: %s", content)
	common.Logger.Debugf("Encrypted content: %s with Token: %s", string(encryptedContent), encryptToken)

	// Save to file
	fileNameEncrypted := path.Join(path.Dir(fileName), filePrefix+path.Base(fileName)+fileSuffix)
	err = os.WriteFile(fileNameEncrypted, encryptedContent, 0644)
	if err != nil {
		return err
	}
	common.Logger.Debugf("Saved encrypted alert file: %s", fileNameEncrypted)
	return nil
}

// encryptData simulates encryption; replace with actual encryption logic
func encryptData(data []byte, token string) ([]byte, error) {
	// Simulated logic; can be replaced with real encryption
	key := common.GenerateKey(token)
	encryptedStr, err := common.UuidEncrypt(data, []byte(key))
	if err != nil {
		return nil, err
	}
	return []byte(encryptedStr), nil
}
