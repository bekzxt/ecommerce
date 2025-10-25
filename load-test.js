import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
    vus: 100, // 100 виртуальных пользователей
    duration: '40s', // в течение 30 секунд
};

export default function () {
    const url = 'http://localhost:8080/orders';
    const payload = JSON.stringify({
        user_id: "12345",
        items: [
            { product_id: 11, quantity: 2, price: 1500.50 },
            { product_id: 11, quantity: 1, price: 3200.00 }
        ]
    });

    const params = { headers: { 'Content-Type': 'application/json' } };

    const res = http.post(url, payload, params);

    check(res, {
        'status is 200 or 201': (r) => r.status === 200 || r.status === 201,
    });

    // sleep не обязателен, но помогает имитировать реальных пользователей
    sleep(0.5);
}
