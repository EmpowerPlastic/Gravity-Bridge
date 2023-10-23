//! This file parses the Gravity contract ethereum events. Note that there is no Ethereum ABI unpacking implementation. Instead each event
//! is parsed directly from it's binary representation. This is technical debt within this implementation. It's quite easy to parse any
//! individual event manually but a generic decoder can be quite challenging to implement. A proper implementation would probably closely
//! mirror Serde and perhaps even become a serde crate for Ethereum ABI decoding
//! For now reference the ABI encoding document here https://docs.soliditylang.org/en/v0.8.3/abi-spec.html

// TODO this file needs static assertions that prevent it from compiling on 16 bit systems.
// we assume a system bit width of at least 32

use crate::error::GravityError;
use crate::num_conversion::downcast_uint256;
use clarity::Address as EthAddress;
use deep_space::utils::bytes_to_hex_str;
use deep_space::{Address as CosmosAddress, Address, Msg};
use gravity_proto::gravitynft::{
    MsgSendNftToCosmosClaim
};
use num256::Uint256;
use web30::types::Log;
use super::ethereum_events::EthereumEvent;

const ONE_MEGABYTE: usize = 1000usize.pow(3);

// gravitynft msg type urls
pub const MSG_SEND_NFT_TO_COSMOS_CLAIM_TYPE_URL: &str = "/gravitynft.v1.MsgSendNFTToCosmosClaim";

/// Used to limit the length of variable length user provided inputs like
/// ERC20 names and deposit destination strings
// const ONE_MEGABYTE: usize = 1000usize.pow(3);

/// A parsed struct representing the Ethereum event fired when someone makes an erc721 deposit
/// on the GravityERC721 contract
#[derive(Serialize, Deserialize, Debug, Clone, Eq, PartialEq, Hash)]
pub struct SendERC721ToCosmosEvent {
    /// The token contract address for the deposit
    pub erc721: EthAddress,
    /// The Ethereum Sender
    pub sender: EthAddress,
    /// The Cosmos destination, this is a raw value from the Ethereum contract
    /// and therefore could be provided by an attacker. If the string is valid
    /// utf-8 it will be included here, if it is invalid utf8 we will provide
    /// an empty string. Values over 1mb of text are not permitted and will also
    /// be presented as empty
    pub destination: String,
    /// the validated destination is the destination string parsed and interpreted
    /// as a valid Bech32 Cosmos address, if this is not possible the value is none
    pub validated_destination: Option<CosmosAddress>,
    /// The id of ERC721 token that is being sent
    pub token_id: Uint256,
    /// URI of the ERC721 token
    pub token_uri: String,
    /// The transaction's nonce, used to make sure there can be no accidental duplication
    pub event_nonce: u64,
    /// The block height this event occurred at
    pub block_height: Uint256,
}

/// struct for holding the data encoded fields
/// of a send erc721 to Cosmos event for unit testing
#[derive(Eq, PartialEq, Debug)]
struct SendERC721ToCosmosEventData {
    /// The Cosmos destination, None for an invalid deposit address
    pub destination: String,
    /// The id of ERC721 token that is being sent
    pub token_id: Uint256,
    /// The transaction's nonce, used to make sure there can be no accidental duplication
    pub event_nonce: Uint256,
    /// URI of the ERC721 token
    pub token_uri: String,
}

impl SendERC721ToCosmosEvent {
    fn decode_data_bytes(data: &[u8]) -> Result<SendERC721ToCosmosEventData, GravityError> {
        if data.len() < 4 * 32 {
            return Err(GravityError::InvalidEventLogError(
                "too short for SendERC721ToCosmosEventData".to_string(),
            ));
        }

        let token_id = Uint256::from_be_bytes(&data[32..64]);
        let event_nonce = Uint256::from_be_bytes(&data[64..96]);

        // discard words three and four which contain the data type and length
        let destination_str_len_start = 4 * 32;
        let destination_str_len_end = 5 * 32;
        let destination_str_len =
            Uint256::from_be_bytes(&data[destination_str_len_start..destination_str_len_end]);

        if destination_str_len > u32::MAX.into() {
            return Err(GravityError::InvalidEventLogError(
                "denom length overflow, probably incorrect parsing".to_string(),
            ));
        }
        let destination_str_len: usize = destination_str_len.to_string().parse().unwrap();

        let destination_str_start = 5 * 32;
        let destination_str_end = destination_str_start + destination_str_len;

        if data.len() < destination_str_end {
            return Err(GravityError::InvalidEventLogError(
                "Incorrect length for dynamic data".to_string(),
            ));
        }

        let destination = &data[destination_str_start..destination_str_end];

        let dest = String::from_utf8(destination.to_vec());

        let mut destination_str_in_bytes = destination_str_len / 32;
        if destination_str_len % 32 != 0 {
            destination_str_in_bytes = destination_str_in_bytes + 1;
        }
        let token_uri_str_len_start = (5 + destination_str_in_bytes) * 32;
        let token_uri_str_len_end = (6 + destination_str_in_bytes) * 32;
        let token_uri_str_len =
            Uint256::from_be_bytes(&data[token_uri_str_len_start..token_uri_str_len_end]);

        if token_uri_str_len > u32::MAX.into() {
            return Err(GravityError::InvalidEventLogError(
                "denom length overflow, probably incorrect parsing".to_string(),
            ));
        }
        let token_uri_str_len: usize = token_uri_str_len.to_string().parse().unwrap();
        let token_uri: String;

        if token_uri_str_len > 0 {
            let token_uri_str_start = (6 + destination_str_in_bytes) * 32;
            let token_uri_str_end = token_uri_str_start + token_uri_str_len;

            if data.len() < token_uri_str_end {
                return Err(GravityError::InvalidEventLogError(
                    "Incorrect length for dynamic data".to_string(),
                ));
            }

            let uri = &data[token_uri_str_start..token_uri_str_end];

            let uri = String::from_utf8(uri.to_vec());

            if uri.is_err() {
                return Err(GravityError::InvalidEventLogError(
                    "Token URI parsing error, probably incorrect parsing".to_string(),
                ));
            }

            token_uri = uri.unwrap();
        } else {
            token_uri = String::new();
        }

        if dest.is_err() {
            if destination.len() < 1000 {
                warn!("Event nonce {} sends token to {} which is invalid utf-8, this token will be allocated to the community pool", event_nonce, bytes_to_hex_str(destination));
            } else {
                warn!("Event nonce {} sends token to a destination that is invalid utf-8, this token will be allocated to the community pool", event_nonce);
            }
            return Ok(SendERC721ToCosmosEventData {
                destination: String::new(),
                event_nonce,
                token_id,
                token_uri,
            });
        }
        // whitespace can not be a valid part of a bech32 address, so we can safely trim it
        let dest = dest.unwrap().trim().to_string();

        if dest.as_bytes().len() > ONE_MEGABYTE {
            warn!("Event nonce {} sends token to a destination that exceeds the length limit, this token will be allocated to the community pool", event_nonce);
            Ok(SendERC721ToCosmosEventData {
                destination: String::new(),
                event_nonce,
                token_id,
                token_uri,
            })
        } else {
            Ok(SendERC721ToCosmosEventData {
                destination: dest,
                event_nonce,
                token_id,
                token_uri,
            })
        }
    }
}
impl EthereumEvent for SendERC721ToCosmosEvent {
    fn get_block_height(&self) -> u64 {
        downcast_uint256(self.block_height.clone()).unwrap()
    }

