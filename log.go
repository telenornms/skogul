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

import "github.com/sirupsen/logrus"

// ConfigureLogger sets up the logger based on calling parameters
func ConfigureLogger(requestedLoglevel, logformat string) {
	loglevel := getLogLevelFromString(requestedLoglevel)
	logrus.SetLevel(loglevel)

	textFormatter := getLogTextFormatter(logformat)
	logrus.SetFormatter(textFormatter)
}

func getLogTextFormatter(requestedFormatter string) *logrus.TextFormatter {
	switch requestedFormatter {
	case "syslog":
		return &syslogFormat
	default:
		return &logrus.TextFormatter{}
	}
}

var syslogFormat = logrus.TextFormatter{
	DisableTimestamp: true,
	DisableColors:    true,
}

func getLogLevelFromString(requestedLevel string) logrus.Level {
	switch requestedLevel {
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
