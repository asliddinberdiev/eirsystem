export default defineNuxtRouteMiddleware((to) => {
  const authStore = useAuthStore();

  // Define public routes that don't require authentication
  const publicRoutes = ["/auth/login"];

  // If the user is logged in
  if (authStore.isLoggedIn) {
    // If trying to access a public auth route (like login), redirect to dashboard
    if (publicRoutes.includes(to.path)) {
      return navigateTo("/");
    }
  } else {
    // If the user is NOT logged in and trying to access a protected route
    if (!publicRoutes.includes(to.path)) {
      return navigateTo("/auth/login");
    }
  }
});
