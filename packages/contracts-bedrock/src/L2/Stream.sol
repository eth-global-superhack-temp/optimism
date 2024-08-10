// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

import { Predeploys } from "../libraries/Predeploys.sol";
import { ISemver } from "../universal/ISemver.sol";
import { Ownable } from "@openzeppelin/contracts/access/Ownable.sol";

interface IStream {
    function stream() external;
}

/**
 * @custom:proxied
 * @custom:predeploy 0x42000000000000000000000000000000000000A0
 * @title Stream
 * @notice The Stream predeploy streams tokens for the grants.
 */
contract Stream is ISemver, Ownable, IStream {
    /**
     * @notice Address of the special depositor account.
     */
    address public constant DEPOSITOR_ACCOUNT = 0xDeaDDEaDDeAdDeAdDEAdDEaddeAddEAdDEAd0001;

    /**
     * @notice Address of the stream contract to be called.
     */
    address public target;

    /// @notice Semantic version.
    /// @custom:semver 2.4.0
    string public constant version = "2.4.0";

    /**
     * @param _owner Address that will initially own this contract.
     */
    constructor(address _owner) Ownable() {
        transferOwnership(_owner);
    }

    /**
     * @notice Allows the owner to modify the target address.
     * @param _target New target address.
     */
    function setTarget(address _target) public onlyOwner {
        target = _target;
    }

    /**
     * @notice Calls the stream function in the target contract.
     */
    function stream() external {
        require(msg.sender == DEPOSITOR_ACCOUNT, "Stream: only the depositor account can stream");
        if (target == address(0)) {
            return;
        }
        IStream(target).stream();
    }
}