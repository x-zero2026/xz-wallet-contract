// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

/**
 * @title XZToken
 * @dev XZ Token (XZT) - ERC20 token for XZ Wallet task escrow system
 * 
 * Initial Supply: 10,000 XZT
 * Decimals: 18
 * Network: Sepolia Testnet
 */
contract XZToken is ERC20, Ownable {
    
    /**
     * @dev Constructor mints initial supply to deployer
     */
    constructor() ERC20("XZ Token", "XZT") Ownable(msg.sender) {
        // Mint 10,000 XZT to deployer (system wallet)
        _mint(msg.sender, 10000 * 10**decimals());
    }
    
    /**
     * @dev Allows owner to mint additional tokens if needed
     * @param to Address to receive minted tokens
     * @param amount Amount of tokens to mint (in wei)
     */
    function mint(address to, uint256 amount) external onlyOwner {
        _mint(to, amount);
    }
    
    /**
     * @dev Burns tokens from caller's balance
     * @param amount Amount of tokens to burn (in wei)
     */
    function burn(uint256 amount) external {
        _burn(msg.sender, amount);
    }
}
