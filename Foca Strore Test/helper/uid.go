package helper

import (
	"fmt"
	"strconv"
	"time"
	"voca-store/models"

	"gorm.io/gorm"
)

func GenerateAddressUID(tx *gorm.DB) (string, error) {
    today := time.Now().Format("20060102")
    prefix := fmt.Sprintf("ADDR-%s-", today)

    var lastAddress models.Address
    
    err := tx.Set("gorm:query_option", "FOR UPDATE").
        Unscoped().
        Where("uid LIKE ?", prefix+"%").
        Order("uid DESC").
        First(&lastAddress).Error

    nextSeq := 1

    if err == nil {
        uidStr := lastAddress.UID
        if len(uidStr) >= 4 {
            lastSeqStr := uidStr[len(uidStr)-4:]
            seq, parseErr := strconv.Atoi(lastSeqStr)
            if parseErr == nil {
                nextSeq = seq + 1
            }
        }
    } else if err != gorm.ErrRecordNotFound {
        return "", fmt.Errorf("failed to query last address UID: %w", err)
    }

    uid := fmt.Sprintf("ADDR-%s-%04d", today, nextSeq)
    return uid, nil
}