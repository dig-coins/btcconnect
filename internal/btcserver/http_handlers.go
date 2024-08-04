package btcserver

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dig-coins/btcconnect/pkg/btctx"
	"github.com/dig-coins/btcconnect/pkg/helper"
	"github.com/dig-coins/btcconnect/pkg/share"
	"github.com/gin-gonic/gin"
	"github.com/sgostarter/libeasygo/ptl"
	"github.com/spf13/cast"
	"golang.org/x/exp/slices"
)

func (s *BTCServer) getCacheWallets() (wallets []string) {
	i, ok := s.cacheData.Get("wallets")
	if ok {
		return cast.ToStringSlice(i)
	}

	wallets, err := s.listWallets()
	if err == nil {
		s.cacheData.Set("wallets", wallets, time.Second*10)
	}

	return
}

func (s *BTCServer) handlerGetWallets(c *gin.Context) {
	wallets, code, msg := s.doHTTPGetWallets(c)

	var resp ptl.ResponseWrapper

	if resp.Apply(code, msg) {
		resp.Resp = wallets
	}

	c.JSON(http.StatusOK, resp)
}

func (s *BTCServer) doHTTPGetWallets(_ *gin.Context) (wallets []string, code ptl.Code, msg string) {
	wallets, err := s.listWallets()
	if err != nil {
		code = ptl.CodeErrCommunication

		msg = err.Error()

		return
	}

	return
}

func (s *BTCServer) handlerGetBalance(c *gin.Context) {
	balance, code, msg := s.doGetBalance(c)

	var resp ptl.ResponseWrapper

	if resp.Apply(code, msg) {
		resp.Resp = balance
	}

	c.JSON(http.StatusOK, resp)
}

func (s *BTCServer) doGetBalance(c *gin.Context) (allBalance int64, code ptl.Code, msg string) {
	var m map[string]any

	err := c.ShouldBindJSON(&m)
	if err != nil {
		code = ptl.CodeErrCommunication

		msg = err.Error()

		return
	}

	var wallets []string

	if m["wallets"] != nil {
		wallets, err = cast.ToStringSliceE(m["wallets"])
		if err != nil {
			msg = err.Error()

			return
		}
	}

	if len(wallets) == 0 {
		wallets = s.getCacheWallets()
	}

	if len(wallets) == 0 {
		msg = "no wallets"

		return
	}

	for _, wallet := range wallets {
		var balance float64

		balance, err = s.getBalance(wallet)
		if err != nil {
			msg = err.Error()

			return
		}

		allBalance += helper.UnitBTC2SatoshiBTC(balance)
	}

	return
}

type SimplePayParams struct {
	Wallets            []string      `json:"wallets,omitempty"`
	Outputs            []TransOutput `json:"outputs"`
	ChangeAddress      string        `json:"change_address"`
	TotalFee           int64         `json:"total_fee"`
	ConfirmationTarget int           `json:"confirmation_target,omitempty"`
	FeeSatoshiPerKB    int64         `json:"fee_satoshi_per_kb,omitempty"`
	MinTransFlag       bool          `json:"min_trans_flag,omitempty"`
}

type UnsignedTxResponse struct {
	UnsignedTxHex   string               `json:"unsigned_tx_hex"`
	UnsignedTx      *btctx.UnsignedTx    `json:"unsigned_tx"`
	WalletUnspent   map[string][]Unspent `json:"wallet_unspent"`
	FeeSatoshiPerKB int64                `json:"fee_satoshi_per_kb"`
	Fee             int64                `json:"fee"`
	CommandJSON     string               `json:"command_json"`
}

func (resp *UnsignedTxResponse) think() {
	if resp.UnsignedTx == nil {
		return
	}

	var totalAmount, outAmount int64

	for _, input := range resp.UnsignedTx.Inputs {
		totalAmount += input.Amount
	}

	for _, output := range resp.UnsignedTx.Outputs {
		outAmount += output.Amount
	}

	resp.Fee = totalAmount - outAmount
	resp.CommandJSON, _ = share.MarshalCommandToJSON(share.Command{
		CommandType: share.CommandTypeGenTx,
		Input:       resp.UnsignedTxHex,
	})
}

func (s *BTCServer) handlerSimplePay(c *gin.Context) {
	unsignedTx, code, msg := s.doSimplePay(c)

	var resp ptl.ResponseWrapper

	if resp.Apply(code, msg) {
		resp.Resp = unsignedTx
	}

	c.JSON(http.StatusOK, resp)
}

