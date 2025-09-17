package util

import (
	"database/sql"
	"fmt"
	"net/smtp"
	"os"
)

// SendVerificationEmail sends the verification email using Gmail SMTP instead of SendGrid
func SendVerificationEmail(userEmail, userName, verificationLink string) error {
	from := os.Getenv("GMAIL_ADDRESS")
	password := os.Getenv("GMAIL_APP_PASSWORD")

	if from == "" || password == "" {
		return fmt.Errorf("GMAIL_ADDRESS and GMAIL_APP_PASSWORD must be set in environment variables")
	}

	to := []string{userEmail}

	// Gmail SMTP server config
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Subject + MIME headers for HTML email
	subject := "Subject: Potvrdite svoj email\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	// ✅ Full HTML version (copied from your SendGrid version)
	htmlContent := fmt.Sprintf(`
<!DOCTYPE html>
<html lang="bs">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Potvrdite svoj email</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            line-height: 1.6;
            color: #333;
            background-color: #f4f4f4;
        }
        .container {
            max-width: 600px;
            margin: 0 auto;
            background-color: #ffffff;
            border-radius: 10px;
            box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
            overflow: hidden;
        }
        .header {
            background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%);
            color: white;
            padding: 40px 30px;
            text-align: center;
        }
        .header h1 {
            font-size: 28px;
            margin-bottom: 10px;
            font-weight: 300;
        }
        .header p {
            font-size: 16px;
            opacity: 0.9;
        }
        .content {
            padding: 40px 30px;
        }
        .greeting {
            font-size: 18px;
            color: #2c3e50;
            margin-bottom: 20px;
        }
        .message {
            font-size: 16px;
            color: #555;
            margin-bottom: 30px;
            line-height: 1.8;
        }
        .cta-button {
            display: inline-block;
            background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%);
            color: white;
            padding: 15px 30px;
            text-decoration: none;
            border-radius: 50px;
            font-size: 16px;
            font-weight: 600;
            text-align: center;
            transition: transform 0.3s ease;
            box-shadow: 0 4px 15px rgba(102, 126, 234, 0.3);
        }
        .cta-button:hover {
            transform: translateY(-2px);
            box-shadow: 0 6px 20px rgba(102, 126, 234, 0.4);
        }
        .cta-container {
            text-align: center;
            margin: 30px 0;
        }
        .security-note {
            background-color: #f8f9fa;
            border-left: 4px solid #667eea;
            padding: 15px;
            margin: 30px 0;
            border-radius: 5px;
        }
        .security-note h3 {
            color: #2c3e50;
            font-size: 16px;
            margin-bottom: 8px;
        }
        .security-note p {
            color: #666;
            font-size: 14px;
        }
        .footer {
            background-color: #2c3e50;
            color: white;
            padding: 30px;
            text-align: center;
        }
        .footer p {
            margin-bottom: 10px;
            font-size: 14px;
        }
        .footer .company-name {
            font-weight: 600;
            color: #667eea;
        }
        .divider {
            height: 2px;
            background: linear-gradient(to right, #667eea, #764ba2);
            margin: 20px 0;
        }
        @media (max-width: 600px) {
            .container {
                margin: 0 10px;
            }
            .header,
            .content,
            .footer {
                padding: 20px;
            }
            .header h1 {
                font-size: 24px;
            }
            .cta-button {
                padding: 12px 25px;
                font-size: 14px;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Dobrodošli!</h1>
            <p>Samo jedan korak do završetka registracije</p>
        </div>
        <div class="content">
            <div class="greeting">
                Zdravo %s! 👋
            </div>
            <div class="message">
                Hvala vam što ste se registrovali na platformu eDnevnik. Da biste završili registraciju, potrebno je da potvrdite svoju email adresu.
            </div>
            <div class="cta-container">
                <a href="%s" target="_blank" class="cta-button" style="color:white;">Potvrdite Email</a>
            </div>
            <div class="divider"></div>
            <div class="security-note">
                <h3>🔒 Sigurnosna napomena</h3>
                <p>Ovaj link za potvrdu će biti aktivan 24 sata.
            </div>
        </div>
        <div class="footer">
            <p style="margin-top: 20px; font-size: 12px; opacity: 0.8;">
                Ovaj email je automatski generisan. Molimo vas da ne odgovarate direktno na ovu adresu.
            </p>
        </div>
    </div>
</body>
</html>
`, userName, verificationLink)

	// Build final message with Subject + MIME + HTML
	message := []byte(subject + mime + htmlContent)

	// SMTP authentication
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Send it!
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		return fmt.Errorf("failed to send verification email: %v", err)
	}

	return nil
}

