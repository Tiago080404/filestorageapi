import { createRouter, createWebHistory } from "vue-router";
import LoginView from "../views/LoginView.vue";
import HomeView from "../views/HomeView.vue";

const routes = [
  { path: "/login", component: LoginView },
  { path: "/", component: HomeView },
];

export const router = createRouter({
  history: createWebHistory(),
  routes,
});

router.beforeEach(async function (to) {
  try {
    const response = await fetch(`${import.meta.env.VITE_API_URL}api/me`, {
      credentials: "include",
    });
    const loggedIn = response.ok;

    if (!loggedIn && to.path !== "/login") {
      return "/login";
    }

    if (loggedIn && to.path === "/login") {
      return "/";
    }
  } catch (err) {
    return "/login";
  }
});
