package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
)

func main() {
	tmp := NewEmailSender("smtp.yandex.ru", 465, "pyaterochka.team@yandex.ru", "пароль из лички", logrus.New().WithField("service", "mail"))
	go func() {
		tmp.Run()
	}()
	tmp.SendMessage([]string{"vet_v2002@mail.ru"}, HTMLString(htmlBody))
	str := ""
	fmt.Scanf("%s", &str)
	tmp.Stop()
}
