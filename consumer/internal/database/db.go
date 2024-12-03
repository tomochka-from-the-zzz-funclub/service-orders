package database

import (
	"consumer/internal/config"
	myErrors "consumer/internal/errors"
	myLog "consumer/internal/logger"
	"consumer/internal/models"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/lib/pq"
)

type Postgres struct {
	Connection *sql.DB
}

func NewPostgres(cfg config.Config) *Postgres {
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=%s", cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBHost, cfg.DBPort, cfg.SslMode)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		myLog.Log.Fatalf("Failed to connect to PostgreSQL: %v", err)
		return nil //, myErrors.ErrCreatePostgresConnection
	}
	time.Sleep(time.Minute)
	err = db.Ping()
	if err != nil {
		myLog.Log.Fatalf("Failed to ping PostgreSQL: %v", err)
		return nil //, myErrors.ErrPing
	} else {
		myLog.Log.Debugf("ping success")
	}
	time.Sleep(time.Minute)
	query :=
		` CREATE TABLE IF NOT EXISTS orders (
    order_uid VARCHAR(255) PRIMARY KEY,
    track_number VARCHAR(255),
    entry VARCHAR(255),
    delivery UUID,
    payment UUID,
    items UUID[],
    locale VARCHAR(3),
    internal_signature VARCHAR(255),
    customer_id VARCHAR(255),
    delivery_service VARCHAR(255),
    shardkey INT,
    sm_id INT,
    date_created VARCHAR(255),
    oof_shard INT
);

CREATE TABLE IF NOT EXISTS payment (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    transaction VARCHAR(255),
    request_id VARCHAR(255),
    currency VARCHAR(3),
    provider VARCHAR(255),
    amount INT,
    payment_dt INT,
    bank VARCHAR(255),
    delivery_cost INT,
    goods_total INT,
    custom_fee INT
);

CREATE TABLE IF NOT EXISTS items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    chrt_id INT,
    track_number VARCHAR(255),
    price INT,
    rid VARCHAR(255),
    name VARCHAR(255),
    sale INT,
    size INT,
    total_price INT,
    nm_id INT,
    brand VARCHAR(255),
    status VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS delivery (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255),
    phone VARCHAR(255),
    zip INT,
    city VARCHAR(255),
    address VARCHAR(255),
    region VARCHAR(255),
    email VARCHAR(255)
);
`
	_, err = db.Exec(query)
	if err != nil {
		myLog.Log.Errorf(err.Error())
	}
	return &Postgres{
		Connection: db,
	}
}

func (db *Postgres) AddOrderStruct(order models.Order) (string, error) {
	myLog.Log.Debugf("Go to bd in Set")

	id_delivery, err := db.AddDelivery(order.Delivery)
	if err != nil {
		myLog.Log.Errorf("Error CreateDelivery: %v", err.Error())
		return "", err
	}

	id_payment, err := db.AddPayment(order.Payment)
	if err != nil {
		myLog.Log.Errorf("Error CreatePayment: %v", err.Error())
		return "", err
	}

	id_item, err := db.AddItems(order.Items)
	if err != nil {
		myLog.Log.Errorf("Error CreateItems: %v", err.Error())
		return "", err
	}

	err = db.AddOrder(order, id_item, id_delivery, id_payment)
	if err != nil {
		if err == sql.ErrNoRows {
			myLog.Log.Debugf("The entry was not added. No data returned.")
			return "", err
		}
		myLog.Log.Errorf("Error CreateOrder: %v", err.Error())
		return "", err
	}

	myLog.Log.Infof("Entry successfully added with ID: ", order.OrderUID)

	return order.OrderUID, nil
}

func (db *Postgres) AddDelivery(delivery models.Delivery) (string, error) {
	query_add_delivery :=
		`WITH insert_return AS (
		INSERT INTO  delivery (name, phone, zip, city, address, region,	email)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	)
	SELECT id FROM insert_return`
	var id_delivery string
	err := db.Connection.QueryRow(query_add_delivery, delivery.Name, delivery.Phone, delivery.Zip, delivery.City, delivery.Address, delivery.Region, delivery.Email).Scan(&id_delivery)
	if err != nil {
		return "", err
	}
	return id_delivery, nil
}

