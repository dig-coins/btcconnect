package redistorage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/rueidis"
	"github.com/spf13/cast"
)

const (
	redisKeyProcessedBlock        = "processed_block"
	redisKeyProcessedBlockTxIndex = "processed_block_tx_index"
)

type Storage interface {
	StoreProcessedBlockPosition(block, txIndex int64) (err error)
	LoadProcessedBlockPosition() (block, txIndex int64, err error)

	NewTXO(txID string, vOut int, scriptPubKeyType, scriptPubKeyHex string, address string) (err error)
	SpendTXO(txID string, vOut int) (err error)
	GetUnspentTXO(txID string, vOut int) (scriptPubKeyType, scriptPubKeyHex, address string, err error)
}

func NewRedisStorage() (Storage, error) {
	stg := &rStorage{}

	err := stg.init()
	if err != nil {
		return nil, err
	}

	return stg, nil
}

type rStorage struct {
	rClient rueidis.Client
}

func (impl *rStorage) init() (err error) {
	rClient, err := rueidis.NewClient(rueidis.ClientOption{
		InitAddress: []string{"127.0.0.1:8300"}, Password: "redis_default_pass"})
	if err != nil {
		return
	}

	impl.rClient = rClient

	return
}

//
//
//

type RedisCommandProc func(ctx context.Context, client rueidis.Client)

func (impl *rStorage) doRedisCommand(proc RedisCommandProc) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	proc(ctx, impl.rClient)
}

func (impl *rStorage) doRedisCommandDedicated(proc RedisCommandProc) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_ = impl.rClient.Dedicated(func(_ rueidis.DedicatedClient) error {
		proc(ctx, impl.rClient)

		return nil
	})
}

//
//
//

func (impl *rStorage) StoreProcessedBlockPosition(block, txIndex int64) (err error) {
	impl.doRedisCommandDedicated(func(ctx context.Context, client rueidis.Client) {
		client.Do(ctx, client.B().Watch().Key(redisKeyProcessedBlock, redisKeyProcessedBlockTxIndex).Build())

		rrs := client.DoMulti(ctx,
			client.B().Multi().Build(),
			client.B().Set().Key(redisKeyProcessedBlock).Value(cast.ToString(block)).Build(),
			client.B().Set().Key(redisKeyProcessedBlockTxIndex).Value(cast.ToString(txIndex)).Build(),
			client.B().Exec().Build(),
		)

		for _, rr := range rrs {
			if rr.Error() != nil {
				err = rr.Error()

				break
			}
		}
	})

	return
}

func (impl *rStorage) LoadProcessedBlockPosition() (block, txIndex int64, err error) {
	impl.doRedisCommand(func(ctx context.Context, client rueidis.Client) {
		rrs := client.DoMulti(ctx,
			client.B().Get().Key(redisKeyProcessedBlock).Build(),
			client.B().Get().Key(redisKeyProcessedBlockTxIndex).Build(),
		)

		if len(rrs) != 2 {
			err = errors.New("internal")

			return
		}

		err = rrs[0].Error()
		if err != nil {
			if rueidis.IsRedisNil(err) {
				block = -1
				txIndex = -1

				err = nil
			}

			return
		}

		block, err = rrs[0].AsInt64()
		if err != nil {
			return
		}

		err = rrs[1].Error()
		if err != nil {
			if rueidis.IsRedisNil(err) {
				txIndex = -1

				err = nil
			}

			return
		}

		txIndex, err = rrs[1].AsInt64()
	})

	return
}

func (impl *rStorage) utxoRedisKey(txID string, vOut int) string {
	return fmt.Sprintf("utxo:%s:%d", txID, vOut)
}

func (impl *rStorage) NewTXO(txID string, vOut int, scriptPubKeyType, scriptPubKeyHex string, address string) (err error) {
	impl.doRedisCommand(func(ctx context.Context, client rueidis.Client) {
		err = client.Do(ctx, client.B().Hset().Key(impl.utxoRedisKey(txID, vOut)).FieldValue().
			FieldValue("scriptPubKeyType", scriptPubKeyType).
			FieldValue("scriptPubKeyHex", scriptPubKeyHex).
			FieldValue("address", address).Build()).Error()
	})

	return
}

func (impl *rStorage) SpendTXO(txID string, vOut int) (err error) {
	impl.doRedisCommand(func(ctx context.Context, client rueidis.Client) {
		var delCnt int64

		delCnt, err = client.Do(ctx, client.B().Del().Key(impl.utxoRedisKey(txID, vOut)).Build()).ToInt64()

		if err == nil && delCnt != 1 {
			err = errors.New("internal")
		}
	})

	return
}

func (impl *rStorage) GetUnspentTXO(txID string, vOut int) (scriptPubKeyType, scriptPubKeyHex string, address string, err error) {
	impl.doRedisCommand(func(ctx context.Context, client rueidis.Client) {
		var vs []string

		vs, err = client.Do(ctx, client.B().Hmget().Key(impl.utxoRedisKey(txID, vOut)).
			Field("scriptPubKeyType", "scriptPubKeyHex", "address").Build()).AsStrSlice()
		if err != nil {
			return
		}

		if len(vs) != 3 {
			err = errors.New("internal")

			return
		}

		scriptPubKeyType = vs[0]
		scriptPubKeyHex = vs[1]
		address = vs[2]
	})

	return
}
