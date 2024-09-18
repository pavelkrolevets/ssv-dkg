package integration_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/bloxapp/ssv/logging"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	eth_crypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"

	"github.com/ssvlabs/dkg-spec/testing/stubs"
	cli_initiator "github.com/ssvlabs/ssv-dkg/cli/initiator"
)

func TestLoadProofsFromFileResign(t *testing.T) {
	err := os.RemoveAll("./output/")
	require.NoError(t, err)
	err = logging.SetGlobalLogger("info", "capital", "console", nil)
	require.NoError(t, err)
	version := "test.version"

	// Open ethereum keystore
	jsonBytes, err := os.ReadFile("../examples/initiator/UTC--2024-06-14T14-05-12.366668334Z--dcc846fa10c7cfce9e6eb37e06ed93b666cfc5e9")
	require.NoError(t, err)
	keyStorePassword, err := os.ReadFile(filepath.Clean("../examples/initiator/password"))
	require.NoError(t, err)
	sk, err := keystore.DecryptKey(jsonBytes, string(keyStorePassword))
	require.NoError(t, err)
	owner := eth_crypto.PubkeyToAddress(sk.PrivateKey.PublicKey)

	stubClient := &stubs.Client{
		CallContractF: func(call ethereum.CallMsg) ([]byte, error) {
			return nil, nil
		},
	}
	servers, ops := createOperatorsFromExamplesFolder(t, version, stubClient)
	operators, err := json.Marshal(ops)
	require.NoError(t, err)
	RootCmd := &cobra.Command{
		Use:   "ssv-dkg",
		Short: "CLI for running Distributed Key Generation protocol",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
		},
	}
	RootCmd.AddCommand(cli_initiator.StartResigning)
	RootCmd.Short = "ssv-dkg-test"
	RootCmd.Version = version
	cli_initiator.StartResigning.Version = version
	args := []string{"resign",
		"--proofsFilePath", "./stubs/4/000001-0xaa57eab07f1a740672d0c106867d366c798d3b932d373c88cf047da1a3c16d0816ac58bab5a9d6f6f4b63a07608f8f39/proofs.json",
		"--operatorsInfo", string(operators),
		"--owner", owner.Hex(),
		"--withdrawAddress", "0x81592c3de184a3e2c0dcb5a261bc107bfa91f494",
		"--operatorIDs", "11,22,33,44",
		"--nonce", "10",
		"--ethKeystorePath", "./stubs/UTC--2024-06-14T14-05-12.366668334Z--dcc846fa10c7cfce9e6eb37e06ed93b666cfc5e9",
		"--ethKeystorePass", "./stubs/password"}
	RootCmd.SetArgs(args)
	err = RootCmd.Execute()
	require.NoError(t, err)
	err = os.RemoveAll("./output/")
	require.NoError(t, err)
	for _, srv := range servers {
		srv.HttpSrv.Close()
	}
}

