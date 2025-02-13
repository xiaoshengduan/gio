// SPDX-License-Identifier: Unlicense OR MIT

package unit_test

import (
	"testing"

	"github.com/xiaoshengduan/gio-fly/unit"
)

func TestMetric_DpToSp(t *testing.T) {
	m := unit.Metric{
		PxPerDp: 2,
		PxPerSp: 3,
	}

	{
		exp := m.Dp(5)
		got := m.Sp(m.DpToSp(5))
		if got != exp {
			t.Errorf("DpToSp conversion mismatch %v != %v", exp, got)
		}
	}

	{
		exp := m.Sp(5)
		got := m.Dp(m.SpToDp(5))
		if got != exp {
			t.Errorf("SpToDp conversion mismatch %v != %v", exp, got)
		}
	}
}
