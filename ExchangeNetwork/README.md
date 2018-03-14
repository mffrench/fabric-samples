# Créer votre premier réseau Fabric


## Génération PKI

```
$ cd $FLAB
$ ./byfn.sh -m gen_pki
Generating certs for with channel 'BanksCo' and CLI timeout of '10'
Continue (y/n)? y
proceeding ...

##########################################################
##### Generate certificates using cryptogen tool #########
##########################################################
bank1.banksco.com
bank2.banksco.com
```

## Generation du genesis block

```
$ ./byfn.sh -m gen_cha
Generating genesis block for with channel 'BanksCo' and CLI timeout of '10'
Continue (y/n)? y
proceeding ...
##########################################################
#########  Generating Orderer Genesis block ##############
##########################################################
2018-03-14 13:48:20.002 UTC [common/configtx/tool] main -> INFO 001 Loading configuration
2018-03-14 13:48:20.067 UTC [common/configtx/tool] doOutputBlock -> INFO 002 Generating genesis block
2018-03-14 13:48:20.069 UTC [common/configtx/tool] doOutputBlock -> INFO 003 Writing genesis block

#################################################################
### Generating channel configuration transaction 'channel.tx' ###
#################################################################
2018-03-14 13:48:21.363 UTC [common/configtx/tool] main -> INFO 001 Loading configuration
2018-03-14 13:48:21.368 UTC [common/configtx/tool] doOutputChannelCreateTx -> INFO 002 Generating new channel configtx
2018-03-14 13:48:21.369 UTC [common/configtx/tool] doOutputChannelCreateTx -> INFO 003 Writing new channel tx

#################################################################
#######    Generating anchor peer update for Bank1MSP   ##########
#################################################################
2018-03-14 13:48:22.661 UTC [common/configtx/tool] main -> INFO 001 Loading configuration
2018-03-14 13:48:22.667 UTC [common/configtx/tool] doOutputAnchorPeersUpdate -> INFO 002 Generating anchor peer update
2018-03-14 13:48:22.667 UTC [common/configtx/tool] doOutputAnchorPeersUpdate -> INFO 003 Writing anchor peer update

#################################################################
#######    Generating anchor peer update for Bank2MSP   ##########
#################################################################
2018-03-14 13:48:23.931 UTC [common/configtx/tool] main -> INFO 001 Loading configuration
2018-03-14 13:48:23.940 UTC [common/configtx/tool] doOutputAnchorPeersUpdate -> INFO 002 Generating anchor peer update
2018-03-14 13:48:23.940 UTC [common/configtx/tool] doOutputAnchorPeersUpdate -> INFO 003 Writing anchor peer update
```

## Démarrage et interconnection entre les composants

Dans terminator, lancer les commandes suivantes : 

```
$ docker-compose -f docker-compose-cli.yaml up -d
```

Puis vérifiez votre infrastructure : 

````
$ docker ps -a
CONTAINER ID        IMAGE                                     COMMAND             CREATED              STATUS                   PORTS                                              NAMES
4a551023ffa9        hyperledger/fabric-tools:x86_64-1.0.2     "/bin/bash"         About a minute ago   Up About a minute                                                           cliBank2
181127bd9400        hyperledger/fabric-tools:x86_64-1.0.2     "/bin/bash"         About a minute ago   Up About a minute                                                           cliBank1
63cb35a23194        hyperledger/fabric-peer:x86_64-1.0.2      "peer node start"   About a minute ago   Up About a minute        0.0.0.0:7051->7051/tcp, 0.0.0.0:7053->7053/tcp     peer0.bank1.banksco.com
0cc375cf5e60        hyperledger/fabric-peer:x86_64-1.0.2      "peer node start"   About a minute ago   Up About a minute        0.0.0.0:10051->7051/tcp, 0.0.0.0:10053->7053/tcp   peer1.bank2.banksco.com
328c178a796f        hyperledger/fabric-peer:x86_64-1.0.2      "peer node start"   About a minute ago   Up About a minute        0.0.0.0:8051->7051/tcp, 0.0.0.0:8053->7053/tcp     peer1.bank1.banksco.com
cf2d318a1ac6        hyperledger/fabric-orderer:x86_64-1.0.2   "orderer"           About a minute ago   Up About a minute        0.0.0.0:7050->7050/tcp                             orderer.banksco.com
774ee0805ef4        hyperledger/fabric-peer:x86_64-1.0.2      "peer node start"   About a minute ago   Up About a minute        0.0.0.0:9051->7051/tcp, 0.0.0.0:9053->7053/tcp     peer0.bank2.banksco.com
````