func (s *BTCServer) doSimplePay(c *gin.Context) (unsignedTx UnsignedTxResponse, code ptl.Code, msg string) {
	var req SimplePayParams

	err := c.BindJSON(&req)
	if err != nil {
		code = ptl.CodeErrInvalidArgs

		msg = err.Error()

		return
	}

	if len(req.Outputs) == 0 && req.ChangeAddress == "" {
		code = ptl.CodeErrInvalidArgs

		return
	}

	unsignedTx.WalletUnspent, err = s.getUnspent(req.Wallets, nil)
	if err != nil {
		code = ptl.CodeErrInternal

		msg = err.Error()

		return
	}

	feeSatoshiPerKB, err := s.fixFeeSatoshiPerKB(req.FeeSatoshiPerKB, req.ConfirmationTarget)
	if err != nil {
		msg = err.Error()

		return
	}

	unsignedTx.FeeSatoshiPerKB = feeSatoshiPerKB

	unspents := make([]Unspent, 0, 10)
	for _, wUnspents := range unsignedTx.WalletUnspent {
		unspents = append(unspents, wUnspents...)
	}

	unsignedTx.UnsignedTx, unsignedTx.UnsignedTxHex, err = s.genToUnsignedTx(feeSatoshiPerKB,
		req.TotalFee, unspents, req.Outputs, req.ChangeAddress, req.MinTransFlag)
	if err != nil {
		msg = err.Error()

		return
	}

	if req.TotalFee > 0 {
		unsignedTx.FeeSatoshiPerKB = 0
	}

	unsignedTx.think()

	return
}

type PayTx struct {
	ID   string `json:"id"`
	VOut uint32 `json:"v_out"`
}

type PayInput struct {
	Wallet    string   `json:"wallet,omitempty"`
	Addresses []string `json:"addresses,omitempty"`
	PayTxs    []PayTx  `json:"pay_txs,omitempty"`
}

type PayParams struct {
	Inputs  []PayInput    `json:"inputs"`
	Outputs []TransOutput `json:"outputs"`

	ChangeAddress      string `json:"change_address"`
	TotalFee           int64  `json:"total_fee"`
	ConfirmationTarget int    `json:"confirmation_target,omitempty"`
	FeeSatoshiPerKB    int64  `json:"fee_satoshi_per_kb,omitempty"`
	MinTransFlag       bool   `json:"min_trans_flag,omitempty"`
}

func (s *BTCServer) handlerPay(c *gin.Context) {
	unsignedTx, code, msg := s.doPay(c)

	var resp ptl.ResponseWrapper

	if resp.Apply(code, msg) {
		resp.Resp = unsignedTx
	}

	c.JSON(http.StatusOK, resp)
}

func (s *BTCServer) doPay(c *gin.Context) (unsignedTx UnsignedTxResponse, code ptl.Code, msg string) {
	var req PayParams

	err := c.BindJSON(&req)
	if err != nil {
		code = ptl.CodeErrInvalidArgs

		msg = err.Error()

		return
	}

	if len(req.Outputs) == 0 && req.ChangeAddress == "" {
		code = ptl.CodeErrInvalidArgs

		return
	}

	unsignedTx.WalletUnspent, err = s.getUnspent(nil, req.Inputs)
	if err != nil {
		code = ptl.CodeErrInternal

		msg = err.Error()

		return
	}

	feeSatoshiPerKB, err := s.fixFeeSatoshiPerKB(req.FeeSatoshiPerKB, req.ConfirmationTarget)
	if err != nil {
		msg = err.Error()

		return
	}

	unsignedTx.FeeSatoshiPerKB = feeSatoshiPerKB

	unspents := make([]Unspent, 0, 10)
	for _, wUnspents := range unsignedTx.WalletUnspent {
		unspents = append(unspents, wUnspents...)
	}

	unsignedTx.UnsignedTx, unsignedTx.UnsignedTxHex, err = s.genToUnsignedTx(feeSatoshiPerKB, req.TotalFee,
		unspents, req.Outputs, req.ChangeAddress, req.MinTransFlag)
	if err != nil {
		msg = err.Error()

		return
	}

	if req.TotalFee > 0 {
		unsignedTx.FeeSatoshiPerKB = 0
	}

	unsignedTx.think()

	return
}

func (s *BTCServer) fixFeeSatoshiPerKB(feeSatoshiPerKB int64, confirmationTarget int) (int64, error) {
	if feeSatoshiPerKB > 0 {
		return feeSatoshiPerKB, nil
	}

	if confirmationTarget <= 0 {
		confirmationTarget = 2
	}

	return s.estimateSmartFee(confirmationTarget)
}

