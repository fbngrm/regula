package regula

import (
	"testing"

	"github.com/heetch/regula/rule"
	"github.com/stretchr/testify/require"
)

func TestGetString(t *testing.T) {
	p := Params{
		"string": "string",
		"bool":   true,
	}

	t.Run("GetString - OK", func(t *testing.T) {
		v, err := p.GetString("string")
		require.NoError(t, err)
		require.Equal(t, "string", v)
	})

	t.Run("GetString - NOK - ErrParamNotFound", func(t *testing.T) {
		_, err := p.GetString("badkey")
		require.Error(t, err)
		require.Equal(t, err, rule.ErrParamNotFound)
	})

	t.Run("GetString - NOK - ErrParamTypeMismatch", func(t *testing.T) {
		_, err := p.GetString("bool")
		require.Error(t, err)
		require.Equal(t, err, rule.ErrParamTypeMismatch)
	})
}

func TestGetBool(t *testing.T) {
	p := Params{
		"bool":   true,
		"string": "string",
	}

	t.Run("GetBool - OK", func(t *testing.T) {
		v, err := p.GetBool("bool")
		require.NoError(t, err)
		require.Equal(t, true, v)
	})

	t.Run("GetBool - NOK - ErrParamNotFound", func(t *testing.T) {
		_, err := p.GetBool("badkey")
		require.Error(t, err)
		require.Equal(t, err, rule.ErrParamNotFound)
	})

	t.Run("GetBool - NOK - ErrParamTypeMismatch", func(t *testing.T) {
		_, err := p.GetBool("string")
		require.Error(t, err)
		require.Equal(t, err, rule.ErrParamTypeMismatch)
	})
}

func TestGetInt64(t *testing.T) {
	p := Params{
		"int64":  int64(42),
		"string": "string",
	}

	t.Run("GetInt64 - OK", func(t *testing.T) {
		v, err := p.GetInt64("int64")
		require.NoError(t, err)
		require.Equal(t, int64(42), v)
	})

	t.Run("GetInt64 - NOK - ErrParamNotFound", func(t *testing.T) {
		_, err := p.GetInt64("badkey")
		require.Error(t, err)
		require.Equal(t, err, rule.ErrParamNotFound)
	})

	t.Run("GetInt64 - NOK - ErrParamTypeMismatch", func(t *testing.T) {
		_, err := p.GetInt64("string")
		require.Error(t, err)
		require.Equal(t, err, rule.ErrParamTypeMismatch)
	})
}

func TestGetFloat64(t *testing.T) {
	p := Params{
		"float64": 42.42,
		"string":  "string",
	}

	t.Run("GetFloat64 - OK", func(t *testing.T) {
		v, err := p.GetFloat64("float64")
		require.NoError(t, err)
		require.Equal(t, 42.42, v)
	})

	t.Run("GetFloat64 - NOK - ErrParamNotFound", func(t *testing.T) {
		_, err := p.GetFloat64("badkey")
		require.Error(t, err)
		require.Equal(t, err, rule.ErrParamNotFound)
	})

	t.Run("GetFloat64 - NOK - ErrParamTypeMismatch", func(t *testing.T) {
		_, err := p.GetFloat64("string")
		require.Error(t, err)
		require.Equal(t, err, rule.ErrParamTypeMismatch)
	})
}