### Création du channel

D'abord nous devons nous 'connecter' à l'un des container client. Nous allons commencer par le client `cliBank1` 
```
$ docker exec -it cliBank1 bash
```

Une fois connecté à ce container, nous pouvons lancer les commandes pour créer le channel HLF `bankscochannel`.  
```
# peer channel create -o orderer.banksco.com:7050 -c $CHANNEL_NAME -f ./channel-artifacts/channel.tx --tls true --cafile $ORDERER_TLS_CA
```

*NOTE:* vous pouvez suivre le comportement de cette commande via les logs du container orderer : `docker logs -f orderer.banksco.com`

### Connection des peers sur le channel (JOIN)

Une fois le channel HLF crée, nous allons connecter les peers de l'organisation `Bank1` au channel `bankscochannel`. Nous restons donc 'connecté' au client de l'organisation `Bank1`.

Commencons par connecter le peer `peer0.bank1` au channel : 
```
# export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_BANK1_TLS_CA
# export CORE_PEER_ADDRESS=$PEER0_BANK1_PRADDR
# peer channel join -b bankscochannel.block
```

*NOTE:* depuis la VM et sur une autre console, vous pouvez vérifier la connection au peer0.bank1 vers le channel porté par l'orderer en 
regardant les logs du peer : 

```
$ docker logs -f peer0.bank1.banksco.com
2018-03-14 14:22:44.845 UTC [deliveryClient] StartDeliverForChannel -> DEBU 2d0 This peer will pass blocks from orderer service to other peers for channel bankscochannel
2018-03-14 14:22:44.849 UTC [deliveryClient] connect -> DEBU 2d1 Connected to orderer.banksco.com:7050
2018-03-14 14:22:44.849 UTC [deliveryClient] connect -> DEBU 2d2 Establishing gRPC stream with orderer.banksco.com:7050 ...
2018-03-14 14:22:44.850 UTC [deliveryClient] afterConnect -> DEBU 2d3 Entering
2018-03-14 14:22:44.850 UTC [deliveryClient] RequestBlocks -> DEBU 2d4 Starting deliver with block [1] for channel bankscochannel
```

Puis nous continuons par le peer `peer1.bank1` :
```
# export CORE_PEER_TLS_ROOTCERT_FILE=$PEER1_BANK1_TLS_CA
# export CORE_PEER_ADDRESS=$PEER1_BANK1_PRADDR
# peer channel join -b bankscochannel.block
```

Nous allons maintenant connecter les peers de l'organisation `Bank2` au channel `bankscochannel`. Pour ce faire nous nous 'deconnectons' du client de l'organisation `Bank1` 
et nous nous connectons sur le client de l'organisation `Bank2`:

```
# exit
$ docker exec -it cliBank2 bash
```

A la différence de l'organisation `Bank1` ou nous avons utilisé le client pour créer le channel, sur le client de l'organisation `Bank2` nous allons 'fetch' le channel pour
récupérer le `genesis.block` qui va nous permettre de connecter les peers au channel : 

```
# peer channel fetch newest $CHANNEL_NAME.block -o orderer.banksco.com:7050 -c $CHANNEL_NAME --tls --cafile $ORDERER_TLS_CA
```

Nous pouvons ensuite rejouer les opérations pour les peer de l'organisation `Bank2`. Commencons par le peer `peer0.bank2` : 
```
# export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_BANK2_TLS_CA
# export CORE_PEER_ADDRESS=$PEER0_BANK2_PRADDR
# peer channel join -b bankscochannel.block
```

