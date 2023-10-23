use crate::query::{get_last_erc721_event_nonce_for_validator};
use deep_space::{Address as CosmosAddress};
use gravity_proto::gravitynft::query_client::QueryClient as GravityNftQueryClient;
use gravity_utils::get_with_retry::RETRY_TIME;
use tokio::time::sleep;
use tonic::transport::Channel;


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
