package handler

import "errors"

var errMissingGuild = errors.New("X-Guild-ID header is required")
