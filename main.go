package main

import (
	"fmt"
	"log"
	"net/http"
  "net/smtp"
	"os"
)

func UnsubscribeFromEmail(w http.ResponseWriter, r *http.Request){
  to := r.FormValue("email")
  if to == ""{
    http.Error(w, "Email field not provided", http.StatusBadRequest)
    return
  }

  if err := SendUnsubscribeRequest(to); err != nil{
    http.Error(w, "Failed to send unsubscribe email request", http.StatusInternalServerError)
    return
  }

  response := fmt.Sprintf("Unsubscribe email sent to %s", to)
  w.WriteHeader(http.StatusOK)
  w.Write([]byte(response))
}

func SendUnsubscribeRequest(to string) error{
  from := os.Getenv("EMAIL")
  fromPassword := os.Getenv("EMAIL_PASSWORD")
  smtpPort := "mail.google.com" 
  smtpHost := "587"

  msg := []byte(
    "To: "+ to + "\r\n"+
    "Subject: Unsubscription Confirmation" + "\r\n"+
    "r\n" +
    "You have successfully unsubscribed!\r\n",
  )

  auth := smtp.PlainAuth("", from, fromPassword, smtpHost)
  return smtp.SendMail(smtpHost+":"+smtpPort, auth,from, []string{to}, msg)
}

func main(){
  http.HandleFunc("/unsubscribe", UnsubscribeFromEmail)

  port := ":8080"
  fmt.Printf("Server listening on port %s", port)

  if err := http.ListenAndServe(port, nil); err != nil{
    log.Fatal(err)
  }
}