    fn get_event_nonce(&self) -> u64 {
        self.event_nonce
    }

    fn from_log(input: &Log) -> Result<SendERC721ToCosmosEvent, GravityError> {
        let topics = (input.topics.get(1), input.topics.get(2));
        if let (Some(erc721_data), Some(sender_data)) = topics {
            let erc721 = EthAddress::from_slice(&erc721_data[12..32])?;
            let sender = EthAddress::from_slice(&sender_data[12..32])?;
            let block_height = if let Some(bn) = input.block_number.clone() {
                if bn > u64::MAX.into() {
                    return Err(GravityError::InvalidEventLogError(
                        "Block height overflow! probably incorrect parsing".to_string(),
                    ));
                } else {
                    bn
                }
            } else {
                return Err(GravityError::InvalidEventLogError(
                    "Log does not have block number, we only search logs already in blocks?"
                        .to_string(),
                ));
            };

            let data = SendERC721ToCosmosEvent::decode_data_bytes(&input.data)?;
            if data.event_nonce > u64::MAX.into() || block_height > u64::MAX.into() {
                Err(GravityError::InvalidEventLogError(
                    "Event nonce overflow, probably incorrect parsing".to_string(),
                ))
            } else {
                let event_nonce: u64 = data.event_nonce.to_string().parse().unwrap();
                let validated_destination = match data.destination.parse() {
                    Ok(v) => Some(v),
                    Err(_) => {
                        if data.destination.len() < 1000 {
                            warn!("Event nonce {} sends token to {} which is invalid bech32, this token will be allocated to the community pool", event_nonce, data.destination);
                        } else {
                            warn!("Event nonce {} sends token to a destination which is invalid bech32, this token will be allocated to the community pool", event_nonce);
                        }
                        None
                    }
                };
                Ok(SendERC721ToCosmosEvent {
                    erc721,
                    sender,
                    destination: data.destination,
                    validated_destination,
                    token_id: data.token_id,
                    token_uri: data.token_uri,
                    event_nonce,
                    block_height,
                })
            }
        } else {
            Err(GravityError::InvalidEventLogError(
                "Too few topics".to_string(),
            ))
        }
    }

    fn from_logs(input: &[Log]) -> Result<Vec<SendERC721ToCosmosEvent>, GravityError> {
        let mut res = Vec::new();
        for item in input {
            res.push(SendERC721ToCosmosEvent::from_log(item)?);
        }
        Ok(res)
    }
    /// returns all values in the array with event nonces greater
    /// than the provided value
    fn filter_by_event_nonce(event_nonce: u64, input: &[Self]) -> Vec<Self> {
        let mut ret = Vec::new();
        for item in input {
            if item.event_nonce > event_nonce {
                ret.push(item.clone())
            }
        }
        ret
    }

    // gets the Ethereum block for the given nonce
    fn get_block_for_nonce(event_nonce: u64, input: &[Self]) -> Option<Uint256> {
        for item in input {
            if item.event_nonce == event_nonce {
                return Some(item.block_height.clone());
            }
        }
        None
    }

    fn to_claim_msg(self, orchestrator: Address) -> Msg {
        let claim = MsgSendNftToCosmosClaim {
            event_nonce: self.event_nonce,
            eth_block_height: self.get_block_height(),
            token_contract: self.erc721.to_string(),
            token_id: self.token_id.to_string(),
            token_uri: self.token_uri,
            cosmos_receiver: self.destination,
            ethereum_sender: self.sender.to_string(),
            orchestrator: orchestrator.to_string(),
        };
        Msg::new(MSG_SEND_NFT_TO_COSMOS_CLAIM_TYPE_URL, claim)
    }
}
