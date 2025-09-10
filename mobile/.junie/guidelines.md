SubTracker Mobile – Development Guidelines

Audience: Advanced developers working on this repository.
Scope: Project-specific guidance for build/configuration, testing (unit/integration/E2E), and development practices. This document captures realities of this codebase as of 2025‑09‑10.

1) Build and Configuration

1.1 Runtime and Tooling
- React Native + Expo SDK 53, RN 0.79.x, React 19.
- TypeScript strict mode; path aliases via tsconfig.json ("@/*" -> src/*).
- NativeWind
- Expo Router 5 (typedRoutes enabled in app.config.js experiments).

1.2 Environment Configuration
- Required env vars (dev and prod):
  - EXPO_PUBLIC_CLERK_PUBLISHABLE_KEY: Clerk publishable key (required for production builds; validated by app.config.js).
  - EXPO_PUBLIC_BACKEND_URL: Backend API base URL (defaults to https://api.subtracker.mistribe.com when missing).
- Optional:
  - EXPO_PUBLIC_DEBUG: "true" to enable extra debug logging.
  - EXPO_PUBLIC_UPDATES_URL, EAS_PROJECT_ID: for OTA/EAS if used.
- Files: .env, .env.development, .env.production, .env.local (gitignored).
- Validation:
  - npm run config:validate prints a report and fails when EXPO_PUBLIC_CLERK_PUBLISHABLE_KEY or EXPO_PUBLIC_BACKEND_URL are missing (we verified this).
- App configuration (app.config.js):
  - Injects env into config.extra: { clerkPublishableKey, backendUrl, environment, enableDebugLogging }.
  - Production build guards: throws if EXPO_PUBLIC_CLERK_PUBLISHABLE_KEY is not set.
  - Clerk plugin is currently commented out; re-enable after resolving dependency issues (search “@clerk/clerk-expo/plugin” in app.config.js).

1.3 Local Run
- Install: npm install
- Set up env:
  - cp .env.development .env.local
  - Edit .env.local to include EXPO_PUBLIC_CLERK_PUBLISHABLE_KEY and EXPO_PUBLIC_BACKEND_URL.
  - Validate: npm run config:validate
- Start:
  - npm start (alias: expo start)
  - Open with Expo Go, iOS Simulator, or Android Emulator.

1.4 Native/EAS Considerations
- Detox E2E requires native binaries under ios/ and android/. This Expo project is currently managed. To generate native projects:
  - npx expo prebuild --platform ios
  - npx expo prebuild --platform android
  - Review and adjust detox.config.js binaryPath/scheme/Gradle tasks if names differ.

1.5 Aliases and Module Resolution
- tsconfig.json provides aliases; Jest mirrors with moduleNameMapper "^@/(.*)$": "<rootDir>/src/$1". Use @/… imports across the codebase for consistency.

2) Testing

2.1 Test Stacks Available
- Unit/Integration: Jest 29 with preset 'react-native'. Setup file at jest.setup.js mocks AsyncStorage, Expo SQLite, SecureStore, and basic RN Appearance/Platform.
- E2E: Detox configured to run via Jest (e2e/jest.config.js). Requires native apps (see 1.4).

2.2 Running Tests
- Full unit suite: npm run test:unit
  - Note: Currently, the full suite contains failing tests depending on not-yet-implemented modules/mocks. Prefer targeted runs while developing.
- Single/targeted unit test (verified):
  - npx jest src/__tests__/currency-converter-smoke.test.ts -i
  - -i (runInBand) avoids concurrency issues with RN mocks.
- Coverage:
  - npm run test:coverage (collects from src/**/*.{ts,tsx}, excluding d.ts and test folders; see jest.config.js).
- Detox:
  - Build (example iOS sim): DETOX_CONFIGURATION=ios.sim.debug npm run test:e2e:build
  - Test: DETOX_CONFIGURATION=ios.sim.debug npm run test:e2e
  - Adjust DETOX_CONFIGURATION to one of: ios.sim.debug, ios.sim.release, android.emu.debug, etc. See detox.config.js.

