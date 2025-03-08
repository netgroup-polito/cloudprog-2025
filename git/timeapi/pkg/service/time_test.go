package service_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/cheina97/timeserver/pkg/service"
)

var _ = Describe("Time", func() {
	Context("GetTimeWithTimezone", func() {
		It("should return the current time in the specified timezone", func() {
			location := "America/New_York"
			t, err := service.GetTimeWithTimezone(location)
			Expect(err).ToNot(HaveOccurred())
			loc, _ := time.LoadLocation(location)
			Expect(t.Location()).To(Equal(loc))
		})

		It("should return an error for an invalid timezone", func() {
			location := "Invalid/Timezone"
			_, err := service.GetTimeWithTimezone(location)
			Expect(err).To(HaveOccurred())
		})

		It("should return the current time in UTC when no timezone is provided", func() {
			location := ""
			t, err := service.GetTimeWithTimezone(location)
			Expect(err).ToNot(HaveOccurred())
			Expect(t.Location()).To(Equal(time.UTC))
		})
	})
})
