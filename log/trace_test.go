package log

import (
	"context"
	"regexp"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTraceCtx(t *testing.T) {
	Convey("Test trace context", t, func() {
		So(traceInit, ShouldNotBeNil)

		Convey("Test nextTraceID", func() {
			rexID := regexp.MustCompile(`^[0-9a-f]{16}-[0-9a-f]{32}-[0-9a-f]{16}$`)
			for i := 0; i < 10; i++ {
				id := nextTraceID()
				So(id, ShouldNotBeEmpty)
				t.Logf("trace id: %s\n  sum:%s - hex:%s",
					id.String(), id.Sum(), id.Hex())
				So(rexID.MatchString(id.String()), ShouldBeTrue)
				So(len(id.Sum()), ShouldEqual, 16)
				So(len(id.Hex()), ShouldEqual, 64)
			}
		})

		Convey("Test traceCtx", func() {
			basectx := context.Background()

			l1ctx := WithTrace(basectx, "test1")
			So(l1ctx, ShouldNotBeNil)
			So(l1ctx.Value(traceKey), ShouldNotBeNil)
			t1val := GetTrace(l1ctx)
			So(t1val, ShouldNotBeNil)
			So(len(t1val), ShouldEqual, 1)
			t.Logf("trace value L1: %s", t1val)
			t.Logf("trace detail L1: %s", t1val.Format("ID-", "-END", ">"))
			l2ctx := WithTrace(l1ctx, "test2")
			l2ctxwc, _ := context.WithCancel(l2ctx)
			l3ctx := WithTrace(l2ctxwc, "test3")
			So(l3ctx, ShouldNotBeNil)
			tval := GetTrace(l3ctx)
			So(tval, ShouldNotBeNil)
			So(len(tval), ShouldEqual, 3)
			So(tval[0].Name, ShouldEqual, "test1")
			So(tval[1].Name, ShouldEqual, "test2")
			So(tval[2].Name, ShouldEqual, "test3")
			t.Logf("trace value L3: %s", tval)
			t.Logf("trace detail L3: %s", tval.Format("T.", "-ID", "\n    >>>> "))
		})
	})
}
