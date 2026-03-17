REST API для управления подписками + подсчёт общей стоимости за период


Запуск через docker:
docker-compose up --build -d

Endpoints:
POST	/subscriptions	Создать подписку
GET	/subscriptions	Список всех подписок
GET	/subscriptions/total	Общая стоимость за период
GET	/subscriptions/:id	Подписка по ID	Subscription
PUT	/subscriptions/:id	Обновить подписку	200 OK
DELETE	/subscriptions/:id	Удалить подписку	204 No Content


