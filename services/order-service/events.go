package main

import (
	"encoding/json"
	"log"

	"common"
)

func handlePaymentProcessed(data []byte) error {
	var event common.Event
	if err := json.Unmarshal(data, &event); err != nil {
		log.Printf("[ORDER] ERROR: Failed to unmarshal PaymentProcessed event: %v", err)
		return err
	}

	log.Printf("[ORDER] Received PaymentProcessed for order: %s", event.OrderID)
	err := UpdateOrderStatus(event.OrderID, OrderStatusConfirmed)
	if err != nil {
		log.Printf("[ORDER] ERROR: Failed to update order status to CONFIRMED: %v", err)
		return err
	}
	return nil
}

func handleDeliveryAssigned(data []byte) error {
	var event common.Event
	if err := json.Unmarshal(data, &event); err != nil {
		log.Printf("[ORDER] ERROR: Failed to unmarshal DeliveryAssigned event: %v", err)
		return err
	}

	log.Printf("[ORDER] Received DeliveryAssigned for order: %s", event.OrderID)
	err := UpdateOrderStatus(event.OrderID, OrderStatusDelivered)
	if err != nil {
		log.Printf("[ORDER] ERROR: Failed to update order status to DELIVERED: %v", err)
		return err
	}
	return nil
}

func startEventConsumers(mq *common.RabbitMQ) {
	if mq == nil {
		return
	}

	go func() {
		err := mq.ConsumeWithRetry("PaymentProcessed", handlePaymentProcessed)
		if err != nil {
			log.Printf("[ORDER] Failed to consume PaymentProcessed: %v", err)
		}
	}()

	go func() {
		err := mq.ConsumeWithRetry("DeliveryAssigned", handleDeliveryAssigned)
		if err != nil {
			log.Printf("[ORDER] Failed to consume DeliveryAssigned: %v", err)
		}
	}()
}
