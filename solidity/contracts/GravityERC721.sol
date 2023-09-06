//SPDX-License-Identifier: Apache-2.0
pragma solidity 0.8.10;

import "@openzeppelin/contracts/security/ReentrancyGuard.sol";
import "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";
import { ERC721Holder } from "@openzeppelin/contracts/token/ERC721/utils/ERC721Holder.sol";
import "./interfaces/IGravity.sol";
import "./CosmosNFT.sol";

error InvalidSignature();
error InsufficientPower(uint256 cumulativePower, uint256 powerThreshold);
error IncorrectCheckpoint();
error MalformedCurrentValidatorSet();
error InvalidWithdrawalNonce(uint256 newNonce, uint256 currentNonce);
error WithdrawalTimedOut();
error MintTimedOut();

// This is used purely to avoid stack too deep errors
// represents everything about a given validator set
struct ValsetArgs {
	// the validators in this set, represented by an Ethereum address
	address[] validators;
	// the powers of the given validators in the same order as above
	uint256[] powers;
	// the nonce of this validator set
	uint256 valsetNonce;
	// the reward amount denominated in the below reward token, can be
	// set to zero
	uint256 rewardAmount;
	// the reward token, should be set to the zero address if not being used
	address rewardToken;
}

// This represents a validator signature
struct Signature {
	uint8 v;
	bytes32 r;
	bytes32 s;
}

