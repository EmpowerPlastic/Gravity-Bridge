use clarity::{Address, Uint256};
use cosmos_gravitynft::utils::get_last_nft_event_nonce_with_retry;
use deep_space::address::Address as CosmosAddress;
use gravity_proto::gravitynft::query_client::QueryClient as GravityNftQueryClient;
use gravity_utils::get_with_retry::RETRY_TIME;
use gravity_utils::types::{event_signatures::*, EthereumEvent, SendERC721ToCosmosEvent};

use metrics_exporter::metrics_errors_counter;
use tokio::time::sleep as delay_for;
use tonic::transport::Channel;
use web30::client::Web3;

use crate::ethereum_event_watcher::get_latest_safe_block;

/// This is roughly the maximum number of blocks a reasonable Ethereum node
/// can search in a single request before it starts timing out or behaving badly
pub const BLOCKS_TO_SEARCH: u128 = 5_000u128;

/// This function retrieves the last event nonce this oracle has relayed to Cosmos
/// it then uses the Ethereum indexes to determine what block the last entry
pub async fn get_last_checked_block_nft(
    grpc_client: GravityNftQueryClient<Channel>,
    our_cosmos_address: CosmosAddress,
    prefix: String,
    gravity_contract_address: Address,
    web3: &Web3,
) -> Uint256 {
    let mut grpc_client = grpc_client;

    let latest_block = get_latest_safe_block(web3).await;
    let mut last_event_nonce =
            get_last_nft_event_nonce_with_retry(&mut grpc_client, our_cosmos_address, prefix)
                .await
                .into();

    // zero indicates this oracle has never submitted an event before since there is no
    // zero event nonce (it's pre-incremented in the solidity contract) we have to go
    // and look for event nonce one.
    if last_event_nonce == 0u8.into() {
        last_event_nonce = 1u8.into();
    }

    let current_block: Uint256 = latest_block.clone();

    while current_block.clone() > 0u8.into() {
        info!(
            "Oracle is resyncing, looking back into the history to find our last event nonce {}, on block {}",
            last_event_nonce, current_block
        );
        let end_search = if current_block.clone() < BLOCKS_TO_SEARCH.into() {
            0u8.into()
        } else {
            current_block.clone() - BLOCKS_TO_SEARCH.into()
        };
            let send_erc721_to_cosmos_events = web3
                .check_for_events(
                    end_search.clone(),
                    Some(current_block.clone()),
                    vec![gravity_contract_address],
                    vec![SENT_ERC721_TO_COSMOS_EVENT_SIG],
                )
                .await;
            if send_erc721_to_cosmos_events.is_err() {
                error!("Failed to get blockchain events while resyncing, is your Eth node working? If you see only one of these it's fine",);
                delay_for(RETRY_TIME).await;
                metrics_errors_counter(1, "Failed to get blockchain events while resyncing");
                continue;
            }
            let send_erc721_to_cosmos_events = send_erc721_to_cosmos_events.unwrap();
            for event in send_erc721_to_cosmos_events {
                match SendERC721ToCosmosEvent::from_log(&event) {
                    Ok(send) => {
                        debug!(
                            "{} send event nonce {} last event nonce",
                            send.event_nonce, last_event_nonce
                        );
                        if upcast(send.event_nonce) == last_event_nonce
                            && event.block_number.is_some()
                        {
                            return event.block_number.unwrap();
                        }
                    }
                    Err(e) => {
                        error!("Got SendERC721ToCosmos event that we can't parse {}", e);
                        metrics_errors_counter(
                            3,
                            "Got SendERC721ToCosmos event that we can't parse",
                        );
                    }
                }
            }
            let gravityerc721_deployed_events = web3
                .check_for_events(
                    end_search.clone(),
                    Some(current_block.clone()),
                    vec![gravity_contract_address],
                    vec![GRAVITYERC721_DEPLOYED_EVENT_SIG],
                )
                .await;
            if gravityerc721_deployed_events.is_err() {
                error!("Failed to get blockchain events while resyncing, is your Eth node working? If you see only one of these it's fine",);
                delay_for(RETRY_TIME).await;
                metrics_errors_counter(1, "Failed to get blockchain events while resyncing");
                continue;
            }
            let gravityerc721_deployed_events = gravityerc721_deployed_events.unwrap();
            if gravityerc721_deployed_events.len() > 1 {
                error!("More than one GravityERC721DeployedEvent found! This should never happen!");
                metrics_errors_counter(
                    2,
                    "More than one GravityERC721DeployedEvent found! This should never happen!",
                );
            }
            if gravityerc721_deployed_events.len() == 1 {
                let gravityerc721_deploy_height = gravityerc721_deployed_events[0]
                    .clone()
                    .block_number
                    .unwrap();

                info!(
                    "Found GravityERC721DeployedEvent at block {}",
                    gravityerc721_deploy_height
                );
                return gravityerc721_deploy_height;
            }
    }

    // we should exit above when we find the zero valset, if we have the wrong contract address through we could be at it a while as we go over
    // the entire history to 'prove' it.
    panic!("You have reached the end of block history without finding the Gravity contract deploy event! You must have the wrong contract address!");
}

fn upcast(input: u64) -> Uint256 {
    input.into()
}
