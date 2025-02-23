package bitcoin

import (
	"os"

	bitcoinv1alpha1 "github.com/kotalco/kotal/apis/bitcoin/v1alpha1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Bitcoin core client", func() {

	node := &bitcoinv1alpha1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "bitcoin-node",
			Namespace: "default",
		},
		Spec: bitcoinv1alpha1.NodeSpec{
			Network:          "mainnet",
			RPC:              true,
			P2PPort:          8888,
			P2PHost:          "127.0.0.1",
			RPCPort:          7777,
			RPCHost:          "127.0.0.1",
			Wallet:           false,
			TransactionIndex: true,
		},
	}

	node.Default()
	// nil is passed because there's no reconciler client
	// TODO: create test for rpcUsers where client is not nil
	client := NewClient(node, nil)

	It("Should get correct image", func() {
		// default image
		img := client.Image()
		Expect(img).To(Equal(DefaultBitcoinCoreImage))
		// after setting custom image
		testImage := "kotalco/bitcoin-core:spec"
		node.Spec.Image = &testImage
		Expect(client.Image()).To(Equal(testImage))
		// after setting custom image
		testImage = "kotalco/bitcoin-core:test"
		os.Setenv(EnvBitcoinCoreImage, testImage)
		Expect(client.Image()).To(Equal(testImage))
	})

	It("Should get correct command", func() {
		Expect(client.Command()).To(Equal([]string{
			"bitcoind",
		}))
	})

	It("Should get correct home directory", func() {
		Expect(client.HomeDir()).To(Equal(BitcoinCoreHomeDir))
	})

	It("Should generate correct client arguments", func() {
		Expect(client.Args()).To(ContainElements([]string{
			"-chain=main",
			"-datadir=/home/bitcoin/kotal-data",
			"-server=1",
			"-bind=127.0.0.1:8888",
			"-rpcport=7777",
			"-rpcbind=127.0.0.1",
			"-rpcallowip=0.0.0.0/0",
			"-disablewallet",
			"-txindex=1",
		}))
	})

})
