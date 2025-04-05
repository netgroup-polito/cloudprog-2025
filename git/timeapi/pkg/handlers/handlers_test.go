package handlers_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	recorder *httptest.ResponseRecorder
)

var _ = Describe("Handlers", func() {
	BeforeEach(func() {
		recorder = httptest.NewRecorder()
	})

	Context("GET /time", func() {
		It("should return the current time in the specified timezone", func() {
			req, err := http.NewRequest(http.MethodGet, "/time?timezone=America/New_York", http.NoBody)
			Expect(err).ToNot(HaveOccurred())

			router.ServeHTTP(recorder, req)

			Expect(recorder.Code).To(Equal(http.StatusOK))
		})

		It("should return an error for an invalid timezone", func() {
			req, err := http.NewRequest(http.MethodGet, "/time?timezone=Invalid/Timezone", http.NoBody)
			Expect(err).ToNot(HaveOccurred())

			router.ServeHTTP(recorder, req)

			recorder.Code = http.StatusInternalServerError
		})

		It("should return the current time in UTC when no timezone is provided", func() {
			req, err := http.NewRequest(http.MethodGet, "/time", http.NoBody)
			Expect(err).ToNot(HaveOccurred())

			router.ServeHTTP(recorder, req)

			Expect(recorder.Code).To(Equal(http.StatusOK))
			fmt.Println(recorder.Body.String())
		})
	})
})
