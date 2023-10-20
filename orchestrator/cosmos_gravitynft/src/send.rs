
use deep_space::address::Address as CosmosAddress;
use deep_space::error::CosmosGrpcError;
use deep_space::private_key::PrivateKey;
use deep_space::Contact;
use deep_space::Msg;
use deep_space::{coin::Coin};
use gravity_proto::cosmos_sdk_proto::cosmos::base::abci::v1beta1::TxResponse;
use gravity_utils::types::*;
use std::{time::Duration};

pub const TIMEOUT: Duration = Duration::from_secs(60);

pub async fn send_erc721_claims(
    contact: &Contact,
    our_cosmos_key: impl PrivateKey,
    erc721_deposits: Vec<SendERC721ToCosmosEvent>,
    fee: Coin,
) -> Result<TxResponse, CosmosGrpcError> {
    let our_cosmos_address = our_cosmos_key.to_address(&contact.get_prefix()).unwrap();

    let mut erc721_deposit_nonces_msgs: Vec<(u64, Msg)> =
        create_claim_msgs(erc721_deposits, our_cosmos_address);

    erc721_deposit_nonces_msgs.sort_unstable_by(|a, b| a.0.cmp(&b.0));

    const MAX_ORACLE_MESSAGES: usize = 1000;
    while erc721_deposit_nonces_msgs.len() > MAX_ORACLE_MESSAGES {
        // pops messages off of the end
        erc721_deposit_nonces_msgs.pop();
    }
    let mut msgs = Vec::new();
    for i in erc721_deposit_nonces_msgs {
        msgs.push(i.1);
    }
    contact
        .send_message(&msgs, None, &[fee], Some(TIMEOUT), our_cosmos_key)
        .await
}

/// Creates the `Msg`s needed for `orchestrator` to attest to `events`
/// Returns a Vec of (event_nonce: u64, Msg), which will contain one (nonce, msg) per event
fn create_claim_msgs(
    events: Vec<impl EthereumEvent>,
    orchestrator: CosmosAddress,
) -> Vec<(u64, Msg)> {
    let mut msgs = vec![];
    for event in events {
        // Create msg
        msgs.push((event.get_event_nonce(), event.to_claim_msg(orchestrator)));
    }
    msgs
}

