import { createRouter, createWebHistory } from "vue-router";
import Conn from "@/views/Conn.vue";
import DB from "@/views/DB.vue";

const router = createRouter({
    history: createWebHistory(),
    routes: [
        {
            path: '/conn',
            component: Conn,
        },
        {
            path: '/db',
            component: DB,
        },
    ]
})

export default router