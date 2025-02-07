import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
    vus: 10,
    duration: '15s',
};

function getRandomIP() {
    const o = Math.floor(Math.random() * 256);

    return `83.150.59.250.${o}`;
}

export default function () {
    const baseURL = __ENV.API_URL || 'http://localhost:5000';
    const payload = JSON.stringify({
        "timestamp": "2020-06-24T15:27:00.123456Z",
        "ip": getRandomIP(),
        "url": "http://example.com",
    });
    const params = {
        headers: {
            'Content-Type': 'application/json',
        },
    };
    const res = http.post(`${baseURL}/logs`, payload, params);

    check(res, {
        "Post status is 200": (r) => res.status === 200,
    });

    sleep(1);
}
