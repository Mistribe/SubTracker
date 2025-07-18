# Authentication Implementation

## Overview

This document describes the authentication implementation for the Recurrent Payment Tracker web application. The authentication system uses Kinde for user authentication and implements protected routes with redirection for unauthenticated users.

## Components

### Authentication Provider

The application uses the `KindeProvider` from `@kinde-oss/kinde-auth-react` to provide authentication context to the entire application. The provider is configured in `App.tsx` with the following parameters:

- `clientId`: The Kinde client ID
- `domain`: The Kinde domain URL
- `redirectUri`: The URL to redirect to after successful authentication (set to the application's root URL)
- `logoutUri`: The URL to redirect to after logout (set to the application's root URL)

### Protected Route Component

The `ProtectedRoute` component (`src/components/ProtectedRoute.tsx`) is a wrapper component that:

1. Checks if the user is authenticated using the `useKindeAuth` hook
2. Redirects unauthenticated users to the sign-in page
3. Shows a loading indicator while checking authentication status
4. Renders its children only if the user is authenticated

### Sign-In Page

The `SignIn` component (`src/pages/SignIn.tsx`) provides a user interface for authentication with:

1. A sign-in button that calls Kinde's `login()` method
2. A sign-up button that calls Kinde's `register()` method
3. Automatic redirection to the dashboard if the user is already authenticated

### Routing System

The routing system (`src/Routes.tsx`) uses `react-router-dom` to define routes and handle navigation:

1. Public routes (like the home page and sign-in page) are accessible to all users
2. Protected routes (like the dashboard) are wrapped with the `ProtectedRoute` component
3. A catch-all route redirects to either the dashboard or sign-in page based on authentication status

## Authentication Flow

1. When a user tries to access a protected route (e.g., `/dashboard`), the `ProtectedRoute` component checks if they are authenticated
2. If not authenticated, they are redirected to the sign-in page (`/signin`)
3. On the sign-in page, the user can choose to sign in or create a new account
4. After successful authentication, the user is redirected to the dashboard
5. If an authenticated user tries to access the sign-in page, they are automatically redirected to the dashboard

## Testing

To test the authentication flow:

1. Start the application with `npm run dev`
2. Try accessing the dashboard directly (http://localhost:5173/dashboard)
   - You should be redirected to the sign-in page if not authenticated
3. Sign in or create an account
   - You should be redirected to the dashboard after successful authentication
4. Try accessing the sign-in page after authentication
   - You should be automatically redirected to the dashboard

## Future Improvements

1. Add a navigation bar with a logout button
2. Implement user profile management
3. Add more granular role-based access control
4. Enhance the authentication UI with more branding and customization