func TestLoadProofsFromRawJSONResign(t *testing.T) {
	err := os.RemoveAll("./output/")
	require.NoError(t, err)
	proofsRawJSON := `[{"proof":{"validator":"aa57eab07f1a740672d0c106867d366c798d3b932d373c88cf047da1a3c16d0816ac58bab5a9d6f6f4b63a07608f8f39","encrypted_share":"959b0501cc8792851cbbb6fadbc1d983e9a91c8bbc3e92fd05498b817c7d4061c22a8a61e67d097083f70176a2f62e1d4c34428a5a7d5e1bd4d4164d178de8b45606d5953a915351dace666314668d44fba95cee65bf9b9d3e611dffc558f0769332f63cdf101a2a487d05885e52402529e11849c6f2f60d305111addd37746d3e47892e636d6d8ce2d61b45750cb625dcf2faeedf673b104bdc2c7dc862a723354f2ca2867517853278952f917e6e74846938a8567588dca731590f585037968658c6a998c90d1e6595ceb39a52410641531f0a14137faf5d89894e9912d91b57652a7103c96c7b8747a84628ffb901c2d2f967ee6130d7851544477b24835d","share_pub":"8ce17b6af1fc18aa66298ef3f722cfa3e891d7cb5bdf116051a0af5aeb6f6c38a7860cd3d47427fb584240706c3a4acb","owner":"dcc846fa10c7cfce9e6eb37e06ed93b666cfc5e9"},"signature":"25defc7041c2bd425a2c405298d7a153f653bb6f49068938dc04ed81bd03ab7330c049dec3ff8389f7ae8b96cac5559ac39631e343f027cb41826af1b8cdee0fa9deb0b6210de725084396046904a7b2e9816decb2b33c88729ba1c814e59d04749c0a133e6190df830a919292cb8bd33535657c88936dfa05e83894517132c1fa090efe333b569dcbb7959477806bcdce797835e3a8e990f5a1ba091e9a6c5c4ba0807546882befd63273ab4a2857c0153b5adfe6342088d65d3192e782e30b5441f9aaf63d3ee6a0603bdc0be489202e347302aef61873593236102666e9a9bcbdb3b8322d18ffcc4a8bb41754207ec243a18099c6271ec8593f913d60bfac"},{"proof":{"validator":"aa57eab07f1a740672d0c106867d366c798d3b932d373c88cf047da1a3c16d0816ac58bab5a9d6f6f4b63a07608f8f39","encrypted_share":"2723a6b922ece432b6ad497ba50969534a5c3f9d29cc05739d951efc986f8ccf25d6d2d6331aeedb3a7f30c3f29760aeba44f9c42bc83f91ce3148ccca64e46405e627d02a4215d5bf1e166eb540838c8ac9cb79551e04a177ff765b4663ade090276606fba5c1cb91b731c8ef25d4e541dca541061eb1c1a5772f38445a65fbe40e6b2f072892e0a9982aebe70f27ce5ca241ffc7c0013f99d04ef586bb4ff04bbfbd1f3eaa6eddab998e355d3e160effd078baf364d187b59f1d74e75f891f2708f65326284f692dea5f82fa736f0b2b625dbf3bbda31047ed6d7b37f9138d9cc63d502282d926d59a9c825e58b15cd271497778d5f3d102e7012d1b8eac9e","share_pub":"8f476e4cb557aede4137696d63b2910a88bc55123e9a8a6b03cbc338f8fec3bf8bfffaf6dea5607b76c7d00915b4dd25","owner":"dcc846fa10c7cfce9e6eb37e06ed93b666cfc5e9"},"signature":"2dc3a34a648277a3e94d972d37b54b2a43c4754b8e934640064ab3274b2c43ea75998ba0b95457c7f8e518187779cf469c3a3850388921d00a2e1de31e98ccc8e22a5ece097bb904addc439c4f4de5dbe0f741fb554347807a853e94ae0c9271a22eaa9f909d97ffc9b16a84a1f6194fbf3aa032479e4fd1538a49cfd9d9a93d5269a8705a92bc0f92aa354012289eafc98e2774cf9a39e556e29b33bc9e0c18315bdb83a2b5310a5e8c556fa2888e89577c813e86e2fe90d74a9c67db2abde81e68cfa052e5e1519c62600f0f6044f1be48ebd665e42b480e74b54d19b91b3e1284e70fa6d0d268dabf0b4240eaae49d0928711bd7fe8e709375a68624c1c14"},{"proof":{"validator":"aa57eab07f1a740672d0c106867d366c798d3b932d373c88cf047da1a3c16d0816ac58bab5a9d6f6f4b63a07608f8f39","encrypted_share":"02a0a1e9cccf8aaaa2035fdc848fb4ca298a76a3462394391c500c68ef963e33ef91bacd2caa6e654ab05f50d49c38a8e0019ec41b15e63399defeaf6616461ddf646fd59d0830b54f99365446fcd8fad0059d7c4c5ea500e7d3e2d7cf739f5378a24b66ee2271878876b0db263db6c0391589deffba9e4d56f56ff23b9ee3dcafdfbcc1c3b0a8f365482393981b8b3d231840c32ee2be9fdceb4915544deec6fcd1803e75b0159883cd7eeb87bee17885724037b0fbfd93878b7b17c32260a41e491962b67fb405f0c39b47f9f93785f1852cb8fa75d133a025e64092b4ed18059a60abb24e9372a3ca0e08838abb9b42524a74ae3eab5410ee2314b6ab118e","share_pub":"a99bf34e0aef11749080ed4cf36acf1711460f5bdbd098b28de1afac618cc273f457613ca454262d8f5e46a8d010dd00","owner":"dcc846fa10c7cfce9e6eb37e06ed93b666cfc5e9"},"signature":"15050d609d2e256fa934578e57c8e97df895fb2e8a887317800c9797460bae1356c274ebeea1b2e96daedac9c24ecc4629acf0cff618fac8df6b4f90851e7e176b636d07db4eb204c02a752a69418e640f1b08d7e04310c95a91529421b4359b31f29ea6db08c3d2e0c7f0ca2f3f8f159812dcd48b613a4e706e02ce9164efcdf385fcbedbd465b5448b677c66bcfe433500e05ed97514c031c20e1ee54a1ee7b154b2aa7c0bdef97ad9db3d90e2dbe7b6ec0fda5b780ac6d4ac78f7e002d845d244e8ed43f6bda63280306e0c014b232b760dd1b16969fb9b0b5a139e3a41742d0870ec2f03c3606d2d36b782740fb4b2edad79082254bbda275060c64591d6"},{"proof":{"validator":"aa57eab07f1a740672d0c106867d366c798d3b932d373c88cf047da1a3c16d0816ac58bab5a9d6f6f4b63a07608f8f39","encrypted_share":"1078b9a1f64cd8efbdd44cd04ba262f78f79d7e1ba99b2c87a580492410314ae26cda143ffc09462335250f03e478fcd90c3a213f127d905e4ceef67b476d9ef649fc6e0f2603b9e1d8dac5801f962882797dfba972f96ca664ba51269b1a3534560a1313586a8f42c97431581a12bddaf6460c27248e31b5746d117da31f073f0862ea5a893852d0fe2174c73c6313dd3596f514e7e5758bf6f5195e80310f9c1f92374784d98c3503b1806b24e7a6f4e1bcf012e7d5f582d66899487247c3666ad131e5dcdd43812e4c4aa3e9f5e0f0120425e018b800e51b6207a3f4d4f8969c177603a3e360c95b18375d9c5b1405bf39a90bc24aa111bba009f772360e3","share_pub":"84faff90116da6fc1222fc2060ea24cae6f135c5d2684b09e5f9b0253d682df0a96025e81a606aafb2d434d774ec9340","owner":"dcc846fa10c7cfce9e6eb37e06ed93b666cfc5e9"},"signature":"7fddfeb183641468c6799448c6462e687a3278b6336aa48ec658340ef41a68910495053dc2fb0033cf42f7a9c035a0deb8b01eb184b4268fc3868dd4c2188727c1340dac988bf7c2917bc5b970dcb6d675446308e9f1d10d4079324a7ded67bf5efab995cd6168132d8b8e73d29062c70808b73b7c75e6efb292871c6710d4baee4dfc6c7043f011f5a38175fbbb062f2b49696fdb522c9a312eb319ee723619f86c0d253a4ae59c7968b3bd879538af052b49ab18aff88ea76ae60da29c9f2e327a76a1e35b357c45afa3c0435ece28f407946b26c44a1839dffb5c061689425fc6ccaeac9af5cfe72061ac02e13824f8d136aeeea6a292eb9738f83a57a1a1"}]`
	err = logging.SetGlobalLogger("info", "capital", "console", nil)
	require.NoError(t, err)
	version := "test.version"
	// Open ethereum keystore
	jsonBytes, err := os.ReadFile("../examples/initiator/UTC--2024-06-14T14-05-12.366668334Z--dcc846fa10c7cfce9e6eb37e06ed93b666cfc5e9")
	require.NoError(t, err)
	keyStorePassword, err := os.ReadFile(filepath.Clean("../examples/initiator/password"))
	require.NoError(t, err)
	sk, err := keystore.DecryptKey(jsonBytes, string(keyStorePassword))
	require.NoError(t, err)
	owner := eth_crypto.PubkeyToAddress(sk.PrivateKey.PublicKey)

	stubClient := &stubs.Client{
		CallContractF: func(call ethereum.CallMsg) ([]byte, error) {
			return nil, nil
		},
	}
	servers, ops := createOperatorsFromExamplesFolder(t, version, stubClient)
	operators, err := json.Marshal(ops)
	require.NoError(t, err)
	RootCmd := &cobra.Command{
		Use:   "ssv-dkg",
		Short: "CLI for running Distributed Key Generation protocol",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
		},
	}
	RootCmd.AddCommand(cli_initiator.StartResigning)
	RootCmd.Short = "ssv-dkg-test"
	RootCmd.Version = version
	cli_initiator.StartResigning.Version = version
	args := []string{"resign",
		"--proofsRawJSON", proofsRawJSON,
		"--operatorsInfo", string(operators),
		"--owner", owner.Hex(),
		"--withdrawAddress", "0x81592c3de184a3e2c0dcb5a261bc107bfa91f494",
		"--operatorIDs", "11,22,33,44",
		"--nonce", "10",
		"--ethKeystorePath", "./stubs/UTC--2024-06-14T14-05-12.366668334Z--dcc846fa10c7cfce9e6eb37e06ed93b666cfc5e9",
		"--ethKeystorePass", "./stubs/password"}
	RootCmd.SetArgs(args)
	err = RootCmd.Execute()
	require.NoError(t, err)

	for _, srv := range servers {
		srv.HttpSrv.Close()
	}
}

