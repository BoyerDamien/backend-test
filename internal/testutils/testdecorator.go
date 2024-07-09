package testutils

import (
	"context"
	"os"
	"testing"
	"time"

	charmLog "github.com/charmbracelet/log"
	"github.com/japhy-tech/backend-test/internal/gateways"
	"github.com/japhy-tech/backend-test/internal/gateways/mysql"
	"github.com/maxatome/go-testdeep/td"
)

const (
	MysqlDSN = "root:root@(localhost:53306)/core?parseTime=true"
)

func TestDecorator(t *testing.T, test func(context.Context, gateways.IDatastore, *td.T)) {
	var (
		logger = charmLog.NewWithOptions(os.Stderr, charmLog.Options{
			Formatter:       charmLog.TextFormatter,
			ReportCaller:    true,
			ReportTimestamp: true,
			TimeFormat:      time.Kitchen,
			Prefix:          "🧑‍💻 backend-test",
			Level:           charmLog.DebugLevel,
		})
		datastore = mysql.New(MysqlDSN, logger)
		require   = td.Require(t)
		ctx       = context.Background()
	)
	defer func() {
		require.CmpNoError(datastore.Reset(ctx))
		require.CmpNoError(datastore.Close())
	}()

	test(ctx, datastore, require)
}
