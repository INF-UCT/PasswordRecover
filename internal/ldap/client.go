package ldap

import (
	"errors"
	"fmt"

	"chpassword/internal/config"

	goldap "github.com/go-ldap/ldap/v3"
)

var ErrPasswordPolicy = errors.New("la contraseña no cumple con la política del servidor")

type Client struct {
	cfg config.LDAPConfig
}

func New(cfg config.LDAPConfig) *Client {
	return &Client{cfg: cfg}
}

func (c *Client) connect() (*goldap.Conn, error) {
	addr := fmt.Sprintf("%s:%s", c.cfg.Host, c.cfg.Port)

	var (
		conn *goldap.Conn
		err  error
	)

	if c.cfg.UseTLS {
		conn, err = goldap.DialTLS("tcp", addr, nil)
	} else {
		conn, err = goldap.DialURL("ldap://" + addr)
	}
	if err != nil {
		return nil, err
	}

	if err := conn.Bind(c.cfg.BindDN, c.cfg.BindPassword); err != nil {
		conn.Close()
		return nil, err
	}

	return conn, nil
}

// UserExists busca un usuario por su correo en LDAP.
// Retorna si existe, su DN, y un posible error.
func (c *Client) UserExists(email string) (bool, string, error) {
	conn, err := c.connect()
	if err != nil {
		return false, "", err
	}
	defer conn.Close()

	filter := fmt.Sprintf("(mail=%s)", goldap.EscapeFilter(email))
	req := goldap.NewSearchRequest(
		c.cfg.BaseDN,
		goldap.ScopeWholeSubtree,
		goldap.NeverDerefAliases,
		0, 0, false,
		filter,
		[]string{"dn"},
		nil,
	)

	result, err := conn.Search(req)
	if err != nil {
		return false, "", err
	}

	if len(result.Entries) == 0 {
		return false, "", nil
	}

	return true, result.Entries[0].DN, nil
}

// ChangePassword cambia la contraseña de un usuario identificado por su DN.
func (c *Client) ChangePassword(userDN, newPassword string) error {
	conn, err := c.connect()
	if err != nil {
		return err
	}
	defer conn.Close()

	req := goldap.NewPasswordModifyRequest(userDN, "", newPassword)
	_, err = conn.PasswordModify(req)
	if err != nil {
		var ldapErr *goldap.Error
		if errors.As(err, &ldapErr) && ldapErr.ResultCode == goldap.LDAPResultConstraintViolation {
			return ErrPasswordPolicy
		}
		return err
	}

	return nil
}