func TestLoadProofsFromFileReshare(t *testing.T) {
	err := os.RemoveAll("./output/")
	require.NoError(t, err)
	err = logging.SetGlobalLogger("info", "capital", "console", nil)
	require.NoError(t, err)
	version := "test.version"

	// Open ethereum keystore
	jsonBytes, err := os.ReadFile("../examples/initiator/UTC--2024-06-14T14-05-12.366668334Z--dcc846fa10c7cfce9e6eb37e06ed93b666cfc5e9")
	require.NoError(t, err)
	keyStorePassword, err := os.ReadFile(filepath.Clean("../examples/initiator/password"))
	require.NoError(t, err)
	sk, err := keystore.DecryptKey(jsonBytes, string(keyStorePassword))
	require.NoError(t, err)
	owner := eth_crypto.PubkeyToAddress(sk.PrivateKey.PublicKey)

	stubClient := &stubs.Client{
		CallContractF: func(call ethereum.CallMsg) ([]byte, error) {
			return nil, nil
		},
	}
	servers, ops := createOperatorsFromExamplesFolder(t, version, stubClient)
	operators, err := json.Marshal(ops)
	require.NoError(t, err)
	RootCmd := &cobra.Command{
		Use:   "ssv-dkg",
		Short: "CLI for running Distributed Key Generation protocol",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
		},
	}
	RootCmd.AddCommand(cli_initiator.StartReshare)
	RootCmd.Short = "ssv-dkg-test"
	RootCmd.Version = version
	cli_initiator.StartReshare.Version = version
	args := []string{"reshare",
		"--proofsFilePath", "./stubs/4/000001-0xaa57eab07f1a740672d0c106867d366c798d3b932d373c88cf047da1a3c16d0816ac58bab5a9d6f6f4b63a07608f8f39/proofs.json",
		"--operatorsInfo", string(operators),
		"--owner", owner.Hex(),
		"--withdrawAddress", "0x81592c3de184a3e2c0dcb5a261bc107bfa91f494",
		"--operatorIDs", "11,22,33,44",
		"--newOperatorIDs", "55,66,77,88",
		"--nonce", strconv.Itoa(10),
		"--ethKeystorePath", "./stubs/UTC--2024-06-14T14-05-12.366668334Z--dcc846fa10c7cfce9e6eb37e06ed93b666cfc5e9",
		"--ethKeystorePass", "./stubs/password",
		"--network", "holesky"}
	RootCmd.SetArgs(args)
	err = RootCmd.Execute()
	require.NoError(t, err)
	err = os.RemoveAll("./output/")
	require.NoError(t, err)
	for _, srv := range servers {
		srv.HttpSrv.Close()
	}
}

