package database

type DBUser struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
	Role         string `json:"role"`
	TOTPSecret   string `json:"totp_secret"`
	TOTPEnabled  int    `json:"totp_enabled"`
}

func (d *DB) GetUserByUsername(username string) (*DBUser, error) {
	var u DBUser
	err := d.QueryRow("SELECT id, username, password_hash, role, totp_secret, totp_enabled FROM users WHERE username = ?", username).
		Scan(&u.ID, &u.Username, &u.PasswordHash, &u.Role, &u.TOTPSecret, &u.TOTPEnabled)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (d *DB) GetUserByID(id int) (*DBUser, error) {
	var u DBUser
	err := d.QueryRow("SELECT id, username, password_hash, role, totp_secret, totp_enabled FROM users WHERE id = ?", id).
		Scan(&u.ID, &u.Username, &u.PasswordHash, &u.Role, &u.TOTPSecret, &u.TOTPEnabled)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (d *DB) UpdatePassword(userID int, newHash string) error {
	_, err := d.Exec("UPDATE users SET password_hash = ? WHERE id = ?", newHash, userID)
	return err
}

func (d *DB) SetTOTPSecret(userID int, secret string) error {
	_, err := d.Exec("UPDATE users SET totp_secret = ?, totp_enabled = 0 WHERE id = ?", secret, userID)
	return err
}

func (d *DB) Enable2FA(userID int) error {
	_, err := d.Exec("UPDATE users SET totp_enabled = 1 WHERE id = ?", userID)
	return err
}

func (d *DB) Disable2FA(userID int) error {
	_, err := d.Exec("UPDATE users SET totp_secret = '', totp_enabled = 0 WHERE id = ?", userID)
	return err
}
