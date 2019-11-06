/*
 * skogul, logging utilities
 *
 * Copyright (c) 2019 Telenor Norge AS
 * Author(s):
 *  - hakon.solbjorg@telenor.com
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

package skogul

import (
	"strings"

	"github.com/sirupsen/logrus"
)

// ConfigureLogger sets up the logger based on calling parameters
func ConfigureLogger(requestedLoglevel string, logtimestamp bool) {
	loglevel := getLogLevelFromString(requestedLoglevel)
	logrus.SetLevel(loglevel)

	logrus.SetFormatter(&logrus.TextFormatter{
		DisableTimestamp: !logtimestamp,
	})
}

func getLogLevelFromString(requestedLevel string) logrus.Level {
	switch strings.ToLower(requestedLevel) {
	case "e", "error":
		return logrus.ErrorLevel
	case "w", "warn":
		return logrus.WarnLevel
	case "i", "info":
		return logrus.InfoLevel
	case "d", "debug":
		return logrus.DebugLevel
	case "v", "verbose", "t", "trace":
		return logrus.TraceLevel
	default:
		return logrus.WarnLevel
	}
}
