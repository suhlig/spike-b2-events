package b2_test

import (
	"encoding/json"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/suhlig/spike-b2-events/b2"
)

type Timestamped struct {
	Timestamp b2.Timestamp
}

var _ = Describe("Timestamp", func() {
	var (
		err         error
		timestamped Timestamped
	)

	JustBeforeEach(func() {
		err = json.Unmarshal([]byte(`{"timestamp": 1716482076260}`), &timestamped)
	})

	It("parses", func() {
		Expect(err).ToNot(HaveOccurred())
	})

	It("has the expected timestamp", func() {
		Expect(time.Time(timestamped.Timestamp)).To(Equal(time.Unix(1716482076, 260000000).UTC()))
	})
})
