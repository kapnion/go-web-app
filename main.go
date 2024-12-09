package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
)

type Order struct {
    ID          string  `json:"id"`
    Description string  `json:"description"`
    Amount      float64 `json:"amount"`
}

type OrderConfirmation struct {
    OrderID     string `json:"order_id"`
    Status      string `json:"status"`
    ConfirmedAt string `json:"confirmed_at"`
}

type OrderComparison struct {
    Order        Order
    Confirmation *OrderConfirmation
}

func main() {
    http.HandleFunc("/", indexHandler)
    http.HandleFunc("/api/orders", ordersHandler)
    http.HandleFunc("/api/order-confirmations", orderConfirmationsHandler)

    log.Println("Server started at http://localhost:8086")
    log.Fatal(http.ListenAndServe(":8086", nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
    tmpl, err := template.ParseFiles("templates/index.html")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    orders := getOrders()
    confirmations := getOrderConfirmations()

    var comparisons []OrderComparison
    for _, order := range orders {
        comparison := OrderComparison{Order: order}
        for _, conf := range confirmations {
            if conf.OrderID == order.ID {
                comparison.Confirmation = &conf
                break
            }
        }
        comparisons = append(comparisons, comparison)
    }

    err = tmpl.Execute(w, comparisons)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func ordersHandler(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(getOrders())
}

func orderConfirmationsHandler(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(getOrderConfirmations())
}

func getOrders() []Order {
    return []Order{
        {ID: "1", Description: "Order 1", Amount: 100.0},
        {ID: "2", Description: "Order 2", Amount: 200.0},
    }
}

func getOrderConfirmations() []OrderConfirmation {
    return []OrderConfirmation{
        {OrderID: "1", Status: "Confirmed", ConfirmedAt: "2023-01-01"},
        {OrderID: "2", Status: "Pending", ConfirmedAt: ""},
    }
}