func (db *Postgres) AddPayment(payment models.Payment) (string, error) {
	query_add_payment := `
        WITH insert_return AS (
            INSERT INTO payment (transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee)
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
            RETURNING id
        )
        SELECT id FROM insert_return
    `
	var id_payment string
	err := db.Connection.QueryRow(query_add_payment, payment.Transaction, payment.RequestID, payment.Currency, payment.Provider,
		strconv.Itoa(payment.Amount), strconv.Itoa(payment.PaymentDT), payment.Bank, strconv.Itoa(payment.DeliveryCost),
		strconv.Itoa(payment.GoodsTotal), strconv.Itoa(payment.CustomFee)).Scan(&id_payment)
	if err != nil {
		return "", err
	}
	return id_payment, nil
}

func (db *Postgres) AddItems(items []models.Item) ([]string, error) {

	query_add_items := `
        WITH insert_return AS (
            INSERT INTO items (chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status)
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
            RETURNING id
        )
        SELECT id FROM insert_return
    `
	var id_items []string
	for i := 0; i < len(items); i++ {
		var id string
		err := db.Connection.QueryRow(query_add_items, strconv.Itoa(items[i].ChrtID), items[i].TrackNumber, strconv.Itoa(items[i].Price),
			items[i].Rid, items[i].Name, strconv.Itoa(items[i].Sale), items[i].Size, strconv.Itoa(items[i].TotalPrice), strconv.Itoa(items[i].NmID),
			items[i].Brand, strconv.Itoa(items[i].Status)).Scan(&id)
		if err != nil {
			myLog.Log.Errorf("Error CreateItems: %v", err.Error())
			return id_items, err
		}
		myLog.Log.Debugf("Insert items %v", strconv.Itoa(items[i].ChrtID))
		id_items = append(id_items, id)
	}
	return id_items, nil
}

func (db *Postgres) AddOrder(order models.Order, items []string, delivery string, payment string) error {
	query_add_order := `
	INSERT INTO orders (order_uid, track_number, entry, delivery, payment, items, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard)
	VALUES ($1, $2, $3, $4::uuid, $5::uuid, $6::uuid[], $7, $8, $9, $10, $11, $12, $13, $14)
`
	err := db.Connection.QueryRow(query_add_order, order.OrderUID, order.TrackNumber, order.Entry, delivery, payment, pq.Array(items), order.Locale,
		order.InternalSignature, order.CustomerID, order.DeliveryService, order.ShardKey, order.SmID, order.DateCreated, order.OofShard).Err()
	return err
}

func (db *Postgres) GetOrder(order_uuid string) (models.Order, error) {
	myLog.Log.Debugf("Go to db in GetOrder with id: %+v", order_uuid)
	var order models.Order
	var delivery models.Delivery
	var payment models.Payment
	query_get_order := `
    SELECT 
        o.track_number, 
        o.entry, o.items, o.locale, o.internal_signature, o.customer_id, o.delivery_service, o.shardkey, o.sm_id, o.date_created, o.oof_shard,
        d.name, d.phone, d.zip, d.city, d.address, d.region, d.email,
        p.transaction, p.request_id, p.currency, p.provider, p.amount, p.payment_dt, p.bank, p.delivery_cost, p.goods_total, p.custom_fee
    FROM 
        orders o
    LEFT JOIN 
        delivery d ON o.delivery = d.id
    LEFT JOIN 
        payment p ON o.payment = p.id
    WHERE 
        o.order_uid = $1`
	rows, err := db.Connection.Query(query_get_order, order_uuid)
	if err != nil {
		myLog.Log.Errorf("Error GetOrder: %v", err.Error())
		return models.Order{}, err
	}
	defer rows.Close()
	//itemMap := make(map[string]models.Item)
	var id_items []string
	if !rows.Next() { // Проверка на наличие записи
		myLog.Log.Errorf("No order found with uuid: %v", order_uuid)
		return models.Order{}, myErrors.ErrNotFoundOrder
	}
	for rows.Next() {
		//var item models.Item
		err = rows.Scan(
			&order.TrackNumber,
			&order.Entry,
			pq.Array(&id_items),
			&order.Locale,
			&order.InternalSignature,
			&order.CustomerID,
			&order.DeliveryService,
			&order.ShardKey,
			&order.SmID,
			&order.DateCreated,
			&order.OofShard,
			&delivery.Name,
			&delivery.Phone,
			&delivery.Zip,
			&delivery.City,
			&delivery.Address,
			&delivery.Region,
			&delivery.Email,
			&payment.Transaction,
			&payment.RequestID,
			&payment.Currency,
			&payment.Provider,
			&payment.Amount,
			&payment.PaymentDT,
			&payment.Bank,
			&payment.DeliveryCost,
			&payment.GoodsTotal,
			&payment.CustomFee,
		)
		if err != nil {
			myLog.Log.Errorf("Error scanning row: %v", err.Error())
			return models.Order{}, err
		}
	}

	query_get_item :=
		`SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status
	WHERE id = $1`
	var item models.Item
	var items []models.Item
	for i := 0; i < len(order.Items); i++ {
		err = db.Connection.QueryRow(query_get_item, order.Items[i]).Scan(&item.ChrtID, &item.TrackNumber, &item.Price, &item.Rid, &item.Name, &item.Sale, &item.Size,
			&item.TotalPrice, &item.NmID, &item.Brand, &item.Status)
		if err != nil {
			myLog.Log.Errorf("Error GetItems: %v", err.Error())
			return models.Order{}, err
		}
		items = append(items, item)
	}
	order.Items = items
	order.Delivery = delivery
	order.Payment = payment
	return order, nil
}