Puis continuons avec le peer `peer1.bank2` et deconnectons nous : 
```
# export CORE_PEER_TLS_ROOTCERT_FILE=$PEER1_BANK2_TLS_CA
# export CORE_PEER_ADDRESS=$PEER1_BANK2_PRADDR
# peer channel join -b bankscochannel.block
# exit
```

## Installation et instantiation de la chaincode

Nous pouvons à présent déployer la chaincode sur le réseau et donc les peers *endorsers* :

+ peer0.bank1
+ peer0.bank2
+ peer1.bank2

Commencons par l'installation de la chaincode sur `peer0.bank1`:  
```
$ docker exec -it cliBank1 bash
# export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_BANK1_TLS_CA
# export CORE_PEER_ADDRESS=$PEER0_BANK1_PRADDR
# peer chaincode install -n excc -v 1.0 -p github.com/hyperledger/fabric/examples/chaincode/go/exchangeCC-0
# exit
```

Vérifions maintenant que la chaincode a bien été installée : 
```
$ docker exec -ti peer0.bank1.banksco.com /bin/bash
# ls /var/hyperledger/production/chaincodes/ 
excc.1.0
# exit
``` 

Installons maintenant sur sur `peer0.bank2` :
```
$ docker exec -it cliBank2 bash
# export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_BANK2_TLS_CA
# export CORE_PEER_ADDRESS=$PEER0_BANK2_PRADDR
# peer chaincode install -n excc -v 1.0 -p github.com/hyperledger/fabric/examples/chaincode/go/exchangeCC-0
```

Puis sur `peer1.bank2` : 
```
# export CORE_PEER_TLS_ROOTCERT_FILE=$PEER1_BANK2_TLS_CA
# export CORE_PEER_ADDRESS=$PEER1_BANK2_PRADDR
# peer chaincode install -n excc -v 1.0 -p github.com/hyperledger/fabric/examples/chaincode/go/exchangeCC-0
```

Nous pouvons désormais instancier la chaincode excc v1.0 en initialisant les comptes des companies `a` à 100 et `b` à 200. 
Prenons par exemple le peer `peer0.bank2`.
```
# export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_BANK2_TLS_CA
# export CORE_PEER_ADDRESS=$PEER0_BANK2_PRADDR
# peer chaincode instantiate -o orderer.banksco.com:7050 --tls true --cafile $ORDERER_TLS_CA -C $CHANNEL_NAME -n excc -v 1.0 \
-c '{"Args":["init","a","100","b","200"]}' -P "OR ('Bank1MSP.member','Bank2MSP.member')"
# exit
```

*NOTE 1: * nous pouvons vérifier que les peers synchronisent leurs blocks à la fin de l'operation en vérifiant les logs : 

