use crate::query::{get_last_erc721_event_nonce_for_validator};
use deep_space::client::ChainStatus;
use deep_space::error::CosmosGrpcError;
use deep_space::utils::encode_any;
use deep_space::{Address as CosmosAddress, Contact};
use gravity_proto::gravity::query_client::QueryClient as GravityQueryClient;
use gravity_proto::gravitynft::query_client::QueryClient as GravityNftQueryClient;
use gravity_proto::gravity::OutgoingLogicCall as ProtoLogicCall;
use gravity_proto::gravity::OutgoingTxBatch as ProtoBatch;
use gravity_proto::gravity::Valset as ProtoValset;
use gravity_utils::get_with_retry::RETRY_TIME;
use gravity_utils::types::LogicCall;
use gravity_utils::types::TransactionBatch;
use gravity_utils::types::Valset;
use num256::Uint256;
use prost_types::Any;
use std::convert::TryFrom;
use std::ops::Mul;
use tokio::time::sleep;
use tonic::metadata::AsciiMetadataValue;
use tonic::transport::Channel;
use tonic::{IntoRequest, Request};

/// gets the Cosmos last nft event nonce, no matter how long it takes.
pub async fn get_last_nft_event_nonce_with_retry(
    client: &mut GravityNftQueryClient<Channel>,
    our_cosmos_address: CosmosAddress,
    prefix: String,
) -> u64 {
    let mut res =
        get_last_erc721_event_nonce_for_validator(client, our_cosmos_address, prefix.clone()).await;
    while res.is_err() {
        error!(
            "Failed to get last event nonce, is the Cosmos GRPC working? {:?}",
            res
        );
        sleep(RETRY_TIME).await;
        res = get_last_erc721_event_nonce_for_validator(client, our_cosmos_address, prefix.clone())
            .await;
    }
    res.unwrap()
}
