/*
 * Copyright (C) 2021 The poly network Authors
 * This file is part of The poly network library.
 *
 * The  poly network  is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The  poly network  is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 * You should have received a copy of the GNU Lesser General Public License
 * along with The poly network .  If not, see <http://www.gnu.org/licenses/>.
 */

package manager

import (
	"fmt"
	"github.com/urfave/cli/v2"
)

const (
	BATCH_CREATE_ACCOUNT = "batchcreateaccount"
	BATCH_TRANSFER_TOKEN = "batchtransfertoken"
)

var _Handlers = map[string]func(*cli.Context) error{}

func init() {
	_Handlers[BATCH_CREATE_ACCOUNT] = BatchCreateAccount
	_Handlers[BATCH_TRANSFER_TOKEN] = BatchTransferToken
}

func HandleCommand(method string, ctx *cli.Context) error {
	h, ok := _Handlers[method]
	if !ok {
		return fmt.Errorf("Unsupported subcommand %s", method)
	}
	return h(ctx)
}
