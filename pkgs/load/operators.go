package load

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"net/url"

	"github.com/bloxapp/ssv-dkg/pkgs/initiator"
)

func LoadOperatorsJson(operatorsMetaData []byte) (initiator.Operators, error) {
	opmap := make(map[uint64]initiator.Operator)
	var operators []initiator.OperatorDataJson
	err := json.Unmarshal(bytes.TrimSpace(operatorsMetaData), &operators)
	if err != nil {
		return nil, err
	}
	for _, opdata := range operators {
		_, err := url.ParseRequestURI(opdata.Addr)
		if err != nil {
			return nil, fmt.Errorf("invalid operator URL")
		}
		operatorKeyByte, err := base64.StdEncoding.DecodeString(opdata.PubKey)
		if err != nil {
			return nil, err
		}
		pemBlock, _ := pem.Decode(operatorKeyByte)
		if pemBlock == nil {
			return nil, fmt.Errorf("wrong pub key string")
		}
		pbKey, err := x509.ParsePKIXPublicKey(pemBlock.Bytes)
		if err != nil {
			return nil, err
		}

		opmap[opdata.ID] = initiator.Operator{
			Addr:   opdata.Addr,
			ID:     opdata.ID,
			PubKey: pbKey.(*rsa.PublicKey),
		}
	}
	return opmap, nil
}
