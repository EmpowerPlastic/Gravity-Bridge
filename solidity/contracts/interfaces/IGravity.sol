//SPDX-License-Identifier: Apache-2.0
pragma solidity 0.8.10;

interface IGravity {

    function state_lastValsetCheckpoint() external view returns (bytes32);
    
}