package main

import (
	"blockchain/blockchain"
	"blockchain/chaincfg"
	"blockchain/mining"
	miner "blockchain/mining"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

// server provides a bitcoin server for handling communications to and from
// bitcoin peers.
type server struct {
	// The following variables must only be used atomically.
	// Putting the uint64s first makes them 64-bit aligned for 32-bit systems.
	bytesReceived uint64 // Total bytes received from all peers since start.
	bytesSent     uint64 // Total bytes sent by all peers since start.
	started       int32
	shutdown      int32
	shutdownSched int32
	startupTime   int64

	chainParams *chaincfg.Params
	//	addrManager          *addrmgr.AddrManager
	//connManager          *connmgr.ConnManager
	//sigCache             *txscript.SigCache
	//hashCache            *txscript.HashCache
	//rpcServer            *rpcServer
	//syncManager          *netsync.SyncManager
	chain *blockchain.Blockchain
	//txMemPool            *mempool.TxPool
	Miner *miner.Miner
	//modifyRebroadcastInv chan interface{}
	//newPeers          chan *serverPeer
	//donePeers         chan *serverPeer
	//banPeers          chan *serverPeer
	//query             chan interface{}
	//relayInv          chan relayMsg
	//broadcast         chan broadcastMsg
	//peerHeightsUpdate chan updatePeerHeightsMsg
	wg   sync.WaitGroup
	quit chan struct{}
	//nat        NAT
	//db         database.DB
	//timeSource blockchain.MedianTimeSource
	//services wire.ServiceFlag
	// The following fields are used for optional indexes.  They will be nil
	// if the associated index is not enabled.  These fields are set during
	// initial creation of the server and never changed afterwards, so they
	// do not need to be protected for concurrent access.
	//txIndex   *indexers.TxIndex
	//addrIndex *indexers.AddrIndex
	//cfIndex *indexers.CfIndex
	// The fee estimator keeps track of how long transactions are left in
	// the mempool before they are mined into blocks.
	//feeEstimator *mempool.FeeEstimator
	// cfCheckptCaches stores a cached slice of filter headers for cfcheckpt
	// messages for each filter type.
	//	cfCheckptCaches    map[wire.FilterType][]cfHeaderKV
	//cfCheckptCachesMtx sync.RWMutex
}

// newServer returns a new server configured to listen on addr for the
//  network type specified by chainParams.  Use start to begin accepting
// connections from peers.
func newServer(chainParams *chaincfg.Params, interrupt <-chan struct{}) (*server, error) {
	//	db database.DB, chainParams *chaincfg.Params,
	//services := defaultServices
	//if cfg.NoPeerBloomFilters {
	//	services &^= wire.SFNodeBloom
	//	}
	//	if cfg.NoCFilters {
	//		services &^= wire.SFNodeCF
	//	}
	//	amgr := addrmgr.New(cfg.DataDir, btcdLookup)
	var listeners []net.Listener
	_ = listeners
	//	var nat NAT
	if !cfg.DisableListen {
		//var err error
		fmt.Println("listening")
		//	listeners, nat, err = initListeners(amgr, listenAddrs, services)
		//if err != nil {
		//	return nil, err
		//}
		//if len(listeners) == 0 {
		//	return nil, errors.New("no valid listen address")
		//}
	}
	//	if len(agentBlacklist) > 0 {
	//		srvrLog.Infof("User-agent blacklist %s", agentBlacklist)
	//	}
	//	if len(agentWhitelist) > 0 {
	//		srvrLog.Infof("User-agent whitelist %s", agentWhitelist)
	//	}
	s := server{
		//		chainParams:          chainParams,
		//		addrManager:          amgr,
		//		newPeers:             make(chan *serverPeer, cfg.MaxPeers),
		//		donePeers:            make(chan *serverPeer, cfg.MaxPeers),
		//		banPeers:             make(chan *serverPeer, cfg.MaxPeers),
		//		query:                make(chan interface{}),
		//		relayInv:             make(chan relayMsg, cfg.MaxPeers),
		//		broadcast:            make(chan broadcastMsg, cfg.MaxPeers),
		quit: make(chan struct{}),
		//		modifyRebroadcastInv: make(chan interface{}),
		//		peerHeightsUpdate:    make(chan updatePeerHeightsMsg),
		//		nat:                  nat,
		//		db:                   db,
		//		timeSource:           blockchain.NewMedianTime(),
		//services: services,
		//		sigCache:             txscript.NewSigCache(cfg.SigCacheMaxSize),
		//		hashCache:            txscript.NewHashCache(cfg.SigCacheMaxSize),
		//		cfCheckptCaches:      make(map[wire.FilterType][]cfHeaderKV),
		//		agentBlacklist:       agentBlacklist,
		//		agentWhitelist:       agentWhitelist,
	}
	// Create the transaction and address indexes if needed.
	//
	// CAUTION: the txindex needs to be first in the indexes array because
	// the addrindex uses data from the txindex during catchup.  If the
	// addrindex is run first, it may not have the transactions from the
	// current block indexed.
	//	var indexes []indexers.Indexer
	//	if cfg.TxIndex || cfg.AddrIndex {
	// Enable transaction index if address index is enabled since it
	//		// requires it.
	//		if !cfg.TxIndex {
	//			indxLog.Infof("Transaction index enabled because it " +
	//				"is required by the address index")
	//			cfg.TxIndex = true
	//		} else {
	//			indxLog.Info("Transaction index is enabled")
	//		}
	//		s.txIndex = indexers.NewTxIndex(db)
	//		indexes = append(indexes, s.txIndex)
	//	}
	//	if cfg.AddrIndex {
	//		indxLog.Info("Address index is enabled")
	//		s.addrIndex = indexers.NewAddrIndex(db, chainParams)
	//		indexes = append(indexes, s.addrIndex)
	//	}
	//	if !cfg.NoCFilters {
	//		indxLog.Info("Committed filter index is enabled")
	//		s.cfIndex = indexers.NewCfIndex(db, chainParams)
	//		indexes = append(indexes, s.cfIndex)
	//	}

	// Create an index manager if any of the optional indexes are enabled.
	//	var indexManager blockchain.IndexManager
	//	if len(indexes) > 0 {
	//		indexManager = indexers.NewManager(db, indexes)
	//	}

	// Merge given checkpoints with the default ones unless they are disabled.
	//	var checkpoints []chaincfg.Checkpoint
	//	if !cfg.DisableCheckpoints {
	//		checkpoints = mergeCheckpoints(s.chainParams.Checkpoints, cfg.addCheckpoints)
	//	}

	// Create a new block chain instance with the appropriate configuration.
	var err error
	s.chain, err = blockchain.New(&blockchain.Config{
		//		DB:           s.db,
		//		Interrupt:    interrupt,
		ChainParams: s.chainParams,
		//		Checkpoints:  checkpoints,
		//		TimeSource:   s.timeSource,
		//		SigCache:     s.sigCache,
		//		IndexManager: indexManager,
		//		HashCache:    s.hashCache,
	})
	if err != nil {
		return nil, err
	}

	// Search for a FeeEstimator state in the database. If none can be found
	// or if it cannot be loaded, create a new one.
	//	db.Update(func(tx database.Tx) error {
	//		metadata := tx.Metadata()
	//		feeEstimationData := metadata.Get(mempool.EstimateFeeDatabaseKey)
	//		if feeEstimationData != nil {
	// delete it from the database so that we don't try to restore the
	//			// same thing again somehow.
	//			metadata.Delete(mempool.EstimateFeeDatabaseKey)

	// If there is an error, log it and make a new fee estimator.
	//			var err error
	//			s.feeEstimator, err = mempool.RestoreFeeEstimator(feeEstimationData)

	//			if err != nil {
	//				peerLog.Errorf("Failed to restore fee estimator %v", err)
	//			}
	//		}

	//		return nil
	//	})

	// If no feeEstimator has been found, or if the one that has been found
	// is behind somehow, create a new one and start over.
	//	if s.feeEstimator == nil || s.feeEstimator.LastKnownHeight() != s.chain.BestSnapshot().Height {
	//		s.feeEstimator = mempool.NewFeeEstimator(
	//			mempool.DefaultEstimateFeeMaxRollback,
	//			mempool.DefaultEstimateFeeMinRegisteredBlocks)
	//	}

	//		txC := mempool.Config{
	//		Policy: mempool.Policy{
	//			DisableRelayPriority: cfg.NoRelayPriority,
	//			AcceptNonStd:         cfg.RelayNonStd,
	//			FreeTxRelayLimit:     cfg.FreeTxRelayLimit,
	//			MaxOrphanTxs:         cfg.MaxOrphanTxs,
	//			MaxOrphanTxSize:      defaultMaxOrphanTxSize,
	//			MaxSigOpCostPerTx:    blockchain.MaxBlockSigOpsCost / 4,
	//			MinRelayTxFee:        cfg.minRelayTxFee,
	//			MaxTxVersion:         2,
	//			RejectReplacement:    cfg.RejectReplacement,
	//		},
	//		ChainParams:    chainParams,
	//		FetchUtxoView:  s.chain.FetchUtxoView,
	//		BestHeight:     func() int32 { return s.chain.BestSnapshot().Height },
	//		MedianTimePast: func() time.Time { return s.chain.BestSnapshot().MedianTime },
	//		CalcSequenceLock: func(tx *btcutil.Tx, view *blockchain.UtxoViewpoint) (*blockchain.SequenceLock, error) {
	//			return s.chain.CalcSequenceLock(tx, view, true)
	//		},
	//		IsDeploymentActive: s.chain.IsDeploymentActive,
	//		SigCache:           s.sigCache,
	//		HashCache:          s.hashCache,
	//		AddrIndex:          s.addrIndex,
	//		FeeEstimator:       s.feeEstimator,
	//	}
	//s.txMemPool = mempool.New(&txC)

	/*s.syncManager, err = netsync.New(&netsync.Config{
		PeerNotifier:       &s,
		Chain:              s.chain,
		TxMemPool:          s.txMemPool,
		ChainParams:        s.chainParams,
		DisableCheckpoints: cfg.DisableCheckpoints,
		MaxPeers:           cfg.MaxPeers,
		FeeEstimator:       s.feeEstimator,
	})
	if err != nil {
		return nil, err
	}*/

	// Create the mining policy and block template generator based on the
	// configuration options.
	//
	// NOTE: The miner relies on the mempool, so the mempool has to be
	// created before calling the function to create the CPU miner.
	policy := miner.Policy{
		BlockMinWeight:    cfg.BlockMinWeight,
		BlockMaxWeight:    cfg.BlockMaxWeight,
		BlockMinSize:      cfg.BlockMinSize,
		BlockMaxSize:      cfg.BlockMaxSize,
		BlockPrioritySize: cfg.BlockPrioritySize,
		//	TxMinFreeFee:      cfg.minRelayTxFee,
	}
	blockTemplateGenerator := mining.NewBlkTmplGenerator(&policy, s.chainParams, s.chain) //s.chainParams, s.txMemPool, s.chain, s.timeSource,
	//s.sigCache, s.hashCache)

	s.Miner = miner.New(&miner.Config{
		ChainParams:            chainParams,
		BlockTemplateGenerator: blockTemplateGenerator,
		//MiningAddrs:            cfg.miningAddrs,
		//ProcessBlock:   s.syncManager.ProcessBlock,
		//ConnectedCount: s.ConnectedCount,
		//IsCurrent: s.syncManager.IsCurrent,
	})

	// Only setup a function to return new addresses to connect to when
	// not running in connect-only mode.  The simulation network is always
	// in connect-only mode since it is only intended to connect to
	// specified peers and actively avoid advertising and connecting to
	// discovered peers in order to prevent it from becoming a public test
	// network.
	/*	var newAddressFunc func() (net.Addr, error)
		if !cfg.SimNet && len(cfg.ConnectPeers) == 0 {
			newAddressFunc = func() (net.Addr, error) {
				for tries := 0; tries < 100; tries++ {
					addr := s.addrManager.GetAddress()
					if addr == nil {
						break
					}

					// Address will not be invalid, local or unroutable
					// because addrmanager rejects those on addition.
					// Just check that we don't already have an address
					// in the same group so that we are not connecting
					// to the same network segment at the expense of
					// others.
					key := addrmgr.GroupKey(addr.NetAddress())
					if s.OutboundGroupCount(key) != 0 {
						continue
					}

					// only allow recent nodes (10mins) after we failed 30
					// times
					if tries < 30 && time.Since(addr.LastAttempt()) < 10*time.Minute {
						continue
					}

					// allow nondefault ports after 50 failed tries.
					if tries < 50 && fmt.Sprintf("%d", addr.NetAddress().Port) !=
						activeNetParams.DefaultPort {
						continue
					}

					// Mark an attempt for the valid address.
					//		s.addrManager.Attempt(addr.NetAddress())

					//				addrString := addrmgr.NetAddressKey(addr.NetAddress())
					//			return addrStringToNetAddr(addrString)
				}

				return nil, errors.New("no valid connect address")
			}
		}

		// Create a connection manager.
		targetOutbound := defaultTargetOutbound
		if cfg.MaxPeers < targetOutbound {
			targetOutbound = cfg.MaxPeers
		}
		cmgr, err := connmgr.New(&connmgr.Config{
			Listeners:      listeners,
			OnAccept:       s.inboundPeerConnected,
			RetryDuration:  connectionRetryInterval,
			TargetOutbound: uint32(targetOutbound),
			Dial:           btcdDial,
			OnConnection:   s.outboundPeerConnected,
			GetNewAddress:  newAddressFunc,
		})
		if err != nil {
			return nil, err
		}
		s.connManager = cmgr

		// Start up persistent peers.
		permanentPeers := cfg.ConnectPeers
		if len(permanentPeers) == 0 {
			permanentPeers = cfg.AddPeers
		}
		for _, addr := range permanentPeers {
			netAddr, err := addrStringToNetAddr(addr)
			if err != nil {
				return nil, err
			}

			go s.connManager.Connect(&connmgr.ConnReq{
				Addr:      netAddr,
				Permanent: true,
			})
		}
	*/
	/*if !cfg.DisableRPC {
		// Setup listeners for the configured RPC listen addresses and
		// TLS settings.
		rpcListeners, err := setupRPCListeners()
		if err != nil {
			return nil, err
		}
		if len(rpcListeners) == 0 {
			return nil, errors.New("RPCS: No valid listen address")
		}

		s.rpcServer, err = newRPCServer(&rpcserverConfig{
			Listeners:    rpcListeners,
			StartupTime:  s.startupTime,
			ConnMgr:      &rpcConnManager{&s},
			SyncMgr:      &rpcSyncMgr{&s, s.syncManager},
			TimeSource:   s.timeSource,
			Chain:        s.chain,
			ChainParams:  chainParams,
			DB:           db,
			TxMemPool:    s.txMemPool,
			Generator:    blockTemplateGenerator,
			CPUMiner:     s.cpuMiner,
			TxIndex:      s.txIndex,
			AddrIndex:    s.addrIndex,
			CfIndex:      s.cfIndex,
			FeeEstimator: s.feeEstimator,
		})
		if err != nil {
			return nil, err
		}

		// Signal process shutdown when the RPC server requests it.
		go func() {
			<-s.rpcServer.RequestedProcessShutdown()
			shutdownRequestChannel <- struct{}{}
		}()
	}*/

	return &s, nil
}

// Start begins accepting connections from peers.
func (s *server) Start() {
	// Already started?
	if atomic.AddInt32(&s.started, 1) != 1 {
		return
	}

	//srvrLog.Trace("Starting server")

	// Server startup time. Used for the uptime command for uptime calculation.
	s.startupTime = time.Now().Unix()

	// Start the peer handler which in turn starts the address and block
	// managers.
	s.wg.Add(1)
	//go s.peerHandler()

	//if s.nat != nil {
	//	s.wg.Add(1)
	//	go s.upnpUpdateThread()
	//}

	if !cfg.DisableRPC {
		s.wg.Add(1)

		// Start the rebroadcastHandler, which ensures user tx received by
		// the RPC server are rebroadcast until being included in a block.
		//	go s.rebroadcastHandler()

		//	s.rpcServer.Start()
	}
	// Start the CPU miner if generation is enabled.
	if cfg.Generate {
		s.Miner.Start()
	}
}
