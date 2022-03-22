/*
 * Copyright 2022 Business Process Technologies. All rights reserved.
 */

package accounts

import (
	"context"
	"database/sql"

	accountstypes "git.ooo.ua/vipcoin/chain/x/accounts/types"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type (
	// repository - defines a repository for accounts repository
	Repository struct {
		db *sqlx.DB
	}
)

// NewRepository constructor.
func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r Repository) SaveAccounts(accounts ...*accountstypes.Account) error {
	if len(accounts) == 0 {
		return nil
	}

	tx, err := r.db.BeginTxx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	defer tx.Rollback()

	query := `INSERT INTO vipcoin_chain_accounts_accounts 
			(address, hash, public_key, kinds, state, extras, affiliates, wallets) 
		VALUES 
			(:address, :hash, :public_key, :kinds, :state, :extras, :affiliates, :wallets)`

	for _, acc := range accounts {
		accountDB := toAccountDatabase(acc)
		if accountDB.Affiliates, err = saveAffiliates(tx, acc.Affiliates); err != nil {
			return err
		}

		if _, err := tx.NamedExec(query, accountDB); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func saveAffiliates(tx *sqlx.Tx, affiliates []*accountstypes.Affiliate) (pq.Int64Array, error) {
	if len(affiliates) == 0 {
		return pq.Int64Array{}, nil
	}

	query := `INSERT INTO vipcoin_chain_accounts_affiliates
		(address, affiliation_kind, extras)
	VALUES
		(:address, :affiliation_kind, :extras)
	RETURNING id`

	resultID := make(pq.Int64Array, len(affiliates))
	for index, affiliate := range affiliates {
		resp, err := tx.NamedQuery(query, toAffiliatesDatabase(affiliate))
		if err != nil {
			return pq.Int64Array{}, err
		}

		for resp.Next() {
			if err := resp.Scan(&resultID[index]); err != nil {
				return pq.Int64Array{}, err
			}
		}

		if err := resp.Err(); err != nil {
			return pq.Int64Array{}, err
		}
	}

	return resultID, nil
}

// func updateAffiliates(tx *sqlx.Tx, affiliates []*accountstypes.Affiliate) error {
// 	if len(affiliates) == 0 {
// 		return nil
// 	}

// 	query := `UPDATE vipcoin_chain_accounts_affiliates SET
// 				address = :address, affiliation_kind = :affiliation_kind, extras = :extras
// 			WHERE id = :id`

// 	resultID := make(pq.Int64Array, len(affiliates))
// 	for index, affiliate := range affiliates {
// 		if _, err := tx.NamedExec(query, toAffiliatesDatabase(affiliate)); err != nil {
// 			return err
// 		}

// 		resp, err := tx.NamedQuery(query, toAffiliatesDatabase(affiliate))
// 		if err != nil {
// 			return err
// 		}

// 		for resp.Next() {
// 			if err := resp.Scan(&resultID[index]); err != nil {
// 				return err
// 			}
// 		}

// 		if err := resp.Err(); err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }
