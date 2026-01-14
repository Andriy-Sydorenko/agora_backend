package email

import "fmt"

const (
	PasswordResetEmailTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Reset Your Password</title>
</head>
<body style="margin: 0; padding: 0; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif; background-color: #f4f4f5;">
    <table role="presentation" style="width: 100%%; border-collapse: collapse; background-color: #f4f4f5;">
        <tr>
            <td style="padding: 40px 20px;">
                <table role="presentation" style="max-width: 600px; margin: 0 auto; background-color: #ffffff; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1);">
                    <!-- Header -->
                    <tr>
                        <td style="padding: 40px 40px 20px; text-align: center; border-bottom: 1px solid #e5e7eb;">
                            <h1 style="margin: 0; font-size: 24px; font-weight: 600; color: #111827;">Reset Your Password</h1>
                        </td>
                    </tr>
                    
                    <!-- Content -->
                    <tr>
                        <td style="padding: 32px 40px;">
                            <p style="margin: 0 0 16px; font-size: 16px; line-height: 24px; color: #374151;">
                                We received a request to reset your password. Click the button below to create a new password:
                            </p>
                            
                            <!-- CTA Button -->
                            <table role="presentation" style="width: 100%%; margin: 24px 0;">
                                <tr>
                                    <td style="text-align: center;">
                                        <a href="%s" style="display: inline-block; padding: 12px 32px; background-color: #2563eb; color: #ffffff; text-decoration: none; border-radius: 6px; font-weight: 500; font-size: 16px;">Reset Password</a>
                                    </td>
                                </tr>
                            </table>
                            
                            <p style="margin: 24px 0 0; font-size: 14px; line-height: 20px; color: #6b7280;">
                                Or copy and paste this link into your browser:
                            </p>
                            <p style="margin: 8px 0 0; font-size: 12px; line-height: 18px; color: #9ca3af; word-break: break-all;">
                                %s
                            </p>
                        </td>
                    </tr>
                    
                    <!-- Security Notice -->
                    <tr>
                        <td style="padding: 0 40px 32px;">
                            <div style="padding: 16px; background-color: #fef3c7; border-left: 4px solid #f59e0b; border-radius: 4px;">
                                <p style="margin: 0 0 8px; font-size: 14px; font-weight: 600; color: #92400e;">
                                    ⚠️ Security Notice
                                </p>
                                <ul style="margin: 0; padding-left: 20px; font-size: 13px; line-height: 20px; color: #78350f;">
                                    <li>This link will expire in <strong>30 minutes</strong></li>
                                    <li>If you didn't request this, please ignore this email</li>
                                    <li>Your password will not change until you click the link above</li>
                                </ul>
                            </div>
                        </td>
                    </tr>
                    
                    <!-- Footer -->
                    <tr>
                        <td style="padding: 24px 40px; background-color: #f9fafb; border-top: 1px solid #e5e7eb; border-radius: 0 0 8px 8px;">
                            <p style="margin: 0 0 8px; font-size: 12px; line-height: 18px; color: #6b7280; text-align: center;">
                                This is an automated message. Please do not reply to this email.
                            </p>
                            <p style="margin: 0; font-size: 12px; line-height: 18px; color: #9ca3af; text-align: center;">
                                © %d Agora. All rights reserved.
                            </p>
                        </td>
                    </tr>
                </table>
            </td>
        </tr>
    </table>
</body>
</html>
`
	PasswordResetEmailSubject = "Reset Your Password"
)

// generatePasswordResetEmailHTML generates the HTML content for password reset email
func generatePasswordResetEmailHTML(resetURL string, year int) string {
	return fmt.Sprintf(PasswordResetEmailTemplate, resetURL, resetURL, year)
}
