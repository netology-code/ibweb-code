-- если указываем только имя таблицы, то обязаны заполнять все поля
INSERT INTO clients
VALUES (1, 'petr', 'password-hash', 'Пётр Николаевич Иванов', 'ACTIVE', NOW());

SELECT *
FROM clients;

DELETE
FROM clients
WHERE id = 1;

-- если указываем имя таблицы и столбцы,
-- то обязаны заполнять только то, что указали
-- остальное примет значение по умолчанию либо NULL
-- (если не установлено NOT NULL)
-- ключевое слово DEFAULT позволяет использовать значения по умолчанию при необходимости
INSERT INTO clients(login, password, name, status)
VALUES ('petr', 'password-hash', 'Пётр Николаевич Иванов', DEFAULT);


-- вставка сразу нескольких записей
INSERT INTO clients(login, password, name, status)
VALUES ('vasya', 'password-hash', 'Василий Николаевич Иванов', 'ACTIVE'),
       ('masha', 'password-hash', 'Мария Ивановна Петрова', 'ACTIVE'),
       ('dasha', 'password-hash', 'Дарья Ивановна Крылова', 'ACTIVE')
;


