package main

import (
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"
)

func UnsubscribeFromEmail(w http.ResponseWriter, r *http.Request) {
	// Only allow POST method
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	to := r.FormValue("email")
	if to == "" {
		http.Error(w, "Email field not provided", http.StatusBadRequest)
		return
	}

	if err := SendUnsubscribeRequest(to); err != nil {
		log.Printf("Failed to send unsubscribe email: %v", err)
		http.Error(w, "Failed to send unsubscribe email request", http.StatusInternalServerError)
		return
	}

	response := fmt.Sprintf("Unsubscribe email sent to %s", to)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}

func SendUnsubscribeRequest(to string) error {
	from := os.Getenv("EMAIL")
	fromPassword := os.Getenv("EMAIL_PASSWORD")
	
	// Fixed: Swapped host and port values
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Fixed: Proper email formatting with headers
	msg := []byte(fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: Unsubscription Confirmation\r\n"+
		"\r\n"+
		"You have successfully unsubscribed!\r\n",
		from, to))

	auth := smtp.PlainAuth("", from, fromPassword, smtpHost)
	
	// Fixed: Properly concatenated host and port
	return smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, msg)
}

func main() {
	// Check for required environment variables
	required := []string{"EMAIL", "EMAIL_PASSWORD"}
	for _, env := range required {
		if os.Getenv(env) == "" {
			log.Fatalf("Required environment variable %s is not set", env)
		}
	}

	http.HandleFunc("/unsubscribe", UnsubscribeFromEmail)
	port := ":8080"
	fmt.Printf("Server listening on port %s\n", port) // Added newline
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
