package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"golang.org/x/crypto/bcrypt"
)





func HashPassword(password string)(string,error){
	bytes,err:=bcrypt.GenerateFromPassword([]byte (password),bcrypt.DefaultCost)
	if err!=nil{
		return "",err
	}

	return string(bytes),nil
}

func GenerateOtp()string{
	max:=big.NewInt(10000)
	n,err:=rand.Int(rand.Reader,max)

	if err!=nil{
		return "0000"
	}
	return fmt.Sprintf("%04d",n.Int64())


}

