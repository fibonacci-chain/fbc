res=$(fbchaincli tx wasm store ./wasm/escrow/artifacts/cw_escrow-aarch64.wasm --fees 0.01fibo --from captain --gas=2000000 -b block -y)
echo "store code..."
echo $res
code_id=$(echo "$res" | jq '.logs[0].events[1].attributes[0].value' | sed 's/\"//g')
res=$(fbchaincli tx wasm instantiate "$code_id" '{"arbiter":"fb19z5mcfz22amvmthldejx7z3k3w4ngj7j85fl2m","end_height":100000,"recipient":"ex190227rqaps5nplhg2tg8hww7slvvquzy0qa0l0"}' --label test1 --admin fb19z5mcfz22amvmthldejx7z3k3w4ngj7j85fl2m --fees 0.001fibo --from captain -b block -y)
contractAddr=$(echo "$res" | jq '.logs[0].events[0].attributes[0].value' | sed 's/\"//g')
echo "instantiate contract..."
echo $res
#fbchaincli tx send fb19z5mcfz22amvmthldejx7z3k3w4ngj7j85fl2m $contractAddr 999fibo --fees 0.01fibo -y -b block
fbchaincli tx wasm execute "$contractAddr" '{"approve":{"quantity":[{"amount":"1","denom":"fibo"}]}}' --amount 888fibo --fees 0.001fibo --from captain -b block -y
