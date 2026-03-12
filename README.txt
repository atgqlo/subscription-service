REST API для управления подписками + подсчёт общей стоимости за период


POST	/subscriptions	Создать подписку
GET	/subscriptions	Список всех подписок
GET	/subscriptions/total	Общая стоимость за период
GET	/subscriptions/:id	Подписка по ID	Subscription
PUT	/subscriptions/:id	Обновить подписку	200 OK
DELETE	/subscriptions/:id	Удалить подписку	204 No Content


swagger: http://localhost:8080/swagger/index.html