package main

import (
	"fmt"
	"log"
	"os"
	"time"

	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
	"github.com/joho/godotenv"
)

func main() {
	var Bot *tgbotapi.BotAPI
	if _, err := os.Stat(".env"); os.IsNotExist(err) {
		// 如果 .env 文件不存在，创建并写入默认值
		fmt.Printf(".env 文件不存在，正在创建...")

		// 创建并打开 .env 文件
		file, err := os.Create(".env")
		if err != nil {
			log.Fatalf("创建 .env 文件失败: %v", err)
		}
		defer file.Close()

		// 写入默认的环境变量内容
		defaultEnv := `Token=
`
		if _, err := file.WriteString(defaultEnv); err != nil {
			log.Fatalf("写入 .env 文件失败: %v", err)
		}
		fmt.Printf(".env 文件已创建，并写入默认内容.")
	}
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("%s", err)
	}
	Token := os.Getenv("Token")
	Bot, Err := tgbotapi.NewBotAPI(Token)
	if Err != nil {
	}
	Bot.Debug = true
	updatecfg := tgbotapi.NewUpdate(0)
	updatecfg.Timeout = 60
	updates := Bot.GetUpdatesChan(updatecfg)
	for update := range updates {
		if update.Message != nil { // 如果有消息
			// 检查是否是新成员加入
			if update.Message.NewChatMembers != nil {
				for _, user := range update.Message.NewChatMembers {
					chatID := update.Message.Chat.ID

					// 封禁新用户
					_, err := Bot.BanChatMember(chatID, user.ID)
					if err != nil {
						log.Printf("封禁用户 %s 失败: %s", user.UserName, err)
						continue
					}
					log.Printf("封禁用户 %s", user.UserName)

					// 5 秒后解除封禁
					go func(userID int64) {
						time.Sleep(5 * time.Second)
						_, err := Bot.UnbanChatMember(chatID, userID)
						if err != nil {
							log.Printf("解除封禁用户 %s 失败: %s", user.UserName, err)
						} else {
							log.Printf("解除封禁用户 %s", user.UserName)
						}
					}(user.ID)
				}
			}
		}
	}
}