contract GravityERC721 is ERC721Holder, ReentrancyGuard {
	
	// The number of 'votes' required to execute NFT withdrawal
	uint256 constant constant_powerThreshold = 2863311530;

	address public immutable state_gravitySolAddress;
	bytes32 public immutable state_gravityId;
	
	uint256 public state_lastERC721EventNonce = 1;
	mapping(address => uint256) public state_lastWithdrawalNonces;

	event SendERC721ToCosmosEvent(
		address indexed _tokenContract,
		address indexed _sender,
		string _destination,
		uint256 _tokenId,
		uint256 _eventNonce,
		string _tokenURI
	);

	event ERC721WithdrawnEvent(
		uint256 _withdrawNonce,
		address indexed _tokenContract,
		uint256 _eventNonce
	);

	event ERC721DeployedEvent(
		string _cosmosClass,
		address indexed _tokenContract,
		string _name,
		string _symbol,
		uint256 _eventNonce
	);

	event ERC721MintedEvent(
		address indexed _tokenContract,
		uint256 _tokenId,
		string _tokenURI,
		uint256 _eventNonce
	);

	event GravityERC721DeployedEvent();

	constructor(
		// reference gravity.sol for many functions peformed here
		address _gravitySolAddress,
		bytes32 _gravityId
	) {
		state_gravitySolAddress = _gravitySolAddress;
		state_gravityId = _gravityId;
		emit GravityERC721DeployedEvent();
		}

	function sendERC721ToCosmos(
		address _tokenContract,
		string calldata _destination,
		uint256 _tokenId
	) external nonReentrant {
		ERC721(_tokenContract).safeTransferFrom(msg.sender, address(this), _tokenId);

		emit SendERC721ToCosmosEvent(
			_tokenContract,
			msg.sender,
			_destination,
			_tokenId,
			state_lastERC721EventNonce,
			ERC721(_tokenContract).tokenURI(_tokenId)
		);
		state_lastERC721EventNonce = state_lastERC721EventNonce + 1;
	}

	function withdrawERC721 (
		ValsetArgs calldata _currentValset,
		Signature[] calldata _sigs,
		address _tokenContract,
		uint256[] calldata _tokenIds,
		address[] calldata _destinations,
		uint256 _withdrawNonce,
		uint256 _withdrawTimeout
	) external {

		// CHECKS
		{
			// Check that the batch nonce is higher than the last nonce for this token
			if (_withdrawNonce <= state_lastWithdrawalNonces[_tokenContract]) {
				revert InvalidWithdrawalNonce({
					newNonce: _withdrawNonce,
					currentNonce: state_lastWithdrawalNonces[_tokenContract]
				});
			}

			// Check that the withdrawal nonce is less than one million nonces forward from the old one
			// this makes it difficult for an attacker to lock out the contract by getting a single
			// bad withdrawal through with uint256 max nonce
			if (_withdrawNonce > state_lastWithdrawalNonces[_tokenContract] + 1000000) {
				revert InvalidWithdrawalNonce({
					newNonce: _withdrawNonce,
					currentNonce: state_lastWithdrawalNonces[_tokenContract]
				});
			}

			// Check that the block height is less than the timeout height
			if (block.number >= _withdrawTimeout) {
				revert WithdrawalTimedOut();
			}
			validateValset(_currentValset, _sigs);

			validateCheckpoint(makeCheckpoint(_currentValset, state_gravityId));

			bytes32 hash = keccak256(
					abi.encode(
						state_gravityId,
						// bytes 32 encoding of "withdrawERC721"
						0x7769746864726177455243373231000000000000000000000000000000000000,
						_tokenIds,
						_destinations,
						_withdrawNonce,
						_withdrawTimeout,
						_tokenContract
					)
				);
			checkValidatorSignatures(
				_currentValset,
				_sigs,
				hash,
				constant_powerThreshold
			);
		
			// ACTIONS

			for (uint256 i = 0; i < _tokenIds.length; i++) {
				ERC721(_tokenContract).safeTransferFrom(address(this), _destinations[i], _tokenIds[i]);
			}
		}
		{
			state_lastERC721EventNonce = state_lastERC721EventNonce + 1;
			emit ERC721WithdrawnEvent(
				_withdrawNonce,
				_tokenContract,
				state_lastERC721EventNonce
			);
		}
	}

	function mintERC721(
		ValsetArgs calldata _currentValset,
		Signature[] calldata _sigs,
		address _tokenContract,
		address _destination,
		uint256 _tokenId,
		string calldata _tokenURI,
		uint256 _mintTimeout
	) external {

		// CHECKS
		{
			// Check that the block height is less than the timeout height
			if (block.number >= _mintTimeout) {
				revert MintTimedOut();
			}
			validateValset(_currentValset, _sigs);

			validateCheckpoint(makeCheckpoint(_currentValset, state_gravityId));

			bytes32 hash = keccak256(
					abi.encode(
						state_gravityId,
						// bytes 32 encoding of "mintERC721"
						0x6d696e7445524337323100000000000000000000000000000000000000000000,
						_tokenId,
						_destination,
						_mintTimeout,
						_tokenContract
					)
				);
			checkValidatorSignatures(
				_currentValset,
				_sigs,
				hash,
				constant_powerThreshold
			);

			// ACTIONS

			CosmosERC721(_tokenContract).safeMint(_destination, _tokenId, _tokenURI);
		}
		{
			state_lastERC721EventNonce = state_lastERC721EventNonce + 1;
			emit ERC721MintedEvent(
				_tokenContract,
				_tokenId,
				_tokenURI,
				state_lastERC721EventNonce
			);
		}
	}

	function deployERC721(
		string calldata _cosmosClass,
		string calldata _name,
		string calldata _symbol
	) external {
		// Deploy an ERC721 and grant ownership to Gravity.sol
		CosmosERC721 erc721 = new CosmosERC721(_name, _symbol);

		// Fire an event to let the Cosmos module know
		state_lastERC721EventNonce = state_lastERC721EventNonce + 1;
		emit ERC721DeployedEvent(
			_cosmosClass,
			address(erc721),
			_name,
			_symbol,
			state_lastERC721EventNonce
		);
	}

	function validateCheckpoint(bytes32 checkpoint)
		private
		view	
	{
		if( IGravity(state_gravitySolAddress).state_lastValsetCheckpoint() != checkpoint) {
			revert IncorrectCheckpoint();
		}
	}

	function makeCheckpoint(ValsetArgs memory _valsetArgs, bytes32 _gravityId)
		private
		pure	
		returns (bytes32)
	{
		// bytes32 encoding of the string "checkpoint"
		bytes32 methodName = 0x636865636b706f696e7400000000000000000000000000000000000000000000;

		bytes32 checkpoint = keccak256(
			abi.encode(
				_gravityId,
				methodName,
				_valsetArgs.valsetNonce,
				_valsetArgs.validators,
				_valsetArgs.powers,
				_valsetArgs.rewardAmount,
				_valsetArgs.rewardToken
			)
		);

		return checkpoint;
	}

	function validateValset(ValsetArgs calldata _valset, Signature[] calldata _sigs) private pure {
		if (
			_valset.validators.length != _valset.powers.length ||
			_valset.validators.length != _sigs.length
		) {
			revert MalformedCurrentValidatorSet();
		}
	}

	function checkValidatorSignatures(
		// The current validator set and their powers
		ValsetArgs calldata _currentValset,
		// The current validator's signatures
		Signature[] calldata _sigs,
		// This is what we are checking they have signed
		bytes32 _theHash,
		uint256 _powerThreshold
	) private pure {
		uint256 cumulativePower = 0;

		for (uint256 i = 0; i < _currentValset.validators.length; i++) {
			// If v is set to 0, this signifies that it was not possible to get a signature from this validator and we skip evaluation
			// (In a valid signature, it is either 27 or 28)
			if (_sigs[i].v != 0) {
				// Check that the current validator has signed off on the hash
				if (!verifySig(_currentValset.validators[i], _theHash, _sigs[i])) {
					revert InvalidSignature();
				}

				// Sum up cumulative power
				cumulativePower = cumulativePower + _currentValset.powers[i];

				// Break early to avoid wasting gas
				if (cumulativePower > _powerThreshold) {
					break;
				}
			}
		}

		// Check that there was enough power
		if (cumulativePower <= _powerThreshold) {
			revert InsufficientPower(cumulativePower, _powerThreshold);
		}
		// Success
	}

	// Utility function to verify geth style signatures
	function verifySig(
		address _signer,
		bytes32 _theHash,
		Signature calldata _sig
	) private pure returns (bool) {
		bytes32 messageDigest = keccak256(
			abi.encodePacked("\x19Ethereum Signed Message:\n32", _theHash)
		);
		return _signer == ECDSA.recover(messageDigest, _sig.v, _sig.r, _sig.s);
	}
}