func TestLoadProofsFromRawJSONReshare(t *testing.T) {
	err := os.RemoveAll("./output/")
	require.NoError(t, err)
	proofsRawJSON := `[{"proof":{"validator":"aa57eab07f1a740672d0c106867d366c798d3b932d373c88cf047da1a3c16d0816ac58bab5a9d6f6f4b63a07608f8f39","encrypted_share":"959b0501cc8792851cbbb6fadbc1d983e9a91c8bbc3e92fd05498b817c7d4061c22a8a61e67d097083f70176a2f62e1d4c34428a5a7d5e1bd4d4164d178de8b45606d5953a915351dace666314668d44fba95cee65bf9b9d3e611dffc558f0769332f63cdf101a2a487d05885e52402529e11849c6f2f60d305111addd37746d3e47892e636d6d8ce2d61b45750cb625dcf2faeedf673b104bdc2c7dc862a723354f2ca2867517853278952f917e6e74846938a8567588dca731590f585037968658c6a998c90d1e6595ceb39a52410641531f0a14137faf5d89894e9912d91b57652a7103c96c7b8747a84628ffb901c2d2f967ee6130d7851544477b24835d","share_pub":"8ce17b6af1fc18aa66298ef3f722cfa3e891d7cb5bdf116051a0af5aeb6f6c38a7860cd3d47427fb584240706c3a4acb","owner":"dcc846fa10c7cfce9e6eb37e06ed93b666cfc5e9"},"signature":"25defc7041c2bd425a2c405298d7a153f653bb6f49068938dc04ed81bd03ab7330c049dec3ff8389f7ae8b96cac5559ac39631e343f027cb41826af1b8cdee0fa9deb0b6210de725084396046904a7b2e9816decb2b33c88729ba1c814e59d04749c0a133e6190df830a919292cb8bd33535657c88936dfa05e83894517132c1fa090efe333b569dcbb7959477806bcdce797835e3a8e990f5a1ba091e9a6c5c4ba0807546882befd63273ab4a2857c0153b5adfe6342088d65d3192e782e30b5441f9aaf63d3ee6a0603bdc0be489202e347302aef61873593236102666e9a9bcbdb3b8322d18ffcc4a8bb41754207ec243a18099c6271ec8593f913d60bfac"},{"proof":{"validator":"aa57eab07f1a740672d0c106867d366c798d3b932d373c88cf047da1a3c16d0816ac58bab5a9d6f6f4b63a07608f8f39","encrypted_share":"2723a6b922ece432b6ad497ba50969534a5c3f9d29cc05739d951efc986f8ccf25d6d2d6331aeedb3a7f30c3f29760aeba44f9c42bc83f91ce3148ccca64e46405e627d02a4215d5bf1e166eb540838c8ac9cb79551e04a177ff765b4663ade090276606fba5c1cb91b731c8ef25d4e541dca541061eb1c1a5772f38445a65fbe40e6b2f072892e0a9982aebe70f27ce5ca241ffc7c0013f99d04ef586bb4ff04bbfbd1f3eaa6eddab998e355d3e160effd078baf364d187b59f1d74e75f891f2708f65326284f692dea5f82fa736f0b2b625dbf3bbda31047ed6d7b37f9138d9cc63d502282d926d59a9c825e58b15cd271497778d5f3d102e7012d1b8eac9e","share_pub":"8f476e4cb557aede4137696d63b2910a88bc55123e9a8a6b03cbc338f8fec3bf8bfffaf6dea5607b76c7d00915b4dd25","owner":"dcc846fa10c7cfce9e6eb37e06ed93b666cfc5e9"},"signature":"2dc3a34a648277a3e94d972d37b54b2a43c4754b8e934640064ab3274b2c43ea75998ba0b95457c7f8e518187779cf469c3a3850388921d00a2e1de31e98ccc8e22a5ece097bb904addc439c4f4de5dbe0f741fb554347807a853e94ae0c9271a22eaa9f909d97ffc9b16a84a1f6194fbf3aa032479e4fd1538a49cfd9d9a93d5269a8705a92bc0f92aa354012289eafc98e2774cf9a39e556e29b33bc9e0c18315bdb83a2b5310a5e8c556fa2888e89577c813e86e2fe90d74a9c67db2abde81e68cfa052e5e1519c62600f0f6044f1be48ebd665e42b480e74b54d19b91b3e1284e70fa6d0d268dabf0b4240eaae49d0928711bd7fe8e709375a68624c1c14"},{"proof":{"validator":"aa57eab07f1a740672d0c106867d366c798d3b932d373c88cf047da1a3c16d0816ac58bab5a9d6f6f4b63a07608f8f39","encrypted_share":"02a0a1e9cccf8aaaa2035fdc848fb4ca298a76a3462394391c500c68ef963e33ef91bacd2caa6e654ab05f50d49c38a8e0019ec41b15e63399defeaf6616461ddf646fd59d0830b54f99365446fcd8fad0059d7c4c5ea500e7d3e2d7cf739f5378a24b66ee2271878876b0db263db6c0391589deffba9e4d56f56ff23b9ee3dcafdfbcc1c3b0a8f365482393981b8b3d231840c32ee2be9fdceb4915544deec6fcd1803e75b0159883cd7eeb87bee17885724037b0fbfd93878b7b17c32260a41e491962b67fb405f0c39b47f9f93785f1852cb8fa75d133a025e64092b4ed18059a60abb24e9372a3ca0e08838abb9b42524a74ae3eab5410ee2314b6ab118e","share_pub":"a99bf34e0aef11749080ed4cf36acf1711460f5bdbd098b28de1afac618cc273f457613ca454262d8f5e46a8d010dd00","owner":"dcc846fa10c7cfce9e6eb37e06ed93b666cfc5e9"},"signature":"15050d609d2e256fa934578e57c8e97df895fb2e8a887317800c9797460bae1356c274ebeea1b2e96daedac9c24ecc4629acf0cff618fac8df6b4f90851e7e176b636d07db4eb204c02a752a69418e640f1b08d7e04310c95a91529421b4359b31f29ea6db08c3d2e0c7f0ca2f3f8f159812dcd48b613a4e706e02ce9164efcdf385fcbedbd465b5448b677c66bcfe433500e05ed97514c031c20e1ee54a1ee7b154b2aa7c0bdef97ad9db3d90e2dbe7b6ec0fda5b780ac6d4ac78f7e002d845d244e8ed43f6bda63280306e0c014b232b760dd1b16969fb9b0b5a139e3a41742d0870ec2f03c3606d2d36b782740fb4b2edad79082254bbda275060c64591d6"},{"proof":{"validator":"aa57eab07f1a740672d0c106867d366c798d3b932d373c88cf047da1a3c16d0816ac58bab5a9d6f6f4b63a07608f8f39","encrypted_share":"1078b9a1f64cd8efbdd44cd04ba262f78f79d7e1ba99b2c87a580492410314ae26cda143ffc09462335250f03e478fcd90c3a213f127d905e4ceef67b476d9ef649fc6e0f2603b9e1d8dac5801f962882797dfba972f96ca664ba51269b1a3534560a1313586a8f42c97431581a12bddaf6460c27248e31b5746d117da31f073f0862ea5a893852d0fe2174c73c6313dd3596f514e7e5758bf6f5195e80310f9c1f92374784d98c3503b1806b24e7a6f4e1bcf012e7d5f582d66899487247c3666ad131e5dcdd43812e4c4aa3e9f5e0f0120425e018b800e51b6207a3f4d4f8969c177603a3e360c95b18375d9c5b1405bf39a90bc24aa111bba009f772360e3","share_pub":"84faff90116da6fc1222fc2060ea24cae6f135c5d2684b09e5f9b0253d682df0a96025e81a606aafb2d434d774ec9340","owner":"dcc846fa10c7cfce9e6eb37e06ed93b666cfc5e9"},"signature":"7fddfeb183641468c6799448c6462e687a3278b6336aa48ec658340ef41a68910495053dc2fb0033cf42f7a9c035a0deb8b01eb184b4268fc3868dd4c2188727c1340dac988bf7c2917bc5b970dcb6d675446308e9f1d10d4079324a7ded67bf5efab995cd6168132d8b8e73d29062c70808b73b7c75e6efb292871c6710d4baee4dfc6c7043f011f5a38175fbbb062f2b49696fdb522c9a312eb319ee723619f86c0d253a4ae59c7968b3bd879538af052b49ab18aff88ea76ae60da29c9f2e327a76a1e35b357c45afa3c0435ece28f407946b26c44a1839dffb5c061689425fc6ccaeac9af5cfe72061ac02e13824f8d136aeeea6a292eb9738f83a57a1a1"}]`
	err = logging.SetGlobalLogger("info", "capital", "console", nil)
	require.NoError(t, err)
	version := "test.version"
	// Open ethereum keystore
	jsonBytes, err := os.ReadFile("../examples/initiator/UTC--2024-06-14T14-05-12.366668334Z--dcc846fa10c7cfce9e6eb37e06ed93b666cfc5e9")
	require.NoError(t, err)
	keyStorePassword, err := os.ReadFile(filepath.Clean("../examples/initiator/password"))
	require.NoError(t, err)
	sk, err := keystore.DecryptKey(jsonBytes, string(keyStorePassword))
	require.NoError(t, err)
	owner := eth_crypto.PubkeyToAddress(sk.PrivateKey.PublicKey)
	stubClient := &stubs.Client{
		CallContractF: func(call ethereum.CallMsg) ([]byte, error) {
			return nil, nil
		},
	}
	servers, ops := createOperatorsFromExamplesFolder(t, version, stubClient)
	operators, err := json.Marshal(ops)
	require.NoError(t, err)
	RootCmd := &cobra.Command{
		Use:   "ssv-dkg",
		Short: "CLI for running Distributed Key Generation protocol",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
		},
	}
	RootCmd.AddCommand(cli_initiator.StartReshare)
	RootCmd.Short = "ssv-dkg-test"
	RootCmd.Version = version
	cli_initiator.StartReshare.Version = version
	args := []string{"reshare",
		"--proofsRawJSON", proofsRawJSON,
		"--operatorsInfo", string(operators),
		"--owner", owner.Hex(),
		"--withdrawAddress", "0x81592c3de184a3e2c0dcb5a261bc107bfa91f494",
		"--operatorIDs", "11,22,33,44",
		"--newOperatorIDs", "55,66,77,88",
		"--nonce", strconv.Itoa(10),
		"--ethKeystorePath", "./stubs/UTC--2024-06-14T14-05-12.366668334Z--dcc846fa10c7cfce9e6eb37e06ed93b666cfc5e9",
		"--ethKeystorePass", "./stubs/password",
		"--network", "holesky"}
	RootCmd.SetArgs(args)
	err = RootCmd.Execute()
	require.NoError(t, err)

	for _, srv := range servers {
		srv.HttpSrv.Close()
	}
}

