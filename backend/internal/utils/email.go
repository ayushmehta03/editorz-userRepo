package utils

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "os"
)

func SendOTPEmail(toEmail string, otp string) error {

    apiKey := os.Getenv("RESEND_API_KEY")
    from := os.Getenv("EMAIL_FROM")

    if apiKey == "" || from == "" {
        return fmt.Errorf("resend env variables not set")
    }

    payload := map[string]interface{}{
        "from": from,
        "to":   []string{toEmail},
        "subject": "Editorzzz • Verify your email",
        "html": fmt.Sprintf(`
<!DOCTYPE html>
<html>
<body style="background:#09090b; font-family:'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; padding:40px;">
  <div style="max-width:520px; margin:auto; background:#18181b; padding:40px; border-radius:24px; color:#ffffff; border: 1px solid #27272a;">
    <h2 style="text-align:center; margin:0; font-size:28px; letter-spacing:-0.5px;">
      Editor<span style="color:#8b5cf6;">zzz</span>
    </h2>

    <p style="text-align:center; color:#a1a1aa; margin-top:10px; font-size:16px;">
      Verify your account
    </p>

    <div style="margin:40px 0; text-align:center;">
      <p style="color:#e4e4e7; margin-bottom:20px; font-size:14px;">Your secure verification code:</p>

      <div style="
        font-size:32px;
        font-weight:800;
        letter-spacing:8px;
        color:#ffffff;
        background:linear-gradient(135deg, #8b5cf6 0%%, #6d28d9 100%%);
        padding:20px 40px;
        border-radius:16px;
        display:inline-block;
        box-shadow: 0 10px 15px -3px rgba(139, 92, 246, 0.3);
      ">
        %s
      </div>

      <p style="font-size:13px; color:#71717a; margin-top:20px;">
        This code expires in <span style="color:#a1a1aa; font-weight:600;">10 minutes</span>
      </p>
    </div>

    <div style="border-top:1px solid #27272a; padding-top:24px;">
      <p style="font-size:12px; color:#52525b; text-align:center; line-height:1.6;">
        If you didn’t request this code, you can safely ignore this email. <br/> 
        &copy; 2026 Editorzzz Inc.
      </p>
    </div>
  </div>
</body>
</html>
`, otp),
    }

    body, _ := json.Marshal(payload)

    req, err := http.NewRequest(
        "POST",
        "https://api.resend.com/emails",
        bytes.NewBuffer(body),
    )
    if err != nil {
        return err
    }

    req.Header.Set("Authorization", "Bearer "+apiKey)
    req.Header.Set("Content-Type", "application/json")

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode >= 300 {
        return fmt.Errorf("resend failed with status %d", resp.StatusCode)
    }

    return nil
}