```
$ docker logs peer0.bank2.banksco.com
...
2018-03-14 18:40:24.084 UTC [historyleveldb] Commit -> DEBU 405 Channel [bankscochannel]: Updating history database for blockNo [1] with [1] transactions
2018-03-14 18:40:24.084 UTC [historyleveldb] Commit -> DEBU 406 Channel [bankscochannel]: Updates committed to history database for blockNo [1]
2018-03-14 18:40:24.084 UTC [eventhub_producer] SendProducerBlockEvent -> DEBU 407 Entry
2018-03-14 18:40:24.084 UTC [eventhub_producer] SendProducerBlockEvent -> DEBU 408 Channel [bankscochannel]: Block event for block number [1] contains transaction id: f37164fd6b2a7d3c1fb7d841e62fd3d74aaffb11dccc1a197f0fcd56f3471a72
2018-03-14 18:40:24.084 UTC [eventhub_producer] SendProducerBlockEvent -> INFO 409 Channel [bankscochannel]: Sending event for block number [1]
2018-03-14 18:40:24.084 UTC [eventhub_producer] Send -> DEBU 40a Entry
2018-03-14 18:40:24.084 UTC [eventhub_producer] Send -> DEBU 40b Event processor timeout > 0
2018-03-14 18:40:24.085 UTC [eventhub_producer] Send -> DEBU 40c Event sent successfully
2018-03-14 18:40:24.085 UTC [eventhub_producer] Send -> DEBU 40d Exit
2018-03-14 18:40:24.085 UTC [eventhub_producer] SendProducerBlockEvent -> DEBU 40e Exit

$ docker logs peer0.bank1.banksco.com
...
2018-03-14 18:40:24.082 UTC [lockbasedtxmgr] Commit -> DEBU 348 Updates committed to state database
2018-03-14 18:40:24.085 UTC [historyleveldb] Commit -> DEBU 349 Channel [bankscochannel]: Updating history database for blockNo [1] with [1] transactions
2018-03-14 18:40:24.085 UTC [historyleveldb] Commit -> DEBU 34a Channel [bankscochannel]: Updates committed to history database for blockNo [1]
2018-03-14 18:40:24.085 UTC [eventhub_producer] SendProducerBlockEvent -> DEBU 34b Entry
2018-03-14 18:40:24.085 UTC [eventhub_producer] SendProducerBlockEvent -> DEBU 34c Channel [bankscochannel]: Block event for block number [1] contains transaction id: f37164fd6b2a7d3c1fb7d841e62fd3d74aaffb11dccc1a197f0fcd56f3471a72
2018-03-14 18:40:24.086 UTC [eventhub_producer] SendProducerBlockEvent -> INFO 34d Channel [bankscochannel]: Sending event for block number [1]
2018-03-14 18:40:24.086 UTC [eventhub_producer] Send -> DEBU 34e Entry
2018-03-14 18:40:24.086 UTC [eventhub_producer] Send -> DEBU 34f Event processor timeout > 0
2018-03-14 18:40:24.086 UTC [eventhub_producer] Send -> DEBU 350 Event sent successfully
2018-03-14 18:40:24.086 UTC [eventhub_producer] Send -> DEBU 351 Exit
2018-03-14 18:40:24.086 UTC [eventhub_producer] SendProducerBlockEvent -> DEBU 352 Exit
```  

*NOTE 2: * nous pouvons aussi constater la présence de la chaincode instanciée : 

```
$ docker ps -a
  CONTAINER ID        IMAGE                                                                                                   COMMAND                  CREATED             STATUS              PORTS                                              NAMES
  947d62cf2574        dev-peer0.bank2.banksco.com-excc-1.0-decfd737e1e178040a5e3859c01d55477dcf32bcf7964eb00c086b56b64a1893   "chaincode -peer.add…"   3 minutes ago       Up 3 minutes                                                           dev-peer0.bank2.banksco.com-excc-1.0
  089ec3f79dc2        hyperledger/fabric-tools:x86_64-1.0.2                                                                   "/bin/bash"              9 minutes ago       Up 10 minutes                                                          cliBank1
  9461b215d906        hyperledger/fabric-tools:x86_64-1.0.2                                                                   "/bin/bash"              9 minutes ago       Up 10 minutes                                                          cliBank2
  d3d92a0a665f        hyperledger/fabric-peer:x86_64-1.0.2                                                                    "peer node start"        9 minutes ago       Up 10 minutes       0.0.0.0:10051->7051/tcp, 0.0.0.0:10053->7053/tcp   peer1.bank2.banksco.com
  4a46affbdd64        hyperledger/fabric-orderer:x86_64-1.0.2                                                                 "orderer"                9 minutes ago       Up 10 minutes       0.0.0.0:7050->7050/tcp                             orderer.banksco.com
  0c84be2cc5dc        hyperledger/fabric-peer:x86_64-1.0.2                                                                    "peer node start"        9 minutes ago       Up 10 minutes       0.0.0.0:8051->7051/tcp, 0.0.0.0:8053->7053/tcp     peer1.bank1.banksco.com
  c534b37d5083        hyperledger/fabric-peer:x86_64-1.0.2                                                                    "peer node start"        9 minutes ago       Up 10 minutes       0.0.0.0:7051->7051/tcp, 0.0.0.0:7053->7053/tcp     peer0.bank1.banksco.com
  5508c04b6c27        hyperledger/fabric-peer:x86_64-1.0.2                                                                    "peer node start"        9 minutes ago       Up 10 minutes       0.0.0.0:9051->7051/tcp, 0.0.0.0:9053->7053/tcp     peer0.bank2.banksco.com
```



