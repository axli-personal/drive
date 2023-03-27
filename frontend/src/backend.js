import env from "/env-dev.json";
import axios from "axios";

const site = "http://127.0.0.1:80";

const userService = axios.create({
    baseURL: env.USER_SERVICE_URL,
    withCredentials: true,
});

const driveService = axios.create({
    baseURL: env.DRIVE_SERVICE_URL,
    withCredentials: true,
});

const storageService = axios.create({
    baseURL: env.STORAGE_SERVICE_URL,
    withCredentials: true,
})

export {
    site,
    userService,
    driveService,
    storageService,
}
