// Copyright 2020 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package ethash

import (
    "math/big"

    "github.com/ethereum/go-ethereum/core/types"
    "github.com/holiman/uint256"
)

const (
    // frontierDurationLimit is the target block time in seconds.
    frontierDurationLimit = 5

    // minimumDifficulty is the minimum that the difficulty may ever be.
    minimumDifficulty = 131072

    // expDiffPeriodUint is the exponential difficulty period.
    // Set to a high value to effectively disable the difficulty bomb.
    expDiffPeriodUint = 999999999

    // difficultyBoundDivisor is the bound divisor of the difficulty.
    // Set to 5 for a divisor of 1024 (2^5) to ensure smoother adjustments.
    difficultyBoundDivisor = 5 // Divisor of 1024 (2^5)
)

// CalcDifficultyFrontierU256 calculates the difficulty using Frontier rules.
func CalcDifficultyFrontierU256(time uint64, parent *types.Header) *big.Int {
    pDiff, _ := uint256.FromBig(parent.Difficulty)
    adjust := pDiff.Clone()
    adjust.Rsh(adjust, difficultyBoundDivisor) // adjust = pDiff / 1024

    timeDiff := int64(time - parent.Time)

    if timeDiff < frontierDurationLimit {
        // Increase difficulty by adjust
        pDiff.Add(pDiff, adjust)
    } else {
        // Decrease difficulty by adjust
        if pDiff.Gt(adjust) {
            pDiff.Sub(pDiff, adjust)
        } else {
            pDiff.SetUint64(minimumDifficulty)
        }
    }

    // Enforce minimum difficulty
    if pDiff.LtUint64(minimumDifficulty) {
        pDiff.SetUint64(minimumDifficulty)
    }

    // Disable the difficulty bomb by not adding the exponential factor
    return pDiff.ToBig()
}

// CalcDifficultyHomesteadU256 calculates the difficulty using Homestead rules.
func CalcDifficultyHomesteadU256(time uint64, parent *types.Header) *big.Int {
    pDiff, _ := uint256.FromBig(parent.Difficulty)
    adjust := pDiff.Clone()
    adjust.Rsh(adjust, difficultyBoundDivisor) // adjust = pDiff / 1024

    timeDiff := int64(time - parent.Time)
    var sign int64

    if timeDiff < frontierDurationLimit {
        // Increase difficulty
        sign = 1
    } else {
        // Decrease difficulty
        sign = -1
    }

    if sign > 0 {
        pDiff.Add(pDiff, adjust)
    } else {
        if pDiff.Gt(adjust) {
            pDiff.Sub(pDiff, adjust)
        } else {
            pDiff.SetUint64(minimumDifficulty)
        }
    }

    // Enforce minimum difficulty
    if pDiff.LtUint64(minimumDifficulty) {
        pDiff.SetUint64(minimumDifficulty)
    }

    // Disable the difficulty bomb by not adding the exponential factor
    return pDiff.ToBig()
}

// MakeDifficultyCalculatorU256 creates a difficulty calculator with the given bomb delay.
func MakeDifficultyCalculatorU256(bombDelay *big.Int) func(time uint64, parent *types.Header) *big.Int {
    // bombDelay is effectively disabled; no changes needed related to the bomb
    return func(time uint64, parent *types.Header) *big.Int {
        pDiff, _ := uint256.FromBig(parent.Difficulty)
        adjust := pDiff.Clone()
        adjust.Rsh(adjust, difficultyBoundDivisor) // adjust = pDiff / 1024

        timeDiff := int64(time - parent.Time)
        var sign int64

        if timeDiff < frontierDurationLimit {
            // Increase difficulty
            sign = 1
        } else {
            // Decrease difficulty
            sign = -1
        }

        if sign > 0 {
            pDiff.Add(pDiff, adjust)
        } else {
            if pDiff.Gt(adjust) {
                pDiff.Sub(pDiff, adjust)
            } else {
                pDiff.SetUint64(minimumDifficulty)
            }
        }

        // Enforce minimum difficulty
        if pDiff.LtUint64(minimumDifficulty) {
            pDiff.SetUint64(minimumDifficulty)
        }

        // Disable the difficulty bomb by not adding the exponential factor
        return pDiff.ToBig()
    }
}
