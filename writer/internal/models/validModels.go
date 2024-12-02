package models

func CheckValidOrder(order Order) bool {

	if order.OrderUID == "" ||
		order.TrackNumber == "" ||
		order.Entry == "" ||
		order.Delivery.Name == "" ||
		order.Delivery.Phone == "" ||
		order.Delivery.Zip == "" ||
		order.Delivery.City == "" ||
		order.Delivery.Address == "" ||
		order.Delivery.Region == "" ||
		order.Delivery.Email == "" ||
		order.Payment.Transaction == "" ||
		order.Payment.Currency == "" ||
		order.Payment.Provider == "" ||
		order.Payment.Amount < 0 || // допускаем 0
		order.Payment.PaymentDT <= 0 || // допускаем 0
		order.Payment.Bank == "" ||
		order.Payment.DeliveryCost < 0 || // допускаем 0
		order.Payment.GoodsTotal < 0 || // допускаем 0
		order.Payment.CustomFee < 0 { // допускаем 0
		return false
	}

	k := 0
	for i := 0; i < len(order.Items); i++ {
		if order.Items[i].ChrtID <= 0 || // допускаем 0
			order.Items[i].TrackNumber == "" ||
			order.Items[i].Price < 0 || // допускаем 0
			order.Items[i].Rid == "" ||
			order.Items[i].Name == "" ||
			order.Items[i].Sale < 0 || // допускаем 0
			order.Items[i].Size == "" ||
			order.Items[i].TotalPrice < 0 || // допускаем 0
			order.Items[i].NmID <= 0 || // допускаем 0
			order.Items[i].Brand == "" ||
			order.Items[i].Status <= 0 {
			k++
		}
	}
	if k != 0 {
		return false
	}
	return true
}
