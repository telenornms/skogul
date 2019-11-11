/*
 * skogul, switch transformer
 *
 * Copyright (c) 2019 Telenor Norge AS
 * Author(s):
 *  - Håkon Solbjørg <hakon.solbjorg@telenor.com>
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

package transformer

import (
	log "github.com/sirupsen/logrus"
	"github.com/telenornms/skogul"
)

type Case struct {
	When   string `doc:"Used as a conditional statement on a field"`
	Is     string `doc:"Used for the specific value of the stated field"`
	Invert bool   `doc:"Invert the result of the conditional"`
	// Transformers []skogul.Transformer `doc:"The transformers to run when the defined conditional is true"`
	Transformers []string `doc:"The transformers to run when the defined conditional is true"`
}

type Switch struct {
	Cases []Case `doc:"A list of switch cases "`
}

func (sw *Switch) Transform(c *skogul.Container) error {
	for _, cas := range sw.Cases {
		log.Warningf("cases: %+v", cas)

		// c.
	}

	return nil
}