// GetPendingAccountVerificationToken TODO: Add description
func GetPendingAccountVerificationToken(accountID int, workspaceDB *sql.DB) (string, error) {
	var token string
	query := "SELECT verification_token FROM pending_accounts WHERE id = ?"
	err := workspaceDB.QueryRow(query, accountID).Scan(&token)
	if err != nil {
		return "", fmt.Errorf("failed to get verification token: %v", err)
	}
	return token, nil
}

// VerifyAccount TODO: Add description
func VerifyAccount(verificationToken string, workspaceDB *sql.DB) error {

	type PendingAccountData struct {
		ID          int
		AccountType string
		Email       string
	}
	var pendingAccount PendingAccountData

	// Get the pending account data using the verification token
	query := `SELECT id, account_type, email FROM pending_accounts WHERE verification_token = ?;`
	err := workspaceDB.QueryRow(query, verificationToken).Scan(
		&pendingAccount.ID,
		&pendingAccount.AccountType,
		&pendingAccount.Email,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("no pending account found with the provided verification token")
		}
		return fmt.Errorf("error retrieving pending account: %v", err)
	}

	accountInsertQuery := `INSERT INTO accounts (email, password, account_type)
    SELECT email, password, account_type
    FROM pending_accounts
    WHERE verification_token = ?;`

	res, err := workspaceDB.Exec(accountInsertQuery, verificationToken)
	if err != nil {
		return fmt.Errorf("error inserting account: %v", err)
	}

	accountID, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("error getting last insert ID: %v", err)
	}

	if pendingAccount.AccountType == "teacher" {
		insertTeacherQuery := `INSERT INTO teachers (name, last_name, phone, account_id,
        contractions, title)
        SELECT name, last_name, phone, ?, contractions, title
        FROM pending_teachers
        WHERE account_id = ?;`

		_, err = workspaceDB.Exec(insertTeacherQuery, accountID, pendingAccount.ID)
		if err != nil {
			return fmt.Errorf("error inserting teacher: %v", err)
		}
	} else if pendingAccount.AccountType == "pupil" {
		insertPupilQuery := `INSERT INTO pupil_global (name, last_name, jmbg,
        gender, address, guardian_name, phone_number, guardian_number, date_of_birth,
        religion, account_id, place_of_birth)
        SELECT name, last_name, jmbg, gender, address, guardian_name, phone_number,
        guardian_number, date_of_birth, religion, ?, place_of_birth
        FROM pending_pupil_global
        WHERE account_id = ?;`

		_, err = workspaceDB.Exec(insertPupilQuery, accountID, pendingAccount.ID)
		if err != nil {
			return fmt.Errorf("error inserting pupil: %v", err)
		}
	} else {
		return fmt.Errorf("unsupported account type: %s", pendingAccount.AccountType)
	}

	pendingCleanupQuery := `DELETE FROM pending_accounts WHERE verification_token = ?;`
	_, err = workspaceDB.Exec(pendingCleanupQuery, verificationToken)
	if err != nil {
		return fmt.Errorf("error cleaning up pending accounts: %v", err)
	}

	return nil
}
