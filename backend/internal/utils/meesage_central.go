package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

var mcClient = &http.Client{
	Timeout: 10 * time.Second,
}

// message central otp service for phone verification


func MessageCentralSendOTP(mobile string) (string, error) {

	authToken := os.Getenv("MESSAGE_CENTRAL_AUTH_TOKEN")
	customerId := os.Getenv("MESSAGE_CENTRAL_CUSTOMER_ID")
	countryCode := os.Getenv("MESSAGE_CENTRAL_COUNTRY_CODE")

	if authToken == "" || customerId == "" || countryCode == "" {
		return "", errors.New("message central env vars not set")
	}

	url := fmt.Sprintf(
		"https://cpaas.messagecentral.com/verification/v3/send?customerId=%s&countryCode=%s&flowType=SMS&mobileNumber=%s",
		customerId,
		countryCode,
		mobile,
	)

	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("authToken", authToken)
	req.Header.Set("accept", "application/json")

	resp, err := mcClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("send otp failed: %s", string(body))
	}

	var response struct {
		Data struct {
			VerificationID string `json:"verificationId"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return "", err
	}

	if response.Data.VerificationID == "" {
		return "", errors.New("verificationId missing in response")
	}

	return response.Data.VerificationID, nil
}


func MessageCentralVerifyOTP(verificationId, otp string) error {

	authToken := os.Getenv("MESSAGE_CENTRAL_AUTH_TOKEN")
	if authToken == "" {
		return errors.New("message central auth token not set")
	}

	url := fmt.Sprintf(
		"https://cpaas.messagecentral.com/verification/v3/validateOtp?verificationId=%s&code=%s&flowType=SMS",
		verificationId,
		otp,
	)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("authToken", authToken)
	req.Header.Set("accept", "application/json")

	resp, err := mcClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("otp validation failed: %s", string(body))
	}

	return nil
}

