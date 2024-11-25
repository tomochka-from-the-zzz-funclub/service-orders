package database

import (
	"consumer/internal/config"
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

	return &Postgres{
		Connection: db,
	} //, nil
}

func (db *Postgres) AddOrder(order models.Order) (string, error) {
	myLog.Log.Debugf("Go to bd in Set")
	query_add_delivery :=
		`WITH insert_return AS (
            INSERT INTO  delivery (name, phone, zip, city, address, region,	email)
            VALUES ($1, $2, $3, $4, $5, $6, $7)
            RETURNING id
        )
        SELECT id FROM insert_return`
	var id_delivery string
	err := db.Connection.QueryRow(query_add_delivery, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip, order.Delivery.City, order.Delivery.Address, order.Delivery.Region, order.Delivery.Email).Scan(&id_delivery)
	if err != nil {
		myLog.Log.Errorf("Error CreateDelivery: %v", err.Error())
		return "", err
	}
	query_add_payment := `
        WITH insert_return AS (
            INSERT INTO payment (transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee)
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
            RETURNING id
        )
        SELECT id FROM insert_return
    `
	var id_payment string
	err = db.Connection.QueryRow(query_add_payment, order.Payment.Transaction, order.Payment.RequestID, order.Payment.Currency, order.Payment.Provider,
		strconv.Itoa(order.Payment.Amount), strconv.Itoa(order.Payment.PaymentDT), order.Payment.Bank, strconv.Itoa(order.Payment.DeliveryCost), strconv.Itoa(order.Payment.GoodsTotal), strconv.Itoa(order.Payment.CustomFee)).Scan(&id_payment)
	if err != nil {
		myLog.Log.Errorf("Error CreatePayment: %v", err.Error())
		return "", err
	}
	query_add_items := `
        WITH insert_return AS (
            INSERT INTO items (chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status)
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
            RETURNING id
        )
        SELECT id FROM insert_return
    `
	var id_item []string
	for i := 0; i < len(order.Items); i++ {
		var id string
		err = db.Connection.QueryRow(query_add_items, strconv.Itoa(order.Items[i].ChrtID), order.Items[i].TrackNumber, strconv.Itoa(order.Items[i].Price), order.Items[i].Rid, order.Items[i].Name,
			strconv.Itoa(order.Items[i].Sale), order.Items[i].Size, strconv.Itoa(order.Items[i].TotalPrice), strconv.Itoa(order.Items[i].NmID), order.Items[i].Brand, strconv.Itoa(order.Items[i].Status)).Scan(&id)
		if err != nil {
			myLog.Log.Errorf("Error CreateItems: %v", err.Error())
			return "", err
		}
		myLog.Log.Debugf("Insert items %v", strconv.Itoa(order.Items[i].ChrtID))
		id_item = append(id_item, id)
	}
	query_add_order := `
		INSERT INTO orders (order_uid, track_number, entry, delivery, payment, items, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard)
		VALUES ($1, $2, $3, $4::uuid, $5::uuid, $6::uuid[], $7, $8, $9, $10, $11, $12, $13, $14)
`
	err = db.Connection.QueryRow(query_add_order, order.OrderUID, order.TrackNumber, order.Entry, id_delivery, id_payment, pq.Array(id_item), order.Locale,
		order.InternalSignature, order.CustomerID, order.DeliveryService, order.ShardKey, order.SmID, order.DateCreated, order.OofShard).Err()
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

func (db *Postgres) GetOrder(order_uuid string) (models.Order, error) {
	myLog.Log.Debugf("Go to db in GetOrder")

	var order models.Order
	var delivery models.Delivery
	var payment models.Payment

	query_get_order := `
    SELECT 
        o.track_number, 
        o.entry,
        d.name, d.phone, d.zip, d.city, d.address, d.region, d.email,
        p.transaction, p.request_id, p.currency, p.provider, p.amount, p.payment_dt, p.bank, p.delivery_cost, p.goods_total, p.custom_fee,
        i.chrt_id, i.track_number, i.price, i.rid, i.name, i.sale, i.size, i.total_price, i.nm_id, i.brand, i.status
    FROM 
        orders o
    LEFT JOIN 
        delivery d ON o.delivery = d.id
    LEFT JOIN 
        payment p ON o.payment = p.id
    LEFT JOIN 
        items i ON i.id = ANY(o.items)
    WHERE 
        o.order_uid = $1`

	rows, err := db.Connection.Query(query_get_order, order_uuid)
	if err != nil {
		myLog.Log.Errorf("Error GetOrder: %v", err.Error())
		return models.Order{}, err
	}
	defer rows.Close()

	itemMap := make(map[string]models.Item)

	for rows.Next() {
		var item models.Item

		err = rows.Scan(
			&order.TrackNumber,
			&order.Entry,
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
			&item.ChrtID,
			&item.TrackNumber,
			&item.Price,
			&item.Rid,
			&item.Name,
			&item.Sale,
			&item.Size,
			&item.TotalPrice,
			&item.NmID,
			&item.Brand,
			&item.Status,
		)

		if err != nil {
			myLog.Log.Errorf("Error scanning row: %v", err.Error())
			return models.Order{}, err
		}
	}
	for _, item := range itemMap {
		order.Items = append(order.Items, item)
	}

	return order, nil
}
