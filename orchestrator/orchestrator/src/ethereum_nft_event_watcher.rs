//! Ethereum Event watcher watches for events such as a deposit to the Gravity Ethereum contract or a validator set update
//! or a transaction batch update. It then responds to these events by performing actions on the Cosmos chain if required

use clarity::{Address as EthAddress, Uint256};
use cosmos_gravitynft::{
    query::get_last_erc721_event_nonce_for_validator,
    send::send_erc721_claims,
};
use deep_space::Contact;
use deep_space::{
    coin::Coin,
    private_key::{CosmosPrivateKey, PrivateKey},
};
use gravity_proto::gravitynft::query_client::QueryClient as GravityNftQueryClient;
use gravity_utils::{error::GravityError, types::{event_signatures::*, EthereumEvent, SendERC721ToCosmosEvent}};

use metrics_exporter::metrics_errors_counter;
use tonic::transport::Channel;
use web30::client::Web3;
use web30::jsonrpc::error::Web3Error;

use super::ethereum_event_watcher::{CheckedNonces, get_latest_safe_block};

use crate::oracle_resync::{BLOCKS_TO_SEARCH};



#[allow(clippy::too_many_arguments)]
pub async fn check_for_events_nft(
    web3: &Web3,
    contact: &Contact,
    grpc_client: &mut GravityNftQueryClient<Channel>,
    gravity_contract_address: EthAddress,
    our_private_key: CosmosPrivateKey,
    fee: Coin,
    starting_block: Uint256,
) -> Result<CheckedNonces, GravityError> {
    let our_cosmos_address = our_private_key.to_address(&contact.get_prefix()).unwrap();
    let latest_block = get_latest_safe_block(web3).await;
    trace!(
        "Checking for events starting {} safe {}",
        starting_block,
        latest_block
    );

    // if the latest block is more than BLOCKS_TO_SEARCH ahead do not search the full history
    // comparison only to prevent panic on underflow.
    let latest_block = if latest_block > starting_block
        && latest_block.clone() - starting_block.clone() > BLOCKS_TO_SEARCH.into()
    {
        starting_block.clone() + BLOCKS_TO_SEARCH.into()
    } else {
        latest_block
    };

        let erc721_deposits = web3
            .check_for_events(
                starting_block.clone(),
                Some(latest_block.clone()),
                vec![gravity_contract_address],
                vec![SENT_ERC721_TO_COSMOS_EVENT_SIG],
            )
            .await;
        trace!("ERC721 deposits {:?}", erc721_deposits);

        if let (Ok(erc721_deposits),) = (erc721_deposits,) {
            let erc721_deposits = SendERC721ToCosmosEvent::from_logs(&erc721_deposits)?;
            trace!("parsed erc721 deposits {:?}", erc721_deposits);

            // note that starting block overlaps with our last checked block, because we have to deal with
            // the possibility that the relayer was killed after relaying only one of multiple events in a single
            // block, so we also need this routine so make sure we don't send in the first event in this hypothetical
            // multi event block again. In theory we only send all events for every block and that will pass of fail
            // atomicly but lets not take that risk.
            let last_event_nonce = get_last_erc721_event_nonce_for_validator(
                grpc_client,
                our_cosmos_address,
                contact.get_prefix(),
            )
            .await?;
            let erc721_deposits =
                SendERC721ToCosmosEvent::filter_by_event_nonce(last_event_nonce, &erc721_deposits);
            if !erc721_deposits.is_empty() {
                info!(
                "Oracle observed erc721 deposit with sender {}, destination {:?}, token id {}, token uri {}, and event nonce {}",
                erc721_deposits[0].sender, erc721_deposits[0].validated_destination, erc721_deposits[0].token_id, erc721_deposits[0].token_uri, erc721_deposits[0].event_nonce
            );
                let res =
                    send_erc721_claims(contact, our_private_key, erc721_deposits.clone(), fee)
                        .await;
                if res.is_err() {
                    error!("Failed to process GravityERC721 claims");
                    metrics_errors_counter(2, "Failed process erc721 claims");
                    return Err(GravityError::CosmosGrpcError(res.unwrap_err()));
                }

                let new_event_nonce = get_last_erc721_event_nonce_for_validator(
                    grpc_client,
                    our_cosmos_address,
                    contact.get_prefix(),
                )
                .await?;

                info!("Current gravityerc721 event nonce is {}", new_event_nonce);

                // since we can't actually trust that the above txresponse is correct we have to check here
                // we may be able to trust the tx response post grpc
                if new_event_nonce == last_event_nonce {
                    return Err(GravityError::InvalidBridgeStateError(
                format!("GravityERC721 claims did not process, trying to update but still on {}, trying again in a moment, check txhash {:?} for errors", last_event_nonce, res),
            ));
                } else {
                    info!(
                        "GravityERC721 claims processed, new nonce {}",
                        new_event_nonce
                    );
                }

                // find the eth block for our newest event nonce
                let erc721_deposits =
                    SendERC721ToCosmosEvent::get_block_for_nonce(new_event_nonce, &erc721_deposits);

                Ok(CheckedNonces {
                    block_number: erc721_deposits.unwrap(),
                    event_nonce: new_event_nonce.into(),
                })
            } else {
                // no changes
                Ok(CheckedNonces {
                    block_number: latest_block,
                    event_nonce: last_event_nonce.into(),
                })
            }
        } else {
            error!("Failed to get events");
            metrics_errors_counter(1, "Failed to get events");
            Err(GravityError::EthereumRestError(Web3Error::BadResponse(
                "Failed to get logs!".to_string(),
            )))
        }
}
