package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"

	abcitypes "github.com/tendermint/tendermint/abci/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// AssetsApplication mongodb connection.
type AssetsApplication struct {
	db         *mongo.Client
	validators []abcitypes.ValidatorUpdate
}

var _ abcitypes.Application = (*AssetsApplication)(nil)

// NewAssetsApplication mongodb connection come from main.go .
func NewAssetsApplication(db *mongo.Client) *AssetsApplication {
	return &AssetsApplication{db: db}
}

// Info interface.
func (app *AssetsApplication) Info(req abcitypes.RequestInfo) abcitypes.ResponseInfo {
	/*return abcitypes.ResponseInfo{Data: "codechainv0.0.1",
		Version:    version.ABCIVersion,
		AppVersion: ProtocolVersion.Uint64(),
	}*/
	return abcitypes.ResponseInfo{}
}

// SetOption interface.
func (AssetsApplication) SetOption(req abcitypes.RequestSetOption) abcitypes.ResponseSetOption {
	return abcitypes.ResponseSetOption{}
}

// isValid
func (app *AssetsApplication) isValid(tx []byte) (code uint32) {
	parts := bytes.Split(tx, []byte("="))
	if len(parts) != 2 {
		return 1
	}
	return 0
}

// DeliverTx check it and save to mongodb.
func (app *AssetsApplication) DeliverTx(req abcitypes.RequestDeliverTx) abcitypes.ResponseDeliverTx {
	fmt.Println(string(req.Tx))
	code := app.isValid(req.Tx)
	if code != 0 {
		return abcitypes.ResponseDeliverTx{Code: code}
	}
	parts := bytes.Split(req.Tx, []byte("="))
	key, value := string(parts[0]), string(parts[1])
	collection := app.db.Database("chain").Collection("assets")
	assets := bson.M{"key": string(key), "value": string(value)}
	fmt.Println(assets)
	insertResult, err := collection.InsertOne(context.TODO(), assets)
	if err != nil {
		panic(err)
	}
	fmt.Println("Inserted a single document:", insertResult.InsertedID)
	return abcitypes.ResponseDeliverTx{Code: 0}
}

// CheckTx check tx format .
func (app *AssetsApplication) CheckTx(req abcitypes.RequestCheckTx) abcitypes.ResponseCheckTx {
	code := app.isValid(req.Tx)
	return abcitypes.ResponseCheckTx{Code: code, GasWanted: 1}
}

// Commit interface .
func (app *AssetsApplication) Commit() abcitypes.ResponseCommit {
	return abcitypes.ResponseCommit{Data: []byte{}}
}

// Query  query document from mongledb.
func (app *AssetsApplication) Query(reqQuery abcitypes.RequestQuery) (resQuery abcitypes.ResponseQuery) {
	parts := bytes.Split(reqQuery.Data, []byte("="))
	value := string(parts[1])
	filter := bson.M{"key": string(value)}
	collection := app.db.Database("chain").Collection("assets")
	assets := bson.M{}
	err := collection.FindOne(context.TODO(), filter).Decode(&assets)
	if err != nil {
		error := fmt.Sprintf("%s", err)
		resQuery.Code = 1
		resQuery.Log = error
		resQuery.Value = nil
	} else {
		if value, ok := assets["value"].(string); ok {
			resQuery.Value = []byte(value)
			resQuery.Info = value
			resQuery.Code = 0
			resQuery.Log = ""
		} else {
			resQuery.Value = nil
			resQuery.Code = 1
			resQuery.Log = "error type"
		}
	}
	return
}

// InitChain drop collection .
func (app *AssetsApplication) InitChain(req abcitypes.RequestInitChain) abcitypes.ResponseInitChain {
	app.validators = req.Validators
	collection := app.db.Database("chain").Collection("assets")
	collection.Drop(context.TODO())
	return abcitypes.ResponseInitChain{}
}

// BeginBlock interface.
func (app *AssetsApplication) BeginBlock(req abcitypes.RequestBeginBlock) abcitypes.ResponseBeginBlock {
	return abcitypes.ResponseBeginBlock{}
}

// EndBlock interface.
func (app *AssetsApplication) EndBlock(req abcitypes.RequestEndBlock) abcitypes.ResponseEndBlock {
	//fmt.Printf("%+v", req) test dynamic add validator
	if len(app.validators) == 0 || req.Height <= 21 {
		return abcitypes.ResponseEndBlock{}
	}
	fmt.Println(len(app.validators))
	var v abcitypes.ValidatorUpdate
	// test new validator's public key
	v.Power = 10
	v.PubKey.Type = "ed25519"
	v.PubKey.Data, _ = base64.StdEncoding.DecodeString("BsY96CRY2RK+vcVbMFpOiGQSLJARQTlDB00BbyZuwM0=")
	//
	keyExists := false
	for i := 0; i < len(app.validators); i++ {
		if bytes.Compare(app.validators[i].PubKey.Data, v.PubKey.Data) == 0 {
			keyExists = true
			break
		}
	}
	if keyExists {
		return abcitypes.ResponseEndBlock{}
	}
	app.validators = append(app.validators, v)
	return abcitypes.ResponseEndBlock{ValidatorUpdates: app.validators}
}
