package log

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestQuickLogger(t *testing.T) {
	Convey("TestQuickLogger", t, func() {
		Convey("GetQuickLogger", func() {
			logtest := GetQuickLogger("test")
			So(logtest, ShouldNotBeNil)
			So(quickLogger["test"], ShouldEqual, logtest)
			logtest.SetLevel(LLDebug)
			logtest.Debug("debug")
			logtest.Info("info")
			logtest.Warn("warn")
			logtest.Error("error")
			logt2 := GetQuickLogger("t2log")
			So(logt2, ShouldNotBeNil)
			So(quickLogger["t2log"], ShouldEqual, logt2)
			logt2.Debug("debug")
			logt2.Info("info")
			logt2.Warn("warn")
			logt2.Error("error")

			prt := logtest.Debugxt()
			So(prt == nil, ShouldBeFalse)
			prt("debug xt")
			prt = logt2.Debugxt()
			So(prt == nil, ShouldBeTrue)

			logtvfy := GetQuickLogger("test")
			So(logtvfy, ShouldEqual, logtest)
			logt2vfy := GetQuickLogger("t2log")
			So(logt2vfy, ShouldEqual, logt2)

			ctx := WithTrace(context.Background(), "testTrace")
			ltt, err := logtest.FromTrace(ctx)
			So(err, ShouldBeNil)
			So(ltt, ShouldNotBeNil)
			ltt.Debug("debug")
			ltt.Info("info")
			ltt.Warn("warn")
			ltt.Error("error")

			ctx2 := WithTrace(ctx, "testTraceII")
			ltt2, err := ltt.FromTrace(ctx2)
			So(err, ShouldBeNil)
			So(ltt2, ShouldNotBeNil)
			ltt2.SetLevel(LLWarn)
			ltt2.Debug("debug")
			ltt2.Info("info")
			ltt2.Warn("warn")
			ltt2.Error("error")
			ltt.Debug("debug")
			ltt.Info("info")

			_, err = ltt.FromTrace(context.Background())
			So(err, ShouldNotBeNil)
		})
	})
}