2.3 Adding New Unit Tests
- Location: Either src/**/__tests__/**/*.test.ts[x] or co-located *.test.ts[x]. Jest testMatch covers both patterns.
- Imports: Prefer @/ aliases for source modules (Jest maps to <rootDir>/src).
- React Native/Expo Modules: If your new test pulls in RN TurboModules (DevMenu, etc.), extend jest.setup.js to mock them or isolate logic in pure TS services to avoid native bindings in unit tests.
- Transform ignore for RN libs: jest.config.js allows (expo|@expo|expo-constants|@react-native|react-native). If you add a library that ships untranspiled code under node_modules, include it in transformIgnorePatterns as needed.

2.4 React Component Tests
- @testing-library/react-native is installed. When adding component tests:
  - Import from '@testing-library/react-native'.
  - Mock RN/Expo modules in jest.setup.js as needed (AsyncStorage/SQLite/SecureStore already mocked).
  - Avoid rendering components that require native navigation containers or Clerk until appropriate mocks/providers are added.

3) Additional Development Information

3.1 Code Style and Linting
- ESLint configured via eslint-config-expo (eslint.config.js). Run npm run lint.
- Prefer functional components, hooks, and file-based routing under app/ with Expo Router conventions.

3.2 Project Structure (feature-first, core foundations)
- app/ — File-based routes (Expo Router): route groups like (tabs), (auth), (modals); _layout.tsx; +not-found.tsx
- src/features/* — Feature-owned code (each feature contains components, hooks, services, store, screens, types, index.ts)
- src/core/* — Cross-cutting foundations (api client factory, config, db/migrations, state, ui primitives, utils, analytics, navigation)
- src/lib/* — Third-party adapters/wrappers (e.g., auth adapter, React Query config)
- src/test/* — Shared test utilities and mocks
- assets/* — Images, fonts, icons, animations
- scripts/* — Project scripts (migrations, generators, env validation)

3.3 Configuration Scripts
- npm run config:validate – validates env presence/format (we verified behavior).
- npm run config:inspect – calls src/config/ConfigInspector.ts directly via Node. If your Node runtime cannot import .ts, use ts-node/register or compile first.
- npm run config:template:dev|prod – prints an env template (same note about .ts in Node applies).

3.4 Clerk Integration Status
- app.config.js has Clerk Expo plugin commented out to avoid a dependency issue. Until re-enabled, use Clerk features guarded in code; ensure EXPO_PUBLIC_CLERK_PUBLISHABLE_KEY is set when you re-enable the plugin.

3.5 Currency/Analytics Notes
- CurrencyConverter supports cached and inverse rates, approximate flag, and formatting.
- When working offline, provide fallback rates and handle null from converter.convert.

3.6 Testing Pitfalls and Remedies
- RN TurboModules (e.g., DevMenu) can surface when importing 'react-native' directly in tests. Prefer testing pure TS services or expand jest.setup.js with targeted mocks.
- Use -i / --runInBand for flaky suites interacting with global mocks.
- If you add new RN packages, you may need to extend transformIgnorePatterns in jest.config.js to ensure they’re transformed.

4) References
- Design: .kiro/specs/subtracker-mobile-app/design.md
- Requirements: .kiro/specs/subtracker-mobile-app/requirements.md
- README: ./README.md
- Detox config: ./detox.config.js and ./e2e/jest.config.js
- Jest config: ./jest.config.js and ./jest.setup.js
- TS config: ./tsconfig.json

5) Verified Commands (as of this commit)
- npm run config:validate – executed; produced errors for missing EXPO_PUBLIC_* variables as expected in development when not set.
- npx jest src/__tests__/currency-converter-smoke.test.ts -i – executed and PASSED, verifying the Jest stack and module resolution.

6) Do no modify this files in this folders
- `services/api`
- `app-example`

- End of guidelines.