func (db *Postgres) GetAllOrders() (map[string]models.Order, error) {
	result := make(map[string]models.Order)

	query := `SELECT o.order_uid, o.track_number, 
        o.entry, o.items, o.locale, o.internal_signature, o.customer_id, o.delivery_service, o.shardkey, o.sm_id, o.date_created, o.oof_shard,
        d.name, d.phone, d.zip, d.city, d.address, d.region, d.email,
        p.transaction, p.request_id, p.currency, p.provider, p.amount, p.payment_dt, p.bank, p.delivery_cost, p.goods_total, p.custom_fee

	FROM  orders o
    LEFT JOIN 
        delivery d ON o.delivery = d.id
    LEFT JOIN 
        payment p ON o.payment = p.id`

	rows, err := db.Connection.Query(query)
	if err != nil {
		return result, err
	}
	defer rows.Close()
	var id_items []string
	var order models.Order
	var delivery models.Delivery
	var payment models.Payment
	for rows.Next() {
		err = rows.Scan(
			&order.OrderUID,
			&order.TrackNumber,
			&order.Entry,
			pq.Array(&id_items),
			&order.Locale,
			&order.InternalSignature,
			&order.CustomerID,
			&order.DeliveryService,
			&order.ShardKey,
			&order.SmID,
			&order.DateCreated,
			&order.OofShard,
			&delivery.Name,
			&delivery.Phone,
			&delivery.Zip,
			&delivery.City,
			&delivery.Address,
			&delivery.Region,
			&delivery.Email,
			&payment.Transaction,
			&payment.RequestID,
			&payment.Currency,
			&payment.Provider,
			&payment.Amount,
			&payment.PaymentDT,
			&payment.Bank,
			&payment.DeliveryCost,
			&payment.GoodsTotal,
			&payment.CustomFee,
		)
		if err != nil {
			myLog.Log.Errorf("Error scanning row: %v", err.Error())
			return result, err
		}
		query_get_item :=
			`SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status
	WHERE id = $1`
		var item models.Item
		var items []models.Item
		for i := 0; i < len(order.Items); i++ {
			err = db.Connection.QueryRow(query_get_item, order.Items[i]).Scan(&item.ChrtID, &item.TrackNumber, &item.Price, &item.Rid, &item.Name, &item.Sale, &item.Size,
				&item.TotalPrice, &item.NmID, &item.Brand, &item.Status)
			if err != nil {
				myLog.Log.Errorf("Error GetItems: %v", err.Error())
				return result, err
			}
			items = append(items, item)
		}
		order.Items = items
		order.Delivery = delivery
		order.Payment = payment
		result[order.OrderUID] = order
		fmt.Println(order.OrderUID)
	}
	return result, nil
}
