# Perepisochnik - backend для приватного общения в реальном времени
## Requests:
- "/" (с заголовком 'username') - подключиться к серверу
- "/chat{id}" - написать сообщение в чат
  
### Приложение принимает и отдаёт JSON объекты:
__Message__  
{
 "id"
 "reciver"
 "recipient"
 "text"
 "chat_id"
}  
__User__  
{"id" "username" "chats"}  
__Chat__  
{"id" "users" "messages"}

#### Приложение использует базу данных Postgres, сохраняя пользователей, чаты и сообщения
