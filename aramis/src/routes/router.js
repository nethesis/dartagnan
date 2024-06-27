import Vue from "vue";
import Router from "vue-router";
import VueAnalytics from "vue-analytics";

import Login from "@/components/Login";
import Callback from "@/components/Callback";
import NotFound from "@/components/NotFound";

import Dashboard from "@/components/Dashboard";
import Alerts from "@/components/Alerts";
import Servers from "@/components/Servers";
import Services from "@/components/Services";
import ServerDetail from "@/components/ServerDetails";
import Profile from "@/components/Profile";

Vue.use(Router);

const router = new Router({
  mode: "history",
  routes: [
    {
      path: "/login",
      name: "Login",
      component: Login
    },
    {
      path: "/callback",
      name: "Callback",
      component: Callback
    },
    {
      path: "/dashboard",
      name: "Dashboard",
      component: Dashboard
    },
    {
      path: "/alerts",
      name: "Alerts",
      component: Alerts
    },
    {
      path: "/servers",
      name: "Servers",
      component: Servers
    },
    {
      path: "/services",
      name: "Services",
      component: Services
    },
    {
      path: "/servers/:id",
      name: "ServerDetail",
      component: ServerDetail
    },
    {
      path: "/profile",
      name: "Profile",
      component: Profile
    },
    {
      path: "/notfound",
      name: "NotFound",
      component: NotFound
    },
    {
      path: "*",
      name: "NotFound",
      component: NotFound
    }
  ]
});

Vue.use(VueAnalytics, {
  id: CONFIG.UA_ANALYTICS,
  router
});

export default router;