func (s *BTCServer) selectWalletsByPayInputs(inputs []PayInput) []string {
	if len(inputs) == 0 {
		return s.getCacheWallets()
	}

	walletsM := make(map[string]int)

	for _, input := range inputs {
		if input.Wallet == "" {
			return s.getCacheWallets()
		}

		walletsM[input.Wallet]++
	}

	wallets := make([]string, 0, len(walletsM))

	for wallet := range walletsM {
		wallets = append(wallets, wallet)
	}

	return wallets
}

func (s *BTCServer) getUnspent(wallets []string, inputs []PayInput) (walletUnspent map[string][]Unspent, err error) {
	walletTxs := make(map[string][]Unspent)

	if len(wallets) == 0 {
		if len(inputs) > 0 {
			wallets = s.selectWalletsByPayInputs(inputs)
		} else {
			wallets = s.getCacheWallets()
		}
	}

	for _, wallet := range wallets {
		walletTxs[wallet], err = s.listWalletUnspent(wallet)
		if err != nil {
			return
		}
	}

	if len(inputs) == 0 {
		walletUnspent = walletTxs

		return
	}

	walletUnspent = make(map[string][]Unspent)
	selectIDs := make(map[string]any)

	fnSelectUnspent := func(input PayInput) {
		for wallet, unspents := range walletTxs {
			if input.Wallet != "" && input.Wallet != wallet {
				continue
			}

			for _, unspent := range unspents {
				if len(input.Addresses) > 0 && !slices.Contains(input.Addresses, unspent.Address) {
					continue
				}

				if len(input.PayTxs) == 0 {
					walletUnspent[wallet] = append(walletUnspent[wallet], unspent)
					selectIDs[fmt.Sprintf("%s:%d", unspent.TxID, unspent.VOut)] = true

					continue
				}

				for _, tx := range input.PayTxs {
					if tx.ID == unspent.TxID && tx.VOut == unspent.VOut {
						walletUnspent[wallet] = append(walletUnspent[wallet], unspent)
						selectIDs[fmt.Sprintf("%s:%d", unspent.TxID, unspent.VOut)] = true

						break
					}
				}
			}
		}
	}

	for _, input := range inputs {
		fnSelectUnspent(input)
	}

	return
}

func (s *BTCServer) genToUnsignedTx(feeSatoshiPerKB, totalFee int64, unspent []Unspent, rawOutputs []TransOutput,
	changeAddress string, minTransFlag bool) (unsignedTx *btctx.UnsignedTx, unsignedTxHex string, err error) {
	inputs, outputs, err := s.selectUnspentInputs(unspent, nil, feeSatoshiPerKB, totalFee, rawOutputs,
		changeAddress, minTransFlag)
	if err != nil {
		return
	}

	unsignedTx, unsignedTxHex, err = s.GetUnsignedTxEx(inputs, outputs)

	return
}

type NetworkFeeResponse struct {
	CoreFee2  int64 `json:"core_fee_2"`
	CoreFee6  int64 `json:"core_fee_6"`
	CoreFee20 int64 `json:"core_fee_20"`
	CoreFee40 int64 `json:"core_fee_40"`
	CoinExFee int64 `json:"coin_ex_fee"`
}

func (s *BTCServer) handlerGetNetworkFee(c *gin.Context) {
	var feeResp NetworkFeeResponse

	if s.cfg.CoinType == share.CoinTypeBTC {
		feeResp.CoinExFee, _, _ = s.getNetworkFee4CoinExCom(c.Request.Context())
	}

	feeResp.CoreFee2, _ = s.estimateSmartFee(2)
	feeResp.CoreFee6, _ = s.estimateSmartFee(6)
	feeResp.CoreFee20, _ = s.estimateSmartFee(20)
	feeResp.CoreFee40, _ = s.estimateSmartFee(40)

	c.JSON(http.StatusOK, ptl.ResponseWrapper{
		Code: ptl.CodeSuccess,
		Resp: feeResp,
	})
}

func (s *BTCServer) getNetworkFee4CoinExCom(ctx context.Context) (bestTxFee int64, code ptl.Code, msg string) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://explorer.coinex.com/res/btc/network", nil)
	if err != nil {
		msg = err.Error()

		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		msg = err.Error()

		return
	}

	if resp == nil {
		msg = "no resp"

		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		msg = fmt.Sprintf("%d", resp.StatusCode)

		return
	}

	var m map[string]any

	err = json.NewDecoder(resp.Body).Decode(&m)
	if err != nil {
		msg = err.Error()

		return
	}

	if cast.ToInt(m["code"]) != 0 {
		msg = "invalid code"

		return
	}

	di, ok := m["data"]
	if !ok {
		msg = "invalid data"

		return
	}

	feeBTC, err := cast.ToFloat64E(cast.ToStringMap(di)["best_tx_fee"])
	if err != nil {
		msg = err.Error()

		return
	}

	bestTxFee = helper.UnitBTC2SatoshiBTC(feeBTC)

	return
}

