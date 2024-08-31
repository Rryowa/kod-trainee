## Привет! Меня зовут Антон и я очень надеюсь, что вам понравится мое решение!
## Если нет - пожалуйста, дайте знать что нужно исправить!

### Конфигурация в .env файле

### Запуск:
    docker-compose up -d
    make all
    используй коллекцию постман приведенную ниже

### Остановка:
    ctrl + с
    make down
    docker-compose down -v

### Заметки: NoteService - internal/service/note.go
    1) AddNote - Валидирует заголовок и текст заметки с помощью Yandex Speller.
        Затем добавляет заметку в базу.
        Использует интерфейс Storage для взаимодействия с бд.
    2) GetNotes - Выводит список заметок пользователя постранично и отсортированный по дате.
        Последние доабвленные являются первыми выведенными.
        Использует интерфейс Storage для взаимодействия с бд.

### Авторизация: UserService - internal/service/user.go
    1) SingUp - Создает юзера.
    2) LogIn - Проверяет юзера и создает куки
        Использует SessionService
        Который создает jwt токен и записывает его в куки.
    3) LogOut - Удаляет куки
        Использует SessionService
        Который создает такой же куки но с пустым value.

### Аутентификация: Middleware - internal/middleware
    Без валидного JWT токена в куки не получится ничего сделать.
    Использует SessionService - internal/service/session.go.
    Который достает JWT токен из куки и проверяет его.
    Затем записывает данные пользователя из JWT токена в Context

### Трейсинг - http://localhost:16686/ service - kod

### Postman коллекция - https://www.postman.com/rryowa/workspace/kod/collection/27242165-ca26f13d-a4e4-4104-990d-3512e8f03c77?action=share&creator=27242165