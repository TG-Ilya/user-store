# User store

## Тестовое задание.
Написать веб сервис, принимающий запросы по http api. 
Методы CRUD(create | read | update | delete). Для создания пользователей.
Пользователь имеет id, name, birth_date. Входными данными для RUD методов должен быть id пользователя. В create методе все кроме id. В методах просто создавать юзера в базе и что-нибудь выводить в лог(если потребуется).
Использовать подход specification first. В качестве спецификации использовать swagger 2.0. По спецификации генерировать серверный код. Использовать https://github.com/go-swagger/go-swagger.
В качестве базы использовать sqlite. Не использовать ORM.
В проекте должен быть dockerfile,  с помощью которого должен собираться образ с сервисом и впоследствии запускаться.
Все необходимые параметры для конфигурации приложения передавать через env.
Тесты по желанию, но будет не лишним.
Не забыть про миграции в базе.

Исходный код выложить на public гитхаб / bitbucket.
