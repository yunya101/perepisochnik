package projlib

import "github.com/yunya101/perepisochnik/internal/models"

func RemoveElementFromSlice(slice []*models.Message, s int) []*models.Message {
	return append(slice[:s], slice[s+1:]...)
}

func InsertMsg(slice []*models.Message, msg *models.Message) []*models.Message {
	newSlice := []*models.Message{msg}

	newSlice = append(newSlice, slice...)

	return newSlice
}