func TestLoadProofsErrorResign(t *testing.T) {
	err := os.RemoveAll("./output/")
	require.NoError(t, err)
	err = logging.SetGlobalLogger("info", "capital", "console", nil)
	require.NoError(t, err)
	version := "test.version"
	// Open ethereum keystore
	jsonBytes, err := os.ReadFile("../examples/initiator/UTC--2024-06-14T14-05-12.366668334Z--dcc846fa10c7cfce9e6eb37e06ed93b666cfc5e9")
	require.NoError(t, err)
	keyStorePassword, err := os.ReadFile(filepath.Clean("../examples/initiator/password"))
	require.NoError(t, err)
	sk, err := keystore.DecryptKey(jsonBytes, string(keyStorePassword))
	require.NoError(t, err)
	owner := eth_crypto.PubkeyToAddress(sk.PrivateKey.PublicKey)
	stubClient := &stubs.Client{
		CallContractF: func(call ethereum.CallMsg) ([]byte, error) {
			return nil, nil
		},
	}
	servers, ops := createOperatorsFromExamplesFolder(t, version, stubClient)
	operators, err := json.Marshal(ops)
	require.NoError(t, err)
	RootCmd := &cobra.Command{
		Use:   "ssv-dkg",
		Short: "CLI for running Distributed Key Generation protocol",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
		},
	}
	RootCmd.AddCommand(cli_initiator.StartResigning)
	RootCmd.Short = "ssv-dkg-test"
	RootCmd.Version = version
	cli_initiator.StartResigning.Version = version
	proofsRawJSON := `[{"proof":{"validator":"aa57eab07f1a740672d0c106867d366c798d3b932d373c88cf047da1a3c16d0816ac58bab5a9d6f6f4b63a07608f8f39","encrypted_share":"959b0501cc8792851cbbb6fadbc1d983e9a91c8bbc3e92fd05498b817c7d4061c22a8a61e67d097083f70176a2f62e1d4c34428a5a7d5e1bd4d4164d178de8b45606d5953a915351dace666314668d44fba95cee65bf9b9d3e611dffc558f0769332f63cdf101a2a487d05885e52402529e11849c6f2f60d305111addd37746d3e47892e636d6d8ce2d61b45750cb625dcf2faeedf673b104bdc2c7dc862a723354f2ca2867517853278952f917e6e74846938a8567588dca731590f585037968658c6a998c90d1e6595ceb39a52410641531f0a14137faf5d89894e9912d91b57652a7103c96c7b8747a84628ffb901c2d2f967ee6130d7851544477b24835d","share_pub":"8ce17b6af1fc18aa66298ef3f722cfa3e891d7cb5bdf116051a0af5aeb6f6c38a7860cd3d47427fb584240706c3a4acb","owner":"dcc846fa10c7cfce9e6eb37e06ed93b666cfc5e9"},"signature":"25defc7041c2bd425a2c405298d7a153f653bb6f49068938dc04ed81bd03ab7330c049dec3ff8389f7ae8b96cac5559ac39631e343f027cb41826af1b8cdee0fa9deb0b6210de725084396046904a7b2e9816decb2b33c88729ba1c814e59d04749c0a133e6190df830a919292cb8bd33535657c88936dfa05e83894517132c1fa090efe333b569dcbb7959477806bcdce797835e3a8e990f5a1ba091e9a6c5c4ba0807546882befd63273ab4a2857c0153b5adfe6342088d65d3192e782e30b5441f9aaf63d3ee6a0603bdc0be489202e347302aef61873593236102666e9a9bcbdb3b8322d18ffcc4a8bb41754207ec243a18099c6271ec8593f913d60bfac"},{"proof":{"validator":"aa57eab07f1a740672d0c106867d366c798d3b932d373c88cf047da1a3c16d0816ac58bab5a9d6f6f4b63a07608f8f39","encrypted_share":"2723a6b922ece432b6ad497ba50969534a5c3f9d29cc05739d951efc986f8ccf25d6d2d6331aeedb3a7f30c3f29760aeba44f9c42bc83f91ce3148ccca64e46405e627d02a4215d5bf1e166eb540838c8ac9cb79551e04a177ff765b4663ade090276606fba5c1cb91b731c8ef25d4e541dca541061eb1c1a5772f38445a65fbe40e6b2f072892e0a9982aebe70f27ce5ca241ffc7c0013f99d04ef586bb4ff04bbfbd1f3eaa6eddab998e355d3e160effd078baf364d187b59f1d74e75f891f2708f65326284f692dea5f82fa736f0b2b625dbf3bbda31047ed6d7b37f9138d9cc63d502282d926d59a9c825e58b15cd271497778d5f3d102e7012d1b8eac9e","share_pub":"8f476e4cb557aede4137696d63b2910a88bc55123e9a8a6b03cbc338f8fec3bf8bfffaf6dea5607b76c7d00915b4dd25","owner":"dcc846fa10c7cfce9e6eb37e06ed93b666cfc5e9"},"signature":"2dc3a34a648277a3e94d972d37b54b2a43c4754b8e934640064ab3274b2c43ea75998ba0b95457c7f8e518187779cf469c3a3850388921d00a2e1de31e98ccc8e22a5ece097bb904addc439c4f4de5dbe0f741fb554347807a853e94ae0c9271a22eaa9f909d97ffc9b16a84a1f6194fbf3aa032479e4fd1538a49cfd9d9a93d5269a8705a92bc0f92aa354012289eafc98e2774cf9a39e556e29b33bc9e0c18315bdb83a2b5310a5e8c556fa2888e89577c813e86e2fe90d74a9c67db2abde81e68cfa052e5e1519c62600f0f6044f1be48ebd665e42b480e74b54d19b91b3e1284e70fa6d0d268dabf0b4240eaae49d0928711bd7fe8e709375a68624c1c14"},{"proof":{"validator":"aa57eab07f1a740672d0c106867d366c798d3b932d373c88cf047da1a3c16d0816ac58bab5a9d6f6f4b63a07608f8f39","encrypted_share":"02a0a1e9cccf8aaaa2035fdc848fb4ca298a76a3462394391c500c68ef963e33ef91bacd2caa6e654ab05f50d49c38a8e0019ec41b15e63399defeaf6616461ddf646fd59d0830b54f99365446fcd8fad0059d7c4c5ea500e7d3e2d7cf739f5378a24b66ee2271878876b0db263db6c0391589deffba9e4d56f56ff23b9ee3dcafdfbcc1c3b0a8f365482393981b8b3d231840c32ee2be9fdceb4915544deec6fcd1803e75b0159883cd7eeb87bee17885724037b0fbfd93878b7b17c32260a41e491962b67fb405f0c39b47f9f93785f1852cb8fa75d133a025e64092b4ed18059a60abb24e9372a3ca0e08838abb9b42524a74ae3eab5410ee2314b6ab118e","share_pub":"a99bf34e0aef11749080ed4cf36acf1711460f5bdbd098b28de1afac618cc273f457613ca454262d8f5e46a8d010dd00","owner":"dcc846fa10c7cfce9e6eb37e06ed93b666cfc5e9"},"signature":"15050d609d2e256fa934578e57c8e97df895fb2e8a887317800c9797460bae1356c274ebeea1b2e96daedac9c24ecc4629acf0cff618fac8df6b4f90851e7e176b636d07db4eb204c02a752a69418e640f1b08d7e04310c95a91529421b4359b31f29ea6db08c3d2e0c7f0ca2f3f8f159812dcd48b613a4e706e02ce9164efcdf385fcbedbd465b5448b677c66bcfe433500e05ed97514c031c20e1ee54a1ee7b154b2aa7c0bdef97ad9db3d90e2dbe7b6ec0fda5b780ac6d4ac78f7e002d845d244e8ed43f6bda63280306e0c014b232b760dd1b16969fb9b0b5a139e3a41742d0870ec2f03c3606d2d36b782740fb4b2edad79082254bbda275060c64591d6"},{"proof":{"validator":"aa57eab07f1a740672d0c106867d366c798d3b932d373c88cf047da1a3c16d0816ac58bab5a9d6f6f4b63a07608f8f39","encrypted_share":"1078b9a1f64cd8efbdd44cd04ba262f78f79d7e1ba99b2c87a580492410314ae26cda143ffc09462335250f03e478fcd90c3a213f127d905e4ceef67b476d9ef649fc6e0f2603b9e1d8dac5801f962882797dfba972f96ca664ba51269b1a3534560a1313586a8f42c97431581a12bddaf6460c27248e31b5746d117da31f073f0862ea5a893852d0fe2174c73c6313dd3596f514e7e5758bf6f5195e80310f9c1f92374784d98c3503b1806b24e7a6f4e1bcf012e7d5f582d66899487247c3666ad131e5dcdd43812e4c4aa3e9f5e0f0120425e018b800e51b6207a3f4d4f8969c177603a3e360c95b18375d9c5b1405bf39a90bc24aa111bba009f772360e3","share_pub":"84faff90116da6fc1222fc2060ea24cae6f135c5d2684b09e5f9b0253d682df0a96025e81a606aafb2d434d774ec9340","owner":"dcc846fa10c7cfce9e6eb37e06ed93b666cfc5e9"},"signature":"7fddfeb183641468c6799448c6462e687a3278b6336aa48ec658340ef41a68910495053dc2fb0033cf42f7a9c035a0deb8b01eb184b4268fc3868dd4c2188727c1340dac988bf7c2917bc5b970dcb6d675446308e9f1d10d4079324a7ded67bf5efab995cd6168132d8b8e73d29062c70808b73b7c75e6efb292871c6710d4baee4dfc6c7043f011f5a38175fbbb062f2b49696fdb522c9a312eb319ee723619f86c0d253a4ae59c7968b3bd879538af052b49ab18aff88ea76ae60da29c9f2e327a76a1e35b357c45afa3c0435ece28f407946b26c44a1839dffb5c061689425fc6ccaeac9af5cfe72061ac02e13824f8d136aeeea6a292eb9738f83a57a1a1"}]`
	args := []string{"resign",
		"--proofsRawJSON", proofsRawJSON,
		"--proofsFilePath", "./stubs/4/000001-0xaa57eab07f1a740672d0c106867d366c798d3b932d373c88cf047da1a3c16d0816ac58bab5a9d6f6f4b63a07608f8f39/proofs.json",
		"--operatorsInfo", string(operators),
		"--owner", owner.Hex(),
		"--withdrawAddress", "0x81592c3de184a3e2c0dcb5a261bc107bfa91f494",
		"--operatorIDs", "11,22,33,44",
		"--nonce", "10",
		"--ethKeystorePath", "./stubs/UTC--2024-06-14T14-05-12.366668334Z--dcc846fa10c7cfce9e6eb37e06ed93b666cfc5e9",
		"--ethKeystorePass", "./stubs/password"}
	RootCmd.SetArgs(args)
	err = RootCmd.Execute()
	require.ErrorContains(t, err, "please provide either proofsRaw flag or proofsFilePath, not both")
	for _, srv := range servers {
		srv.HttpSrv.Close()
	}
}

