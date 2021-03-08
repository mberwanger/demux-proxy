package cmd

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRootCmd(t *testing.T) {
	var mem = &exitMemento{}
	Execute("1.2.3", mem.Exit, []string{"-h"})
	require.Equal(t, 0, mem.code)
}

func TestRootCmdHelp(t *testing.T) {
	var mem = &exitMemento{}
	var cmd = newRootCmd("", mem.Exit).cmd
	cmd.SetArgs([]string{"-h"})
	require.NoError(t, cmd.Execute())
	require.Equal(t, 0, mem.code)
}

func TestRootCmdVersion(t *testing.T) {
	var b bytes.Buffer
	var mem = &exitMemento{}
	var cmd = newRootCmd("1.2.3", mem.Exit).cmd
	cmd.SetOut(&b)
	cmd.SetArgs([]string{"-v"})
	require.NoError(t, cmd.Execute())
	require.Equal(t, "demux-proxy version 1.2.3\n", b.String())
	require.Equal(t, 0, mem.code)
}
