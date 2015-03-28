package townsita

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTownsita(t *testing.T) {
	Convey("Given new Townsita", t, func() {
		tw := New(NewConfig(), nil)
		Convey("Townsita is not nil", func() {
			So(tw, ShouldNotBeNil)
		})
	})
}
