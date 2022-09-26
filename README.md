# note-service

Это REST-api, позволяющий пользователям регистрироваться, авторизовываться, создавать, редактировать и удалять сообщения. 

# Архитектура

Проект разделен на две части - взаимодействия с пользователями и заметками, каждая часть имеет три раздела: router, service и store.

Router отвечает за вазимодействие с клиентской частью, он получает запрос, проверяет его на валидность, после чего вызывают service, и выдает ответ на основе полученных результатов.

Service отвечает за проверку правильности запроса. NoteService проверяет может ли пользователь увидеть заметки,и имеет ли он право редактирования/удаления. UserService отвечает за шифровку и расшифровки пароля.

Store отвечает за взаимодействие с данными, если операция проверена сервисом, то она выполняется store, который возвращает результат своей работу сервису, который в свою очередь возвращает результат роутеру.

Так же существует expiration service. Он запускается каждые десять секунд и удаляет сообщения ttl которых уже прошло.

Структура проекта сделана на основе https://github.com/golang-standards/project-layout

# Requests

Примеры всех запросов находятся в папке [postman](https://github.com/VorobevNickolay/note-service/tree/master/postman)

## User router

### Sign-up

'POST /user'

Позволяет зарегистрировать пользователя, получает имя и пароль пользователя, после чего создает пользователя и возвращает его.
### Sign-up

'POST /user/:id'

Позволяет войти пользователю, получает имя и пароль пользователя, после чего создает jwt-токен и возвращает его, все действия.

## Note router

Каждый из методов NoteRouter вызывает перед собой middle-ware функцию, которая получает jwt-токен и расшивровывает его в id пользователя. Таким образом, действия с заметками могут совершить только вошедшие пользователи

### GetNotes

'GET /notes'

Возвращает пользователю все его заметки, отсортированные по выбранному параметру, если параметр не был выбран, возвращает массив, отсортированный по ID.

### GetNoteByID

'GET /note/:id'

Возвращает пользователю заметку, если он создатель или если заметка публичная, или он включен в массив пользователей, кому дан доступ

### PostNote

'POST /note'

Позволяет создать заметку, обязательным параметром является только текст, другие параметры пользователь может указать при желании, или не указывать их вовсе

### UpdateNote

'PUT /note/:id'

Позволяет пользователю обновить свою заметку

### DeleteNote

'DELETE /note/:id'

Позволяет пользователю удалить свою заметку
