// Copyright © 2022 Obol Labs Inc.
//
// This program is free software: you can redistribute it and/or modify it
// under the terms of the GNU General Public License as published by the Free
// Software Foundation, either version 3 of the License, or (at your option)
// any later version.
//
// This program is distributed in the hope that it will be useful, but WITHOUT
// ANY WARRANTY; without even the implied warranty of  MERCHANTABILITY or
// FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for
// more details.
//
// You should have received a copy of the GNU General Public License along with
// this program.  If not, see <http://www.gnu.org/licenses/>.

// Code generated by "stringer -type=OrderStop -trimprefix=Stop"; DO NOT EDIT.

package lifecycle

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[StopTracing-8]
	_ = x[StopScheduler-9]
	_ = x[StopP2PPeerDB-10]
	_ = x[StopP2PTCPNode-11]
	_ = x[StopP2PUDPNode-12]
	_ = x[StopMonitoringAPI-13]
	_ = x[StopBeaconMock-14]
	_ = x[StopValidatorAPI-15]
	_ = x[StopRetryer-16]
}

const _OrderStop_name = "TracingSchedulerP2PPeerDBP2PTCPNodeP2PUDPNodeMonitoringAPIBeaconMockValidatorAPIRetryer"

var _OrderStop_index = [...]uint8{0, 7, 16, 25, 35, 45, 58, 68, 80, 87}

func (i OrderStop) String() string {
	i -= 8
	if i < 0 || i >= OrderStop(len(_OrderStop_index)-1) {
		return "OrderStop(" + strconv.FormatInt(int64(i+8), 10) + ")"
	}
	return _OrderStop_name[_OrderStop_index[i]:_OrderStop_index[i+1]]
}
