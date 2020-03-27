export const sidebar = state => state.app.sidebar
export const device = state => state.app.device
export const token = state => state.user.token
export const avatar = state => state.user.avatar
export const name = state => state.user.name

// blockchain
export const blockchain = state => state.blockchain
export const blocks = state => state.blockchain.blocks
export const config = state => state.config
export const consensusState = state => state.blockchain.consensusState
export const dumpConsensusState = state => state.blockchain.dumpConsensusState
export const materials = state => state.materials
export const nodes = state => state.blockchain.nodes
export const roundStep = state => state.blockchain.roundStep
export const validators = state => state.blockchain.validators
export const bc = state => state.blockchain

// TODO rename to lastBlockHeader? such as it doesn't include block txs
export const latestBlock = (state, getters) => {
  let { blocks } = getters
  if (blocks && blocks.length >= 1) {
    return blocks[0].header
  } else {
    return {
      height: 0,
      time: "",
      last_commit_hash: "",
      num_txs: 0, // txs in this block
      total_txs: 0 // total txs in blockchain at the moment of this block
    }
  }
}

export const totalBlocks = (state, getters) => parseInt(getters.latestBlock.height)
export const totalTxs = (state, getters) => parseInt(getters.latestBlock.total_txs)

// Abount account
export const account = state => state.account.tem_new_account
