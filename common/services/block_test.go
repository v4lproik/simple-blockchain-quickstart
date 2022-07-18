package services

import (
	"context"
	"errors"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/v4lproik/simple-blockchain-quickstart/common/models"
	"github.com/v4lproik/simple-blockchain-quickstart/test"
)

func init() {
	test.InitTestContext()
}

var acc models.Account
var ctx, _ = context.WithTimeout(context.Background(), 1*time.Millisecond)

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
		want    interface{}
		wantErr bool
	}{
		{
			name: "mining a block should return a block with a nonce",
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
		{
			name: "mining a block with context error should return error",
			fields: fields{
				miningComplexity: 10,
			},
			args: args{
				ctx: ctx,
				pb: PendingBlock{
					models.Hash{},
					1,
					uint64(time.Now().Unix()),
					acc,
					[]models.Transaction{*models.NewTransaction(acc, acc, 10, models.SELF_REWARD)},
				},
			},
			wantErr: true,
			want:    errors.New("Mine: mining task has been shutdown"),
		},
	}

	// define variables
	ss, _ := NewEthKeystore(test.KeystoreDirPath)
	ethAccount, _ := ss.NewKeystoreAccount("password")
	acc = models.Account(ethAccount.String())

	// run tests
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
			if (err != nil) && err.Error() != tt.want.(error).Error() {
				t.Errorf("Mine() error = %v, wantErr %v", err, tt.want)
				return
			}
		})
	}
}