func (s *BTCServer) handlerGetUnspentOfGroupWalletAddress(c *gin.Context) {
	unspent, code, msg := s.getUnspentOfGroupOfWalletAddress(c)

	var resp ptl.ResponseWrapper

	if resp.Apply(code, msg) {
		resp.Resp = unspent
	}

	c.JSON(http.StatusOK, resp)
}

type UnspentVO struct {
	TxID          string `json:"tx_id"`
	VOut          uint32 `json:"v_out"`
	Label         string `json:"label"`
	Address       string `json:"address"`
	Confirmations int    `json:"confirmations"`
	Amount        int64  `json:"amount"`
}

type AddressUnspentListVO struct {
	Address   string      `json:"address"`
	Label     string      `json:"label"`
	Amount    int64       `json:"amount"`
	UnspentVO []UnspentVO `json:"unspent"`
}

type WalletUnspentVO struct {
	Wallet  string                 `json:"wallet"`
	Amount  int64                  `json:"amount"`
	Unspent []AddressUnspentListVO `json:"unspent_list"`
}

func (s *BTCServer) getUnspentOfGroupOfWalletAddress(c *gin.Context) (
	unspentListVO []*WalletUnspentVO, code ptl.Code, msg string) {
	var m map[string]any

	err := c.ShouldBindJSON(&m)
	if err != nil {
		code = ptl.CodeErrCommunication

		msg = err.Error()

		return
	}

	var wallets []string

	if m["wallets"] != nil {
		wallets, err = cast.ToStringSliceE(m["wallets"])
		if err != nil {
			msg = err.Error()

			return
		}
	}

	if len(wallets) == 0 {
		wallets = s.getCacheWallets()
	}

	if len(wallets) == 0 {
		msg = "no wallets"

		return
	}

	var unspents []Unspent

	for _, wallet := range wallets {
		walletUnspent := &WalletUnspentVO{
			Wallet: wallet,
		}

		unspentListVO = append(unspentListVO, walletUnspent)

		unspents, _ = s.listWalletUnspent(wallet)

		addressUnspent := make(map[string]*AddressUnspentListVO)

		for _, unspent := range unspents {
			if !unspent.Safe {
				continue
			}

			u, ok := addressUnspent[unspent.Address]
			if !ok {
				u = &AddressUnspentListVO{
					Address: unspent.Address,
					Label:   unspent.Label,
				}

				addressUnspent[unspent.Address] = u
			}

			u.Amount += helper.UnitBTC2SatoshiBTC(unspent.Amount)
			u.UnspentVO = append(u.UnspentVO, UnspentVO{
				TxID:          unspent.TxID,
				VOut:          unspent.VOut,
				Label:         unspent.Label,
				Address:       unspent.Address,
				Confirmations: unspent.Confirmations,
				Amount:        helper.UnitBTC2SatoshiBTC(unspent.Amount),
			})
		}

		walletUnspent.Unspent = make([]AddressUnspentListVO, 0, len(addressUnspent))

		for _, vo := range addressUnspent {
			walletUnspent.Amount += vo.Amount
			walletUnspent.Unspent = append(walletUnspent.Unspent, AddressUnspentListVO{
				Address:   vo.Address,
				Label:     vo.Label,
				Amount:    vo.Amount,
				UnspentVO: vo.UnspentVO,
			})
		}
	}

	return
}

func (s *BTCServer) handlerResignedTx(c *gin.Context) {
	unsignedTx, code, msg := s.doResignedTx(c)

	var resp ptl.ResponseWrapper

	if resp.Apply(code, msg) {
		resp.Resp = unsignedTx
	}

	c.JSON(http.StatusOK, resp)
}

type ResignedTxRequest struct {
	UnsignedTx         string `json:"unsigned_tx"`
	TotalFee           int64  `json:"total_fee"`
	ConfirmationTarget int    `json:"confirmation_target,omitempty"`
	FeeSatoshiPerKB    int64  `json:"fee_satoshi_per_kb,omitempty"`
}

