package btcserver

/*
import (
	"context"
	"encoding/hex"
	"strings"

	"github.com/sgostarter/i/commerr"
	"github.com/sgostarter/i/l"
	"github.com/spf13/cast"
)

func (s *BTCServer) dataRoutine() {
	logger := s.logger.WithFields(l.StringField(l.RoutineKey, "dataRoutine"))

	logger.Info("enter")

	defer logger.Info("leave")

	defer s.wg.Done()

	loop := true

	for loop {
		select {
		case <-s.ctx.Done():
			loop = false

			continue
		default:
			s.doDataProcess(s.ctx, logger)
		}
	}
}

func (s *BTCServer) doDataProcess(ctx context.Context, _logger l.Wrapper) {
	cli := s.rpcClient

	var blockHeight int64

	err := cli.CallWrapper("getblockcount", nil, &blockHeight)
	if err != nil {
		_logger.WithFields(l.ErrorField(err)).Error("call getblockcount")

		return
	}

	processedBlock, processedTxIndex, err := s.stg.LoadProcessedBlockPosition()
	if err != nil {
		_logger.WithFields(l.ErrorField(err)).Error("load processed info")

		return
	}

	if processedBlock == -1 {
		processedBlock++
	}

	processedTxIndex++

	if processedBlock > blockHeight {
		return
	}

	for curBlockIndex := processedBlock; curBlockIndex <= blockHeight-10; curBlockIndex++ {
		bLogger := _logger.WithFields(l.Int64Field("blockIndex", curBlockIndex))

		var blockHash string

		err = cli.CallWrapper("getblockhash", []any{curBlockIndex}, &blockHash)
		if err != nil {
			bLogger.WithFields(l.ErrorField(err)).Error("call getblockhash")

			return
		}

		var blockD map[string]any

		err = cli.CallWrapper("getblock", []any{blockHash, 2}, &blockD)
		if err != nil {
			bLogger.WithFields(l.ErrorField(err)).Error("call getblock")

			return
		}

		var txIndex int

		if curBlockIndex == processedBlock {
			txIndex = int(processedTxIndex)
		}

		txS := cast.ToSlice(blockD["tx"])

		for curTxIndex := txIndex; curTxIndex < len(txS); curTxIndex++ {
			txM := cast.ToStringMap(txS[curTxIndex])

			txID := cast.ToString(txM["txid"])

			tLogger := bLogger.WithFields(l.IntField("txIndex", curTxIndex), l.StringField("txID", txID))

			vIns := cast.ToSlice(txM["vin"])

			for vInIndex, vIn := range vIns {
				tILogger := tLogger.WithFields(l.IntField("vInIndex", vInIndex))

				vInM := cast.ToStringMap(vIn)

				if v, ok := vInM["coinbase"]; ok {
					coinBase, _ := hex.DecodeString(cast.ToString(v))

					bLogger.Debugf("coinbase is %s\n", curBlockIndex, string(coinBase))
				} else {
					vTxID := cast.ToString(vInM["txid"])
					vOut := cast.ToInt(vInM["vout"])

					if err = s.stg.SpendTXO(vTxID, vOut); err != nil {
						tILogger.WithFields(l.StringField("txID", vTxID), l.IntField("vOut", vOut)).
							Warn("spend txo failed")
					}
				}
			}

			vOuts := cast.ToSlice(txM["vout"])

			for vOutIdx, vOut := range vOuts {
				tOLogger := tLogger.WithFields(l.IntField("vInIndex", vOutIdx))

				vOutM := cast.ToStringMap(vOut)

				if scriptPubKeyI, ok := vOutM["scriptPubKey"]; ok {
					vs := cast.ToStringMap(scriptPubKeyI)

					asmS := cast.ToString(vs["asm"])
					hexS := cast.ToString(vs["hex"])
					address := cast.ToString(vs["address"])
					vType := cast.ToString(vs["type"])

					switch vType {
					case "pubkey":
						if asmS == "" || hexS == "" {
							tOLogger.Error("no asm or hex")

							err = commerr.ErrInternal

							return
						}

						tOLogger.Infof("pubkey out type is: %s, publicKey is %s\n",
							vType, strings.Split(asmS, " ")[0])
					case "pubkeyhash", "scripthash":
						if address == "" {
							tOLogger.WithFields(l.StringField("vType", vType)).Error("no address")

							err = commerr.ErrInternal

							return
						}

						tOLogger.Infof("%s address is %s\n", vType, address)
					case "nonstandard":
						tOLogger.Warn("nonstandard")
						hexS = ""
					case "multisig":
						if asmS == "" || hexS == "" {
							tOLogger.Error("no asm or hex")

							err = commerr.ErrInternal

							return
						}

						asmPs := strings.Split(asmS, " ")

						dbgInfo := ""
						dbgInfo += asmPs[0]
						dbgInfo += "/"
						dbgInfo += asmPs[len(asmPs)-2]
						dbgInfo += " public key:"
						for x := 1; x < len(asmPs)-2; x++ {
							dbgInfo += asmPs[x]
							dbgInfo += " "
						}

						tOLogger.Info("multisig => " + dbgInfo)
					case "nulldata":
						tOLogger.Warn("nulldata")
						hexS = ""
					case "witness_unknown":
						tOLogger.Warn("witness_unknown")
					default:
						tOLogger.WithFields(l.StringField("vType", vType)).Error("unknown vType")

						err = commerr.ErrInternal

						return
					}

					if hexS != "" {
						err = s.stg.NewTXO(txID, vOutIdx, vType, hexS, address)
						if err != nil {
							tOLogger.WithFields(l.ErrorField(err)).Error("new txo failed")

							return
						}
					}
				} else {
					tOLogger.Error("no scriptPubKey")

					return
				}
			}

			err = s.stg.StoreProcessedBlockPosition(curBlockIndex, int64(curTxIndex))
			if err != nil {
				tLogger.WithFields(l.ErrorField(err)).Error("store processed info")

				return
			}
		}
	}
}
*/
