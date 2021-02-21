-- for simplicity use just names
INSERT INTO clients(login, password, name, status)
VALUES ('vasya', '$2b$12$9VdrEg13oND9Hhf1FpDOaesPOLm3t5OWtSKupTymMeVhGmUDJ7JES', 'Василий', 'ACTIVE'),
       ('masha', '$2b$12$9VdrEg13oND9Hhf1FpDOaesPOLm3t5OWtSKupTymMeVhGmUDJ7JES', 'Мария', 'ACTIVE'),
       ('dasha', '$2b$12$9VdrEg13oND9Hhf1FpDOaesPOLm3t5OWtSKupTymMeVhGmUDJ7JES', 'Дарья', 'ACTIVE'),
       ('sasha', '$2b$12$9VdrEg13oND9Hhf1FpDOaesPOLm3t5OWtSKupTymMeVhGmUDJ7JES', 'Анна', 'ACTIVE'),
       ('petya', '$2b$12$9VdrEg13oND9Hhf1FpDOaesPOLm3t5OWtSKupTymMeVhGmUDJ7JES', 'Пётр', 'ACTIVE')
;

-- for simplicity just last numbers & russian holder names
INSERT INTO cards(number, balance, issuer, holder, owner_id, status)
VALUES ('**01', 10000, 'MIR', 'Василий', 1, 'ACTIVE'),
       ('**02', 10000, 'MIR', 'Мария', 2, 'ACTIVE'),
       ('**03', 10000, 'MIR', 'Дарья', 3, 'ACTIVE'),
       ('**04', 10000, 'MIR', 'Анна', 4, 'ACTIVE'),
       ('**05', 10000, 'MIR', 'Пётр', 5, 'ACTIVE')
;


