package projlib

import "github.com/yunya101/perepisochnik/internal/models"

func RemoveElement(slice []*models.Message, s int) []*models.Message {
	return append(slice[:s], slice[s+1:]...)
}