## Invoke et Query sur la chaincode excc

Nous pouvons désormais faire des requetes sur la chaincode `excc`. Voici un example sur `peer0.bank1` :
```
$ docker exec -it cliBank1 bash
# export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_BANK1_TLS_CA
# export CORE_PEER_ADDRESS=$PEER0_BANK1_PRADDR
# peer chaincode query -C $CHANNEL_NAME -n excc -c '{"Args":["query","a"]}'
...
Query Result: 100
...
```

Nous constatons que du `peer0.bank1`, nous récupérons les information sur le compte de la companie `a` qui est bien 
initialisée avec une valeur à 100. De même nous pouvons constater que la chaincode du `peer0.bank1` a été instancié :
```
$ docker ps -a                                                                                                                                                                                                                        130 ↵
CONTAINER ID        IMAGE                                                                                                   COMMAND                  CREATED             STATUS              PORTS                                              NAMES
d46f123b5c10        dev-peer0.bank1.banksco.com-excc-1.0-d68f0fa1354fdc6d326c0a80417c0c64b6c128b324b7fe81f895f139c307cb8e   "chaincode -peer.add…"   3 minutes ago       Up 3 minutes                                                           dev-peer0.bank1.banksco.com-excc-1.0
947d62cf2574        dev-peer0.bank2.banksco.com-excc-1.0-decfd737e1e178040a5e3859c01d55477dcf32bcf7964eb00c086b56b64a1893   "chaincode -peer.add…"   8 minutes ago       Up 8 minutes                                                           dev-peer0.bank2.banksco.com-excc-1.0
089ec3f79dc2        hyperledger/fabric-tools:x86_64-1.0.2                                                                   "/bin/bash"              15 minutes ago      Up 15 minutes                                                          cliBank1
9461b215d906        hyperledger/fabric-tools:x86_64-1.0.2                                                                   "/bin/bash"              15 minutes ago      Up 15 minutes                                                          cliBank2
d3d92a0a665f        hyperledger/fabric-peer:x86_64-1.0.2                                                                    "peer node start"        15 minutes ago      Up 15 minutes       0.0.0.0:10051->7051/tcp, 0.0.0.0:10053->7053/tcp   peer1.bank2.banksco.com
4a46affbdd64        hyperledger/fabric-orderer:x86_64-1.0.2                                                                 "orderer"                15 minutes ago      Up 15 minutes       0.0.0.0:7050->7050/tcp                             orderer.banksco.com
0c84be2cc5dc        hyperledger/fabric-peer:x86_64-1.0.2                                                                    "peer node start"        15 minutes ago      Up 15 minutes       0.0.0.0:8051->7051/tcp, 0.0.0.0:8053->7053/tcp     peer1.bank1.banksco.com
c534b37d5083        hyperledger/fabric-peer:x86_64-1.0.2                                                                    "peer node start"        15 minutes ago      Up 15 minutes       0.0.0.0:7051->7051/tcp, 0.0.0.0:7053->7053/tcp     peer0.bank1.banksco.com
5508c04b6c27        hyperledger/fabric-peer:x86_64-1.0.2                                                                    "peer node start"        15 minutes ago      Up 15 minutes       0.0.0.0:9051->7051/tcp, 0.0.0.0:9053->7053/tcp     peer0.bank2.banksco.com
```

Nous allons maintenant exécuter notre premiere transaction du compte a vers le compte b (a va payer 10 à b). Cette opération
nécessite une transaction HLF et donc ce que nous appellons une opération d'invocation : 

