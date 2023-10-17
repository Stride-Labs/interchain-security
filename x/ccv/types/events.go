package types

// CCV events
const (
	EventTypeTimeout                    = "timeout"
	EventTypePacket                     = "ccv_packet"
	EventTypeChannelEstablished         = "channel_established"
	EventTypeFeeTransferChannelOpened   = "fee_transfer_channel_opened"
	EventTypeConsumerClientCreated      = "consumer_client_created"
	EventTypeAssignConsumerKey          = "assign_consumer_key"
	EventTypeSubmitConsumerMisbehaviour = "submit_consumer_misbehaviour"
	EventTypeSubmitConsumerDoubleVoting = "submit_consumer_double_voting"
	EventTypeExecuteConsumerChainSlash  = "execute_consumer_chain_slash"
	EventTypeFeeDistribution            = "fee_distribution"
	EventTypeConsumerSlashRequest       = "consumer_slash_request"
	EventTypeVSCMatured                 = "vsc_matured"
	EventTypeAddConsumerRewardDenom     = "add_consumer_reward_denom"
	EventTypeRemoveConsumerRewardDenom  = "remove_consumer_reward_denom"

	AttributeKeyAckSuccess = "success"
	AttributeKeyAck        = "acknowledgement"
	AttributeKeyAckError   = "error"

	AttributeChainID                  = "chain_id"
	AttributeValidatorAddress         = "validator_address"
	AttributeValidatorConsumerAddress = "validator_consumer_address"
	AttributeInfractionType           = "infraction_type"
	AttributeInfractionHeight         = "infraction_height"
	AttributeConsumerHeight           = "consumer_height"
	AttributeValSetUpdateID           = "valset_update_id"
	AttributeTimestamp                = "timestamp"
	AttributeInitialHeight            = "initial_height"
	AttributeInitializationTimeout    = "initialization_timeout"
	AttributeTrustingPeriod           = "trusting_period"
	AttributeUnbondingPeriod          = "unbonding_period"
	AttributeProviderValidatorAddress = "provider_validator_address"
	AttributeConsumerConsensusPubKey  = "consumer_consensus_pub_key"
	AttributeSubmitterAddress         = "submitter_address"
	AttributeConsumerMisbehaviour     = "consumer_misbehaviour"
	AttributeMisbehaviourClientId     = "misbehaviour_client_id"
	AttributeMisbehaviourHeight1      = "misbehaviour_height_1"
	AttributeMisbehaviourHeight2      = "misbehaviour_height_2"
	AttributeConsumerDoubleVoting     = "consumer_double_voting"

	AttributeDistributionCurrentHeight = "current_distribution_height"
	//#nosec G101 -- (false positive) this is not a hardcoded credential
	AttributeDistributionNextHeight = "next_distribution_height"
	AttributeDistributionFraction   = "distribution_fraction"
	AttributeDistributionTotal      = "total"
	AttributeDistributionToProvider = "provider_amount"

	AttributeConsumerRewardDenom = "consumer_reward_denom"
)
