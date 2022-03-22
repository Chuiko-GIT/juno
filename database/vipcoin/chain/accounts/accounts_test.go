/*
 * Copyright 2022 Business Process Technologies. All rights reserved.
 */

package accounts

import (
	"testing"

	accountstypes "git.ooo.ua/vipcoin/chain/x/accounts/types"
	extratypes "git.ooo.ua/vipcoin/chain/x/types"
	"github.com/brianvoe/gofakeit/v6"
	anytype "github.com/cosmos/cosmos-sdk/codec/types"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
)

func TestRepository_SaveAccounts(t *testing.T) {
	// db, err := sqlx.Connect("pgx", "host=10.10.1.79 port=5432 user=postgres dbname=juno password=postgres sslmode=disable")
	db, err := sqlx.Connect("pgx", "host=localhost port=5432 user=postgres dbname=juno password=postgres sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		accounts []*accountstypes.Account
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid",
			args: args{
				accounts: []*accountstypes.Account{
					{
						Hash:    "b8b6cb7629d68b3ecf9ce200f631ffc72232bc798a7db755307332a40add5e37",
						Address: "vcg1qq995wzw6zgqjm8g2twsykl7xsj0apxmtuxfdy",
						PublicKey: &anytype.Any{
							TypeUrl: "/cosmos.crypto.secp256k1.PubKey",
							Value:   []uint8{10, 33, 2, 32, 174, 170, 220, 129, 199, 203, 202, 84, 205, 96, 6, 247, 144, 46, 61, 225, 73, 220, 82, 19, 53, 39, 205, 55, 45, 114, 65, 148, 77, 198, 60},
						},
						Kinds:      []accountstypes.AccountKind{accountstypes.ACCOUNT_KIND_SYSTEM},
						State:      accountstypes.ACCOUNT_STATE_ACTIVE,
						Extras:     []*extratypes.Extra{},
						Affiliates: []*accountstypes.Affiliate{},
						Wallets:    []string{},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "valid",
			args: args{
				accounts: []*accountstypes.Account{
					{
						Hash:    gofakeit.Regex("^0x[a-fA-F0-9]{40}$"),
						Address: gofakeit.Regex("^0x[a-fA-F0-9]{40}$"),
						PublicKey: &anytype.Any{
							TypeUrl: "some type URL",
							Value:   []byte(gofakeit.Regex("^0x[a-fA-F0-9]{40}$")),
						},
						Kinds: []accountstypes.AccountKind{accountstypes.ACCOUNT_KIND_GENERAL, accountstypes.ACCOUNT_KIND_MERCHANT},
						State: accountstypes.ACCOUNT_STATE_ACTIVE,
						Extras: []*extratypes.Extra{
							{
								Kind: extratypes.EXTRA_KIND_COMMENT,
								Data: "some text",
							},
						},
						Affiliates: []*accountstypes.Affiliate{
							{
								Address:     gofakeit.Regex("^0x[a-fA-F0-9]{40}$"),
								Affiliation: accountstypes.AFFILIATION_KIND_REFERRAL,
								Extras: []*extratypes.Extra{
									{
										Kind: extratypes.EXTRA_KIND_EMAIL,
										Data: gofakeit.Email(),
									},
									{
										Kind: extratypes.EXTRA_KIND_PHONE,
										Data: gofakeit.Phone(),
									},
								},
							},
						},
						Wallets: []string{gofakeit.Regex("^0x[a-fA-F0-9]{40}$")},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Repository{
				db: db,
			}
			if err := r.SaveAccounts(tt.args.accounts...); (err != nil) != tt.wantErr {
				t.Errorf("Repository.SaveAccounts() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
