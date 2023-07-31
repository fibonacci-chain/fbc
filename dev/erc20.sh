res=$(fbchaincli tx wasm store ./wasm/erc20/artifacts/cw_erc20-aarch64.wasm --fees 0.01fibo --from local1 --gas=2000000 -b block -y)
code_id=$(echo "$res" | jq '.logs[0].events[1].attributes[0].value' | sed 's/\"//g')
res=$(fbchaincli tx wasm instantiate "$code_id" '{"decimals":10,"initial_balances":[{"address":"fb19z5mcfz22amvmthldejx7z3k3w4ngj7j85fl2m","amount":"100000000"}],"name":"my test token", "symbol":"MTT"}' --label test1 --admin fb19z5mcfz22amvmthldejx7z3k3w4ngj7j85fl2m --fees 0.001fibo --from local1 -b block -y)
contractAddr=$(echo "$res" | jq '.logs[0].events[0].attributes[0].value' | sed 's/\"//g')
fbchaincli tx wasm execute "$contractAddr" '{"transfer":{"amount":"100","recipient":"fb1muxe2dpldu656f3jq052xe0z0har6cqkerqz7c"}}' --fees 0.001fibo --from local1 -b block -y

echo " ========================================================== "
echo "## show all codes uploaded ##"
fbchaincli query wasm list-code

echo " ========================================================== "
echo "## show contract info by contract addr ##"
fbchaincli query wasm contract "$contractAddr"

echo " ========================================================== "
echo "## show contract update history by contract addr ##"
fbchaincli query wasm contract-history "$contractAddr"

echo " ========================================================== "
echo "## query contract state by contract addr ##"
echo "#### all state"
fbchaincli query wasm contract-state all "$contractAddr"
echo "#### raw state"
fbchaincli query wasm contract-state raw "$contractAddr" 0006636F6E666967636F6E7374616E7473
echo "#### smart state"
fbchaincli query wasm contract-state smart "$contractAddr" '{"balance":{"address":"fb19z5mcfz22amvmthldejx7z3k3w4ngj7j85fl2m"}}'
fbchaincli query wasm contract-state smart "$contractAddr" '{"balance":{"address":"fb1muxe2dpldu656f3jq052xe0z0har6cqkerqz7c"}}'