package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"log"
	"os"
	"sync"
	"time"
)

const (
	defaultDSN = "postgres://app:pass@localhost:5432/db"
	delay      = time.Second * 5
	longDelay = delay * 5
)

func main() {
	dsn, ok := os.LookupEnv("DSN")
	if !ok {
		dsn = defaultDSN
	}

	ctx := context.Background()

	cmd := ""
	if len(os.Args) == 2 {
		cmd = os.Args[1]
	}

	switch cmd {
	case "sample1":
		noTx(ctx, dsn)
		break
	case "sample2":
		txErr(ctx, dsn)
	case "sample3":
		txAnomaly(ctx, dsn)
	default:
		log.Printf("wait for postgres start")
		// wait for postgres start
		for i := 0; i < 10; i++ {
			time.Sleep(delay)
			conn, err := pgx.Connect(ctx, dsn)
			if err != nil {
				log.Printf("retry")
				continue
			}
			_, err = conn.Exec(ctx, "SELECT 1 FROM cards")
			if err == nil {
				break
			}
		}
		log.Printf("ok, you can execute commands")
		time.Sleep(time.Hour)
	}
}

func noTx(ctx context.Context, dsn string) {
	var wg sync.WaitGroup
	wg.Add(2)
	defer wg.Wait()

	ownerId := 1

	go func() {
		defer wg.Done()
		conn, err := pgx.Connect(ctx, dsn)
		if err != nil {
			log.Fatal(err)
		}
		// закрытие соединения (игнорируем ошибки для простоты)
		defer conn.Close(ctx)
		printf := First

		// эмуляция первого запроса
		// (например первый клиент делает запрос на web-сервер)
		sum := 10_000
		printf("Клиент переводит деньги со своего счёта на чужой счёт")
		printf("Проверяем, достаточно ли средств")
		var balance int
		conn.QueryRow(ctx, "SELECT balance FROM cards WHERE owner_id = $1", ownerId).Scan(&balance)
		if balance < sum {
			printf("У клиента недостаточно средств")
			os.Exit(-1)
		}
		printf("Средств достаточно, эмулируем задержку")
		time.Sleep(longDelay)
		printf("Списываем средства")
		conn.QueryRow(ctx, "UPDATE cards SET balance = balance - $2 WHERE owner_id = $1 RETURNING balance", ownerId, sum).Scan(&balance)
		printf("Итоговый баланс (не должен был уйти в 0): ", balance)
	}()

	go func() {
		time.Sleep(delay)
		defer wg.Done()
		conn, err := pgx.Connect(ctx, dsn)
		if err != nil {
			log.Fatal(err)
		}
		// закрытие соединения (игнорируем ошибки для простоты)
		defer conn.Close(ctx)
		printf := Second

		// эмуляция первого запроса
		// (например второй клиент делает запрос на web-сервер)
		sum := 60
		printf("В это время параллельно списываем за обслуживание 60 рублей (стартуем позже)")
		printf("Проверяем, достаточно ли средств")
		var balance int
		conn.QueryRow(ctx, "SELECT balance FROM cards WHERE owner_id = $1", ownerId).Scan(&balance)
		if balance < sum {
			printf("У клиента недостаточно средств")
			os.Exit(-1)
		}
		printf("Средств достаточно, списываем средства")
		conn.Exec(ctx, "UPDATE cards SET balance = balance - $2 WHERE owner_id = $1", ownerId, sum)
		printf("Средства списаны")
	}()
}

func txErr(ctx context.Context, dsn string) {
	var wg sync.WaitGroup
	wg.Add(2)

	ownerId := 2

	go func() {
		time.Sleep(delay)
		defer wg.Done()
		conn, err := pgx.Connect(ctx, dsn)
		if err != nil {
			log.Fatal(err)
		}
		// закрытие соединения (игнорируем ошибки для простоты)
		defer conn.Close(ctx)
		printf := First

		tx, err := conn.BeginTx(ctx, pgx.TxOptions{
			IsoLevel: pgx.Serializable,
		})
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			printf("Пробуем зафиксировать транзакцию")
			if cerr := tx.Commit(ctx); cerr != nil {
				printf("Не удалось зафиксировать, откатываемся: ", cerr)
				tx.Rollback(ctx)
				return
			}
			printf("Ok, транзакция зафиксирована")
		}()

		// эмуляция первого запроса
		// (например первый клиент делает запрос на web-сервер)
		sum := 5_000
		printf("Клиент переводит деньги со своего счёта на чужой счёт (стартуем позже)")
		printf("Проверяем, достаточно ли средств (все запросы внутри транзакции)")
		var balance int
		tx.QueryRow(ctx, "SELECT balance FROM cards WHERE owner_id = $1", ownerId).Scan(&balance)
		if balance < sum {
			printf("У клиента недостаточно средств")
			os.Exit(-1)
		}
		printf("Средств достаточно, списываем средства")
		tx.QueryRow(ctx, "UPDATE cards SET balance = balance - $2 WHERE owner_id = $1 RETURNING balance", ownerId, sum).Scan(&balance)
		printf("Итоговый баланс: ", balance)
	}()

	go func() {
		defer wg.Done()
		conn, err := pgx.Connect(ctx, dsn)
		if err != nil {
			log.Fatal(err)
		}
		// закрытие соединения (игнорируем ошибки для простоты)
		defer conn.Close(ctx)
		printf := Second

		tx, err := conn.BeginTx(ctx, pgx.TxOptions{
			IsoLevel: pgx.Serializable,
		})
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			printf("Пробуем зафиксировать транзакцию")
			if cerr := tx.Commit(ctx); cerr != nil {
				printf("Не удалось зафиксировать, откатываемся: ", cerr)
				tx.Rollback(ctx)
				return
			}
			printf("Ok, транзакция зафиксирована")
		}()

		// эмуляция второго запроса
		// (например сервис списания автоплатежей делает запрос на сервер)
		sum := 60
		printf("В транзакции списываем за обслуживание 60 рублей")
		printf("Проверяем, достаточно ли средств")
		var balance int
		tx.QueryRow(ctx, "SELECT balance FROM cards WHERE owner_id = $1", ownerId).Scan(&balance)
		if balance < sum {
			printf("У клиента недостаточно средств")
			os.Exit(-1)
		}
		printf("Средств достаточно, эмулируем задержку")
		time.Sleep(longDelay)
		printf("Списываем средства")
		tx.QueryRow(ctx, "UPDATE cards SET balance = balance - $2 WHERE owner_id = $1 RETURNING balance", ownerId, sum).Scan(&balance)
		printf("Средства списаны, баланс после списания: ", balance)
	}()

	wg.Wait()
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		log.Fatal(err)
	}
	// закрытие соединения (игнорируем ошибки для простоты)
	defer conn.Close(ctx)
	var balance int
	conn.QueryRow(ctx, "SELECT balance FROM cards WHERE owner_id = $1", ownerId).Scan(&balance)
	fmt.Printf("Итоговый баланс, после всех списаний: %d\n", balance)
}

