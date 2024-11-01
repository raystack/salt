package db

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type Client struct {
	*sqlx.DB
	queryTimeOut time.Duration
	cfg          Config
	host         string
}

// NewClient creates a new sqlx database client
func New(cfg Config) (*Client, error) {
	dbURL, err := url.Parse(cfg.URL)
	if err != nil {
		return nil, err
	}
	host := dbURL.Host

	db, err := sqlx.Connect(cfg.Driver, cfg.URL)
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifeTime)

	return &Client{DB: db, queryTimeOut: cfg.MaxQueryTimeout, cfg: cfg, host: host}, err
}

func (c Client) WithTimeout(ctx context.Context, op func(ctx context.Context) error) (err error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, c.queryTimeOut)
	defer cancel()

	return op(ctxWithTimeout)
}

func (c Client) WithTxn(ctx context.Context, txnOptions sql.TxOptions, txFunc func(*sqlx.Tx) error) (err error) {
	txn, err := c.BeginTxx(ctx, &txnOptions)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			switch p := p.(type) {
			case error:
				err = p
			default:
				err = errors.Errorf("%s", p)
			}
			err = txn.Rollback()
			panic(p)
		} else if err != nil {
			if rlbErr := txn.Rollback(); rlbErr != nil {
				err = fmt.Errorf("rollback error: %s while executing: %w", rlbErr, err)
			} else {
				err = fmt.Errorf("rollback: %w", err)
			}
		} else {
			err = txn.Commit()
		}
	}()

	err = txFunc(txn)
	return err
}

// ConnectionURL fetch the database connection url
func (c *Client) ConnectionURL() string {
	return c.cfg.URL
}

// Host fetch the database host information
func (c *Client) Host() string {
	return c.host
}

// Close closes the database connection
func (c *Client) Close() error {
	return c.DB.Close()
}