func (s *BTCServer) doResignedTx(c *gin.Context) (unsignedTx UnsignedTxResponse, code ptl.Code, msg string) {
	var req ResignedTxRequest

	err := c.BindJSON(&req)
	if err != nil {
		code = ptl.CodeErrInvalidArgs

		msg = err.Error()

		return
	}

	if req.UnsignedTx == "" {
		code = ptl.CodeErrInvalidArgs

		msg = "no unsigned tx"

		return
	}

	tx, err := btctx.UnmarshalUnsignedTx(req.UnsignedTx)
	if err != nil {
		code = ptl.CodeErrInvalidArgs

		msg = err.Error()

		return
	}

	feeSatoshiPerKB, err := s.fixFeeSatoshiPerKB(req.FeeSatoshiPerKB, req.ConfirmationTarget)
	if err != nil {
		msg = err.Error()

		return
	}

	var changeAddress string

	for idx := 0; idx < len(tx.Outputs); idx++ {
		if tx.Outputs[idx].ChangeFlag {
			changeAddress = tx.Outputs[idx].Address
			tx.Outputs = slices.Delete(tx.Outputs, idx, idx+1)

			break
		}
	}

	if changeAddress == "" {
		code = ptl.CodeErrInvalidArgs

		msg = "no change address"

		return
	}

	var inputs []TransInput
	for _, input := range tx.Inputs {
		inputs = append(inputs, TransInput{
			TxID:         input.TxID,
			VOut:         input.VOut,
			PrivateKey:   "",
			Address:      input.Address,
			Amount:       input.Amount,
			RedeemScript: input.RedeemScript,
		})
	}

	var outputs []TransOutput
	for _, output := range tx.Outputs {
		outputs = append(outputs, TransOutput{
			Address:    output.Address,
			Amount:     output.Amount,
			ChangeFlag: output.ChangeFlag,
		})
	}

	changeAmount, err := s.genTransToTxCalcChange(inputs, outputs, changeAddress, feeSatoshiPerKB, req.TotalFee)
	if err != nil {
		code = ptl.CodeErrInternal

		msg = err.Error()

		return
	}

	tx.Outputs = append(tx.Outputs, btctx.Output{
		Address:    changeAddress,
		Amount:     changeAmount,
		ChangeFlag: true,
	})

	unsignedTx.UnsignedTx = tx

	unsignedTx.UnsignedTxHex, err = btctx.MarshalUnsignedTx(tx)
	if err != nil {
		code = ptl.CodeErrInternal

		msg = err.Error()

		return
	}

	unsignedTx.FeeSatoshiPerKB = feeSatoshiPerKB

	unsignedTx.think()

	return
}

func (s *BTCServer) handlerLoadUnsignedTx(c *gin.Context) {
	unsignedTx, code, msg := s.unsignedTxLoad(c)

	var resp ptl.ResponseWrapper

	if resp.Apply(code, msg) {
		resp.Resp = unsignedTx
	}

	c.JSON(http.StatusOK, resp)
}

func (s *BTCServer) unsignedTxLoad(c *gin.Context) (unsignedTx UnsignedTxResponse, code ptl.Code, msg string) {
	var m map[string]string

	err := c.BindJSON(&m)
	if err != nil {
		code = ptl.CodeErrInvalidArgs

		msg = err.Error()

		return
	}

	unsignedTxHex := cast.ToString(m["unsigned_tx"])
	if unsignedTxHex == "" {
		code = ptl.CodeErrInvalidArgs

		msg = "no unsigned tx"

		return
	}

	tx, err := btctx.UnmarshalUnsignedTx(unsignedTxHex)
	if err != nil {
		code = ptl.CodeErrInvalidArgs

		msg = err.Error()

		return
	}

	unsignedTx.UnsignedTx = tx

	unsignedTx.UnsignedTxHex, err = btctx.MarshalUnsignedTx(tx)
	if err != nil {
		code = ptl.CodeErrInternal

		msg = err.Error()

		return
	}

	unsignedTx.FeeSatoshiPerKB = 0

	unsignedTx.think()

	return
}

func (s *BTCServer) handleTxBroadcast(c *gin.Context) {
	txID, code, msg := s.handleTxBroadcastInner(c)

	var resp ptl.ResponseWrapper

	if resp.Apply(code, msg) {
		resp.Resp = txID
	}

	c.JSON(http.StatusOK, resp)
}

func (s *BTCServer) handleTxBroadcastInner(c *gin.Context) (txID string, code ptl.Code, msg string) {
	var m map[string]any

	err := c.BindJSON(&m)
	if err != nil {
		code = ptl.CodeErrCommunication
		msg = err.Error()

		return
	}

	tx := cast.ToString(m["tx"])

	if tx == "" {
		code = ptl.CodeErrInvalidArgs
		msg = "no tx"

		return
	}

	txID, err = s.SendRawTransaction(tx)
	if err != nil {
		code = ptl.CodeErrInternal
		msg = err.Error()

		return
	}

	code = ptl.CodeSuccess

	return
}
