package httpreq

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Client", func() {
	Context("test DefaultClient", func() {
		It("test timeout", func() {
			resp, err := Get(ts.URL+"/test/timeout",
				WithTimeout(time.Second/10),
				WithMiddleWares(ExpectStatusCodeMW(200)))
			Expect(err).Should(HaveOccurred())
			Expect(resp).Should(BeNil())
		})
		It("should do a GET", func() {
			resp, err := Get(ts.URL+"/test/get",
				WithBody("abcdeft"),
				WithDebug(true),
				WithMiddleWares(ExpectStatusCodeMW(200)))
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.String()).Should(Equal("this is a get request"))
			Expect(resp.StatusCode()).Should(Equal(200))
			Expect(resp.Status()).Should(Equal("200 OK"))
		})
		It("should do a POST", func() {
			resp, err := Post(ts.URL+"/test/post",
				"application/x-www-form-urlencoded", nil,
				WithForm("a", "a"),
				WithForm("b", "b"),
				WithDebug(true),
				WithMiddleWares(ExpectStatusCodeMW(200)))
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.String()).Should(Equal("this is a post request"))
			Expect(resp.StatusCode()).Should(Equal(200))
			Expect(resp.Status()).Should(Equal("200 OK"))
		})
	})
})
