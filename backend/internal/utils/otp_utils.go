package utils

import "fmt"


// later on we will integerate api for both endpoints


func SendEmailOtp(email,otp string){
	fmt.Println(otp)
	fmt.Println(email)
}

func SendSMSOTP(phone,otp string){
	fmt.Println(phone)
	fmt.Println(otp)

}