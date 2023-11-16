package genlog_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/winebarrel/genlog"
)

func TestParse(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	var buf bytes.Buffer
	fmt.Fprintln(&buf, `2020-05-27T05:03:27.500301Z   11 Query	SET @@sql_log_bin=off
2020-05-27T05:03:27.543379Z   11 Query	select @@session.tx_read_only
2020-05-27T05:03:27.683485Z   11 Query	COMMIT`)

	bs := []*genlog.Block{}

	err := genlog.Parse(&buf, func(block *genlog.Block) {
		bs = append(bs, block)
	})

	require.NoError(err)
	assert.Equal([]*genlog.Block{
		{
			Time:     "2020-05-27T05:03:27.500301Z",
			Id:       "11",
			Command:  "Query",
			Argument: "SET @@sql_log_bin=off",
		},
		{
			Time:     "2020-05-27T05:03:27.543379Z",
			Id:       "11",
			Command:  "Query",
			Argument: "select @@session.tx_read_only",
		},
		{
			Time:     "2020-05-27T05:03:27.683485Z",
			Id:       "11",
			Command:  "Query",
			Argument: "COMMIT",
		},
	}, bs)
}