func TestLoadProofsErrorReshare(t *testing.T) {
	err := os.RemoveAll("./output/")
	require.NoError(t, err)
	err = logging.SetGlobalLogger("info", "capital", "console", nil)
	require.NoError(t, err)
	version := "test.version"
	// Open ethereum keystore
	jsonBytes, err := os.ReadFile("../examples/initiator/UTC--2024-06-14T14-05-12.366668334Z--dcc846fa10c7cfce9e6eb37e06ed93b666cfc5e9")
	require.NoError(t, err)
	keyStorePassword, err := os.ReadFile(filepath.Clean("../examples/initiator/password"))
	require.NoError(t, err)
	sk, err := keystore.DecryptKey(jsonBytes, string(keyStorePassword))
	require.NoError(t, err)
	owner := eth_crypto.PubkeyToAddress(sk.PrivateKey.PublicKey)

	stubClient := &stubs.Client{
		CallContractF: func(call ethereum.CallMsg) ([]byte, error) {
			return nil, nil
		},
	}
	servers, ops := createOperatorsFromExamplesFolder(t, version, stubClient)
	operators, err := json.Marshal(ops)
	require.NoError(t, err)
	RootCmd := &cobra.Command{
		Use:   "ssv-dkg",
		Short: "CLI for running Distributed Key Generation protocol",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
		},
	}
	RootCmd.AddCommand(cli_initiator.StartReshare)
	RootCmd.Short = "ssv-dkg-test"
	RootCmd.Version = version
	cli_initiator.StartReshare.Version = version
	proofsRawJSON := `[{"proof":{"validator":"aa57eab07f1a740672d0c106867d366c798d3b932d373c88cf047da1a3c16d0816ac58bab5a9d6f6f4b63a07608f8f39","encrypted_share":"959b0501cc8792851cbbb6fadbc1d983e9a91c8bbc3e92fd05498b817c7d4061c22a8a61e67d097083f70176a2f62e1d4c34428a5a7d5e1bd4d4164d178de8b45606d5953a915351dace666314668d44fba95cee65bf9b9d3e611dffc558f0769332f63cdf101a2a487d05885e52402529e11849c6f2f60d305111addd37746d3e47892e636d6d8ce2d61b45750cb625dcf2faeedf673b104bdc2c7dc862a723354f2ca2867517853278952f917e6e74846938a8567588dca731590f585037968658c6a998c90d1e6595ceb39a52410641531f0a14137faf5d89894e9912d91b57652a7103c96c7b8747a84628ffb901c2d2f967ee6130d7851544477b24835d","share_pub":"8ce17b6af1fc18aa66298ef3f722cfa3e891d7cb5bdf116051a0af5aeb6f6c38a7860cd3d47427fb584240706c3a4acb","owner":"dcc846fa10c7cfce9e6eb37e06ed93b666cfc5e9"},"signature":"25defc7041c2bd425a2c405298d7a153f653bb6f49068938dc04ed81bd03ab7330c049dec3ff8389f7ae8b96cac5559ac39631e343f027cb41826af1b8cdee0fa9deb0b6210de725084396046904a7b2e9816decb2b33c88729ba1c814e59d04749c0a133e6190df830a919292cb8bd33535657c88936dfa05e83894517132c1fa090efe333b569dcbb7959477806bcdce797835e3a8e990f5a1ba091e9a6c5c4ba0807546882befd63273ab4a2857c0153b5adfe6342088d65d3192e782e30b5441f9aaf63d3ee6a0603bdc0be489202e347302aef61873593236102666e9a9bcbdb3b8322d18ffcc4a8bb41754207ec243a18099c6271ec8593f913d60bfac"},{"proof":{"validator":"aa57eab07f1a740672d0c106867d366c798d3b932d373c88cf047da1a3c16d0816ac58bab5a9d6f6f4b63a07608f8f39","encrypted_share":"2723a6b922ece432b6ad497ba50969534a5c3f9d29cc05739d951efc986f8ccf25d6d2d6331aeedb3a7f30c3f29760aeba44f9c42bc83f91ce3148ccca64e46405e627d02a4215d5bf1e166eb540838c8ac9cb79551e04a177ff765b4663ade090276606fba5c1cb91b731c8ef25d4e541dca541061eb1c1a5772f38445a65fbe40e6b2f072892e0a9982aebe70f27ce5ca241ffc7c0013f99d04ef586bb4ff04bbfbd1f3eaa6eddab998e355d3e160effd078baf364d187b59f1d74e75f891f2708f65326284f692dea5f82fa736f0b2b625dbf3bbda31047ed6d7b37f9138d9cc63d502282d926d59a9c825e58b15cd271497778d5f3d102e7012d1b8eac9e","share_pub":"8f476e4cb557aede4137696d63b2910a88bc55123e9a8a6b03cbc338f8fec3bf8bfffaf6dea5607b76c7d00915b4dd25","owner":"dcc846fa10c7cfce9e6eb37e06ed93b666cfc5e9"},"signature":"2dc3a34a648277a3e94d972d37b54b2a43c4754b8e934640064ab3274b2c43ea75998ba0b95457c7f8e518187779cf469c3a3850388921d00a2e1de31e98ccc8e22a5ece097bb904addc439c4f4de5dbe0f741fb554347807a853e94ae0c9271a22eaa9f909d97ffc9b16a84a1f6194fbf3aa032479e4fd1538a49cfd9d9a93d5269a8705a92bc0f92aa354012289eafc98e2774cf9a39e556e29b33bc9e0c18315bdb83a2b5310a5e8c556fa2888e89577c813e86e2fe90d74a9c67db2abde81e68cfa052e5e1519c62600f0f6044f1be48ebd665e42b480e74b54d19b91b3e1284e70fa6d0d268dabf0b4240eaae49d0928711bd7fe8e709375a68624c1c14"},{"proof":{"validator":"aa57eab07f1a740672d0c106867d366c798d3b932d373c88cf047da1a3c16d0816ac58bab5a9d6f6f4b63a07608f8f39","encrypted_share":"02a0a1e9cccf8aaaa2035fdc848fb4ca298a76a3462394391c500c68ef963e33ef91bacd2caa6e654ab05f50d49c38a8e0019ec41b15e63399defeaf6616461ddf646fd59d0830b54f99365446fcd8fad0059d7c4c5ea500e7d3e2d7cf739f5378a24b66ee2271878876b0db263db6c0391589deffba9e4d56f56ff23b9ee3dcafdfbcc1c3b0a8f365482393981b8b3d231840c32ee2be9fdceb4915544deec6fcd1803e75b0159883cd7eeb87bee17885724037b0fbfd93878b7b17c32260a41e491962b67fb405f0c39b47f9f93785f1852cb8fa75d133a025e64092b4ed18059a60abb24e9372a3ca0e08838abb9b42524a74ae3eab5410ee2314b6ab118e","share_pub":"a99bf34e0aef11749080ed4cf36acf1711460f5bdbd098b28de1afac618cc273f457613ca454262d8f5e46a8d010dd00","owner":"dcc846fa10c7cfce9e6eb37e06ed93b666cfc5e9"},"signature":"15050d609d2e256fa934578e57c8e97df895fb2e8a887317800c9797460bae1356c274ebeea1b2e96daedac9c24ecc4629acf0cff618fac8df6b4f90851e7e176b636d07db4eb204c02a752a69418e640f1b08d7e04310c95a91529421b4359b31f29ea6db08c3d2e0c7f0ca2f3f8f159812dcd48b613a4e706e02ce9164efcdf385fcbedbd465b5448b677c66bcfe433500e05ed97514c031c20e1ee54a1ee7b154b2aa7c0bdef97ad9db3d90e2dbe7b6ec0fda5b780ac6d4ac78f7e002d845d244e8ed43f6bda63280306e0c014b232b760dd1b16969fb9b0b5a139e3a41742d0870ec2f03c3606d2d36b782740fb4b2edad79082254bbda275060c64591d6"},{"proof":{"validator":"aa57eab07f1a740672d0c106867d366c798d3b932d373c88cf047da1a3c16d0816ac58bab5a9d6f6f4b63a07608f8f39","encrypted_share":"1078b9a1f64cd8efbdd44cd04ba262f78f79d7e1ba99b2c87a580492410314ae26cda143ffc09462335250f03e478fcd90c3a213f127d905e4ceef67b476d9ef649fc6e0f2603b9e1d8dac5801f962882797dfba972f96ca664ba51269b1a3534560a1313586a8f42c97431581a12bddaf6460c27248e31b5746d117da31f073f0862ea5a893852d0fe2174c73c6313dd3596f514e7e5758bf6f5195e80310f9c1f92374784d98c3503b1806b24e7a6f4e1bcf012e7d5f582d66899487247c3666ad131e5dcdd43812e4c4aa3e9f5e0f0120425e018b800e51b6207a3f4d4f8969c177603a3e360c95b18375d9c5b1405bf39a90bc24aa111bba009f772360e3","share_pub":"84faff90116da6fc1222fc2060ea24cae6f135c5d2684b09e5f9b0253d682df0a96025e81a606aafb2d434d774ec9340","owner":"dcc846fa10c7cfce9e6eb37e06ed93b666cfc5e9"},"signature":"7fddfeb183641468c6799448c6462e687a3278b6336aa48ec658340ef41a68910495053dc2fb0033cf42f7a9c035a0deb8b01eb184b4268fc3868dd4c2188727c1340dac988bf7c2917bc5b970dcb6d675446308e9f1d10d4079324a7ded67bf5efab995cd6168132d8b8e73d29062c70808b73b7c75e6efb292871c6710d4baee4dfc6c7043f011f5a38175fbbb062f2b49696fdb522c9a312eb319ee723619f86c0d253a4ae59c7968b3bd879538af052b49ab18aff88ea76ae60da29c9f2e327a76a1e35b357c45afa3c0435ece28f407946b26c44a1839dffb5c061689425fc6ccaeac9af5cfe72061ac02e13824f8d136aeeea6a292eb9738f83a57a1a1"}]`
	args := []string{"reshare",
		"--proofsRawJSON", proofsRawJSON,
		"--proofsFilePath", "./stubs/4/000001-0xaa57eab07f1a740672d0c106867d366c798d3b932d373c88cf047da1a3c16d0816ac58bab5a9d6f6f4b63a07608f8f39/proofs.json",
		"--operatorsInfo", string(operators),
		"--owner", owner.Hex(),
		"--withdrawAddress", "0x81592c3de184a3e2c0dcb5a261bc107bfa91f494",
		"--operatorIDs", "11,22,33,44",
		"--newOperatorIDs", "55,66,77,88",
		"--nonce", strconv.Itoa(10),
		"--ethKeystorePath", "./stubs/UTC--2024-06-14T14-05-12.366668334Z--dcc846fa10c7cfce9e6eb37e06ed93b666cfc5e9",
		"--ethKeystorePass", "./stubs/password",
		"--network", "holesky"}
	RootCmd.SetArgs(args)
	err = RootCmd.Execute()
	require.ErrorContains(t, err, "please provide either proofsRaw flag or proofsFilePath, not both")
	for _, srv := range servers {
		srv.HttpSrv.Close()
	}
}
