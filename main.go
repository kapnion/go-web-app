package main

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"

	"github.com/antchfx/xmlquery"
	"github.com/signintech/gopdf"
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
	// Assuming xmlDoc and xslDoc are defined and loaded correctly
	// result, err := xmlquery.Find(xmlDoc, xslDoc)
	// if err != nil {
	//     log.Fatalf("Error querying XML: %v", err)
	// }
	// Process result

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/api/orders", ordersHandler)
	http.HandleFunc("/api/order-confirmations", orderConfirmationsHandler)
	http.HandleFunc("/api/convert-xml", convertXMLHandler)

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
func convertXMLHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the request method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the multipart form data
	err := r.ParseMultipartForm(10 << 20) // 10 MB max
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// Get the file from the request
	file, _, err := r.FormFile("xmlFile")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Read the file content
	xmlContent, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, "Error reading the file", http.StatusInternalServerError)
		return
	}

	// Perform first XSL transformation
	htmlContent, err := performXSLTransformation(xmlContent, "first_transform.xsl")
	if err != nil {
		http.Error(w, "Error in first XSL transformation", http.StatusInternalServerError)
		return
	}

	// Perform second XSL transformation
	finalHTML, err := performSecondXSLTransformation([]byte(htmlContent))
	if err != nil {
		http.Error(w, "Error in second XSL transformation", http.StatusInternalServerError)
		return
	}

	// Convert HTML to PDF
	pdfContent, err := convertHTMLToPDF(finalHTML)
	if err != nil {
		http.Error(w, "Error converting HTML to PDF: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Set headers for PDF download
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=converted_document.pdf")

	// Write PDF content to response
	w.Write(pdfContent)
}

func performXSLTransformation(input []byte, xslFilename string) (string, error) {
	// Read XSL file
	xslPath := filepath.Join("xsl", xslFilename)
	xslContent, err := ioutil.ReadFile(xslPath)
	if err != nil {
		return "", err
	}

	// Create XML document
	xmlDoc, err := xmlquery.Parse(bytes.NewReader(input))
	if err != nil {
		return "", err
	}

	// Create XSLT document
	_, err = xmlquery.Parse(bytes.NewReader(xslContent))
	if err != nil {
		return "", err
	}

	// Perform transformation
	result, err := xmlquery.QueryAll(xmlDoc, "//xsl:template")
	if err != nil {
		return "", err
	}

	// Convert result to string
	var output bytes.Buffer
	for _, node := range result {
		output.WriteString(node.OutputXML(true))
	}

	return output.String(), nil
}

func performSecondXSLTransformation(input []byte) (string, error) {
	// Simply return the input as the output
	return string(input), nil
}

func convertHTMLToPDF(htmlContent string) ([]byte, error) {
	// Initialize PDF
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	pdf.AddPage()

	// Set font
	err := pdf.AddTTFFont("roboto", "fonts/Roboto-Regular.ttf") // Use the relative path to the Roboto-Regular.ttf file
	if err != nil {
		return nil, err
	}
	err = pdf.SetFont("roboto", "", 14)
	if err != nil {
		return nil, err
	}

	// Write HTML content to PDF
	pdf.Cell(nil, htmlContent)

	// Get PDF content
	var buf bytes.Buffer
	err = pdf.Write(&buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
