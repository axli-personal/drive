import { createRouter, createWebHistory } from "vue-router"
import Home from "/src/views/Home.vue"

const routes = [
    {
        path: "/register",
        component: () => import("/src/views/identity/Register.vue"),
    },
    {
        path: "/login",
        component: () => import("/src/views/identity/Login.vue"),
    },
    {
        path: "/drive/plan",
        component: () => import("/src/views/drive/Plan.vue")
    },
    {
        path: "/drive/my-drive",
        alias: "/",
        component: () => import("/src/views/drive/Drive.vue")
    },
    {
        path: "/drive/folders/:folderId",
        component: () => import("/src/views/drive/Folder.vue")
    },
    {
        path: "/drive/files/binary/:fileId",
        component: () => import("/src/views/drive/BinaryFile.vue"),
    },
    {
        path: "/drive/files/text/:fileId",
        component: () => import("/src/views/drive/TextFile.vue"),
    },
    {
        path: "/drive/files/markdown/:fileId",
        component: () => import("/src/views/drive/MarkdownFile.vue"),
    },
    {
        name: "ArticleSearch",
        path: "/drive/search",
        component: () => import("/src/views/drive/Search.vue")
    },
    {
        path: "/identity/dashboard",
        component: () => import("/src/views/identity/Dashboard.vue"),
    },
    {
        path: "/identity/logout",
        component: () => import("/src/views/identity/Logout.vue"),
    }
];

// the class name for active router-link is "link-match" and "link-same".
export default createRouter({
    history: createWebHistory(),
    linkActiveClass: "link-match",
    linkExactActiveClass: "link-same",
    routes: routes,
});
