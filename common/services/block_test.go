package services

import (
	"context"
	"github.com/v4lproik/simple-blockchain-quickstart/common/models"
	"github.com/v4lproik/simple-blockchain-quickstart/test"
	"os"
	"sync"
	"testing"
	"time"
)

func init() {
	test.InitTestContext()
}

var acc models.Account

func TestFileBlockService_Mine(t *testing.T) {
	type fields struct {
		mu               sync.Mutex
		db               *os.File
		miningComplexity uint32
	}
	type args struct {
		ctx context.Context
		pb  PendingBlock
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.Block
		wantErr bool
	}{
		{
			name: "mining a block should return a block with once",
			fields: fields{
				miningComplexity: 1,
			},
			args: args{
				ctx: context.Background(),
				pb: PendingBlock{
					models.Hash{},
					1,
					uint64(time.Now().Unix()),
					acc,
					[]models.Transaction{*models.NewTransaction(acc, acc, 10, models.SELF_REWARD)},
				},
			},
			wantErr: false,
		},
	}

	//define variables
	ethAccount, _ := test.KeyStoreService.NewKeystoreAccount("password")
	acc = models.Account(ethAccount.String())

	//run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &FileBlockService{
				mu:               tt.fields.mu,
				db:               tt.fields.db,
				miningComplexity: tt.fields.miningComplexity,
			}
			_, err := a.Mine(tt.args.ctx, tt.args.pb)
			if (err != nil) != tt.wantErr {
				t.Errorf("Mine() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
