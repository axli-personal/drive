import http from 'k6/http';
import { check } from "k6";
import { randomString, statusCheck } from "./utils.js";

export const options = {
    vus: 500,
    duration: "30s"
};

const USER_SERVICE_URL = "http://127.0.0.1:8080";
const DRIVE_SERVICE_URL = "http://127.0.0.1:8081";
const STORAGE_SERVICE_URL = "http://127.0.0.1:8082";

const jsonHeaders = { 'Content-Type': 'application/json' };
const cookies = { 'SID': '' };
const smallFile = open('small.txt', 'b');

export default function () {
    let response;

    const account = randomString(15);
    const username = randomString(15);
    const password = randomString(15);

    const registerBody = JSON.stringify({
        account,
        username,
        password,
    });

    const loginBody = JSON.stringify({
        account,
        password,
    })

    const uploadBody = {
        parent: "Drive",
        file: http.file(smallFile, 'small.txt')
    }

    response = http.post(`${USER_SERVICE_URL}/register`, registerBody, { headers: jsonHeaders });
    check(response, { 'register': statusCheck });

    response = http.post(`${USER_SERVICE_URL}/login`, loginBody, { headers: jsonHeaders });
    check(response, { 'login': statusCheck });

    cookies.SID = response.cookies.SID[0].value;

    response = http.post(`${DRIVE_SERVICE_URL}/drive/create`, null, { headers: jsonHeaders, cookies });
    check(response, { 'create drive': statusCheck });

    response = http.get(`${DRIVE_SERVICE_URL}/drive`, { headers: jsonHeaders, cookies });
    check(response, { 'get drive': statusCheck });

    response = http.post(`${STORAGE_SERVICE_URL}/upload`, uploadBody, { cookies });
    check(response, { 'upload': statusCheck });

    response = http.get(`${DRIVE_SERVICE_URL}/drive`, { headers: jsonHeaders, cookies });
    check(response, { 'get drive': statusCheck });
}