```
# peer chaincode invoke -o orderer.banksco.com:7050  --tls true --cafile $ORDERER_TLS_CA  -C $CHANNEL_NAME -n excc\
 -c '{"Args":["invoke","a","b","10"]}'
 
 ...
2018-03-14 18:57:05.839 UTC [chaincodeCmd] chaincodeInvokeOrQuery -> DEBU 009 ESCC invoke result: version:1 response:<status:200 message:"OK" > payload:"\n =\303\224\335k\316,\374\365\240_\263\226\222\017\340>E\000R\336\017\332\305z\201;\347\212\007\275\353\022Y\nE\022-\n\004excc\022%\n\007\n\001a\022\002\010\001\n\007\n\001b\022\002\010\001\032\007\n\001a\032\00290\032\010\n\001b\032\003210\022\024\n\004lscc\022\014\n\n\n\004excc\022\002\010\001\032\003\010\310\001\"\013\022\004excc\032\0031.0" endorsement:<endorser:"\n\010Bank1MSP\022\204\006-----BEGIN -----\nMIICHDCCAcKgAwIBAgIQEf7PCXhsHa3pHhbLGcG2wzAKBggqhkjOPQQDAjB1MQsw\nCQYDVQQGEwJVUzETMBEGA1UECBMKQ2FsaWZvcm5pYTEWMBQGA1UEBxMNU2FuIEZy\nYW5jaXNjbzEaMBgGA1UEChMRYmFuazEuYmFua3Njby5jb20xHTAbBgNVBAMTFGNh\nLmJhbmsxLmJhbmtzY28uY29tMB4XDTE4MDMxNDE0MTczNloXDTI4MDMxMTE0MTcz\nNlowXDELMAkGA1UEBhMCVVMxEzARBgNVBAgTCkNhbGlmb3JuaWExFjAUBgNVBAcT\nDVNhbiBGcmFuY2lzY28xIDAeBgNVBAMTF3BlZXIwLmJhbmsxLmJhbmtzY28uY29t\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEDthZXkkuZkINacbZ6BtfbknoUWiV\n7NvH3uVAEJA2lwG8+PnQKJp4am8tLPvjPUxy6zUJUt39r8HKZY0c0OjTZKNNMEsw\nDgYDVR0PAQH/BAQDAgeAMAwGA1UdEwEB/wQCMAAwKwYDVR0jBCQwIoAgrpMC8OQi\nAKQ/J1X2j2mkdAB3c/mkGQDPvQEvxE2WFegwCgYIKoZIzj0EAwIDSAAwRQIhAOV5\nvcUMI6wtHasBCPJAe5M2Ylh1Dc6WV1iy6joVRoxhAiBhiXAw/pTFL/suH47crqou\nNoY/JwXF52UVDOHr950E+Q==\n-----END -----\n" signature:"0D\002 \021\366\343\366\236\363=\353f\2459\220W\203\264\3662C\240\273Tl#\3266\256\366\275\0007>\031\002 \026\362-Dx^\305\371\235\214!;\352\270(\255\257\022\000\350\275,=E\014\356\347`\tn\3345" >
2018-03-14 18:57:05.839 UTC [chaincodeCmd] chaincodeInvokeOrQuery -> INFO 00a Chaincode invoke successful. result: status:200 
 
```

Nous pouvons noter que cette opération d'invocation implique l'écriture d'un nouveau blockchain : 

```
$ docker logs peer0.bank1.banksco.com
...

