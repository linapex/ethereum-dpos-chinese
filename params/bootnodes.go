
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 12:09:45</date>
//</624342662119428096>


package params

//mainnetbootnodes是运行在上的p2p引导节点的enode URL
//以太坊主网络。
var MainnetBootnodes = []string{
//EASTUM基金会Boo节点
"enode://A979FB575495B8D6DB44F75317D0F4622BF4C2A3365D6AF7C284339968EEF29B69AD0DCE72A4D8DB5EBB4968DE0E3BEC910127F134779FBCB0CB6D3331163C@52.16.188.185:30303“，/ie
"enode://3f1d12044546b76342d59d4a0532c14b85aa669704bfe1f864fe079415aa2c02d743e03218e57a3fb94523ab54032871a6c51b2cc5514cb7c7e35b3ed099@13.93.211.84:30303“，//美国西部
"enode://78DE8A0916848093C73790 EAD81D1928BEC737D565119932B98C6B100D944B7A95E94F847F689FC723399D2E31129D182F7EF3863F2B4C820ABF3AB272234D@191.235.84.50:30303“，/br
"enode://1588f8aab45f6d19c6cbf4a089c267054a8da11978a2f90dbf6a502a4a3bab80d288afdb7ec0ef6d92de563767f3b1ea9e8e334ca711e9f8e2df5a0385e6e6@13.75.154.138:30303“，/au
"enode://1118980BF48B0A3640BDBA04E0FE78B1ADD18E1CD99BF22D53DAAC1FD9972AD650DF52176E7C7D89D1114CFEF2BC23A2959AA4998A46AFCF7D91809F0855082@52.74.57.123:30303“，//SG

//EcUnm基金会C++节点
"enode://979B7FA28FEEB35A474166A16076F1943202CB72B6AF70D327F053E248BAB9BA81760F39D0701EF1D8F89CC1FBDB2CACBA0710A12CD5314D5E0C9021AA3637F9@5.1.83.226:30303“，//de
}

//testnetbootnodes是运行在
//Ropsten测试网络。
var TestnetBootnodes = []string{
"enode://30b7ab30a01c124a6cca36863ece12c4f5fa68e3ba9b0b51407ccc002eed3b3102d20a88f1c1d3c154e2449317b8ef95090e7b312d5cc39354f86d5d606@52.176.7.10:30303“，//美国天蓝色Geth
"enode://865A63255B3BB68023B6BFD5095118fcc13e79dcf014fe4e47e065c350c7cc72af2e53eff895f11ba1bbb6a2b3371c116ee870f26618eadfc2e78aa7349c@52.176.100.77:30303“，//us azure parity
"enode://6332792c4a00e4ee0926ed89e0d27ef985424d97b6a45bf0f23e51f0dcb5e66b875777506458aea7af6f9e4fb69f43f3778ee73c81ed9d34c51c4b16b0f@52.232.243.152:30303“，//奇偶校验
"enode://94C15D1B9E2FE7CE56E458B9A3B672EF1894DDDD0C6F247E0F1D3487F52B66208FB4AEB8179FCE6E3749EA93ED147C37976D67AF557508D199D9594C35F09@192.81.208.223:30303“，/@gpip
}

//rinkebybootnodes是运行在
//Rinkeby测试网络。
var RinkebyBootnodes = []string{
"enode://a24ac7c54844ef4ed0c5eb2d36620ba4e4aa13b8c84684e1b4aab0cebea2ae45cb4d375b77eab5616d34bfbd3c1a833fc51296ff084b770b94fb9028c4d25ccf@52.169.42.101:30303“，/ie
"enode://343149e4feefa15d882d9fe4ac7d88f885bd05ebb735e547f12e12080a9fa07c8014ca6fd7f3731488102fe5e3411f8509cf0b7de3f5b44339c9f25e87cb8@52.3.158.184:30303“，//infura
"enode://B6B28890B006743680C52E64E0D16DB57F28124885595FA03A562BE1D2BF0F3A1DA297D56B13DA25FB992888FD556D4C1A27B1F1F39D531BDE7DE1921C90061CC6@159.89.28.211:30303“，//Akasha
}

//discoveryv5bootnodes是用于
//实验性RLPX v5主题发现网络。
var DiscoveryV5Bootnodes = []string{
"enode://06051A5573C81934C9554EF2898EB13B33A34B94CF36B202B69FDE139CA17A85051979867720D4BDAE4323D4943DF9AEB6643633AA65E0BE843659795007A@35.177.226.168:30303“，
"enode://0cc5f5ffb5d9098c8b8c62325f3797f56509bff942704687b6530992ac706e2cb946b90a34f19548cd3c7bacccaeaa354531e5983c7d1bc0dee16ce4b6440b@40.118.3.223:30304“，
"enode://1c7a64d76c034b0418c004af2f67c50e36a3be60b5e4790bdac0439d1603469a85fad36f2473c9a80eb043ae6036df905fa28f1ff614c3e5dc34f15dcd2dc@40.118.3.223:30306“，
"enode://85c85d7143ae8bb96924f2b54f1b3e70d8c4d367af305325d30a61385a432f247d2c75c45c45c45c45c25c6b4a6035060d072d7f5b35dd1d4c45f76941f62a4f83b6e75daaf@40.118.3.223:30307“，
}

