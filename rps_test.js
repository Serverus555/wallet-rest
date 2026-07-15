import http from "k6/http";
import { check } from "k6";

export const options = {
    scenarios: {
        wallet_load: {
            executor: "constant-vus",
            vus: 100,
            duration: "5s",
        },
    },
};

const id = "00000000-0000-0000-0000-000000000000";
const url = "http://app:8080/api/v1";
let params = { headers: { "Content-Type": "application/json" } };

export default function () {
    let json = { walletId: id, amount: Math.floor(Math.random() * 1000) + 1 }
    let res;

    switch (Math.floor(Math.random() * 3)) {
        case 0:
            json.operationType = "DEPOSIT"
            res = http.post(`${url}/wallet`, JSON.stringify(json), params);
            check(res, { "deposit": (r) => r.status === 200 || r.status === 400});
            break;
        case 1:
            json.operationType = "WITHDRAW"
            res = http.post(`${url}/wallet`, JSON.stringify(json), params);
            check(res, { "withdraw": (r) => r.status === 200 || r.status === 400 });
            break;
        case 2:
            res = http.get(`${url}/wallets/${id}`);
            check(res, { "balance": (r) => r.status === 200 });
            break;
    }
}