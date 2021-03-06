package test

import (
	"net/http"
	"time"

	"github.com/nats-io/nats"
	"github.com/onsi/ginkgo"

	. "github.com/onsi/gomega"

	"github.com/cloudfoundry/gorouter/route"
	"github.com/cloudfoundry/gorouter/test/common"
	"github.com/cloudfoundry/gorouter/test_util"
)

func NewWebSocketApp(urls []route.Uri, rPort uint16, mbusClient *nats.Conn, delay time.Duration) *common.TestApp {
	app := common.NewTestApp(urls, rPort, mbusClient, nil, "")
	app.AddHandler("/", func(w http.ResponseWriter, r *http.Request) {
		defer ginkgo.GinkgoRecover()

		Expect(r.Header.Get("Upgrade")).To(Equal("websocket"))
		Expect(r.Header.Get("Connection")).To(Equal("upgrade"))

		conn, _, err := w.(http.Hijacker).Hijack()
		x := test_util.NewHttpConn(conn)

		resp := test_util.NewResponse(http.StatusSwitchingProtocols)
		resp.Header.Set("Upgrade", "websocket")
		resp.Header.Set("Connection", "upgrade")

		time.Sleep(delay)

		x.WriteResponse(resp)
		Expect(err).ToNot(HaveOccurred())

		x.CheckLine("hello from client")
		x.WriteLine("hello from server")
	})

	return app
}
