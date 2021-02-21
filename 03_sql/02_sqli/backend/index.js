const express = require('express');
const cors = require('cors');
const {Pool} = require('pg');

const dsn = process.env.DSN || 'postgres://app:pass@localhost:5432/db';
// named as db for demonstration simplicity
const db = new Pool({ connectionString: dsn });

const app = express();
app.use((req, res, next) => {
    console.log(req.url);
    next();
})
app.use(cors());
app.use(express.json()); // вычитывает JSON из тела запроса и складывает в req.body
app.get('/health', async (req, res) => {
    try {
        await db.query('SELECT 1');
        res.status(200).send();
    } catch (e) {
        res.status(500).send();
    }
});

// запросы вида http://localhost:9999/sqli/getAll
// req - запрос
// res - ответ
app.get('/sqli/getAll', async (req, res) => {
    try {
        const sql = `
                SELECT id, name
                FROM clients
                WHERE status = 'ACTIVE' AND removed = false
                ORDER BY id
            `;
        const data = await db.query(sql);
        // отправляем клиенту данные в формате JSON
        res.send(data.rows);
    } catch (e) {
        console.error(e);
        res.sendStatus(500);
    }
});

// запросы вида http://localhost/sqli/getAllByName?name=Василий
// req - запрос
// res - ответ
app.get('/sqli/getAllByName', async (req, res) => {
    try {
        const name = req.query.name; // name = Василий
        console.log(name)
        const sql = `
                SELECT id, name
                FROM clients
                WHERE status = 'ACTIVE' AND removed = false
                AND name = '${name}'
                ORDER BY id
            `;
        const data = await db.query(sql);
        // отправляем клиенту данные в формате JSON
        res.send(data.rows);
    } catch (e) {
        console.error(e);
        res.sendStatus(500);
    }
});

// запросы вида http://localhost/sqli/getAllByNameSecured?name=Василий
// req - запрос
// res - ответ
app.get('/sqli/getAllByNameSecured/:id', async (req, res) => {
    try {
        const name = req.query.name; // name = Василий
        console.log(name)
        const sql = `
                SELECT id, name
                FROM clients
                WHERE status = 'ACTIVE' AND removed = false
                AND name = $1
                ORDER BY id
            `;
        const data = await db.query(sql, [name]);
        // отправляем клиенту данные в формате JSON
        res.send(data.rows);
    } catch (e) {
        console.error(e);
        res.sendStatus(500);
    }
});

const port = process.env.PORT || 9999;
app.listen(port, () => {
    console.log('server started');
});
