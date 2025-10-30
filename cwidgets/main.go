package cwidgets

import (
	"github.com/bcicen/ctop/models"
)

type WidgetUpdater interface {
	SetMeta(models.Meta)
	SetMetrics(models.Metrics)
}

type NullWidgetUpdater struct{}

// NullWidgetUpdater implements WidgetUpdater
func (wu NullWidgetUpdater) SetMeta(models.Meta) {}

// NullWidgetUpdater implements WidgetUpdater
func (wu NullWidgetUpdater) SetMetrics(models.Metrics) {}
