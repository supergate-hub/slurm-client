// SPDX-FileCopyrightText: Copyright (C) SchedMD LLC.
// SPDX-License-Identifier: Apache-2.0

package diag

import (
	api "github.com/supergate-hub/slurm-client/api/v0040"
)

// DiagInfo represents Slurm database diagnostic information.
type DiagInfo = api.V0040StatsRec
