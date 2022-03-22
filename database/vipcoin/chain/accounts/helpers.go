/*
 * Copyright 2022 Business Process Technologies. All rights reserved.
 */

package accounts

import (
	accountstypes "git.ooo.ua/vipcoin/chain/x/accounts/types"
	extratypes "git.ooo.ua/vipcoin/chain/x/types"
	"github.com/forbole/bdjuno/v2/database/types"
	"github.com/lib/pq"
)

// toAffiliatesDatabase - mapping func to database model
func toAffiliatesDatabase(affiliate *accountstypes.Affiliate) types.DBAffiliates {
	return types.DBAffiliates{
		Address:         affiliate.Address,
		AffiliationKind: accountstypes.AffiliationKind_value[affiliate.Affiliation.String()],
		Extras:          toExtrasDB(affiliate.Extras),
	}

}

// toExtrasDB - mapping func to database model
func toExtrasDB(extras []*extratypes.Extra) types.ExtraDB {
	result := make([]extratypes.Extra, 0, len(extras))
	for _, extra := range extras {
		result = append(result, *extra)
	}

	return types.ExtraDB{Extras: result}
}

// toAccountDatabase - mapping func to database model
func toAccountDatabase(account *accountstypes.Account) types.DBAccount {
	return types.DBAccount{
		Address:   account.Address,
		Hash:      account.Hash,
		PublicKey: account.PublicKey.String(),
		Kinds:     toKindsDB(account.Kinds),
		State:     int32(account.State),
		Extras:    toExtrasDB(account.Extras),
		Wallets:   account.Wallets,
	}
}

// toKindsDB - mapping func to database model
func toKindsDB(kinds []accountstypes.AccountKind) pq.Int32Array {
	result := make(pq.Int32Array, 0, len(kinds))
	for _, kind := range kinds {
		result = append(result, accountstypes.AccountKind_value[kind.String()])
	}

	return result
}
