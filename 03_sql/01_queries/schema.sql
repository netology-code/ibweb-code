CREATE TABLE clients
(
    id        BIGSERIAL PRIMARY KEY,
    login     TEXT      NOT NULL UNIQUE,
    password  TEXT      NOT NULL,
    name      TEXT      NOT NULL,
    status    TEXT      NOT NULL DEFAULT 'INACTIVE' CHECK (status IN ('INACTIVE', 'ACTIVE')),
    created   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- комментарий
-- пример запроса, модифицирующего таблицу (выполнять не нужно)
-- ALTER TABLE clients ADD COLUMN phone TEXT NOT NULL;

CREATE TABLE cards
(
    id       BIGSERIAL PRIMARY KEY,
    number   TEXT      NOT NULL,
    balance  BIGINT    NOT NULL DEFAULT 0,
    issuer   TEXT      NOT NULL CHECK ( issuer IN ('Visa', 'MasterCard', 'MIR') ),
    holder   TEXT      NOT NULL,
    owner_id BIGINT    NOT NULL REFERENCES clients,
    status   TEXT      NOT NULL DEFAULT 'INACTIVE' CHECK (status IN ('INACTIVE', 'ACTIVE')),
    created  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);


