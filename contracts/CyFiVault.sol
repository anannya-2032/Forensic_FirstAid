// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

contract CyFiVault {
    // Mapping from file hash to the timestamp it was recorded
    mapping(string => uint256) private evidenceRegistry;

    event EvidenceAnchored(string fileHash, uint256 timestamp);

    // Function to store the hash on the blockchain
    function anchorEvidence(string memory _fileHash) public {
        // Ensure we don't overwrite an existing record to maintain integrity
        require(evidenceRegistry[_fileHash] == 0, "Evidence already exists on ledger.");
        
        evidenceRegistry[_fileHash] = block.timestamp;
        emit EvidenceAnchored(_fileHash, block.timestamp);
    }

    // Function to verify if a hash exists and return its timestamp
    function verifyEvidence(string memory _fileHash) public view returns (uint256) {
        return evidenceRegistry[_fileHash];
    }
}