func txAnomaly(ctx context.Context, dsn string) {
	var wg sync.WaitGroup
	wg.Add(2)

	ownerId := 3

	go func() {
		defer wg.Done()
		conn, err := pgx.Connect(ctx, dsn)
		if err != nil {
			log.Fatal(err)
		}
		// закрытие соединения (игнорируем ошибки для простоты)
		defer conn.Close(ctx)
		printf := First

		// эмуляция первого запроса
		// (например первый клиент делает запрос на web-сервер)

		tx, err := conn.BeginTx(ctx, pgx.TxOptions{
			IsoLevel: pgx.ReadCommitted, // for nonrepeatable read
		})
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			printf("Пробуем зафиксировать транзакцию")
			if cerr := tx.Commit(ctx); cerr != nil {
				printf("Не удалось зафиксировать, откатываемся: ", cerr)
				tx.Rollback(ctx)
				return
			}
			printf("Ok, транзакция зафиксирована")
		}()

		sum := 5_000
		printf("Клиент переводит деньги со своего счёта на чужой счёт")
		printf("Проверяем, достаточно ли средств (все запросы внутри транзакции)")
		var balance int
		tx.QueryRow(ctx, "SELECT balance FROM cards WHERE owner_id = $1", ownerId).Scan(&balance)
		if balance < sum {
			printf("У клиента недостаточно средств")
			os.Exit(-1)
		}
		printf("Средств достаточно, списываем средства")
		tx.QueryRow(ctx, "UPDATE cards SET balance = balance - $2 WHERE owner_id = $1 RETURNING balance", ownerId, sum).Scan(&balance)
		printf("Эмулируем задержку перед коммитом, эта транзакция видит на счету: ", balance)
		time.Sleep(longDelay)
		printf("Итоговый баланс: ", balance)
	}()

	go func() {
		time.Sleep(delay)
		defer wg.Done()
		conn, err := pgx.Connect(ctx, dsn)
		if err != nil {
			log.Fatal(err)
		}
		// закрытие соединения (игнорируем ошибки для простоты)
		defer conn.Close(ctx)
		printf := Second

		tx, err := conn.BeginTx(ctx, pgx.TxOptions{
			IsoLevel: pgx.ReadCommitted, // for nonrepeatable read
		})
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			printf("Пробуем зафиксировать транзакцию")
			if cerr := tx.Commit(ctx); cerr != nil {
				printf("Не удалось зафиксировать, откатываемся: ", cerr)
				tx.Rollback(ctx)
				return
			}
			printf("Ok, транзакция зафиксирована")
		}()

		// эмуляция второго запроса
		// (например сервис статистики делает запрос)
		printf("В это время параллельно пробуем запросить баланс клиента (стартуем позже)")
		var balance int
		tx.QueryRow(ctx, "SELECT balance FROM cards WHERE owner_id = $1", ownerId).Scan(&balance)
		printf("Получили ответ от сервера о балансе: ", balance)
		printf("Эмулируем задержку перед повторным чтением")
		time.Sleep(longDelay)
		printf("Снова запрашиваем баланс")
		tx.QueryRow(ctx, "SELECT balance FROM cards WHERE owner_id = $1", ownerId).Scan(&balance)
		printf("Получили ответ от сервера о балансе: ", balance)
		printf("Для чистоты эксперимента спросим ещё и в третий раз")
		time.Sleep(longDelay)
		printf("Снова запрашиваем баланс")
		tx.QueryRow(ctx, "SELECT balance FROM cards WHERE owner_id = $1", ownerId).Scan(&balance)
		printf("Получили ответ от сервера о балансе: ", balance)
	}()

	wg.Wait()
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		log.Fatal(err)
	}
	// закрытие соединения (игнорируем ошибки для простоты)
	defer conn.Close(ctx)
	var balance int
	conn.QueryRow(ctx, "SELECT balance FROM cards WHERE owner_id = $1", ownerId).Scan(&balance)
	fmt.Printf("Итоговый баланс, после всех списаний: %d\n", balance)
}

var (
	First  = Color("\033[1;32m%s\033[0m\n")
	Second = Color("\033[1;34m%s\033[0m\n")
)

func Color(color string) func(args ...interface{}) {
	return func(args ...interface{}) {
		fmt.Printf(color, fmt.Sprint(args...))
	}
}
