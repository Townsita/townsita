package townsta

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTownsta(t *testing.T) {
	Convey("Given new Townsta", t, func() {
		tw := NewTownsta()
		Convey("Townsta is not nil", func() {
			So(tw, ShouldNotBeNil)
		})
	})
}