2018-03-14 18:57:07.885 UTC [stateleveldb] ApplyUpdates -> DEBU 4ed Channel [bankscochannel]: Applying key=[[]byte{0x65, 0x78, 0x63, 0x63, 0x0, 0x62}]
2018-03-14 18:57:07.885 UTC [lockbasedtxmgr] Commit -> DEBU 4ee Updates committed to state database
2018-03-14 18:57:07.885 UTC [historyleveldb] Commit -> DEBU 4ef Channel [bankscochannel]: Updating history database for blockNo [2] with [1] transactions
2018-03-14 18:57:07.885 UTC [historyleveldb] Commit -> DEBU 4f0 Channel [bankscochannel]: Updates committed to history database for blockNo [2]
2018-03-14 18:57:07.885 UTC [eventhub_producer] SendProducerBlockEvent -> DEBU 4f1 Entry
2018-03-14 18:57:07.885 UTC [eventhub_producer] SendProducerBlockEvent -> DEBU 4f2 Channel [bankscochannel]: Block event for block number [2] contains transaction id: ba8b4963434b70c8c594bef27bcb36652e67443bb164aa05b52eb5bb416b2b5c
2018-03-14 18:57:07.885 UTC [eventhub_producer] SendProducerBlockEvent -> INFO 4f3 Channel [bankscochannel]: Sending event for block number [2]
2018-03-14 18:57:07.885 UTC [eventhub_producer] Send -> DEBU 4f4 Entry
2018-03-14 18:57:07.885 UTC [eventhub_producer] Send -> DEBU 4f5 Event processor timeout > 0
2018-03-14 18:57:07.885 UTC [eventhub_producer] Send -> DEBU 4f6 Event sent successfully

$ docker logs peer0.bank2.banksco.com
...

18-03-14 18:57:07.859 UTC [lockbasedtxmgr] Commit -> DEBU 450 Committing updates to state database
2018-03-14 18:57:07.859 UTC [lockbasedtxmgr] Commit -> DEBU 451 Write lock acquired for committing updates to state database
2018-03-14 18:57:07.859 UTC [stateleveldb] ApplyUpdates -> DEBU 452 Channel [bankscochannel]: Applying key=[[]byte{0x65, 0x78, 0x63, 0x63, 0x0, 0x61}]
2018-03-14 18:57:07.859 UTC [stateleveldb] ApplyUpdates -> DEBU 453 Channel [bankscochannel]: Applying key=[[]byte{0x65, 0x78, 0x63, 0x63, 0x0, 0x62}]
2018-03-14 18:57:07.859 UTC [lockbasedtxmgr] Commit -> DEBU 454 Updates committed to state database
2018-03-14 18:57:07.859 UTC [historyleveldb] Commit -> DEBU 455 Channel [bankscochannel]: Updating history database for blockNo [2] with [1] transactions
2018-03-14 18:57:07.859 UTC [historyleveldb] Commit -> DEBU 456 Channel [bankscochannel]: Updates committed to history database for blockNo [2]
2018-03-14 18:57:07.859 UTC [eventhub_producer] SendProducerBlockEvent -> DEBU 457 Entry
2018-03-14 18:57:07.859 UTC [eventhub_producer] SendProducerBlockEvent -> DEBU 458 Channel [bankscochannel]: Block event for block number [2] contains transaction id: ba8b4963434b70c8c594bef27bcb36652e67443bb164aa05b52eb5bb416b2b5c
2018-03-14 18:57:07.860 UTC [eventhub_producer] SendProducerBlockEvent -> INFO 459 Channel [bankscochannel]: Sending event for block number [2]
2018-03-14 18:57:07.860 UTC [eventhub_producer] Send -> DEBU 45a Entry
2018-03-14 18:57:07.860 UTC [eventhub_producer] Send -> DEBU 45b Event processor timeout > 0
2018-03-14 18:57:07.860 UTC [eventhub_producer] Send -> DEBU 45c Event sent successfully
2018-03-14 18:57:07.860 UTC [eventhub_producer] Send -> DEBU 45d Exit
2018-03-14 18:57:07.860 UTC [eventhub_producer] SendProducerBlockEvent -> DEBU 45e Exit 
```

Pour finir nous allons faire une derniere requete pour vérifier l'état du compte a sur le peer `peer1.bank2` :


```
$ docker exec -it cliBank2 bash
# export CORE_PEER_TLS_ROOTCERT_FILE=$PEER1_BANK2_TLS_CA
# export CORE_PEER_ADDRESS=$PEER1_BANK2_PRADDR
# peer chaincode query -C $CHANNEL_NAME -n excc -c '{"Args":["query","a"]}'
...
Query Result: 90
...
```