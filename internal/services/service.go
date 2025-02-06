package services

import (
	"github.com/yunya101/perepisochnik/internal/data"
	"github.com/yunya101/perepisochnik/internal/models"
)

type Service struct {
	repo *data.Repository
}

func (s *Service) SetRepo(r *data.Repository) {
	s.repo = r
}

func (s *Service) AddMsg(msg *models.Message) error {
	err := s.repo.InsertMsg(msg)

	if err != nil {
		return err
	}

	return nil
}

func (s *Service) GetUsersChats(username string) (*models.User, error) {

	chats, err := s.repo.GetChatsByUsername(username)

	if err != nil {
		return nil, err
	}

	for i, chat := range chats {
		chat.Users = make([]string, 0)
		chat.Messages = make([]models.Message, 0)

		chat, err := s.repo.GetUsersFromChat(chat)

		if err != nil {
			return nil, err
		}

		chat, err = s.repo.GetMsgsFromChat(chat)

		if err != nil {
			return nil, err
		}

		chats[i] = chat

	}

	user := &models.User{
		Username: username,
		Chats:    chats,
	}

	return user, nil
}

func (s *Service) GetUsersFromChat(chat *models.Chat) (*models.Chat, error) {

	chat, err := s.repo.GetUsersFromChat(chat)

	if err != nil {
		return nil, err
	}

	return chat, nil

}
