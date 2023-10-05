use std::convert::TryFrom;

use clarity::Address as EthAddress;
use deep_space::address::Address;
use deep_space::error::CosmosGrpcError;
use deep_space::Contact;
use gravity_proto::gravitynft::QueryNftAttestationsRequest;
use gravity_proto::gravitynft::QueryLastNftEventNonceByAddrRequest;
use gravity_proto::gravitynft::query_client::QueryClient as GravityNftQueryClient;
use gravity_proto::cosmos_sdk_proto::cosmos::nft::v1beta1::QueryOwnerRequest;
use gravity_proto::cosmos_sdk_proto::cosmos::nft::v1beta1::QueryOwnerResponse;
use gravity_proto::cosmos_sdk_proto::cosmos::nft::v1beta1::query_client::QueryClient as NftQueryClient;
use gravity_proto::gravitynft::{NftAttestation, PendingNftIbcAutoForward, QueryPendingNftIbcAutoForwardsRequest};
use gravity_utils::error::GravityError;
use gravity_utils::types::*;
use tonic::transport::Channel;


/// Gets the last event nonce for gravityerc721.sol contract that a given validator has attested to, this lets us
/// catch up with what the current event nonce should be if a oracle is restarted
pub async fn get_last_erc721_event_nonce_for_validator(
    client: &mut GravityNftQueryClient<Channel>,
    address: Address,
    prefix: String,
) -> Result<u64, GravityError> {
    let request = client
        .last_nft_event_nonce_by_addr(QueryLastNftEventNonceByAddrRequest {
            address: address.to_bech32(prefix).unwrap(),
        })
        .await?;
    Ok(request.into_inner().last_nft_event_nonce)
}

pub async fn get_erc721_attestations(
    client: &mut GravityNftQueryClient<Channel>,
    limit: Option<u64>,
) -> Result<Vec<NftAttestation>, GravityError> {
    let request = client
        .get_nft_attestations(QueryNftAttestationsRequest {
            limit: limit.unwrap_or(1000u64),
            order_by: String::new(),
            claim_type: String::new(),
            nonce: 0,
            height: 0,
            use_v1_key: false,
        })
        .await?;
    let attestations = request.into_inner().nft_attestations;
    Ok(attestations)
}

pub async fn get_nft_owner(
    client: &mut NftQueryClient<Channel>,
    class_id: String,
    token_id: String,
) -> Result<QueryOwnerResponse, GravityError> {
    let request = client
        .owner(QueryOwnerRequest {
            class_id: class_id,
            id: token_id,
        })
        .await?;
    Ok(request.into_inner())
}

/// Queries the Gravity chain for Pending NFT Ibc Auto Forwards, returning an empty vec if there is an error
pub async fn get_all_pending_nft_ibc_auto_forwards(
    grpc_client: &mut GravityNftQueryClient<Channel>,
) -> Vec<PendingNftIbcAutoForward> {
    let pending_forwards = grpc_client
        .get_pending_nft_ibc_auto_forwards(QueryPendingNftIbcAutoForwardsRequest { limit: 0 })
        .await;
    if let Err(status) = pending_forwards {
        // don't print errors during the upgrade test, which involves running
        // a newer orchestrator against an older chain due to current design limitations.
        if !status.message().contains("unknown method") {
            warn!(
                "Received an error when querying for pending ibc auto forwards: {}",
                status.message()
            );
        }
        return vec![];
    }

    let pending_forwards = pending_forwards.unwrap();
    pending_forwards.into_inner().pending_nft_ibc_auto_forwards
}
