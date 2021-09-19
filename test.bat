rd /s /q tmp
md tmp\blocks
main.exe createwallet -nickname Leocao
main.exe createwallet -nickname Krad
main.exe listwallets
main.exe createblockchain -nickname Leocao
main.exe blockchaininfo
main.exe walletinfo -nickname Leocao
main.exe walletinfo -nickname Krad
main.exe send -from Leocao -to Krad -amount 30
main.exe mine
main.exe blockchaininfo
main.exe walletinfo -nickname Leocao
main.exe walletinfo -nickname Krad



