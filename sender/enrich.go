/*
 * skogul, encrichment-updater sender
 *
 * Copyright (c) 2019 Telenor Norge AS
 * Author(s):
 *  - Kristian Lyngstøl <kly@kly.no>
 *
 * This library is free software; you can redistribute it and/or
 * modify it under the terms of the GNU Lesser General Public
 * License as published by the Free Software Foundation; either
 * version 2.1 of the License, or (at your option) any later version.
 *
 * This library is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
 * Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public
 * License along with this library; if not, write to the Free Software
 * Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA
 * 02110-1301  USA
 */

package sender

import (
	"fmt"

	"github.com/telenornms/skogul"
	"github.com/telenornms/skogul/transformer"
)

var enrichLog = skogul.Logger("sender", "enrichment")

// EnrichmentUpdater sends any received container/metric to the
// update-function of the provided transformer, allowing on-the-fly updates
// to enrichment.
type EnrichmentUpdater struct {
	Enricher skogul.TransformerRef `doc:"The enrichment transformer to update."`
}

// Uses received metrics to update the enrichment transformer
func (e *EnrichmentUpdater) Send(c *skogul.Container) error {
	er, _ := e.Enricher.T.(*transformer.Enrich)
	er.Update(c)
	return nil
}

func (e *EnrichmentUpdater) Verify() error {
	_, ok := e.Enricher.T.(*transformer.Enrich)
	if !ok {
		return fmt.Errorf("provided transformer in enrichmentupdater is not an enrichment transformer")
	}
	return nil
}
